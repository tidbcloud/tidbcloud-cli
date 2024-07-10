/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package branch

import (
	"encoding/json"
	"fmt"
)

// V1beta1BranchView View on branch. Pass this enum to control which subsets of fields to get.   - BASIC: Basic response contains basic information for a branch.  - FULL: FULL response contains all detailed information for a branch.
type V1beta1BranchView string

// List of v1beta1BranchView
const (
	BASIC V1beta1BranchView = "BASIC"
	FULL V1beta1BranchView = "FULL"
)

// All allowed values of V1beta1BranchView enum
var AllowedV1beta1BranchViewEnumValues = []V1beta1BranchView{
	"BASIC",
	"FULL",
}

func (v *V1beta1BranchView) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := V1beta1BranchView(value)
	for _, existing := range AllowedV1beta1BranchViewEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid V1beta1BranchView", value)
}

// NewV1beta1BranchViewFromValue returns a pointer to a valid V1beta1BranchView
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewV1beta1BranchViewFromValue(v string) (*V1beta1BranchView, error) {
	ev := V1beta1BranchView(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for V1beta1BranchView: valid values are %v", v, AllowedV1beta1BranchViewEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v V1beta1BranchView) IsValid() bool {
	for _, existing := range AllowedV1beta1BranchViewEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to v1beta1BranchView value
func (v V1beta1BranchView) Ptr() *V1beta1BranchView {
	return &v
}

type NullableV1beta1BranchView struct {
	value *V1beta1BranchView
	isSet bool
}

func (v NullableV1beta1BranchView) Get() *V1beta1BranchView {
	return v.value
}

func (v *NullableV1beta1BranchView) Set(val *V1beta1BranchView) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1BranchView) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1BranchView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1BranchView(val *V1beta1BranchView) *NullableV1beta1BranchView {
	return &NullableV1beta1BranchView{value: val, isSet: true}
}

func (v NullableV1beta1BranchView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1BranchView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

