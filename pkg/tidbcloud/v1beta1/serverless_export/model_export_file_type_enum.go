/*
TiDB Cloud Serverless Export Open API

TiDB Cloud Serverless Export Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// ExportFileTypeEnum  - SQL: SQL type.  - CSV: CSV type.  - PARQUET: PARQUET type.
type ExportFileTypeEnum string

// List of ExportFileType.Enum
const (
	SQL     ExportFileTypeEnum = "SQL"
	CSV     ExportFileTypeEnum = "CSV"
	PARQUET ExportFileTypeEnum = "PARQUET"
)

// All allowed values of ExportFileTypeEnum enum
var AllowedExportFileTypeEnumEnumValues = []ExportFileTypeEnum{
	"SQL",
	"CSV",
	"PARQUET",
}

func (v *ExportFileTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ExportFileTypeEnum(value)
	for _, existing := range AllowedExportFileTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ExportFileTypeEnum", value)
}

// NewExportFileTypeEnumFromValue returns a pointer to a valid ExportFileTypeEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewExportFileTypeEnumFromValue(v string) (*ExportFileTypeEnum, error) {
	ev := ExportFileTypeEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ExportFileTypeEnum: valid values are %v", v, AllowedExportFileTypeEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ExportFileTypeEnum) IsValid() bool {
	for _, existing := range AllowedExportFileTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ExportFileType.Enum value
func (v ExportFileTypeEnum) Ptr() *ExportFileTypeEnum {
	return &v
}

type NullableExportFileTypeEnum struct {
	value *ExportFileTypeEnum
	isSet bool
}

func (v NullableExportFileTypeEnum) Get() *ExportFileTypeEnum {
	return v.value
}

func (v *NullableExportFileTypeEnum) Set(val *ExportFileTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableExportFileTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableExportFileTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExportFileTypeEnum(val *ExportFileTypeEnum) *NullableExportFileTypeEnum {
	return &NullableExportFileTypeEnum{value: val, isSet: true}
}

func (v NullableExportFileTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExportFileTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
