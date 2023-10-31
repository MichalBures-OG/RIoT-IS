// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type NewUserDefinedDeviceTypeInput struct {
	Denotation string                                 `json:"denotation"`
	Parameters []*UserDefinedDeviceTypeParameterInput `json:"parameters"`
}

type UserDefinedDeviceType struct {
	ID         string                            `json:"id"`
	Denotation string                            `json:"denotation"`
	Parameters []*UserDefinedDeviceTypeParameter `json:"parameters"`
}

type UserDefinedDeviceTypeParameter struct {
	ID   string                             `json:"id"`
	Name string                             `json:"name"`
	Type UserDefinedDeviceTypeParameterType `json:"type"`
}

type UserDefinedDeviceTypeParameterInput struct {
	Name string                             `json:"name"`
	Type UserDefinedDeviceTypeParameterType `json:"type"`
}

type UserDefinedDeviceTypeParameterType string

const (
	UserDefinedDeviceTypeParameterTypeString  UserDefinedDeviceTypeParameterType = "STRING"
	UserDefinedDeviceTypeParameterTypeNumber  UserDefinedDeviceTypeParameterType = "NUMBER"
	UserDefinedDeviceTypeParameterTypeBoolean UserDefinedDeviceTypeParameterType = "BOOLEAN"
)

var AllUserDefinedDeviceTypeParameterType = []UserDefinedDeviceTypeParameterType{
	UserDefinedDeviceTypeParameterTypeString,
	UserDefinedDeviceTypeParameterTypeNumber,
	UserDefinedDeviceTypeParameterTypeBoolean,
}

func (e UserDefinedDeviceTypeParameterType) IsValid() bool {
	switch e {
	case UserDefinedDeviceTypeParameterTypeString, UserDefinedDeviceTypeParameterTypeNumber, UserDefinedDeviceTypeParameterTypeBoolean:
		return true
	}
	return false
}

func (e UserDefinedDeviceTypeParameterType) String() string {
	return string(e)
}

func (e *UserDefinedDeviceTypeParameterType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UserDefinedDeviceTypeParameterType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UserDefinedDeviceTypeParameterType", str)
	}
	return nil
}

func (e UserDefinedDeviceTypeParameterType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
