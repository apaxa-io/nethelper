package httphelper

import (
	"errors"
	"github.com/apaxa-io/strconvhelper"
	"net/http"
)

// ScanErrorType define the type of error occurred while scanning form
type ScanErrorType uint8

// Define available ScanError types
const (
	ScanErrorTypeNoSuchField       ScanErrorType = iota // There is not field in form with requested name
	ScanErrorTypeMultipleValues                  = iota // There is more than 1 field in form with requested name
	ScanErrorTypeIncompatibleValue               = iota // Value in form is incompatible with requested field type (i.e. trying to save "one" as int)
	ScanErrorTypeIncompatibleType                = iota // Function unable to handle field with such type (i.e. truing to scan custom type)
)

// ScanError define error occurred while scanning form
type ScanError struct {
	FieldNum  int           // problem field number (beginning from 0)
	FieldName string        // problem field name
	Type      ScanErrorType // type of error
	SubError  error         // child error, used to exactly describe problem with incompatible value (nil for other types of error)
}

func scanErrorNoSuchField(fieldNum int, fieldName string) ScanError {
	return ScanError{FieldNum: fieldNum, FieldName: fieldName, Type: ScanErrorTypeNoSuchField, SubError: nil}
}

func scanErrorMultipleValues(fieldNum int, fieldName string) ScanError {
	return ScanError{FieldNum: fieldNum, FieldName: fieldName, Type: ScanErrorTypeMultipleValues, SubError: nil}
}

func scanErrorIncompatibleValue(fieldNum int, fieldName string, subError error) ScanError {
	return ScanError{FieldNum: fieldNum, FieldName: fieldName, Type: ScanErrorTypeIncompatibleValue, SubError: subError}
}

func scanErrorIncompatibleType(fieldNum int, fieldName string) ScanError {
	return ScanError{FieldNum: fieldNum, FieldName: fieldName, Type: ScanErrorTypeIncompatibleValue, SubError: nil}
}

// Error Implement error interface for ScanError. It returns text representation of error.
func (e ScanError) Error() string {
	prefix := "Scan error in #" + string(e.FieldNum) + "field with name '" + e.FieldName + "': "
	switch e.Type {
	case ScanErrorTypeNoSuchField:
		return prefix + "no field with such name."
	case ScanErrorTypeMultipleValues:
		return prefix + "there is more than 1 field with such name."
	case ScanErrorTypeIncompatibleValue:
		if e.SubError != nil {
			return prefix + e.Error()
		}
		return prefix + "unable to parse string to required type."
	case ScanErrorTypeIncompatibleType:
		return prefix + " type of this field is imcompatible with this function type."
	}
	return prefix + "unknown error"
}

// ScanField stores requested field name and variable to save value for ScanFormData.
type ScanField struct {
	Name  string      // field name
	Value interface{} // variable to store value
}

const scanBoolTrueString = "on"
const scanBoolFalseString = "off"

// ScanFormData scans Request.Form for required fields and save its value.
// Required fields names and variables to store values described by fields.
// This function accept for each required field exactly one value in form. There is error if zero or more than one fields with requested name exists in form.
// This function supports only following types of fields: [u]int[8/16/32/64], bools & strings.
// *int* will be parsed using strconv.ParseInt with base of 10.
// for bools valid values are only "on" & "off" (case sensitive).
// strings accepted as-is.
// Returned error is always of type ScanError or nil.
// Warning: r.ParseForm should be performed before calling this function.
func ScanFormData(r *http.Request, fields ...ScanField) error {
	for i, field := range fields {
		var stringValue string

		if stringValues, ok := r.Form[field.Name]; ok && len(stringValues) == 1 {
			stringValue = stringValues[0]
		} else if !ok {
			return scanErrorNoSuchField(i, field.Name)
		} else {
			return scanErrorMultipleValues(i, field.Name)
		}

		var err error
		switch value := field.Value.(type) {
		case *int:
			if *value, err = strconvhelper.ParseInt(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *int8:
			if *value, err = strconvhelper.ParseInt8(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *int16:
			if *value, err = strconvhelper.ParseInt16(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *int32:
			if *value, err = strconvhelper.ParseInt32(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *int64:
			if *value, err = strconvhelper.ParseInt64(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *uint:
			if *value, err = strconvhelper.ParseUint(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *uint8:
			if *value, err = strconvhelper.ParseUint8(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *uint16:
			if *value, err = strconvhelper.ParseUint16(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *uint32:
			if *value, err = strconvhelper.ParseUint32(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *uint64:
			if *value, err = strconvhelper.ParseUint64(stringValue); err != nil {
				return scanErrorIncompatibleValue(i, field.Name, err)
			}
		case *bool:
			switch stringValue {
			case scanBoolTrueString:
				*value = true
			case scanBoolFalseString:
				*value = false
			default:
				return scanErrorIncompatibleValue(i, field.Name, errors.New("'"+stringValue+"' is not a valid bool value."))
			}
		case *string:
			*value = stringValue
		default:
			return scanErrorIncompatibleType(i, field.Name)
		}
	}
	return nil
}

/*
func Int64FromForm(r *http.Request, name string) (value int64, err error) {
	if svs, ok := r.Form[name]; ok && len(svs) == 1 {
		value, err = strconv.ParseInt(svs[0], 10, 64)
		return
	} else {
		err = errors.New("Not exactly one field with name " + name)
		return
	}
}
*/
