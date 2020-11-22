package dismock

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMocker_Error(t *testing.T) {
	m, s := NewSession(t)
	defer m.Eval()

	sendErr := httputil.HTTPError{
		Status:  http.StatusBadRequest,
		Code:    10011,
		Message: "Unknown guild",
	}

	m.Error(http.MethodGet, "/guilds/123", sendErr)

	_, err := s.Guild(123)
	require.IsType(t, new(httputil.HTTPError), err)

	httpErr := err.(*httputil.HTTPError)

	assert.Equal(t, sendErr.Status, httpErr.Status)
	assert.Equal(t, sendErr.Code, httpErr.Code)
	assert.Equal(t, sendErr.Message, httpErr.Message)
}
