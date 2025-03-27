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

// ImportStateEnum  - PREPARING: The import is preparing.  - IMPORTING: The import is importing.  - COMPLETED: The import is completed.  - FAILED: The import is failed.  - CANCELING: The import is canceling.  - CANCELED: The import is canceled.
type ImportStateEnum string

// List of ImportState.Enum
const (
	IMPORTSTATEENUM_PREPARING ImportStateEnum = "PREPARING"
	IMPORTSTATEENUM_IMPORTING ImportStateEnum = "IMPORTING"
	IMPORTSTATEENUM_COMPLETED ImportStateEnum = "COMPLETED"
	IMPORTSTATEENUM_FAILED    ImportStateEnum = "FAILED"
	IMPORTSTATEENUM_CANCELING ImportStateEnum = "CANCELING"
	IMPORTSTATEENUM_CANCELED  ImportStateEnum = "CANCELED"
)

// All allowed values of ImportStateEnum enum
var AllowedImportStateEnumEnumValues = []ImportStateEnum{
	"PREPARING",
	"IMPORTING",
	"COMPLETED",
	"FAILED",
	"CANCELING",
	"CANCELED",
}

func (v *ImportStateEnum) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ImportStateEnum(value)
	for _, existing := range AllowedImportStateEnumEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	*v = ImportStateEnum(value)
	return nil
}

// NewImportStateEnumFromValue returns a pointer to a valid ImportStateEnum for the value passed as argument
func NewImportStateEnumFromValue(v string) *ImportStateEnum {
	ev := ImportStateEnum(v)
	return &ev
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ImportStateEnum) IsValid() bool {
	for _, existing := range AllowedImportStateEnumEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ImportState.Enum value
func (v ImportStateEnum) Ptr() *ImportStateEnum {
	return &v
}

type NullableImportStateEnum struct {
	value *ImportStateEnum
	isSet bool
}

func (v NullableImportStateEnum) Get() *ImportStateEnum {
	return v.value
}

func (v *NullableImportStateEnum) Set(val *ImportStateEnum) {
	v.value = val
	v.isSet = true
}

func (v NullableImportStateEnum) IsSet() bool {
	return v.isSet
}

func (v *NullableImportStateEnum) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableImportStateEnum(val *ImportStateEnum) *NullableImportStateEnum {
	return &NullableImportStateEnum{value: val, isSet: true}
}

func (v NullableImportStateEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableImportStateEnum) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
