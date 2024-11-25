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

// PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter the model 'PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter'
type PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter string

// List of PrivateEndpointConnectionService_ListPrivateEndpointConnections_cloudProvider_parameter
const (
	PRIVATEENDPOINTCONNECTIONSERVICELISTPRIVATEENDPOINTCONNECTIONSCLOUDPROVIDERPARAMETER_AWS   PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter = "aws"
	PRIVATEENDPOINTCONNECTIONSERVICELISTPRIVATEENDPOINTCONNECTIONSCLOUDPROVIDERPARAMETER_GCP   PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter = "gcp"
	PRIVATEENDPOINTCONNECTIONSERVICELISTPRIVATEENDPOINTCONNECTIONSCLOUDPROVIDERPARAMETER_AZURE PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter = "azure"
)

// All allowed values of PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter enum
var AllowedPrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameterEnumValues = []PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter{
	"aws",
	"gcp",
	"azure",
}

func (v *PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter(value)
	for _, existing := range AllowedPrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameterEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter", value)
}

// NewPrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameterFromValue returns a pointer to a valid PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewPrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameterFromValue(v string) (*PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter, error) {
	ev := PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter: valid values are %v", v, AllowedPrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameterEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) IsValid() bool {
	for _, existing := range AllowedPrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameterEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to PrivateEndpointConnectionService_ListPrivateEndpointConnections_cloudProvider_parameter value
func (v PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) Ptr() *PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter {
	return &v
}

type NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter struct {
	value *PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter
	isSet bool
}

func (v NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) Get() *PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter {
	return v.value
}

func (v *NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) Set(val *PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter(val *PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) *NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter {
	return &NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter{value: val, isSet: true}
}

func (v NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
