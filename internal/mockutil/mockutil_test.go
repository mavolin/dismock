package mockutil

import (
	"io"
	"net/url"
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/utils/json/option"
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

func TestCheckJSON(t *testing.T) {
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

func TestReplaceNullables(t *testing.T) {
	testCases := []struct {
		name   string
		in     interface{}
		expect interface{}
	}{
		{
			name:   "nothing",
			in:     nil,
			expect: nil,
		},
		{
			name: "NullableBool",
			in: &struct {
				B option.NullableBool
			}{
				B: option.NullBool,
			},
			expect: &struct {
				B option.NullableBool
			}{
				B: nil,
			},
		},
		{
			name: "NullableUint",
			in: &struct {
				B option.NullableUint
			}{
				B: option.NullUint,
			},
			expect: &struct {
				B option.NullableUint
			}{
				B: nil,
			},
		},
		{
			name: "NullableInt",
			in: &struct {
				B option.NullableInt
			}{
				B: option.NullInt,
			},
			expect: &struct {
				B option.NullableInt
			}{
				B: nil,
			},
		},
		{
			name: "NullableString",
			in: &struct {
				B option.NullableString
			}{
				B: option.NullString,
			},
			expect: &struct {
				B option.NullableString
			}{
				B: nil,
			},
		},
		{
			name: "NullableColor",
			in: &struct {
				B option.NullableColor
			}{
				B: option.NullColor,
			},
			expect: &struct {
				B option.NullableColor
			}{
				B: nil,
			},
		},
		{
			name: "nested",
			in: &struct {
				Nest struct {
					B option.NullableBool
				}
			}{
				Nest: struct {
					B option.NullableBool
				}{
					B: option.NullBool,
				},
			},
			expect: &struct {
				Nest struct {
					B option.NullableBool
				}
			}{
				Nest: struct {
					B option.NullableBool
				}{
					B: nil,
				},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			val := reflect.ValueOf(c.in)

			replaceNullables(val)

			assert.Equal(t, c.expect, c.in)
		})
	}
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
			name: "unexpected field",
			query: map[string][]string{
				"foo": {"present"},
			},
			falseExpect: map[string]string{},
		},
		{
			name:  "missing field",
			query: map[string][]string{},
			falseExpect: map[string]string{
				"foo": "I expected this to be here",
			},
		},
		{
			name: "empty query",
			query: map[string][]string{
				"foo": {},
			},
			falseExpect: map[string]string{
				"foo": "this should be filled",
			},
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
