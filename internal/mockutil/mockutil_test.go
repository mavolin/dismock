package mockutil

import (
	"io"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type writeRecorder struct {
	data []byte
}

func (r *writeRecorder) Write(p []byte) (n int, err error) {
	r.data = append(r.data, p...)
	return len(p), nil
}

type testStruct struct {
	Field1 int    `json:"field_1"`
	Field2 string `json:"field_2"`
}

func TestWriteJSON(t *testing.T) {
	s := testStruct{
		Field1: 123,
		Field2: "Hello World!",
	}

	expectWrite := []byte("{\"field_1\":123,\"field_2\":\"Hello World!\"}\n")

	rec := &writeRecorder{
		data: make([]byte, 0),
	}

	WriteJSON(t, rec, s)

	assert.Equal(t, expectWrite, rec.data)
}

type mockReader struct {
	src    []byte
	i      int
	closed bool
}

func (r *mockReader) Read(p []byte) (n int, err error) {
	if r.i >= len(r.src) {
		return 0, io.EOF
	}

	n = copy(p, r.src[r.i:])
	r.i += n

	return
}

func (r *mockReader) Close() error {
	r.closed = true
	return nil
}

func TestCheckJSONBody(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := &mockReader{
			src: []byte("{\"field_1\":123,\"field_2\":\"Hello World!\"}\n"),
		}

		expect := &testStruct{
			Field1: 123,
			Field2: "Hello World!",
		}

		CheckJSON(t, r, new(testStruct), expect)

		assert.True(t, r.closed)
	})

	t.Run("failure", func(t *testing.T) {
		r := &mockReader{
			src: []byte("{\"field_1\":321,\"field_2\":\"Bye World!\"}\n"),
		}

		expect := &testStruct{
			Field1: 123,
			Field2: "Hello World!",
		}

		tMock := new(testing.T)

		CheckJSON(tMock, r, &testStruct{}, expect)

		assert.True(t, tMock.Failed())
		assert.True(t, r.closed)
	})
}

func TestCheckQuery(t *testing.T) {
	failureCases := []struct {
		name        string
		query       url.Values
		falseExpect map[string]string
	}{
		{
			name: "unequal",
			query: map[string][]string{
				"foo": {"abc"},
				"bar": {"123"},
			},
			falseExpect: map[string]string{
				"foo": "def",
				"bar": "456",
			},
		},
		{
			name: "field missing",
			query: map[string][]string{
				"foo": {"present"},
			},
			falseExpect: map[string]string{},
		},
	}

	t.Run("success", func(t *testing.T) {
		CheckQuery(t, url.Values{
			"foo": {"present"},
			"bar": {"123"},
		}, map[string]string{
			"foo": "present",
			"bar": "123",
		})
	})

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				CheckQuery(tMock, c.query, c.falseExpect)

				assert.True(t, tMock.Failed())
			})
		}
	})
}
