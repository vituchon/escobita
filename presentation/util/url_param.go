package util

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	NotSupportedType  = errors.New("No parsing implemented for the given type")
	UrlParamNotExists = errors.New("No parameter exists with the given keyname")
)

type DataType int

const (
	Int = iota
	Float
	Bool
	String
)

var dataTypeNames = []string{
	Int:    "Int",
	Float:  "Float",
	Bool:   "Bool",
	String: "String",
	// TODO(vgiordano): implement support for "flag" type (param without value)
}

func (dt DataType) String() string {
	if int(dt) < len(dataTypeNames) {
		return dataTypeNames[dt]
	}
	return "invalid value for DataType"
}

const (
	Optional  = false
	Mandatory = true
)

// URL parameter specification.
// Holds the necessary data to perform a validation in terms of existence and value type.
type UrlQueryParamSpec struct {
	Name        string
	DataType    DataType
	IsMandatory bool
	// TODO (vgiordano, NICE): add a multiplicity for validating quantity
}

func NewUrlQueryParamSpec(name string, kind DataType, isMandatory bool) UrlQueryParamSpec {
	return UrlQueryParamSpec{
		Name:        name,
		DataType:    kind,
		IsMandatory: isMandatory,
	}
}

// URL query/string parsed parameter
type UrlQueryParam struct {
	UrlQueryParamSpec
	rawValues    []string
	parsedValues []interface{}
}

// Gets the raw values
func (p UrlQueryParam) GetRawValues() []string {
	return p.rawValues
}

// Gets gets the first value associated
// If there are no values associated with the key then returns nil.
func (p UrlQueryParam) GetParsedValue() interface{} {
	if p.parsedValues == nil {
		return nil
	}
	return p.parsedValues[0]
}

// Gets the multiples values associated
func (p UrlQueryParam) GetParsedValues() []interface{} {
	return p.parsedValues
}

// Gets an string representation
func (p UrlQueryParam) String() string {
	return fmt.Sprintf("Spec{name=%v,kind=%v,isMandatory=%v};Values{raw=%v(type=%T),parsed=%v(type=%T)}", p.Name, p.DataType, p.IsMandatory, p.rawValues, p.rawValues, p.parsedValues, p.parsedValues)
}

// Signature for parser alike function, they validate and parse a given "string" that acts as the raw value.
type parseRawValueFunc func(rawValue string) (parsedValue interface{}, err error)

func boolParser(rawValue string) (parsedValue interface{}, err error) {
	parsedValue, err = strconv.ParseBool(rawValue)
	return
}

func intParser(rawValue string) (parsedValue interface{}, err error) {
	parsedValue, err = strconv.ParseInt(rawValue, 10, 64)
	return
}

func floatParser(rawValue string) (parsedValue interface{}, err error) {
	parsedValue, err = strconv.ParseFloat(rawValue, 64)
	return
}

func stringParser(rawValue string) (parsedValue interface{}, err error) {
	parsedValue = rawValue
	return
}

func unsupportedTypeParser(rawValue string) (parsedValue interface{}, err error) {
	err = NotSupportedType
	return
}

// Parses raw values delegating the validation and parsing of each element to a given parser function
func parseMultiple(rawValues []string, parseFunc parseRawValueFunc) (parsedValues []interface{}, err error) {
	parsedValues = make([]interface{}, 0, len(rawValues))
	for _, rawValue := range rawValues {
		var parsedValue interface{}
		parsedValue, err = parseFunc(rawValue)
		if err != nil {
			return
		}
		parsedValues = append(parsedValues, parsedValue)
	}
	return
}

func lookupParserFor(dataType DataType) parseRawValueFunc {
	switch dataType {
	case Bool:
		return boolParser
	case Int:
		return intParser
	case Float:
		return floatParser
	case String:
		return stringParser
	default:
		return unsupportedTypeParser
	}
}

// Parse from the request line a query/string parameter using given specifications.
func ParseUrlQueryParam(r *http.Request, ps UrlQueryParamSpec) (urlParam *UrlQueryParam, err error) {
	rawValues := r.URL.Query()[ps.Name]
	if len(rawValues) > 0 && rawValues[0] != "" { // flag parameters cames with rawValues as '[]string{""}' (a single slice with empty string)
		var parsedValues []interface{}
		parseFunc := lookupParserFor(ps.DataType)
		parsedValues, err = parseMultiple(rawValues, parseFunc)
		if err != nil {
			errMsg := fmt.Sprintf("Can not parse url query parameter(name='%s') as type='%v' from string='%s' (error was: '%v')", ps.Name, ps.DataType, rawValues, err)
			err = errors.New(errMsg)
		} else {
			urlParam = &UrlQueryParam{
				UrlQueryParamSpec: ps,
				parsedValues:      parsedValues,
				rawValues:         rawValues,
			}
		}
	} else {
		if ps.IsMandatory {
			errMsg := fmt.Sprintf("Mandatory url query parameter(name='%s') is not present", ps.Name)
			err = errors.New(errMsg)
		} else {

			urlParam = &UrlQueryParam{
				UrlQueryParamSpec: ps,
				parsedValues:      nil,
				rawValues:         rawValues,
			}
		}
	}
	return
}

type URLQueryParamsByName map[string]*UrlQueryParam

func (params URLQueryParamsByName) Get(paramName string) *UrlQueryParam {
	return params[paramName]
}

// Parse from the url request those url query/string paremeters with the given specs *discarting* those that aren't mandatory and, or not exists or has no associated value.
// This behaviour will change when url params admits flags parameters.
// Returns a map by name holding the actual parsed parameters and an error distinct to nil if any exception arises.
func ParseUrlQueryParams(request *http.Request, paramsSpecs []UrlQueryParamSpec) (params URLQueryParamsByName, err error) {
	params = make(map[string]*UrlQueryParam)
	for _, ps := range paramsSpecs {
		var p *UrlQueryParam
		p, err = ParseUrlQueryParam(request, ps)
		if err != nil {
			return
		} else if p.parsedValues != nil { // TODO(vgiordano): When implemented flag params modify this condition to  (p.parsedValues != nil  && p.DataType != Flag)
			params[p.Name] = p
		}
	}
	return
}

// Returns an unix timestamp from the given URL parameter name.
// If it is not present then a default value is used as return value.
func UnixEpochFromQueryParam(request *http.Request, paramName string, defaultValue int64) int64 {
	paramValue := request.URL.Query().Get(paramName)
	if len(paramValue) > 0 {
		date, err := strconv.Atoi(paramValue)
		if err != nil {
			return defaultValue
		}
		return int64(date)
	}
	return defaultValue
}
