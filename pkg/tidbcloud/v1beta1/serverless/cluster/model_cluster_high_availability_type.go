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

// ClusterHighAvailabilityType the model 'ClusterHighAvailabilityType'
type ClusterHighAvailabilityType string

// List of ClusterHighAvailabilityType
const (
	CLUSTERHIGHAVAILABILITYTYPE_ZONAL    ClusterHighAvailabilityType = "ZONAL"
	CLUSTERHIGHAVAILABILITYTYPE_REGIONAL ClusterHighAvailabilityType = "REGIONAL"
)

// All allowed values of ClusterHighAvailabilityType enum
var AllowedClusterHighAvailabilityTypeEnumValues = []ClusterHighAvailabilityType{
	"ZONAL",
	"REGIONAL",
}

func (v *ClusterHighAvailabilityType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ClusterHighAvailabilityType(value)
	for _, existing := range AllowedClusterHighAvailabilityTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	*v = ClusterHighAvailabilityType(value)
	return nil
}

// NewClusterHighAvailabilityTypeFromValue returns a pointer to a valid ClusterHighAvailabilityType for the value passed as argument
func NewClusterHighAvailabilityTypeFromValue(v string) *ClusterHighAvailabilityType {
	ev := ClusterHighAvailabilityType(v)
	return &ev
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ClusterHighAvailabilityType) IsValid() bool {
	for _, existing := range AllowedClusterHighAvailabilityTypeEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ClusterHighAvailabilityType value
func (v ClusterHighAvailabilityType) Ptr() *ClusterHighAvailabilityType {
	return &v
}

type NullableClusterHighAvailabilityType struct {
	value *ClusterHighAvailabilityType
	isSet bool
}

func (v NullableClusterHighAvailabilityType) Get() *ClusterHighAvailabilityType {
	return v.value
}

func (v *NullableClusterHighAvailabilityType) Set(val *ClusterHighAvailabilityType) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterHighAvailabilityType) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterHighAvailabilityType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterHighAvailabilityType(val *ClusterHighAvailabilityType) *NullableClusterHighAvailabilityType {
	return &NullableClusterHighAvailabilityType{value: val, isSet: true}
}

func (v NullableClusterHighAvailabilityType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterHighAvailabilityType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
