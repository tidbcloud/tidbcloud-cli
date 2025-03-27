/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cluster

import (
	"encoding/json"
)

// Serverlessv1beta1ClusterView Enum for the different types of detail view to be returned for a TiDB Cloud Serverless cluster.   - BASIC: Only basic information about the cluster is returned.  - FULL: All details about the cluster are returned.
type Serverlessv1beta1ClusterView string

// List of serverlessv1beta1ClusterView
const (
	SERVERLESSV1BETA1CLUSTERVIEW_BASIC Serverlessv1beta1ClusterView = "BASIC"
	SERVERLESSV1BETA1CLUSTERVIEW_FULL  Serverlessv1beta1ClusterView = "FULL"

	// Unknown value for handling new enum values gracefully
	Serverlessv1beta1ClusterView_UNKNOWN Serverlessv1beta1ClusterView = "UNKNOWN"
)

// All allowed values of Serverlessv1beta1ClusterView enum
var AllowedServerlessv1beta1ClusterViewEnumValues = []Serverlessv1beta1ClusterView{
	"BASIC",
	"FULL",
}

func (v *Serverlessv1beta1ClusterView) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Serverlessv1beta1ClusterView(value)
	for _, existing := range AllowedServerlessv1beta1ClusterViewEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	// Instead of returning an error, assign UNKNOWN value
	*v = Serverlessv1beta1ClusterView_UNKNOWN
	return nil
}

// NewServerlessv1beta1ClusterViewFromValue returns a pointer to a valid Serverlessv1beta1ClusterView
// for the value passed as argument, or UNKNOWN if the value is not in the enum list
func NewServerlessv1beta1ClusterViewFromValue(v string) *Serverlessv1beta1ClusterView {
	ev := Serverlessv1beta1ClusterView(v)
	if ev.IsValid() {
		return &ev
	}
	unknown := Serverlessv1beta1ClusterView_UNKNOWN
	return &unknown
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Serverlessv1beta1ClusterView) IsValid() bool {
	for _, existing := range AllowedServerlessv1beta1ClusterViewEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to serverlessv1beta1ClusterView value
func (v Serverlessv1beta1ClusterView) Ptr() *Serverlessv1beta1ClusterView {
	return &v
}

type NullableServerlessv1beta1ClusterView struct {
	value *Serverlessv1beta1ClusterView
	isSet bool
}

func (v NullableServerlessv1beta1ClusterView) Get() *Serverlessv1beta1ClusterView {
	return v.value
}

func (v *NullableServerlessv1beta1ClusterView) Set(val *Serverlessv1beta1ClusterView) {
	v.value = val
	v.isSet = true
}

func (v NullableServerlessv1beta1ClusterView) IsSet() bool {
	return v.isSet
}

func (v *NullableServerlessv1beta1ClusterView) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServerlessv1beta1ClusterView(val *Serverlessv1beta1ClusterView) *NullableServerlessv1beta1ClusterView {
	return &NullableServerlessv1beta1ClusterView{value: val, isSet: true}
}

func (v NullableServerlessv1beta1ClusterView) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServerlessv1beta1ClusterView) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
