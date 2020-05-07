package mockutil

import (
	"io"
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

		CheckJSONBody(t, r, new(testStruct), expect)

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

		CheckJSONBody(tMock, r, &testStruct{}, expect)

		assert.True(t, tMock.Failed())
		assert.True(t, r.closed)
	})
}
