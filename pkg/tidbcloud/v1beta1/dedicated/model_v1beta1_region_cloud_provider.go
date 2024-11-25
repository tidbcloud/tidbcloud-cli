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

// V1beta1RegionCloudProvider Enum of cloud provider names.   - aws: Amazon Web Services. buf:lint:ignore ENUM_VALUE_UPPER_SNAKE_CASE  - gcp: Google Cloud Platform. buf:lint:ignore ENUM_VALUE_UPPER_SNAKE_CASE  - azure: Microsoft Azure. buf:lint:ignore ENUM_VALUE_UPPER_SNAKE_CASE
type V1beta1RegionCloudProvider string

// List of v1beta1RegionCloudProvider
const (
	V1BETA1REGIONCLOUDPROVIDER_AWS   V1beta1RegionCloudProvider = "aws"
	V1BETA1REGIONCLOUDPROVIDER_GCP   V1beta1RegionCloudProvider = "gcp"
	V1BETA1REGIONCLOUDPROVIDER_AZURE V1beta1RegionCloudProvider = "azure"
)

// All allowed values of V1beta1RegionCloudProvider enum
var AllowedV1beta1RegionCloudProviderEnumValues = []V1beta1RegionCloudProvider{
	"aws",
	"gcp",
	"azure",
}

func (v *V1beta1RegionCloudProvider) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := V1beta1RegionCloudProvider(value)
	for _, existing := range AllowedV1beta1RegionCloudProviderEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid V1beta1RegionCloudProvider", value)
}

// NewV1beta1RegionCloudProviderFromValue returns a pointer to a valid V1beta1RegionCloudProvider
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewV1beta1RegionCloudProviderFromValue(v string) (*V1beta1RegionCloudProvider, error) {
	ev := V1beta1RegionCloudProvider(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for V1beta1RegionCloudProvider: valid values are %v", v, AllowedV1beta1RegionCloudProviderEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v V1beta1RegionCloudProvider) IsValid() bool {
	for _, existing := range AllowedV1beta1RegionCloudProviderEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to v1beta1RegionCloudProvider value
func (v V1beta1RegionCloudProvider) Ptr() *V1beta1RegionCloudProvider {
	return &v
}

type NullableV1beta1RegionCloudProvider struct {
	value *V1beta1RegionCloudProvider
	isSet bool
}

func (v NullableV1beta1RegionCloudProvider) Get() *V1beta1RegionCloudProvider {
	return v.value
}

func (v *NullableV1beta1RegionCloudProvider) Set(val *V1beta1RegionCloudProvider) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1RegionCloudProvider) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1RegionCloudProvider) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1RegionCloudProvider(val *V1beta1RegionCloudProvider) *NullableV1beta1RegionCloudProvider {
	return &NullableV1beta1RegionCloudProvider{value: val, isSet: true}
}

func (v NullableV1beta1RegionCloudProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1RegionCloudProvider) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
