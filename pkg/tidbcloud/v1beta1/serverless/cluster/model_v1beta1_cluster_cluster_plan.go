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

// V1beta1ClusterClusterPlan the model 'V1beta1ClusterClusterPlan'
type V1beta1ClusterClusterPlan string

// List of v1beta1ClusterClusterPlan
const (
	V1BETA1CLUSTERCLUSTERPLAN_FREE     V1beta1ClusterClusterPlan = "FREE"
	V1BETA1CLUSTERCLUSTERPLAN_SCALABLE V1beta1ClusterClusterPlan = "SCALABLE"
	V1BETA1CLUSTERCLUSTERPLAN_STARTER  V1beta1ClusterClusterPlan = "STARTER"
	V1BETA1CLUSTERCLUSTERPLAN_STANDARD V1beta1ClusterClusterPlan = "STANDARD"
)

// All allowed values of V1beta1ClusterClusterPlan enum
var AllowedV1beta1ClusterClusterPlanEnumValues = []V1beta1ClusterClusterPlan{
	"FREE",
	"SCALABLE",
	"STARTER",
	"STANDARD",
}

func (v *V1beta1ClusterClusterPlan) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := V1beta1ClusterClusterPlan(value)
	for _, existing := range AllowedV1beta1ClusterClusterPlanEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid V1beta1ClusterClusterPlan", value)
}

// NewV1beta1ClusterClusterPlanFromValue returns a pointer to a valid V1beta1ClusterClusterPlan
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewV1beta1ClusterClusterPlanFromValue(v string) (*V1beta1ClusterClusterPlan, error) {
	ev := V1beta1ClusterClusterPlan(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for V1beta1ClusterClusterPlan: valid values are %v", v, AllowedV1beta1ClusterClusterPlanEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v V1beta1ClusterClusterPlan) IsValid() bool {
	for _, existing := range AllowedV1beta1ClusterClusterPlanEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to v1beta1ClusterClusterPlan value
func (v V1beta1ClusterClusterPlan) Ptr() *V1beta1ClusterClusterPlan {
	return &v
}

type NullableV1beta1ClusterClusterPlan struct {
	value *V1beta1ClusterClusterPlan
	isSet bool
}

func (v NullableV1beta1ClusterClusterPlan) Get() *V1beta1ClusterClusterPlan {
	return v.value
}

func (v *NullableV1beta1ClusterClusterPlan) Set(val *V1beta1ClusterClusterPlan) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1ClusterClusterPlan) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1ClusterClusterPlan) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1ClusterClusterPlan(val *V1beta1ClusterClusterPlan) *NullableV1beta1ClusterClusterPlan {
	return &NullableV1beta1ClusterClusterPlan{value: val, isSet: true}
}

func (v NullableV1beta1ClusterClusterPlan) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1ClusterClusterPlan) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
