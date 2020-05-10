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

type testSchema struct {
	Needed    string `schema:"needed"`
	Omittable string `schema:"omittable,omitempty"`
}

func TestCheckQuery(t *testing.T) {
	successCases := []struct {
		name   string
		query  url.Values
		expect testSchema
	}{
		{
			name: "all fields filled",
			query: map[string][]string{
				"needed":    {"present"},
				"omittable": {"present as well"},
			},
			expect: testSchema{
				Needed:    "present",
				Omittable: "present as well",
			},
		},
		{
			name: "field omitted",
			query: map[string][]string{
				"needed": {"present"},
			},
			expect: testSchema{
				Needed: "present",
			},
		},
		{
			name: "zero values",
			query: map[string][]string{
				"needed": {""},
			},
			expect: testSchema{},
		},
	}

	failureCases := []struct {
		name        string
		query       url.Values
		falseExpect testSchema
	}{
		{
			name: "unequal",
			query: map[string][]string{
				"needed":    {"abc"},
				"omittable": {"def"},
			},
			falseExpect: testSchema{
				Needed:    "ghi",
				Omittable: "jkl",
			},
		},
		{
			name: "required field missing",
			query: map[string][]string{
				"omittable": {"present"},
			},
			falseExpect: testSchema{},
		},
	}

	t.Run("success", func(t *testing.T) {
		for _, c := range successCases {
			t.Run(c.name, func(t *testing.T) {
				CheckQuery(t, c.query, new(testSchema), &c.expect)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		for _, c := range failureCases {
			t.Run(c.name, func(t *testing.T) {
				tMock := new(testing.T)

				CheckQuery(tMock, c.query, new(testSchema), &c.falseExpect)

				assert.True(t, tMock.Failed())
			})
		}
	})
}
