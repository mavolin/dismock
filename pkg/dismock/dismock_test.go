package dismock

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
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
func TestMocker_New(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("session", func(t *testing.T) {
			m, s := NewSession(t)

			expect := discord.Channel{ID: 123, VideoQualityMode: discord.AutoVideoQuality}
			m.Channel(expect)

			actual, err := s.Channel(123)
			require.NoError(t, err)

			assert.Equal(t, expect, *actual)
		})

		t.Run("state", func(t *testing.T) {
			m, s := NewState(t)

			expect := discord.Channel{ID: 123, VideoQualityMode: discord.AutoVideoQuality}

			m.Channel(expect)

			actual, err := s.Channel(123)
			require.NoError(t, err)

			assert.Equal(t, expect, *actual)
		})
	})

	t.Run("unhandled path", func(t *testing.T) {
		tMock := new(testing.T)
		m := New(tMock)

		url := "https://" + m.Server.Listener.Addr().String() + "/unhandled/path"

		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		_, _ = client.Get(url) // this will error

		//goland:noinspection ALL
		m.eval()

		assert.True(t, tMock.Failed())
	})

	t.Run("unhandled method", func(t *testing.T) {
		tMock := new(testing.T)

		m := New(tMock)

		m.handlers["/handled/path"] = make(map[string][]Handler)

		m.handlers["/handled/path"][http.MethodPost] = append(m.handlers["/handled/path"][http.MethodPost], testHandler)

		url := "https://" + m.Server.Listener.Addr().String() + "/handled/path"

		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		_, _ = client.Get(url) // this will error

		assert.True(t, tMock.Failed())
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("only handler and only method", func(t *testing.T) {
			m := New(t)

			m.handlers["/path"] = make(map[string][]Handler)

			m.handlers["/path"][http.MethodGet] = append(m.handlers["/path"][http.MethodGet], testHandler)

			url := "https://" + m.Server.Listener.Addr().String() + "/path"

			client := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}

			_, err := client.Get(url)
			require.NoError(t, err)

			assert.Empty(t, m.handlers)
		})

		t.Run("only handler multiple methods", func(t *testing.T) {
			m := New(t)

			m.handlers["/path"] = make(map[string][]Handler)

			m.handlers["/path"][http.MethodGet] = append(m.handlers["/path"][http.MethodGet], testHandler)
			m.handlers["/path"][http.MethodPost] = append(m.handlers["/path"][http.MethodPost], testHandler)

			url := "https://" + m.Server.Listener.Addr().String() + "/path"

			client := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}

			_, err := client.Get(url)
			require.NoError(t, err)

			assert.Len(t, m.handlers, 1)
			assert.Len(t, m.handlers["/path"], 1)

			_, ok := m.handlers["/path"][http.MethodPost]
			assert.True(t, ok)

			m.Close() // prevent m.eval from failing
		})

		t.Run("multiple handlers", func(t *testing.T) {
			m := New(t)

			m.handlers["/path"] = make(map[string][]Handler)

			m.handlers["/path"][http.MethodGet] = append(m.handlers["/path"][http.MethodGet], testHandler, testHandler)

			url := "https://" + m.Server.Listener.Addr().String() + "/path"

			client := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}

			_, err := client.Get(url)
			require.NoError(t, err)

			assert.Len(t, m.handlers, 1)
			assert.Len(t, m.handlers["/path"], 1)
			assert.Len(t, m.handlers["/path"][http.MethodGet], 1)

			m.Close() // prevent m.eval from failing
		})
	})
}

func TestMocker_MockAPI(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := New(t)

		invoked := false

		method := http.MethodGet
		path := "path/123"
		f := func(w http.ResponseWriter, r *http.Request, t *testing.T) {
			invoked = true
		}

		m.MockAPI("handler", method, path, f)

		_, ok := m.handlers["/api/v"+api.Version+"/"+path][method]
		require.True(t, ok)

		_, err := m.Client.Get("https://" + m.Server.Listener.Addr().String() + "/api/v" + api.Version + "/" + path)
		require.NoError(t, err)

		assert.True(t, invoked)
	})

	t.Run("nil handler", func(t *testing.T) {
		m := New(t)

		method := http.MethodGet
		path := "path/123"

		m.MockAPI("handler", method, path, nil)

		h, ok := m.handlers["/api/v"+api.Version+"/"+path][method]
		require.True(t, ok)

		r := new(httptest.ResponseRecorder)
		h[0].ServeHTTP(r, new(http.Request)) // if unsuccessful, this would panic

		assert.Equal(t, http.StatusNoContent, r.Code)

		m.Close() // prevent m.eval from failing
	})
}

func TestMocker_Clone(t *testing.T) {
	m1 := New(t)
	m1.handlers["path"] = map[string][]Handler{http.MethodGet: {}}

	m2 := m1.Clone(t)

	assert.NotEqual(t, m1.Client, m2.Client, "clients are the same")
	assert.Equal(t, m1.handlers, m2.handlers)

	m1.handlers["path2"] = map[string][]Handler{http.MethodPatch: {}}

	assert.NotEqual(t, m1.handlers, m2.handlers)

	m2.Close() // prevent m2.eval from failing
}

func TestMocker_CloneSession(t *testing.T) {
	m1 := New(new(testing.T))

	m1.handlers["path"] = map[string][]Handler{
		http.MethodGet: {},
	}

	m2, _ := m1.CloneSession(t)

	assert.NotEqual(t, m1.Client, m2.Client, "clients are the same")
	assert.Equal(t, m1.handlers, m2.handlers)

	m1.handlers["path2"] = map[string][]Handler{http.MethodPatch: {}}

	assert.NotEqual(t, m1.handlers, m2.handlers)

	m2.Close() // prevent m2.eval from failing
}

func TestMocker_CloneState(t *testing.T) {
	m1 := New(new(testing.T))

	m1.handlers["path"] = map[string][]Handler{
		http.MethodGet: {},
	}

	m2, _ := m1.CloneState(t)

	assert.NotEqual(t, m1.Client, m2.Client, "clients are the same")
	assert.Equal(t, m1.handlers, m2.handlers)

	m1.handlers["path2"] = map[string][]Handler{http.MethodPatch: {}}

	assert.NotEqual(t, m1.handlers, m2.handlers)

	m2.Close() // prevent m2.eval from failing
}

func TestMocker_deepCopyHandlers(t *testing.T) {
	m1 := New(t)

	m1.handlers["path"] = map[string][]Handler{http.MethodGet: {}}

	cp := m1.deepCopyHandlers()
	assert.Equal(t, m1.handlers, cp)

	cp["path2"] = map[string][]Handler{http.MethodPatch: {}}
	assert.NotEqual(t, m1.handlers, cp)

	m1.Close() // prevent m1.eval from failing
}

func TestMocker_Eval(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tMock := new(testing.T)

		m := New(tMock)

		//goland:noinspection ALL
		m.eval()

		// although we would never reach this point if tMock.Failed == true, we leave it in for clarity
		assert.False(t, tMock.Failed())
	})

	t.Run("closed", func(t *testing.T) {
		m := New(t)

		m.Mock("", http.MethodGet, "/path", func(http.ResponseWriter, *http.Request, *testing.T) {})
		m.Close()

		// this is here for clarity, but would obviously get called
		// automatically, as it is part of t's cleanup
		//goland:noinspection ALL
		m.eval()
	})

	t.Run("failure", func(t *testing.T) {
		tMock := new(testing.T)

		m := New(tMock)

		m.handlers["path"] = map[string][]Handler{"request0": {}}

		c := make(chan struct{})

		go func() { // prevent failure caused by t.Fatal's runtime.Goexit
			defer func() { c <- struct{}{} }()

			//goland:noinspection ALL
			m.eval()
		}()

		<-c

		assert.True(t, tMock.Failed())
	})
}

func TestMocker_genUninvokedMsg(t *testing.T) {
	t.Run("singular", func(t *testing.T) {
		m := New(new(testing.T))

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
		m := New(new(testing.T))

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
