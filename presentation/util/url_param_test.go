package util

import (
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const urlPrefix = "http://domain.com/path"

const boolParamName = "bool_param"
const intParamName = "int_param"
const floatParamName = "float_param"
const stringParamName = "string_param"

var (
	MandatoryBoolParamSpec   = NewUrlQueryParamSpec(boolParamName, Bool, true)
	MandatoryIntParamSpec    = NewUrlQueryParamSpec(intParamName, Int, true)
	MandatoryFloatParamSpec  = NewUrlQueryParamSpec(floatParamName, Float, true)
	MandatoryStringParamSpec = NewUrlQueryParamSpec(stringParamName, String, true)

	OptionaBoolParamSpec    = NewUrlQueryParamSpec(boolParamName, Bool, false)
	OptionalIntParamSpec    = NewUrlQueryParamSpec(intParamName, Int, false)
	OptionalFloatParamSpec  = NewUrlQueryParamSpec(floatParamName, Float, false)
	OptionalStringParamSpec = NewUrlQueryParamSpec(stringParamName, String, false)
)

func TestParseExistentOptionalParamHasCorrectParsedValue(t *testing.T) {
	asserParseSingleValuedValidParam(t, OptionaBoolParamSpec, []string{"true"}, true)
	asserParseSingleValuedValidParam(t, OptionaBoolParamSpec, []string{"false"}, false)
	asserParseSingleValuedValidParam(t, OptionalIntParamSpec, []string{"5"}, int64(5))
	asserParseSingleValuedValidParam(t, OptionalFloatParamSpec, []string{"5.0"}, float64(5.0))
	asserParseSingleValuedValidParam(t, OptionalStringParamSpec, []string{"abc"}, "abc")
}

func TestParseExistentMandatoryParamHasCorrectParsedValue(t *testing.T) {
	asserParseSingleValuedValidParam(t, MandatoryBoolParamSpec, []string{"true"}, true)
	asserParseSingleValuedValidParam(t, MandatoryBoolParamSpec, []string{"false"}, false)
	asserParseSingleValuedValidParam(t, MandatoryIntParamSpec, []string{"5"}, int64(5))
	asserParseSingleValuedValidParam(t, MandatoryFloatParamSpec, []string{"5.0"}, float64(5.0))
	asserParseSingleValuedValidParam(t, MandatoryStringParamSpec, []string{"abc"}, "abc")
}

func TestParseParamWithInvalidValueProducesError(t *testing.T) {
	assertParseParamWithInvalidValueProducesError(t, OptionaBoolParamSpec, []string{"tru_e"})
	assertParseParamWithInvalidValueProducesError(t, OptionaBoolParamSpec, []string{"fal-se"})
	assertParseParamWithInvalidValueProducesError(t, OptionalIntParamSpec, []string{"a5b"})
	assertParseParamWithInvalidValueProducesError(t, OptionalFloatParamSpec, []string{"jj5.0c"})
}

func TestParseNotExistentOptionalParamHasNoParsedValue(t *testing.T) {
	asserParseSingleValuedValidParam(t, OptionaBoolParamSpec, []string{}, nil)
	asserParseSingleValuedValidParam(t, OptionaBoolParamSpec, []string{}, nil)
	asserParseSingleValuedValidParam(t, OptionalIntParamSpec, []string{}, nil)
	asserParseSingleValuedValidParam(t, OptionalFloatParamSpec, []string{}, nil)
	asserParseSingleValuedValidParam(t, OptionalStringParamSpec, []string{}, nil)
}

func TestParseMultiplesValidOptionalParamHasCorrectParsedValue(t *testing.T) {
	asserParseMultipleValuedValidParam(t, OptionaBoolParamSpec, []string{"true", "false"}, []interface{}{true, false})
	asserParseMultipleValuedValidParam(t, OptionalIntParamSpec, []string{"4", "5"}, []interface{}{int64(4), int64(5)})
	asserParseMultipleValuedValidParam(t, OptionalFloatParamSpec, []string{"1.5", "2.3"}, []interface{}{float64(1.5), float64(2.3)})
	asserParseMultipleValuedValidParam(t, OptionalStringParamSpec, []string{"abc", "bcd"}, []interface{}{"abc", "bcd"})
}

func TestParseMultiplesValidMandatoryParamHasCorrectParsedValue(t *testing.T) {
	asserParseMultipleValuedValidParam(t, MandatoryBoolParamSpec, []string{"true", "false"}, []interface{}{true, false})
	asserParseMultipleValuedValidParam(t, MandatoryIntParamSpec, []string{"4", "5"}, []interface{}{int64(4), int64(5)})
	asserParseMultipleValuedValidParam(t, MandatoryFloatParamSpec, []string{"1.5", "2.3"}, []interface{}{float64(1.5), float64(2.3)})
	asserParseMultipleValuedValidParam(t, MandatoryStringParamSpec, []string{"abc", "bcd"}, []interface{}{"abc", "bcd"})
}

func TestParseMultiplesWithInValidParamProducesError(t *testing.T) {
	assertParseParamWithInvalidValueProducesError(t, MandatoryBoolParamSpec, []string{"true", "but-not-false"})
	assertParseParamWithInvalidValueProducesError(t, MandatoryIntParamSpec, []string{"4", "notInt"})
	assertParseParamWithInvalidValueProducesError(t, MandatoryFloatParamSpec, []string{"1.5", "notFloat"})
}

func TestParseNotExistantMandatoryParamProducesError(t *testing.T) {
	assertParseNotExistantMandatoryParamProducesError(t, MandatoryBoolParamSpec)
	assertParseNotExistantMandatoryParamProducesError(t, MandatoryIntParamSpec)
	assertParseNotExistantMandatoryParamProducesError(t, MandatoryFloatParamSpec)
	assertParseNotExistantMandatoryParamProducesError(t, MandatoryStringParamSpec)
}

func TestUsingNotValidDataTypeProducesError(t *testing.T) {
	ps := NewUrlQueryParamSpec("param", 10, false)
	values := generateValues(ps, []string{"value"})
	r := httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	_, err := ParseUrlQueryParam(r, ps)
	if err == nil {
		t.Error("Assertion error: Expected error to happen by using non existing type for a parameter specification.")
	}
}

func asserParseSingleValuedValidParam(t *testing.T, ps UrlQueryParamSpec, rawValues []string, expectedParsedValue interface{}) {
	values := generateValues(ps, rawValues)
	r := httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	result, err := ParseUrlQueryParam(r, ps)
	if err != nil {
		t.Fatalf("Unexpected error, error is: '%v'", err)
	}
	assertEquals(t, result.GetParsedValue(), expectedParsedValue)
}

func asserParseMultipleValuedValidParam(t *testing.T, ps UrlQueryParamSpec, rawValues []string, expectedParsedValue []interface{}) {
	values := generateValues(ps, rawValues)
	r := httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	result, err := ParseUrlQueryParam(r, ps)
	if err != nil {
		t.Fatalf("Unexpected error, error is: '%v'", err)
	}
	assertEquals(t, result.GetParsedValues(), expectedParsedValue)
}

func assertParseParamWithInvalidValueProducesError(t *testing.T, ps UrlQueryParamSpec, rawValues []string) {
	values := generateValues(ps, rawValues)
	r := httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	_, err := ParseUrlQueryParam(r, ps)
	if err == nil {
		t.Error("Assertion error: Expected error to happen when parsing invald values.")
	}
}

func assertParseNotExistantMandatoryParamProducesError(t *testing.T, ps UrlQueryParamSpec) {
	r := httptest.NewRequest("GET", urlPrefix, nil)
	_, err := ParseUrlQueryParam(r, ps)
	if err == nil {
		t.Error("Assertion error: Expected error to happen when mandatory param doens't exists.")
	}
}

func generateValues(ps UrlQueryParamSpec, rawValues []string) (values url.Values) {
	values = url.Values{}
	for _, rawValue := range rawValues {
		values.Add(ps.Name, rawValue)
	}
	return
}

var SomeParamsSpecs []UrlQueryParamSpec = []UrlQueryParamSpec{
	NewUrlQueryParamSpec(boolParamName, Bool, Optional),
	NewUrlQueryParamSpec(intParamName, Int, Optional),
	NewUrlQueryParamSpec(floatParamName, Float, Optional),
	NewUrlQueryParamSpec(stringParamName, String, Optional),
}

func TestParseUrlParamsGenerateCorrectParsedValues(t *testing.T) {
	values := url.Values{}
	values.Add(boolParamName, "true")
	values.Add(intParamName, "1")
	values.Add(floatParamName, "1.3")
	values.Add(stringParamName, "string")
	r := httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	params, err := ParseUrlQueryParams(r, SomeParamsSpecs)
	if err != nil {
		t.Fatalf("Unexpected error, error is: '%v'", err)
	}
	boolValue := params.Get(boolParamName).GetParsedValue()
	assertEquals(t, boolValue, true)
	intValue := params.Get(intParamName).GetParsedValue()
	assertEquals(t, intValue, int64(1))
	floatValue := params.Get(floatParamName).GetParsedValue()
	assertEquals(t, floatValue, float64(1.3))
	stringValue := params.Get(stringParamName).GetParsedValue()
	assertEquals(t, stringValue, "string")
}

func TestParseUrlParamsReturnsNilForNonExistantParam(t *testing.T) {
	r := httptest.NewRequest("GET", urlPrefix, nil)
	params, err := ParseUrlQueryParams(r, SomeParamsSpecs)
	if err != nil {
		t.Fatalf("Unexpected error, error is: '%v'", err)
	}
	value := params.Get("not-existant-param")
	if value != nil {
		t.Fatal("Assertion error. Expected nil for an unexisting url parameter name as argument")
	}
}

func TestParseUrlParamsReturnsNilForExistingParamBUtWithoutvalue(t *testing.T) {
	flagParamSpec := NewUrlQueryParamSpec("flag", String, Optional)
	r := httptest.NewRequest("GET", urlPrefix+"?flag", nil)
	params, err := ParseUrlQueryParams(r, []UrlQueryParamSpec{flagParamSpec})
	if err != nil {
		t.Fatalf("Unexpected error, error is: '%v'", err)
	}
	value := params.Get("flag") // until flag feature is implemeted (so params without value are allowed) the value will be nil
	if value != nil {
		t.Fatal("Assertion error. Expected nil for an existing url parameter without value")
	}

}

func TestParseUrlParamReturnsNilForExistingParamWithoutvalue(t *testing.T) {
	flagParamSpec := NewUrlQueryParamSpec("flag", String, Optional)
	r := httptest.NewRequest("GET", urlPrefix+"?flag", nil)
	param, err := ParseUrlQueryParam(r, flagParamSpec)
	if err != nil {
		t.Fatalf("Unexpected error, error is: '%v'", err)
	}
	if param.GetParsedValues() != nil {
		t.Fatal("Assertion error. Expected nil as parsed values for an existing url parameter without value")
	}
}

func TestParseUrlParamsProducesErrorForInvalidAnInvalidParamValue(t *testing.T) {
	var specs []UrlQueryParamSpec = []UrlQueryParamSpec{
		NewUrlQueryParamSpec(boolParamName, Bool, Optional),
		NewUrlQueryParamSpec(intParamName, Int, Optional),
	}
	values := url.Values{}
	values.Add(boolParamName, "true")
	values.Add(intParamName, "not-an-int")
	r := httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	_, err := ParseUrlQueryParams(r, specs)
	if err == nil {
		t.Error("Assertion error: Expected error to happen when parsing a set of parameters with at least one of them having an invalud value according to his type.")
	}

	specs = []UrlQueryParamSpec{
		NewUrlQueryParamSpec(boolParamName, Bool, Optional),
		NewUrlQueryParamSpec(intParamName, Int, Optional),
		NewUrlQueryParamSpec("missing-mandatory-string-param", String, Mandatory),
	}
	values = url.Values{}
	values.Add(boolParamName, "true")
	values.Add(intParamName, "1")
	r = httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	_, err = ParseUrlQueryParams(r, specs)
	if err == nil {
		t.Error("Assertion error: Expected error to happen when parsing a set of parameters without a mandatory parameter.")
	}
}

func TestUnixEpochFromParamWork(t *testing.T) {
	const paramName = "unixTimestamp"

	// Test without URL parameter
	r := httptest.NewRequest("GET", urlPrefix, nil)
	expected := int64(0)
	result := UnixEpochFromQueryParam(r, paramName, expected)
	if result != expected {
		t.Fatalf("UnixEpochFromParam fail, result was '%v' and expected is '%v'", result, expected)
	}

	// Test with a valid URL parameter
	expected = time.Now().Unix()
	stamp := strconv.FormatInt(expected, 10)

	values := url.Values{}
	values.Add(paramName, stamp)
	r = httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	result = UnixEpochFromQueryParam(r, paramName, 0)
	if result != expected {
		t.Fatalf("UnixEpochFromParam fail, result was '%v' and expected is '%v'", result, expected)
	}

	// Test with an invalid URL parameter
	values = url.Values{}
	values.Add(paramName, "I am not an stamp")
	r = httptest.NewRequest("GET", urlPrefix+"?"+values.Encode(), nil)
	expected = 123
	result = UnixEpochFromQueryParam(r, paramName, expected)
	if result != expected {
		t.Fatalf("UnixEpochFromParam fail, result was '%v' and expected is '%v'", result, expected)
	}

}

func assertEquals(t *testing.T, result interface{}, expected interface{}) {
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Assertion error: Should be equals. Results is: '%+v' %T and expected is: '%v' %T", result, result, expected, expected)
	}
}
