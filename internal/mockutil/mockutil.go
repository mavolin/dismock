package mockutil

import (
	"encoding/json"
	"io"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WriteJSON writes the passed value to the passed Writer.
func WriteJSON(t *testing.T, w io.Writer, v interface{}) {
	err := json.NewEncoder(w).Encode(v)
	require.NoError(t, err)
}

// CheckJSON checks the body of the passed Request to check against the passed expected value, assuming the body
// contains JSON data.
// v will be used to decode into and should not contain any data.
func CheckJSON(t *testing.T, r io.ReadCloser, v interface{}, expect interface{}) {
	err := json.NewDecoder(r).Decode(v)
	require.NoError(t, err)

	require.NoError(t, r.Close())

	assert.Equal(t, expect, v)
}

// CheckQuery checks if the passed query contains the values found in except.
func CheckQuery(t *testing.T, query url.Values, expect map[string]string) {
	for name, vals := range query {
		if len(vals) == 0 {
			continue
		}
		expectVal, ok := expect[name]
		if !assert.True(t, ok, "unexpected query field: '"+name+"' with value '"+vals[0]+"'") {
			continue
		}

		assert.Equal(t, expectVal, vals[0], "query fields for '"+name+"' don't match")

		delete(expect, name)
	}

	for name := range expect {
		assert.Fail(t, "missing query field: '"+name+"'")
	}
}
