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

// V1beta1NetworkContainerState the model 'V1beta1NetworkContainerState'
type V1beta1NetworkContainerState string

// List of v1beta1NetworkContainerState
const (
	V1BETA1NETWORKCONTAINERSTATE_ACTIVE V1beta1NetworkContainerState = "ACTIVE"
	V1BETA1NETWORKCONTAINERSTATE_INACTIVE V1beta1NetworkContainerState = "INACTIVE"
)

// All allowed values of V1beta1NetworkContainerState enum
var AllowedV1beta1NetworkContainerStateEnumValues = []V1beta1NetworkContainerState{
	"ACTIVE",
	"INACTIVE",
}

func (v *V1beta1NetworkContainerState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := V1beta1NetworkContainerState(value)
	for _, existing := range AllowedV1beta1NetworkContainerStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid V1beta1NetworkContainerState", value)
}

// NewV1beta1NetworkContainerStateFromValue returns a pointer to a valid V1beta1NetworkContainerState
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewV1beta1NetworkContainerStateFromValue(v string) (*V1beta1NetworkContainerState, error) {
	ev := V1beta1NetworkContainerState(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for V1beta1NetworkContainerState: valid values are %v", v, AllowedV1beta1NetworkContainerStateEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v V1beta1NetworkContainerState) IsValid() bool {
	for _, existing := range AllowedV1beta1NetworkContainerStateEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to v1beta1NetworkContainerState value
func (v V1beta1NetworkContainerState) Ptr() *V1beta1NetworkContainerState {
	return &v
}

type NullableV1beta1NetworkContainerState struct {
	value *V1beta1NetworkContainerState
	isSet bool
}

func (v NullableV1beta1NetworkContainerState) Get() *V1beta1NetworkContainerState {
	return v.value
}

func (v *NullableV1beta1NetworkContainerState) Set(val *V1beta1NetworkContainerState) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1NetworkContainerState) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1NetworkContainerState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1NetworkContainerState(val *V1beta1NetworkContainerState) *NullableV1beta1NetworkContainerState {
	return &NullableV1beta1NetworkContainerState{value: val, isSet: true}
}

func (v NullableV1beta1NetworkContainerState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1NetworkContainerState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
