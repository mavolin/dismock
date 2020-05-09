package mockutil

import (
	"encoding/json"
	"io"
	"net/url"
	"testing"

	"github.com/gorilla/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WriteJSON writes the passed value to the passed Writer.
func WriteJSON(t *testing.T, w io.Writer, v interface{}) {
	err := json.NewEncoder(w).Encode(v)
	require.NoError(t, err)
}

// CheckJSONBody checks the body of the passed Request to check against the passed expected value, assuming the body
// contains JSON data.
// v will be used to decode into and should not contain any data.
func CheckJSONBody(t *testing.T, r io.ReadCloser, v interface{}, expect interface{}) {
	err := json.NewDecoder(r).Decode(v)
	require.NoError(t, err)

	require.NoError(t, r.Close())

	assert.Equal(t, expect, v)
}

// CheckQuery checks if the passed query contains the values found in except.
func CheckQuery(t *testing.T, query url.Values, v interface{}, expect interface{}) {
	d := schema.NewDecoder()
	d.ZeroEmpty(true)

	err := d.Decode(v, query)
	require.NoError(t, err)

	assert.Equal(t, expect, v)
}
