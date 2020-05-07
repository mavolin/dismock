package dismock

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		time.Sleep(100 * time.Microsecond)

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
