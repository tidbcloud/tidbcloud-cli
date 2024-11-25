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

// ClusterServiceListNodeInstancesComponentTypeParameter the model 'ClusterServiceListNodeInstancesComponentTypeParameter'
type ClusterServiceListNodeInstancesComponentTypeParameter string

// List of ClusterService_ListNodeInstances_componentType_parameter
const (
	CLUSTERSERVICELISTNODEINSTANCESCOMPONENTTYPEPARAMETER_TIKV    ClusterServiceListNodeInstancesComponentTypeParameter = "TIKV"
	CLUSTERSERVICELISTNODEINSTANCESCOMPONENTTYPEPARAMETER_TIDB    ClusterServiceListNodeInstancesComponentTypeParameter = "TIDB"
	CLUSTERSERVICELISTNODEINSTANCESCOMPONENTTYPEPARAMETER_TIFLASH ClusterServiceListNodeInstancesComponentTypeParameter = "TIFLASH"
	CLUSTERSERVICELISTNODEINSTANCESCOMPONENTTYPEPARAMETER_PD      ClusterServiceListNodeInstancesComponentTypeParameter = "PD"
	CLUSTERSERVICELISTNODEINSTANCESCOMPONENTTYPEPARAMETER_TIPROXY ClusterServiceListNodeInstancesComponentTypeParameter = "TIPROXY"
)

// All allowed values of ClusterServiceListNodeInstancesComponentTypeParameter enum
var AllowedClusterServiceListNodeInstancesComponentTypeParameterEnumValues = []ClusterServiceListNodeInstancesComponentTypeParameter{
	"TIKV",
	"TIDB",
	"TIFLASH",
	"PD",
	"TIPROXY",
}

func (v *ClusterServiceListNodeInstancesComponentTypeParameter) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ClusterServiceListNodeInstancesComponentTypeParameter(value)
	for _, existing := range AllowedClusterServiceListNodeInstancesComponentTypeParameterEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ClusterServiceListNodeInstancesComponentTypeParameter", value)
}

// NewClusterServiceListNodeInstancesComponentTypeParameterFromValue returns a pointer to a valid ClusterServiceListNodeInstancesComponentTypeParameter
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewClusterServiceListNodeInstancesComponentTypeParameterFromValue(v string) (*ClusterServiceListNodeInstancesComponentTypeParameter, error) {
	ev := ClusterServiceListNodeInstancesComponentTypeParameter(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ClusterServiceListNodeInstancesComponentTypeParameter: valid values are %v", v, AllowedClusterServiceListNodeInstancesComponentTypeParameterEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ClusterServiceListNodeInstancesComponentTypeParameter) IsValid() bool {
	for _, existing := range AllowedClusterServiceListNodeInstancesComponentTypeParameterEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ClusterService_ListNodeInstances_componentType_parameter value
func (v ClusterServiceListNodeInstancesComponentTypeParameter) Ptr() *ClusterServiceListNodeInstancesComponentTypeParameter {
	return &v
}

type NullableClusterServiceListNodeInstancesComponentTypeParameter struct {
	value *ClusterServiceListNodeInstancesComponentTypeParameter
	isSet bool
}

func (v NullableClusterServiceListNodeInstancesComponentTypeParameter) Get() *ClusterServiceListNodeInstancesComponentTypeParameter {
	return v.value
}

func (v *NullableClusterServiceListNodeInstancesComponentTypeParameter) Set(val *ClusterServiceListNodeInstancesComponentTypeParameter) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterServiceListNodeInstancesComponentTypeParameter) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterServiceListNodeInstancesComponentTypeParameter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterServiceListNodeInstancesComponentTypeParameter(val *ClusterServiceListNodeInstancesComponentTypeParameter) *NullableClusterServiceListNodeInstancesComponentTypeParameter {
	return &NullableClusterServiceListNodeInstancesComponentTypeParameter{value: val, isSet: true}
}

func (v NullableClusterServiceListNodeInstancesComponentTypeParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterServiceListNodeInstancesComponentTypeParameter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
