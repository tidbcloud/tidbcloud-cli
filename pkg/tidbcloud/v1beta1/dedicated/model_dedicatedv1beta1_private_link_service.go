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

// checks if the Dedicatedv1beta1PrivateLinkService type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1PrivateLinkService{}

// Dedicatedv1beta1PrivateLinkService struct for Dedicatedv1beta1PrivateLinkService
type Dedicatedv1beta1PrivateLinkService struct {
	Name *string `json:"name,omitempty"`
	TidbNodeGroupId *string `json:"tidbNodeGroupId,omitempty"`
	// For AWS, it's the service name of the Private Link Service. For GCP, it's the resource name of the service attachment. For Azure, it's service resource ID of the Private Link Service.
	ServiceName *string `json:"serviceName,omitempty"`
	// For AWS, it's the fully qualified domain name (FQDN) shared for all private endpoints, despite which VPC the endpoint located in. For GCP, it's the zone name (suffix of FQDN) shared for all private endpoints located in a single VPC network. The format of FQDN is `<endpoint_name>.<service_dns_name>`. For Azure, it's the zone name shared across public internet. The format of FQDN is `<endpoint_name>-<random_hash>.<service_dns_name>`.
	ServiceDnsName *string `json:"serviceDnsName,omitempty"`
	// Only available for AWS. Same as the `AvailabilityZones` field in response body of `github.com/aws/aws-sdk-go-v2/service/ec2.DescribeVpcEndpointServices` method.
	AvailableZones []string `json:"availableZones,omitempty"`
	State *V1beta1PrivateLinkServiceState `json:"state,omitempty"`
	RegionId *string `json:"regionId,omitempty"`
	RegionDisplayName *string `json:"regionDisplayName,omitempty"`
	CloudProvider *V1beta1RegionCloudProvider `json:"cloudProvider,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1PrivateLinkService Dedicatedv1beta1PrivateLinkService

// NewDedicatedv1beta1PrivateLinkService instantiates a new Dedicatedv1beta1PrivateLinkService object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1PrivateLinkService() *Dedicatedv1beta1PrivateLinkService {
	this := Dedicatedv1beta1PrivateLinkService{}
	return &this
}

// NewDedicatedv1beta1PrivateLinkServiceWithDefaults instantiates a new Dedicatedv1beta1PrivateLinkService object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1PrivateLinkServiceWithDefaults() *Dedicatedv1beta1PrivateLinkService {
	this := Dedicatedv1beta1PrivateLinkService{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1PrivateLinkService) SetName(v string) {
	o.Name = &v
}

// GetTidbNodeGroupId returns the TidbNodeGroupId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetTidbNodeGroupId() string {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupId
}

// GetTidbNodeGroupIdOk returns a tuple with the TidbNodeGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetTidbNodeGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		return nil, false
	}
	return o.TidbNodeGroupId, true
}

// HasTidbNodeGroupId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasTidbNodeGroupId() bool {
	if o != nil && !IsNil(o.TidbNodeGroupId) {
		return true
	}

	return false
}

// SetTidbNodeGroupId gets a reference to the given string and assigns it to the TidbNodeGroupId field.
func (o *Dedicatedv1beta1PrivateLinkService) SetTidbNodeGroupId(v string) {
	o.TidbNodeGroupId = &v
}

// GetServiceName returns the ServiceName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetServiceName() string {
	if o == nil || IsNil(o.ServiceName) {
		var ret string
		return ret
	}
	return *o.ServiceName
}

// GetServiceNameOk returns a tuple with the ServiceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetServiceNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceName) {
		return nil, false
	}
	return o.ServiceName, true
}

// HasServiceName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasServiceName() bool {
	if o != nil && !IsNil(o.ServiceName) {
		return true
	}

	return false
}

// SetServiceName gets a reference to the given string and assigns it to the ServiceName field.
func (o *Dedicatedv1beta1PrivateLinkService) SetServiceName(v string) {
	o.ServiceName = &v
}

// GetServiceDnsName returns the ServiceDnsName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetServiceDnsName() string {
	if o == nil || IsNil(o.ServiceDnsName) {
		var ret string
		return ret
	}
	return *o.ServiceDnsName
}

// GetServiceDnsNameOk returns a tuple with the ServiceDnsName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetServiceDnsNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceDnsName) {
		return nil, false
	}
	return o.ServiceDnsName, true
}

// HasServiceDnsName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasServiceDnsName() bool {
	if o != nil && !IsNil(o.ServiceDnsName) {
		return true
	}

	return false
}

// SetServiceDnsName gets a reference to the given string and assigns it to the ServiceDnsName field.
func (o *Dedicatedv1beta1PrivateLinkService) SetServiceDnsName(v string) {
	o.ServiceDnsName = &v
}

// GetAvailableZones returns the AvailableZones field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetAvailableZones() []string {
	if o == nil || IsNil(o.AvailableZones) {
		var ret []string
		return ret
	}
	return o.AvailableZones
}

// GetAvailableZonesOk returns a tuple with the AvailableZones field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetAvailableZonesOk() ([]string, bool) {
	if o == nil || IsNil(o.AvailableZones) {
		return nil, false
	}
	return o.AvailableZones, true
}

// HasAvailableZones returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasAvailableZones() bool {
	if o != nil && !IsNil(o.AvailableZones) {
		return true
	}

	return false
}

// SetAvailableZones gets a reference to the given []string and assigns it to the AvailableZones field.
func (o *Dedicatedv1beta1PrivateLinkService) SetAvailableZones(v []string) {
	o.AvailableZones = v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetState() V1beta1PrivateLinkServiceState {
	if o == nil || IsNil(o.State) {
		var ret V1beta1PrivateLinkServiceState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetStateOk() (*V1beta1PrivateLinkServiceState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given V1beta1PrivateLinkServiceState and assigns it to the State field.
func (o *Dedicatedv1beta1PrivateLinkService) SetState(v V1beta1PrivateLinkServiceState) {
	o.State = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *Dedicatedv1beta1PrivateLinkService) SetRegionId(v string) {
	o.RegionId = &v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *Dedicatedv1beta1PrivateLinkService) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetCloudProvider() V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret V1beta1RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetCloudProviderOk() (*V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given V1beta1RegionCloudProvider and assigns it to the CloudProvider field.
func (o *Dedicatedv1beta1PrivateLinkService) SetCloudProvider(v V1beta1RegionCloudProvider) {
	o.CloudProvider = &v
}

func (o Dedicatedv1beta1PrivateLinkService) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1PrivateLinkService) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.TidbNodeGroupId) {
		toSerialize["tidbNodeGroupId"] = o.TidbNodeGroupId
	}
	if !IsNil(o.ServiceName) {
		toSerialize["serviceName"] = o.ServiceName
	}
	if !IsNil(o.ServiceDnsName) {
		toSerialize["serviceDnsName"] = o.ServiceDnsName
	}
	if !IsNil(o.AvailableZones) {
		toSerialize["availableZones"] = o.AvailableZones
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.RegionId) {
		toSerialize["regionId"] = o.RegionId
	}
	if !IsNil(o.RegionDisplayName) {
		toSerialize["regionDisplayName"] = o.RegionDisplayName
	}
	if !IsNil(o.CloudProvider) {
		toSerialize["cloudProvider"] = o.CloudProvider
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1PrivateLinkService) UnmarshalJSON(data []byte) (err error) {
	varDedicatedv1beta1PrivateLinkService := _Dedicatedv1beta1PrivateLinkService{}

	err = json.Unmarshal(data, &varDedicatedv1beta1PrivateLinkService)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1PrivateLinkService(varDedicatedv1beta1PrivateLinkService)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "tidbNodeGroupId")
		delete(additionalProperties, "serviceName")
		delete(additionalProperties, "serviceDnsName")
		delete(additionalProperties, "availableZones")
		delete(additionalProperties, "state")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "cloudProvider")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1PrivateLinkService struct {
	value *Dedicatedv1beta1PrivateLinkService
	isSet bool
}

func (v NullableDedicatedv1beta1PrivateLinkService) Get() *Dedicatedv1beta1PrivateLinkService {
	return v.value
}

func (v *NullableDedicatedv1beta1PrivateLinkService) Set(val *Dedicatedv1beta1PrivateLinkService) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1PrivateLinkService) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1PrivateLinkService) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1PrivateLinkService(val *Dedicatedv1beta1PrivateLinkService) *NullableDedicatedv1beta1PrivateLinkService {
	return &NullableDedicatedv1beta1PrivateLinkService{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1PrivateLinkService) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1PrivateLinkService) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

