package mockutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/diamondburned/arikawa/v2/utils/sendpart"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Field1 int    `json:"field_1"`
	Field2 string `json:"field_2"`
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

type pipe struct {
	data []byte
	i    int
}

func (p *pipe) Write(data []byte) (n int, err error) {
	p.data = append(p.data, data...)

	return len(data), nil
}

func (p *pipe) Read(data []byte) (n int, err error) {
	if p.i >= len(p.data) {
		return 0, io.EOF
	}

	n = copy(data, p.data[p.i:])
	p.i += n

	return
}

func (p *pipe) Close() error { return nil }

func TestWriteJSON(t *testing.T) {
	s := testStruct{
		Field1: 123,
		Field2: "Hello World!",
	}

	expectWrite := []byte("{\"field_1\":123,\"field_2\":\"Hello World!\"}\n")

	rec := httptest.NewRecorder()

	WriteJSON(t, rec, s)

	data, err := ioutil.ReadAll(rec.Body)
	require.NoError(t, err)

	assert.Equal(t, expectWrite, data)
	assert.Equal(t, http.Header{
		"Content-Type": {"application/json"},
	}, rec.Header())
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

func TestCheckMultipart(t *testing.T) {
	successCases := []struct {
		name   string
		body   func(mw *multipart.Writer, t *testing.T)
		expect *testStruct
		f      []sendpart.File
	}{
		{
			name: "files only",
			body: func(mw *multipart.Writer, t *testing.T) {
				data := []byte{1, 2, 4, 8, 16, 32, 64, 128}

				fw, err := mw.CreateFormFile("file0", "abc")
				require.NoError(t, err)

				n, err := fw.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				require.NoError(t, err)

				data = []byte{128, 64, 32, 16, 8, 4, 2, 1}

				fw, err = mw.CreateFormFile("file1", "def")
				require.NoError(t, err)

				n, err = fw.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				require.NoError(t, err)

				require.NoError(t, mw.Close())
			},
			f: []sendpart.File{
				{
					Name:   "abc",
					Reader: bytes.NewBuffer([]byte{1, 2, 4, 8, 16, 32, 64, 128}),
				},
				{
					Name:   "def",
					Reader: bytes.NewBuffer([]byte{128, 64, 32, 16, 8, 4, 2, 1}),
				},
			},
		},
		{
			name: "both",
			body: func(mw *multipart.Writer, t *testing.T) {
				jw, err := mw.CreateFormField("payload_json")
				require.NoError(t, err)

				err = json.NewEncoder(jw).Encode(testStruct{
					Field1: 123,
					Field2: "abc",
				})
				require.NoError(t, err)

				data := []byte{1, 2, 4, 8, 16, 32, 64, 128}

				fw, err := mw.CreateFormFile("file0", "abc")
				require.NoError(t, err)

				n, err := fw.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				require.NoError(t, err)

				require.NoError(t, mw.Close())
			},
			expect: &testStruct{
				Field1: 123,
				Field2: "abc",
			},
			f: []sendpart.File{
				{
					Name:   "abc",
					Reader: bytes.NewBuffer([]byte{1, 2, 4, 8, 16, 32, 64, 128}),
				},
			},
		},
	}

	for _, c := range successCases {
		t.Run(c.name, func(t *testing.T) {
			p := new(pipe)
			mw := multipart.NewWriter(p)
			c.body(mw, t)

			h := http.Header{
				"Content-Type": {mw.FormDataContentType()},
			}

			var expect interface{}
			if c.expect != nil { // *testStruct(nil) != interface{}(nil), so we need this
				expect = c.expect
			}

			CheckMultipart(t, p, h, new(testStruct), expect, c.f)
		})
	}

	failureCases := []struct {
		name   string
		body   func(mw *multipart.Writer, t *testing.T)
		expect *testStruct
		f      []sendpart.File
	}{
		{
			name: "no json_payload",
			body: func(mw *multipart.Writer, t *testing.T) {
				data := []byte{1, 2, 4, 8, 16, 32, 64, 128}

				fw, err := mw.CreateFormFile("file0", "abc")
				require.NoError(t, err)

				n, err := fw.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				require.NoError(t, err)

				require.NoError(t, mw.Close())
			},
			expect: &testStruct{
				Field1: 123,
				Field2: "abc",
			},
			f: []sendpart.File{
				{
					Name:   "abc",
					Reader: bytes.NewBuffer([]byte{1, 2, 4, 8, 16, 32, 64, 128}),
				},
			},
		},
		{
			name: "unexpected json_payload",
			body: func(mw *multipart.Writer, t *testing.T) {
				jw, err := mw.CreateFormField("payload_json")
				require.NoError(t, err)

				err = json.NewEncoder(jw).Encode(testStruct{
					Field1: 123,
					Field2: "abc",
				})
				require.NoError(t, err)

				data := []byte{1, 2, 4, 8, 16, 32, 64, 128}

				fw, err := mw.CreateFormFile("file0", "abc")
				require.NoError(t, err)

				n, err := fw.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				require.NoError(t, err)

				require.NoError(t, mw.Close())
			},
			f: []sendpart.File{
				{
					Name:   "abc",
					Reader: bytes.NewBuffer([]byte{1, 2, 4, 8, 16, 32, 64, 128}),
				},
			},
		},
		{
			name: "too few files",
			body: func(mw *multipart.Writer, t *testing.T) {
				data := []byte{1, 2, 4, 8, 16, 32, 64, 128}

				fw, err := mw.CreateFormFile("file0", "abc")
				require.NoError(t, err)

				n, err := fw.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				require.NoError(t, err)

				require.NoError(t, mw.Close())
			},
			f: []sendpart.File{
				{
					Name:   "abc",
					Reader: bytes.NewBuffer([]byte{1, 2, 4, 8, 16, 32, 64, 128}),
				},
				{
					Name:   "def",
					Reader: bytes.NewBuffer([]byte{128, 64, 32, 16, 8, 4, 2, 1}),
				},
			},
		},
		{
			name: "too many files",
			body: func(mw *multipart.Writer, t *testing.T) {
				data := []byte{1, 2, 4, 8, 16, 32, 64, 128}

				fw, err := mw.CreateFormFile("file0", "abc")
				require.NoError(t, err)

				n, err := fw.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				require.NoError(t, err)

				data = []byte{128, 64, 32, 16, 8, 4, 2, 1}

				fw, err = mw.CreateFormFile("file1", "def")
				require.NoError(t, err)

				n, err = fw.Write(data)
				if err == nil && n < len(data) {
					err = io.ErrShortWrite
				}
				require.NoError(t, err)

				require.NoError(t, mw.Close())
			},
			f: []sendpart.File{
				{
					Name:   "abc",
					Reader: bytes.NewBuffer([]byte{1, 2, 4, 8, 16, 32, 64, 128}),
				},
			},
		},
	}

	for _, c := range failureCases {
		t.Run(c.name, func(t *testing.T) {
			p := new(pipe)
			mw := multipart.NewWriter(p)
			c.body(mw, t)

			h := http.Header{
				"Content-Type": {mw.FormDataContentType()},
			}

			var expect interface{}
			if c.expect != nil { // *testStruct(nil) != interface{}(nil), so we need this
				expect = c.expect
			}

			tMock := new(testing.T)

			CheckMultipart(tMock, p, h, new(testStruct), expect, c.f)

			assert.True(t, tMock.Failed())
		})
	}

	t.Run("unexpected part", func(t *testing.T) {
		p := new(pipe)
		mw := multipart.NewWriter(p)

		_, err := mw.CreateFormField("unexpected")
		require.NoError(t, err)

		require.NoError(t, mw.Close())

		h := http.Header{
			"Content-Type": {mw.FormDataContentType()},
		}

		tMock := new(testing.T)

		CheckMultipart(tMock, p, h, nil, nil, []sendpart.File{})

		assert.True(t, tMock.Failed())
	})
}

func TestCheckQuery(t *testing.T) {
	failureCases := []struct {
		name        string
		query       url.Values
		falseExpect url.Values
	}{
		{
			name: "unequal",
			query: url.Values{
				"foo": {"abc"},
				"bar": {"123"},
			},
			falseExpect: url.Values{
				"foo": {"def"},
				"bar": {"456"},
			},
		},
		{
			name: "unexpected field",
			query: url.Values{
				"foo": {"present"},
			},
			falseExpect: url.Values{},
		},
		{
			name:  "missing field",
			query: url.Values{},
			falseExpect: url.Values{
				"foo": {"I expected this to be here"},
			},
		},
		{
			name: "empty query",
			query: url.Values{
				"foo": {},
			},
			falseExpect: url.Values{
				"foo": {"this should be filled"},
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		CheckQuery(t, url.Values{
			"foo": {"present"},
			"bar": {"123"},
		}, url.Values{
			"foo": {"present"},
			"bar": {"123"},
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

func TestEqualReader(t *testing.T) {
	testCases := []struct {
		name  string
		b1    []byte
		b2    []byte
		equal bool
	}{
		{
			name:  "equal",
			b1:    []byte{0, 1, 2, 4, 8, 16, 32, 64, 128},
			b2:    []byte{0, 1, 2, 4, 8, 16, 32, 64, 128},
			equal: true,
		},
		{
			name:  "not equal",
			b1:    []byte{0, 1, 2, 4, 8, 16, 32, 64, 128},
			b2:    []byte{128, 64, 32, 16, 8, 4, 2, 1, 0},
			equal: false,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := equalReader(bytes.NewBuffer(c.b1), bytes.NewBuffer(c.b2))
			if c.equal {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJoinIntSet(t *testing.T) {
	testCases := []struct {
		name   string
		set    map[int]struct{}
		delim  string
		expect []string // maps are unordered, therefore we expect one of these to be correct
	}{
		{
			name:   "single value",
			set:    map[int]struct{}{0: {}},
			delim:  ", ",
			expect: []string{"0"},
		},
		{
			name:   "multiple values",
			set:    map[int]struct{}{0: {}, 1: {}},
			delim:  ", ",
			expect: []string{"0, 1", "1, 0"},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			join := joinIntSet(c.set, c.delim)

			for _, e := range c.expect {
				if assert.ObjectsAreEqual(e, join) {
					return
				}
			}

			assert.Equal(t, c.expect[0], join) // gen automatic diff
		})
	}
}
