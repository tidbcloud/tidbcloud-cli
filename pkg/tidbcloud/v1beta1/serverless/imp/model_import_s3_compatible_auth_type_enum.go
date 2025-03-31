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

// ImportS3CompatibleAuthTypeEnum  - ACCESS_KEY: The access method is access key.
type ImportS3CompatibleAuthTypeEnum string

// List of ImportS3CompatibleAuthType.Enum
const (
	IMPORTS3COMPATIBLEAUTHTYPEENUM_ACCESS_KEY ImportS3CompatibleAuthTypeEnum = "ACCESS_KEY"
)

// All allowed values of ImportS3CompatibleAuthTypeEnum enum
var AllowedImportS3CompatibleAuthTypeEnumEnumValues = []ImportS3CompatibleAuthTypeEnum{
	"ACCESS_KEY",
}

func (v *ImportS3CompatibleAuthTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ImportS3CompatibleAuthTypeEnum(value)
	for _, existing := range AllowedImportS3CompatibleAuthTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	*v = ImportS3CompatibleAuthTypeEnum(value)
	return nil
}

// NewImportS3CompatibleAuthTypeEnumFromValue returns a pointer to a valid ImportS3CompatibleAuthTypeEnum for the value passed as argument
func NewImportS3CompatibleAuthTypeEnumFromValue(v string) *ImportS3CompatibleAuthTypeEnum {
	ev := ImportS3CompatibleAuthTypeEnum(v)
	return &ev
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ImportS3CompatibleAuthTypeEnum) IsValid() bool {
	for _, existing := range AllowedImportS3CompatibleAuthTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ImportS3CompatibleAuthType.Enum value
func (v ImportS3CompatibleAuthTypeEnum) Ptr() *ImportS3CompatibleAuthTypeEnum {
	return &v
}

type NullableImportS3CompatibleAuthTypeEnum struct {
	value *ImportS3CompatibleAuthTypeEnum
	isSet bool
}

func (v NullableImportS3CompatibleAuthTypeEnum) Get() *ImportS3CompatibleAuthTypeEnum {
	return v.value
}

func (v *NullableImportS3CompatibleAuthTypeEnum) Set(val *ImportS3CompatibleAuthTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableImportS3CompatibleAuthTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableImportS3CompatibleAuthTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableImportS3CompatibleAuthTypeEnum(val *ImportS3CompatibleAuthTypeEnum) *NullableImportS3CompatibleAuthTypeEnum {
	return &NullableImportS3CompatibleAuthTypeEnum{value: val, isSet: true}
}

func (v NullableImportS3CompatibleAuthTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableImportS3CompatibleAuthTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
