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

// BranchView View on branch. Pass this enum to control which subsets of fields to get.   - BASIC: Basic response contains basic information for a branch.  - FULL: FULL response contains all detailed information for a branch.
type BranchView string

// List of BranchView
const (
	BRANCHVIEW_BASIC BranchView = "BASIC"
	BRANCHVIEW_FULL  BranchView = "FULL"
)

// All allowed values of BranchView enum
var AllowedBranchViewEnumValues = []BranchView{
	"BASIC",
	"FULL",
}

func (v *BranchView) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := BranchView(value)
	for _, existing := range AllowedBranchViewEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid BranchView", value)
}

// NewBranchViewFromValue returns a pointer to a valid BranchView
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewBranchViewFromValue(v string) (*BranchView, error) {
	ev := BranchView(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for BranchView: valid values are %v", v, AllowedBranchViewEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v BranchView) IsValid() bool {
	for _, existing := range AllowedBranchViewEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to BranchView value
func (v BranchView) Ptr() *BranchView {
	return &v
}

type NullableBranchView struct {
	value *BranchView
	isSet bool
}

func (v NullableBranchView) Get() *BranchView {
	return v.value
}

func (v *NullableBranchView) Set(val *BranchView) {
	v.value = val
	v.isSet = true
}

func (v NullableBranchView) IsSet() bool {
	return v.isSet
}

func (v *NullableBranchView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBranchView(val *BranchView) *NullableBranchView {
	return &NullableBranchView{value: val, isSet: true}
}

func (v NullableBranchView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBranchView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}