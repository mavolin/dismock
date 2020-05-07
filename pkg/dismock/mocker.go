package dismock

import (
	"context"
	"crypto/tls"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/session"
	"github.com/diamondburned/arikawa/utils/httputil/httpdriver"
	"github.com/stretchr/testify/require"
)

type (
	// Mocker handles the mocking of arikawa's API calls.
	Mocker struct {
		// handlers contains all mock handlers.
		// The first map is sorted by path, the second by method.
		// This ensures that different requests don't share the same Handler array, while still enforcing the call order.
		handlers map[string]map[string][]Handler // map[Path]map[HTTPMethod][]Handler
		// mut is the Mutex used to secure the handlers map, when multiple request come in concurrently.
		// However, adding handlers is not concurrent safe, as there is no point in it, when testing in a single method.
		mut *sync.Mutex
		// server is the Server used to mock the requests.
		server *httptest.Server
		// t is the test type called on error.
		t *testing.T

		// session is the session that is being mocked.
		session *session.Session
	}

	// Handler is a named handler for mocked endpoints.
	Handler struct {
		// Name is the name of the handler.
		Name string
		// Handler is the underlying handler.
		http.Handler
	}

	// MockFunc is the function used to create a mock.
	MockFunc func(w http.ResponseWriter, r *http.Request, t *testing.T)
)

// New creates a new Mocker, starts its test server and returns a manipulated Session using the test server.
func New(t *testing.T) (*Mocker, *session.Session) {
	m := &Mocker{
		handlers: make(map[string]map[string][]Handler, 1),
		mut:      new(sync.Mutex),
		t:        t,
	}

	m.server = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.mut.Lock()
		defer m.mut.Unlock()

		h, ok := m.handlers[r.URL.Path]
		require.True(t, ok, "unhandled path '"+r.URL.Path+"'")

		methHandlers, ok := h[r.Method]
		require.True(t, ok, "unhandled method '"+r.Method+"' on path '"+r.URL.Path+"'")

		methHandlers[0].ServeHTTP(w, r)

		if len(methHandlers) == 1 {
			if len(h) == 1 {
				delete(m.handlers, r.URL.Path)
			} else {
				delete(m.handlers[r.URL.Path], r.Method)
			}
		} else {
			m.handlers[r.URL.Path][r.Method] = m.handlers[r.URL.Path][r.Method][1:]
		}
	}))

	m.server.StartTLS()

	return m, newMockedSession(m.server.Listener.Addr().String())
}

// newMockedSession creates a new mocked session using the passed address.
func newMockedSession(addr string) *session.Session {
	gw := gateway.NewCustomGateway("", "")
	s := session.NewWithGateway(gw)

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, _ string) (conn net.Conn, err error) {
				return net.Dial(network, addr)
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	s.Client.Client.Client = httpdriver.WrapClient(client)
	s.Client.Retries = 1

	return s
}

// Mock uses the passed MockFunc to as handler for the passed path and method.
// The path is the path with "/api/v6" stripped.
// If there are already handlers for this path with the same method, the handler will be queued up behind the other
// handlers with the same path and method.
// Queued up handlers are invoked in the same order they were added in.
//
// The MockFunc may be nil, if only testing for invokes is needed.
func (m *Mocker) Mock(name, method, path string, f MockFunc) {
	path = "/api/v6" + path

	if m.handlers[path] == nil {
		m.handlers[path] = make(map[string][]Handler)
	}

	h := Handler{
		Name: name,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if f != nil {
				f(w, r, m.t)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		}),
	}

	m.handlers[path][method] = append(m.handlers[path][method], h)
}

// Clone clones handlers of the Mocker and returns the cloned Mocker and a new Session.
// Useful for multiple tests with the same API calls.
func (m *Mocker) Clone(t *testing.T) (*Mocker, *session.Session) {
	handlersCopy := make(map[string]map[string][]Handler, len(m.handlers))

	for p, sub := range m.handlers {
		subCopy := make(map[string][]Handler, len(sub))

		for m, ssub := range sub {
			ssubCopy := make([]Handler, len(ssub))

			copy(ssubCopy, ssub)

			subCopy[m] = ssubCopy
		}

		handlersCopy[p] = subCopy
	}

	clone, s := New(t)

	clone.handlers = handlersCopy

	return clone, s
}

// Eval closes the server and evaluates if all registered handlers were invoked.
// If not it will Fatal, stating the uninvoked handlers.
// This must be called at the end of every test.
func (m *Mocker) Eval() {
	m.Close()

	if len(m.handlers) == 0 {
		return
	}

	m.t.Fatal("there are uninvoked handlers:\n\n" + m.genUninvokedMsg())
}

// Close shuts down the server and blocks until all outstanding requests on this server have completed.
func (m *Mocker) Close() { m.server.Close() }

// genUninvokedMsg generates an error message stating the unused handlers.
//
// Example
//		/guilds/118456055842734083
// 			Guild: 2 uinvoked handlers
// 		/guilds/118456055842734083/members/256827968133791744
//			ModifyMember: 1 uinvoked handler
func (m *Mocker) genUninvokedMsg() string {
	missing := make(map[string]map[string]int, len(m.handlers))

	for p, methHandlers := range m.handlers {
		missing[p] = make(map[string]int, 1)

		for _, handlers := range methHandlers {
			for _, h := range handlers {
				missing[p][h.Name]++
			}
		}
	}

	n := (len(m.handlers) - 1) * 2 // 2 for the colon and the line feed

	for p, missingRequests := range missing {
		n += len(p)

		for name, qty := range missingRequests {
			// log10(qty) is the number of digits and 19 is the number of characters without placeholders
			n += len(name) + int(math.Log10(float64(qty))) + 19

			if qty > 1 { // handler has to be pluralized
				n++
			}
		}
	}

	var b strings.Builder

	b.Grow(n)

	first := true

	for p, missingRequests := range missing {
		if !first {
			b.WriteRune('\n')
		}

		b.WriteString(p)
		b.WriteRune(':')

		for name, qty := range missingRequests {
			b.WriteString("\n\t")
			b.WriteString(name)
			b.WriteString(": ")
			b.WriteString(strconv.Itoa(qty))
			b.WriteString(" uninvoked handler")

			if qty > 1 {
				b.WriteRune('s')
			}
		}

		first = false
	}

	return b.String()
}
