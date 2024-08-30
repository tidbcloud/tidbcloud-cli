/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cluster

import (
	"encoding/json"
	"fmt"
)

// BranchState Output Only. Branch State.   - CREATING: The branch is being created.  - ACTIVE: The branch is active and running.  - DELETED: The branch is being deleted.  - MAINTENANCE: The branch is under maintenance.
type BranchState string

// List of Branch.State
const (
	BRANCHSTATE_CREATING    BranchState = "CREATING"
	BRANCHSTATE_ACTIVE      BranchState = "ACTIVE"
	BRANCHSTATE_DELETED     BranchState = "DELETED"
	BRANCHSTATE_MAINTENANCE BranchState = "MAINTENANCE"
)

// All allowed values of BranchState enum
var AllowedBranchStateEnumValues = []BranchState{
	"CREATING",
	"ACTIVE",
	"DELETED",
	"MAINTENANCE",
}

func (v *BranchState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := BranchState(value)
	for _, existing := range AllowedBranchStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid BranchState", value)
}

// NewBranchStateFromValue returns a pointer to a valid BranchState
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewBranchStateFromValue(v string) (*BranchState, error) {
	ev := BranchState(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for BranchState: valid values are %v", v, AllowedBranchStateEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v BranchState) IsValid() bool {
	for _, existing := range AllowedBranchStateEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to Branch.State value
func (v BranchState) Ptr() *BranchState {
	return &v
}

type NullableBranchState struct {
	value *BranchState
	isSet bool
}

func (v NullableBranchState) Get() *BranchState {
	return v.value
}

func (v *NullableBranchState) Set(val *BranchState) {
	v.value = val
	v.isSet = true
}

func (v NullableBranchState) IsSet() bool {
	return v.isSet
}

func (v *NullableBranchState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBranchState(val *BranchState) *NullableBranchState {
	return &NullableBranchState{value: val, isSet: true}
}

func (v NullableBranchState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBranchState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
