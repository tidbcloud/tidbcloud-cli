/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package imp

import (
	"encoding/json"
	"fmt"
)

// ImportSourceTypeEnum the model 'ImportSourceTypeEnum'
type ImportSourceTypeEnum string

// List of ImportSourceType.Enum
const (
	IMPORTSOURCETYPEENUM_LOCAL         ImportSourceTypeEnum = "LOCAL"
	IMPORTSOURCETYPEENUM_S3            ImportSourceTypeEnum = "S3"
	IMPORTSOURCETYPEENUM_GCS           ImportSourceTypeEnum = "GCS"
	IMPORTSOURCETYPEENUM_AZURE_BLOB    ImportSourceTypeEnum = "AZURE_BLOB"
	IMPORTSOURCETYPEENUM_S3_COMPATIBLE ImportSourceTypeEnum = "S3_COMPATIBLE"
	IMPORTSOURCETYPEENUM_OSS           ImportSourceTypeEnum = "OSS"
)

// All allowed values of ImportSourceTypeEnum enum
var AllowedImportSourceTypeEnumEnumValues = []ImportSourceTypeEnum{
	"LOCAL",
	"S3",
	"GCS",
	"AZURE_BLOB",
	"S3_COMPATIBLE",
	"OSS",
}

func (v *ImportSourceTypeEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ImportSourceTypeEnum(value)
	for _, existing := range AllowedImportSourceTypeEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ImportSourceTypeEnum", value)
}

// NewImportSourceTypeEnumFromValue returns a pointer to a valid ImportSourceTypeEnum
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewImportSourceTypeEnumFromValue(v string) (*ImportSourceTypeEnum, error) {
	ev := ImportSourceTypeEnum(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ImportSourceTypeEnum: valid values are %v", v, AllowedImportSourceTypeEnumEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ImportSourceTypeEnum) IsValid() bool {
	for _, existing := range AllowedImportSourceTypeEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ImportSourceType.Enum value
func (v ImportSourceTypeEnum) Ptr() *ImportSourceTypeEnum {
	return &v
}

type NullableImportSourceTypeEnum struct {
	value *ImportSourceTypeEnum
	isSet bool
}

func (v NullableImportSourceTypeEnum) Get() *ImportSourceTypeEnum {
	return v.value
}

func (v *NullableImportSourceTypeEnum) Set(val *ImportSourceTypeEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableImportSourceTypeEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableImportSourceTypeEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableImportSourceTypeEnum(val *ImportSourceTypeEnum) *NullableImportSourceTypeEnum {
	return &NullableImportSourceTypeEnum{value: val, isSet: true}
}

func (v NullableImportSourceTypeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableImportSourceTypeEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
