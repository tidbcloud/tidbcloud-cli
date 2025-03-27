/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package imp

import (
	"encoding/json"
)

// ImportGcsAuthTypeEnum  - SERVICE_ACCOUNT_KEY: The access method is service account key.
type ImportGcsAuthTypeEnum string

// List of ImportGcsAuthType.Enum
const (
	IMPORTGCSAUTHTYPEENUM_SERVICE_ACCOUNT_KEY ImportGcsAuthTypeEnum = "SERVICE_ACCOUNT_KEY"

	// Unknown value for handling new enum values gracefully
	ImportGcsAuthTypeEnum_UNKNOWN ImportGcsAuthTypeEnum = "UNKNOWN"
)

// All allowed values of ImportGcsAuthTypeEnum enum
var AllowedImportGcsAuthTypeEnumEnumValues = []ImportGcsAuthTypeEnum{
	"SERVICE_ACCOUNT_KEY",
}

func (v *ImportGcsAuthTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ImportGcsAuthTypeEnum(value)
	for _, existing := range AllowedImportGcsAuthTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	// Instead of returning an error, assign UNKNOWN value
	*v = ImportGcsAuthTypeEnum_UNKNOWN
	return nil
}

// NewImportGcsAuthTypeEnumFromValue returns a pointer to a valid ImportGcsAuthTypeEnum
// for the value passed as argument, or UNKNOWN if the value is not in the enum list
func NewImportGcsAuthTypeEnumFromValue(v string) *ImportGcsAuthTypeEnum {
	ev := ImportGcsAuthTypeEnum(v)
	if ev.IsValid() {
		return &ev
	}
	unknown := ImportGcsAuthTypeEnum_UNKNOWN
	return &unknown
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ImportGcsAuthTypeEnum) IsValid() bool {
	for _, existing := range AllowedImportGcsAuthTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ImportGcsAuthType.Enum value
func (v ImportGcsAuthTypeEnum) Ptr() *ImportGcsAuthTypeEnum {
	return &v
}

type NullableImportGcsAuthTypeEnum struct {
	value *ImportGcsAuthTypeEnum
	isSet bool
}

func (v NullableImportGcsAuthTypeEnum) Get() *ImportGcsAuthTypeEnum {
	return v.value
}

func (v *NullableImportGcsAuthTypeEnum) Set(val *ImportGcsAuthTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableImportGcsAuthTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableImportGcsAuthTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableImportGcsAuthTypeEnum(val *ImportGcsAuthTypeEnum) *NullableImportGcsAuthTypeEnum {
	return &NullableImportGcsAuthTypeEnum{value: val, isSet: true}
}

func (v NullableImportGcsAuthTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableImportGcsAuthTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
