package check

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/mavolin/dismock/v3/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WriteJSON writes the passed value to the passed http.ResponseWriter.
func WriteJSON(t testing.TInterface, w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(v)
	require.NoError(t, err)
}

// JSON checks if body contains the JSON data matching the passed expected
// value.
func JSON(t testing.TInterface, expect interface{}, actualReader io.ReadCloser) {
	checkJSON(t, expect, actualReader)
	require.NoError(t, actualReader.Close())
}

// Multipart checks if the body contains multipart data including the passed
// files and optionally the passed JSON data.
//nolint:funlen
func Multipart(
	t testing.TInterface, body io.ReadCloser, h http.Header, expectJSON interface{}, expectFiles []sendpart.File,
) {
	_, p, err := mime.ParseMediaType(h.Get("Content-Type"))
	require.NoError(t, err)

	bound, ok := p["boundary"]
	require.True(t, ok, "boundary parameter not set")

	mr := multipart.NewReader(body, bound)

	jsonChecked := false
	// we store the numbers of the missingFiles in a set, so that we know later
	// on which missingFiles didn't get sent, if any
	missingFiles := make(map[int]struct{}, len(expectFiles))

	for i := range expectFiles {
		missingFiles[i] = struct{}{}
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)

		name := part.FormName()

		switch {
		case name == "payload_json":
			if expectJSON != nil {
				checkJSON(t, expectJSON, part)
				jsonChecked = true
			} else {
				assert.Failf(t, "error when checking multipart body", "expected no json payload, but got: %+v", part)
			}
		case strings.HasPrefix(name, "file"):
			no, err := strconv.Atoi(strings.TrimLeft(name, "file"))
			require.NoErrorf(t, err, `unexpected part with name "%s"`, name)

			if !assert.Lessf(
				t, no, len(expectFiles), "reading file %d, but expected only %d missingFiles", no, len(expectFiles),
			) {
				break
			}

			assert.Equal(t, expectFiles[no].Name, part.FileName(), "unequal file names")

			err = equalReader(expectFiles[no].Reader, part)
			assert.NoErrorf(t, err, "file %d is not equal to received file", no)

			delete(missingFiles, no)
		default:
			assert.Fail(t, fmt.Sprintf(`unexpected part with name "%s"`, name))
		}
	}

	require.NoError(t, body.Close())

	if !jsonChecked && expectJSON != nil {
		assert.Fail(t, "no json_payload was received, but was expected")
	}

	if len(missingFiles) > 0 {
		s := joinIntSet(missingFiles, ", ")
		assert.Fail(t, fmt.Sprintf("the files %s did not get sent", s))
	}
}

// Query checks if the passed query contains the values found in except.
func Query(t testing.TInterface, expect url.Values, actual url.Values) {
	for name, vals := range actual {
		if len(vals) == 0 {
			continue
		}

		expectVal, ok := expect[name]
		if !assert.Truef(t, ok, "unexpected query field: '%s' with value '%s'", name, vals[0]) {
			continue
		}

		assert.Equal(t, expectVal, vals, "query fields for '"+name+"' don't match")

		delete(expect, name)
	}

	for name := range expect {
		assert.Fail(t, "missing query field: '"+name+"'")
	}
}

// Header checks if the expected http.Header are contained in actual.
func Header(t testing.TInterface, expect http.Header, actual http.Header) {
	for _, expect := range expect {
		assert.Contains(t, actual, expect)
	}
}

// checkJSON checks if body contains the JSON data matching the passed expected
// value.
func checkJSON(t testing.TInterface, expect interface{}, actualReader io.Reader) {
	decodeVal := reflect.New(reflect.TypeOf(expect))

	err := json.NewDecoder(actualReader).Decode(decodeVal.Interface())
	require.NoError(t, err)

	expectVal := reflect.ValueOf(expect)

	// do this so we always have an addressable value
	ptrExpectVal := reflect.New(expectVal.Type())
	ptrExpectVal = ptrExpectVal.Elem()
	ptrExpectVal.Set(expectVal)

	replaceNullables(ptrExpectVal)

	assert.Equal(t, ptrExpectVal.Interface(), decodeVal.Elem().Interface())
}

var (
	nullableBool   = reflect.TypeOf(new(option.NullableBoolData))
	nullableUint   = reflect.TypeOf(new(option.NullableUintData))
	nullableInt    = reflect.TypeOf(new(option.NullableIntData))
	nullableString = reflect.TypeOf(new(option.NullableStringData))

	nilBool   = reflect.ValueOf((*option.NullableBoolData)(nil))
	nilUint   = reflect.ValueOf((*option.NullableUintData)(nil))
	nilInt    = reflect.ValueOf((*option.NullableIntData)(nil))
	nilString = reflect.ValueOf((*option.NullableStringData)(nil))
)

// replacesNullables replaces the values of all nullable types with nil, if
// they represent json null.
//
// This is necessary because json.Unmarshal does not call the custom unmarshal
// methods of pointers and instead just assigns nil to them.
// This will trigger false positives when comparing api.FooData structs, as
// fields that use nullable types and that are set to null will get serialized
// as nil, but the expected value for that field is the null version of the
// type.
// To avoid this, we replace all those null versions with nil, so that
// comparison will pass.
func replaceNullables(val reflect.Value) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)

			t := field.Type()

			elem := field
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}

			if field.Kind() == reflect.Ptr && field.IsNil() {
				continue
			}

			const initField = "Init"

			switch {
			case t.AssignableTo(nullableBool):
				if !elem.FieldByName(initField).Bool() {
					field.Set(nilBool)
				}
			case t.AssignableTo(nullableUint):
				if !elem.FieldByName(initField).Bool() {
					field.Set(nilUint)
				}
			case t.AssignableTo(nullableInt):
				if !elem.FieldByName(initField).Bool() {
					field.Set(nilInt)
				}
			case t.AssignableTo(nullableString):
				if !elem.FieldByName(initField).Bool() {
					field.Set(nilString)
				}
			case elem.Kind() == reflect.Struct:
				replaceNullables(field)
			}
		}
	}
}

// equalReader checks if the values of the two readers contain the same data.
func equalReader(a, b io.Reader) error {
	const size = 16

	b1 := make([]byte, size)
	b2 := make([]byte, size)

	for i := 1; ; i++ {
		_, err1 := a.Read(b1)
		_, err2 := b.Read(b2)

		if !bytes.Equal(b1, b2) {
			return fmt.Errorf("%d. chunk is not equal:\n%v\nvs.\n%v", i, b1, b2)
		}

		switch {
		case err1 == io.EOF && err2 == io.EOF:
			return nil
		case err1 == io.EOF:
			_, err2 = b.Read(b2)
			if err2 == io.EOF {
				return nil
			}

			return errors.New("reader 1's stream ended unexpectedly")
		case err2 == io.EOF:
			_, err1 = a.Read(b1)
			if err1 == io.EOF {
				return nil
			}

			return errors.New("reader 2's stream ended unexpectedly")
		case err1 != nil:
			return err1
		case err2 != nil:
			return err2
		}
	}
}

// strings.Join, but for sets of int.
func joinIntSet(set map[int]struct{}, delim string) string {
	var s string

	first := true

	for no := range set {
		if !first {
			s += delim
		}

		s += strconv.Itoa(no)
		first = false
	}

	return s
}
