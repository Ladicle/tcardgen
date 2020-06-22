package hugo

import (
	"fmt"
	"reflect"
)

func NewFMNotExistError(fmKey string) *FMNotExistError {
	return &FMNotExistError{Key: fmKey}
}

// FMNotExistError is an error that means the front-matter key/value is not exists in the content file.
type FMNotExistError struct {
	Key string
}

func (e *FMNotExistError) Error() string {
	return fmt.Sprintf("%q is not defined or empty", e.Key)
}

func NewFMInvalidTypeError(fmKey, wantType string, got interface{}) *FMInvalidTypeError {
	return &FMInvalidTypeError{Key: fmKey, Got: got, WantType: wantType}
}

// FMInvalidTypeError is an error that means the front-matter type is invalid.
type FMInvalidTypeError struct {
	Key      string
	WantType string
	Got      interface{}
}

func (e *FMInvalidTypeError) Error() string {
	return fmt.Sprintf("%q type of %v can not convert to %s", e.Key, reflect.TypeOf(e.Got), e.WantType)
}
