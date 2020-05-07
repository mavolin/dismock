package dismock

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// DiscordError is the error type returned by the discord API.
type DiscordError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// Error simulates an error response for the given path using the given method.
func (m *Mocker) Error(method, path string, e DiscordError) {
	m.Mock("Error", method, path, func(w http.ResponseWriter, r *http.Request, t *testing.T) {
		err := json.NewEncoder(w).Encode(e)
		require.NoError(t, err)
	})
}
