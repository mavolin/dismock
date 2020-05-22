package mockutil

import (
	"encoding/json"
	"io"
	"net/url"
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/utils/json/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WriteJSON writes the passed value to the passed Writer.
func WriteJSON(t *testing.T, w io.Writer, v interface{}) {
	err := json.NewEncoder(w).Encode(v)
	require.NoError(t, err)
}

var (
	nullableBool   = reflect.TypeOf(new(option.NullableBoolData))
	nullableUint   = reflect.TypeOf(new(option.NullableUintData))
	nullableInt    = reflect.TypeOf(new(option.NullableIntData))
	nullableString = reflect.TypeOf(new(option.NullableStringData))
	nullableColor  = reflect.TypeOf(new(option.NullableColorData))

	nilBool   = reflect.ValueOf((*option.NullableBoolData)(nil))
	nilUint   = reflect.ValueOf((*option.NullableUintData)(nil))
	nilInt    = reflect.ValueOf((*option.NullableIntData)(nil))
	nilString = reflect.ValueOf((*option.NullableStringData)(nil))
	nilColor  = reflect.ValueOf((*option.NullableColorData)(nil))
)

// CheckJSON checks the body of the passed Request to check against the passed expected value, assuming the body
// contains JSON data.
// v will be used to decode into and should not contain any data.
func CheckJSON(t *testing.T, r io.ReadCloser, v interface{}, expect interface{}) {
	err := json.NewDecoder(r).Decode(v)
	require.NoError(t, err)

	require.NoError(t, r.Close())

	val := reflect.ValueOf(expect)
	replaceNullables(val)

	assert.Equal(t, expect, v)
}

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

			isNil := field.Kind() == reflect.Ptr && field.IsNil()

			const initField = "Init"

			// this is a workaround to compensate for json.Unmarshal not calling Unmarshaler
			// functions on JSON null
			switch {
			case t.AssignableTo(nullableBool):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilBool)
				}
			case t.AssignableTo(nullableUint):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilUint)
				}
			case t.AssignableTo(nullableInt):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilInt)
				}
			case t.AssignableTo(nullableString):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilString)
				}
			case t.AssignableTo(nullableColor):
				if !isNil && !elem.FieldByName(initField).Bool() {
					field.Set(nilColor)
				}
			case !isNil && elem.Kind() == reflect.Struct:
				replaceNullables(field)
			}
		}
	}
}

// CheckQuery checks if the passed query contains the values found in except.
func CheckQuery(t *testing.T, query url.Values, expect url.Values) {
	for name, vals := range query {
		if len(vals) == 0 {
			continue
		}
		expectVal, ok := expect[name]
		if !assert.True(t, ok, "unexpected query field: '"+name+"' with value '"+vals[0]+"'") {
			continue
		}

		assert.Equal(t, expectVal, vals, "query fields for '"+name+"' don't match")

		delete(expect, name)
	}

	for name := range expect {
		assert.Fail(t, "missing query field: '"+name+"'")
	}
}
