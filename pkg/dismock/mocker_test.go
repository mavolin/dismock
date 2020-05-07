package dismock

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testHandler = Handler{
	Name: "test",
	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}),
}

// tests the Server started in New.
func TestMocker_Server(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, s := New(t)

		expect := discord.Channel{
			ID: 123,
		}

		m.Channel(expect)

		actual, err := s.Channel(123)
		require.NoError(t, err)

		assert.Equal(t, expect, *actual)
	})

	t.Run("unhandled path", func(t *testing.T) {
		tMock := new(testing.T)

		m, _ := New(tMock)

		url := "https://" + m.server.Listener.Addr().String() + "/unhandled/path"

		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		client.Get(url)

		m.Eval()

		assert.True(t, tMock.Failed())
	})

	t.Run("unhandled method", func(t *testing.T) {
		tMock := new(testing.T)

		m, _ := New(tMock)

		m.handlers["/handled/path"] = make(map[string][]Handler)

		m.handlers["/handled/path"][http.MethodPost] = append(m.handlers["/handled/path"][http.MethodPost], testHandler)

		url := "https://" + m.server.Listener.Addr().String() + "/handled/path"

		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		client.Get(url)

		assert.True(t, tMock.Failed())
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("only handler and only method", func(t *testing.T) {
			m, _ := New(t)

			m.handlers["/path"] = make(map[string][]Handler)

			m.handlers["/path"][http.MethodGet] = append(m.handlers["/path"][http.MethodGet], testHandler)

			url := "https://" + m.server.Listener.Addr().String() + "/path"

			client := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}

			_, err := client.Get(url)
			require.NoError(t, err)

			assert.Len(t, m.handlers, 0)
		})

		t.Run("only handler multiple methods", func(t *testing.T) {
			m, _ := New(t)

			m.handlers["/path"] = make(map[string][]Handler)

			m.handlers["/path"][http.MethodGet] = append(m.handlers["/path"][http.MethodGet], testHandler)
			m.handlers["/path"][http.MethodPost] = append(m.handlers["/path"][http.MethodPost], testHandler)

			url := "https://" + m.server.Listener.Addr().String() + "/path"

			client := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}

			_, err := client.Get(url)
			require.NoError(t, err)

			assert.Len(t, m.handlers, 1)
			assert.Len(t, m.handlers["/path"], 1)

			_, ok := m.handlers["/path"][http.MethodPost]
			assert.True(t, ok)
		})

		t.Run("multiple handlers", func(t *testing.T) {
			m, _ := New(t)

			m.handlers["/path"] = make(map[string][]Handler)

			m.handlers["/path"][http.MethodGet] = append(m.handlers["/path"][http.MethodGet], testHandler, testHandler)

			url := "https://" + m.server.Listener.Addr().String() + "/path"

			client := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			}

			_, err := client.Get(url)
			require.NoError(t, err)

			assert.Len(t, m.handlers, 1)
			assert.Len(t, m.handlers["/path"], 1)
			assert.Len(t, m.handlers["/path"][http.MethodGet], 1)
		})
	})
}

func TestMocker_Mock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, _ := New(t)

		invoked := false

		method := http.MethodPost
		path := "/path/123"
		f := func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			invoked = true
		}

		m.Mock("handler", method, path, f)

		h, ok := m.handlers["/api/v6"+path][method]
		require.True(t, ok)

		h[0].ServeHTTP(new(httptest.ResponseRecorder), new(http.Request))
		assert.True(t, invoked)
	})

	t.Run("nil support", func(t *testing.T) {
		m, _ := New(t)

		method := http.MethodPost
		path := "/path/123"

		m.Mock("handler", method, path, nil)

		h, ok := m.handlers["/api/v6"+path][method]
		require.True(t, ok)

		r := new(httptest.ResponseRecorder)

		h[0].ServeHTTP(r, new(http.Request)) // if unsuccessful, this would panic

		assert.Equal(t, http.StatusNoContent, r.Code)
	})
}

func TestMocker_Clone(t *testing.T) {
	m1, _ := New(new(testing.T))

	m1.handlers["path"] = map[string][]Handler{
		http.MethodGet: {},
	}

	m2, _ := m1.Clone(t)

	assert.Equal(t, m1.handlers, m2.handlers)

	m1.handlers["path2"] = map[string][]Handler{
		http.MethodPatch: {},
	}

	assert.NotEqual(t, m1.handlers, m2.handlers)
}

func TestMocker_Eval(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tMock := new(testing.T)

		m, _ := New(tMock)

		m.Eval()

		// although we would never reach this point if tMock.Failed == true, we leave it in for clarity
		assert.False(t, tMock.Failed())
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m, _ := New(tMock)

		m.handlers["path"] = map[string][]Handler{
			"request0": {},
		}

		go func() { // prevent failure caused by t.Fatal's runtime.Goexit
			m.Eval()
		}()

		// as the started goroutine is terminated, this is necessary to ensure execution
		time.Sleep(500 * time.Microsecond)

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_genUninvokedMsg(t *testing.T) {
	t.Run("singular", func(t *testing.T) {
		m, _ := New(new(testing.T))

		expect := "path:\n\trequest0: 1 uninvoked handler"

		m.handlers["path"] = map[string][]Handler{
			http.MethodGet: {
				{
					Name:    "request0",
					Handler: nil,
				},
			},
		}

		assert.Equal(t, expect, m.genUninvokedMsg())
	})

	t.Run("plural", func(t *testing.T) {
		m, _ := New(new(testing.T))

		expect := "path:\n\trequest0: 2 uninvoked handlers"

		m.handlers["path"] = map[string][]Handler{
			http.MethodGet: {
				{
					Name:    "request0",
					Handler: nil,
				},
				{
					Name:    "request0",
					Handler: nil,
				},
			},
		}

		assert.Equal(t, expect, m.genUninvokedMsg())
	})
}
