/*
TiDB Cloud Serverless Export Open API

TiDB Cloud Serverless Export Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package export

import (
	"encoding/json"
)

// ExportS3AuthTypeEnum  - ROLE_ARN: The access method is role arn.  - ACCESS_KEY: The access method is access key.
type ExportS3AuthTypeEnum string

// List of ExportS3AuthType.Enum
const (
	EXPORTS3AUTHTYPEENUM_ROLE_ARN   ExportS3AuthTypeEnum = "ROLE_ARN"
	EXPORTS3AUTHTYPEENUM_ACCESS_KEY ExportS3AuthTypeEnum = "ACCESS_KEY"

	// Unknown value for handling new enum values gracefully
	ExportS3AuthTypeEnum_UNKNOWN ExportS3AuthTypeEnum = "unknown"
)

// All allowed values of ExportS3AuthTypeEnum enum
var AllowedExportS3AuthTypeEnumEnumValues = []ExportS3AuthTypeEnum{
	"ROLE_ARN",
	"ACCESS_KEY",
	ExportS3AuthTypeEnum_UNKNOWN, // Include unknown
}

func (v *ExportS3AuthTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ExportS3AuthTypeEnum(value)
	for _, existing := range AllowedExportS3AuthTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	// Instead of returning an error, assign UNKNOWN value
	*v = ExportS3AuthTypeEnum_UNKNOWN
	return nil
}

// NewExportS3AuthTypeEnumFromValue returns a pointer to a valid ExportS3AuthTypeEnum
// for the value passed as argument, or UNKNOWN if the value is not in the enum list
func NewExportS3AuthTypeEnumFromValue(v string) *ExportS3AuthTypeEnum {
	ev := ExportS3AuthTypeEnum(v)
	if ev.IsValid() {
		return &ev
	}
	unknown := ExportS3AuthTypeEnum_UNKNOWN
	return &unknown
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ExportS3AuthTypeEnum) IsValid() bool {
	for _, existing := range AllowedExportS3AuthTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ExportS3AuthType.Enum value
func (v ExportS3AuthTypeEnum) Ptr() *ExportS3AuthTypeEnum {
	return &v
}

type NullableExportS3AuthTypeEnum struct {
	value *ExportS3AuthTypeEnum
	isSet bool
}

func (v NullableExportS3AuthTypeEnum) Get() *ExportS3AuthTypeEnum {
	return v.value
}

func (v *NullableExportS3AuthTypeEnum) Set(val *ExportS3AuthTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableExportS3AuthTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableExportS3AuthTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExportS3AuthTypeEnum(val *ExportS3AuthTypeEnum) *NullableExportS3AuthTypeEnum {
	return &NullableExportS3AuthTypeEnum{value: val, isSet: true}
}

func (v NullableExportS3AuthTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExportS3AuthTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
