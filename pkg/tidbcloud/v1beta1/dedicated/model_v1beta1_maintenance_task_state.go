/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"fmt"
)

// V1beta1MaintenanceTaskState the model 'V1beta1MaintenanceTaskState'
type V1beta1MaintenanceTaskState string

// List of v1beta1MaintenanceTaskState
const (
	V1BETA1MAINTENANCETASKSTATE_PENDING V1beta1MaintenanceTaskState = "PENDING"
	V1BETA1MAINTENANCETASKSTATE_RUNNING V1beta1MaintenanceTaskState = "RUNNING"
	V1BETA1MAINTENANCETASKSTATE_DONE V1beta1MaintenanceTaskState = "DONE"
)

// All allowed values of V1beta1MaintenanceTaskState enum
var AllowedV1beta1MaintenanceTaskStateEnumValues = []V1beta1MaintenanceTaskState{
	"PENDING",
	"RUNNING",
	"DONE",
}

func (v *V1beta1MaintenanceTaskState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := V1beta1MaintenanceTaskState(value)
	for _, existing := range AllowedV1beta1MaintenanceTaskStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid V1beta1MaintenanceTaskState", value)
}

// NewV1beta1MaintenanceTaskStateFromValue returns a pointer to a valid V1beta1MaintenanceTaskState
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewV1beta1MaintenanceTaskStateFromValue(v string) (*V1beta1MaintenanceTaskState, error) {
	ev := V1beta1MaintenanceTaskState(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for V1beta1MaintenanceTaskState: valid values are %v", v, AllowedV1beta1MaintenanceTaskStateEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v V1beta1MaintenanceTaskState) IsValid() bool {
	for _, existing := range AllowedV1beta1MaintenanceTaskStateEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to v1beta1MaintenanceTaskState value
func (v V1beta1MaintenanceTaskState) Ptr() *V1beta1MaintenanceTaskState {
	return &v
}

type NullableV1beta1MaintenanceTaskState struct {
	value *V1beta1MaintenanceTaskState
	isSet bool
}

func (v NullableV1beta1MaintenanceTaskState) Get() *V1beta1MaintenanceTaskState {
	return v.value
}

func (v *NullableV1beta1MaintenanceTaskState) Set(val *V1beta1MaintenanceTaskState) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1MaintenanceTaskState) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1MaintenanceTaskState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1MaintenanceTaskState(val *V1beta1MaintenanceTaskState) *NullableV1beta1MaintenanceTaskState {
	return &NullableV1beta1MaintenanceTaskState{value: val, isSet: true}
}

func (v NullableV1beta1MaintenanceTaskState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1MaintenanceTaskState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
