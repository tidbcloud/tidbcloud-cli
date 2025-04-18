/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
)

// EndpointConnectionType  - PUBLIC: The endpoint is a public endpoint.  - VPC_PEERING: The endpoint is a VPC peering endpoint.  - PRIVATE_ENDPOINT: The endpoint is a private link endpoint.
type EndpointConnectionType string

// List of EndpointConnectionType
const (
	ENDPOINTCONNECTIONTYPE_PUBLIC           EndpointConnectionType = "PUBLIC"
	ENDPOINTCONNECTIONTYPE_VPC_PEERING      EndpointConnectionType = "VPC_PEERING"
	ENDPOINTCONNECTIONTYPE_PRIVATE_ENDPOINT EndpointConnectionType = "PRIVATE_ENDPOINT"
)

// All allowed values of EndpointConnectionType enum
var AllowedEndpointConnectionTypeEnumValues = []EndpointConnectionType{
	"PUBLIC",
	"VPC_PEERING",
	"PRIVATE_ENDPOINT",
}

func (v *EndpointConnectionType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := EndpointConnectionType(value)
	for _, existing := range AllowedEndpointConnectionTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	*v = EndpointConnectionType(value)
	return nil
}

// NewEndpointConnectionTypeFromValue returns a pointer to a valid EndpointConnectionType for the value passed as argument
func NewEndpointConnectionTypeFromValue(v string) *EndpointConnectionType {
	ev := EndpointConnectionType(v)
	return &ev
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v EndpointConnectionType) IsValid() bool {
	for _, existing := range AllowedEndpointConnectionTypeEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to EndpointConnectionType value
func (v EndpointConnectionType) Ptr() *EndpointConnectionType {
	return &v
}

type NullableEndpointConnectionType struct {
	value *EndpointConnectionType
	isSet bool
}

func (v NullableEndpointConnectionType) Get() *EndpointConnectionType {
	return v.value
}

func (v *NullableEndpointConnectionType) Set(val *EndpointConnectionType) {
	v.value = val
	v.isSet = true
}

func (v NullableEndpointConnectionType) IsSet() bool {
	return v.isSet
}

func (v *NullableEndpointConnectionType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEndpointConnectionType(val *EndpointConnectionType) *NullableEndpointConnectionType {
	return &NullableEndpointConnectionType{value: val, isSet: true}
}

func (v NullableEndpointConnectionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEndpointConnectionType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
