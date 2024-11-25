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

// checks if the PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{}

// PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest struct for PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest
type PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest struct {
	Name                        *string `json:"name,omitempty"`
	PrivateEndpointConnectionId *string `json:"privateEndpointConnectionId,omitempty"`
	ClusterId                   *string `json:"clusterId,omitempty"`
	ClusterDisplayName          *string `json:"clusterDisplayName,omitempty"`
	// The labels of private link connection. It always contains the `project_id` label.
	Labels *map[string]string `json:"labels,omitempty"`
	// The endpoint ID of the private link connection. For AWS, it's VPC endpoint ID. For GCP, it's private service connect endpoint ID. For Azure, it's private endpoint resource ID.
	EndpointId string `json:"endpointId"`
	// The private IP address of the private endpoint in the user's vNet. TiDB Cloud will setup a public DNS record for this private IP address. So the user can use DNS address to connect to the cluster. Only available for Azure clusters.
	PrivateIpAddress         NullableString                          `json:"privateIpAddress,omitempty"`
	EndpointState            *PrivateEndpointConnectionEndpointState `json:"endpointState,omitempty"`
	Massage                  *string                                 `json:"massage,omitempty"`
	RegionId                 *string                                 `json:"regionId,omitempty"`
	RegionDisplayName        *string                                 `json:"regionDisplayName,omitempty"`
	CloudProvider            *V1beta1RegionCloudProvider             `json:"cloudProvider,omitempty"`
	PrivateLinkServiceName   *string                                 `json:"privateLinkServiceName,omitempty"`
	PrivateLinkServiceState  *V1beta1PrivateLinkServiceState         `json:"privateLinkServiceState,omitempty"`
	TidbNodeGroupDisplayName *string                                 `json:"tidbNodeGroupDisplayName,omitempty"`
	// Only for GCP private service connections. It's GCP project name.
	AccountId            NullableString `json:"accountId,omitempty"`
	Host                 *string        `json:"host,omitempty"`
	Port                 *int32         `json:"port,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest

// NewPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest instantiates a new PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest(endpointId string) *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	this := PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{}
	this.EndpointId = endpointId
	return &this
}

// NewPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequestWithDefaults instantiates a new PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequestWithDefaults() *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	this := PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetName(v string) {
	o.Name = &v
}

// GetPrivateEndpointConnectionId returns the PrivateEndpointConnectionId field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateEndpointConnectionId() string {
	if o == nil || IsNil(o.PrivateEndpointConnectionId) {
		var ret string
		return ret
	}
	return *o.PrivateEndpointConnectionId
}

// GetPrivateEndpointConnectionIdOk returns a tuple with the PrivateEndpointConnectionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateEndpointConnectionIdOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpointConnectionId) {
		return nil, false
	}
	return o.PrivateEndpointConnectionId, true
}

// HasPrivateEndpointConnectionId returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPrivateEndpointConnectionId() bool {
	if o != nil && !IsNil(o.PrivateEndpointConnectionId) {
		return true
	}

	return false
}

// SetPrivateEndpointConnectionId gets a reference to the given string and assigns it to the PrivateEndpointConnectionId field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateEndpointConnectionId(v string) {
	o.PrivateEndpointConnectionId = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}
	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetClusterDisplayName returns the ClusterDisplayName field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetClusterDisplayName() string {
	if o == nil || IsNil(o.ClusterDisplayName) {
		var ret string
		return ret
	}
	return *o.ClusterDisplayName
}

// GetClusterDisplayNameOk returns a tuple with the ClusterDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetClusterDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterDisplayName) {
		return nil, false
	}
	return o.ClusterDisplayName, true
}

// HasClusterDisplayName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasClusterDisplayName() bool {
	if o != nil && !IsNil(o.ClusterDisplayName) {
		return true
	}

	return false
}

// SetClusterDisplayName gets a reference to the given string and assigns it to the ClusterDisplayName field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetClusterDisplayName(v string) {
	o.ClusterDisplayName = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetEndpointId returns the EndpointId field value
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetEndpointId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EndpointId
}

// GetEndpointIdOk returns a tuple with the EndpointId field value
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetEndpointIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EndpointId, true
}

// SetEndpointId sets field value
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetEndpointId(v string) {
	o.EndpointId = v
}

// GetPrivateIpAddress returns the PrivateIpAddress field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateIpAddress() string {
	if o == nil || IsNil(o.PrivateIpAddress.Get()) {
		var ret string
		return ret
	}
	return *o.PrivateIpAddress.Get()
}

// GetPrivateIpAddressOk returns a tuple with the PrivateIpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateIpAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.PrivateIpAddress.Get(), o.PrivateIpAddress.IsSet()
}

// HasPrivateIpAddress returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPrivateIpAddress() bool {
	if o != nil && o.PrivateIpAddress.IsSet() {
		return true
	}

	return false
}

// SetPrivateIpAddress gets a reference to the given NullableString and assigns it to the PrivateIpAddress field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateIpAddress(v string) {
	o.PrivateIpAddress.Set(&v)
}

// SetPrivateIpAddressNil sets the value for PrivateIpAddress to be an explicit nil
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateIpAddressNil() {
	o.PrivateIpAddress.Set(nil)
}

// UnsetPrivateIpAddress ensures that no value is present for PrivateIpAddress, not even an explicit nil
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) UnsetPrivateIpAddress() {
	o.PrivateIpAddress.Unset()
}

// GetEndpointState returns the EndpointState field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetEndpointState() PrivateEndpointConnectionEndpointState {
	if o == nil || IsNil(o.EndpointState) {
		var ret PrivateEndpointConnectionEndpointState
		return ret
	}
	return *o.EndpointState
}

// GetEndpointStateOk returns a tuple with the EndpointState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetEndpointStateOk() (*PrivateEndpointConnectionEndpointState, bool) {
	if o == nil || IsNil(o.EndpointState) {
		return nil, false
	}
	return o.EndpointState, true
}

// HasEndpointState returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasEndpointState() bool {
	if o != nil && !IsNil(o.EndpointState) {
		return true
	}

	return false
}

// SetEndpointState gets a reference to the given PrivateEndpointConnectionEndpointState and assigns it to the EndpointState field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetEndpointState(v PrivateEndpointConnectionEndpointState) {
	o.EndpointState = &v
}

// GetMassage returns the Massage field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetMassage() string {
	if o == nil || IsNil(o.Massage) {
		var ret string
		return ret
	}
	return *o.Massage
}

// GetMassageOk returns a tuple with the Massage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetMassageOk() (*string, bool) {
	if o == nil || IsNil(o.Massage) {
		return nil, false
	}
	return o.Massage, true
}

// HasMassage returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasMassage() bool {
	if o != nil && !IsNil(o.Massage) {
		return true
	}

	return false
}

// SetMassage gets a reference to the given string and assigns it to the Massage field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetMassage(v string) {
	o.Massage = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetRegionId(v string) {
	o.RegionId = &v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetCloudProvider() V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret V1beta1RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetCloudProviderOk() (*V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given V1beta1RegionCloudProvider and assigns it to the CloudProvider field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetCloudProvider(v V1beta1RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetPrivateLinkServiceName returns the PrivateLinkServiceName field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateLinkServiceName() string {
	if o == nil || IsNil(o.PrivateLinkServiceName) {
		var ret string
		return ret
	}
	return *o.PrivateLinkServiceName
}

// GetPrivateLinkServiceNameOk returns a tuple with the PrivateLinkServiceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateLinkServiceNameOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateLinkServiceName) {
		return nil, false
	}
	return o.PrivateLinkServiceName, true
}

// HasPrivateLinkServiceName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPrivateLinkServiceName() bool {
	if o != nil && !IsNil(o.PrivateLinkServiceName) {
		return true
	}

	return false
}

// SetPrivateLinkServiceName gets a reference to the given string and assigns it to the PrivateLinkServiceName field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateLinkServiceName(v string) {
	o.PrivateLinkServiceName = &v
}

// GetPrivateLinkServiceState returns the PrivateLinkServiceState field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateLinkServiceState() V1beta1PrivateLinkServiceState {
	if o == nil || IsNil(o.PrivateLinkServiceState) {
		var ret V1beta1PrivateLinkServiceState
		return ret
	}
	return *o.PrivateLinkServiceState
}

// GetPrivateLinkServiceStateOk returns a tuple with the PrivateLinkServiceState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateLinkServiceStateOk() (*V1beta1PrivateLinkServiceState, bool) {
	if o == nil || IsNil(o.PrivateLinkServiceState) {
		return nil, false
	}
	return o.PrivateLinkServiceState, true
}

// HasPrivateLinkServiceState returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPrivateLinkServiceState() bool {
	if o != nil && !IsNil(o.PrivateLinkServiceState) {
		return true
	}

	return false
}

// SetPrivateLinkServiceState gets a reference to the given V1beta1PrivateLinkServiceState and assigns it to the PrivateLinkServiceState field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateLinkServiceState(v V1beta1PrivateLinkServiceState) {
	o.PrivateLinkServiceState = &v
}

// GetTidbNodeGroupDisplayName returns the TidbNodeGroupDisplayName field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetTidbNodeGroupDisplayName() string {
	if o == nil || IsNil(o.TidbNodeGroupDisplayName) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupDisplayName
}

// GetTidbNodeGroupDisplayNameOk returns a tuple with the TidbNodeGroupDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetTidbNodeGroupDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.TidbNodeGroupDisplayName) {
		return nil, false
	}
	return o.TidbNodeGroupDisplayName, true
}

// HasTidbNodeGroupDisplayName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasTidbNodeGroupDisplayName() bool {
	if o != nil && !IsNil(o.TidbNodeGroupDisplayName) {
		return true
	}

	return false
}

// SetTidbNodeGroupDisplayName gets a reference to the given string and assigns it to the TidbNodeGroupDisplayName field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetTidbNodeGroupDisplayName(v string) {
	o.TidbNodeGroupDisplayName = &v
}

// GetAccountId returns the AccountId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetAccountId() string {
	if o == nil || IsNil(o.AccountId.Get()) {
		var ret string
		return ret
	}
	return *o.AccountId.Get()
}

// GetAccountIdOk returns a tuple with the AccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetAccountIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.AccountId.Get(), o.AccountId.IsSet()
}

// HasAccountId returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasAccountId() bool {
	if o != nil && o.AccountId.IsSet() {
		return true
	}

	return false
}

// SetAccountId gets a reference to the given NullableString and assigns it to the AccountId field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetAccountId(v string) {
	o.AccountId.Set(&v)
}

// SetAccountIdNil sets the value for AccountId to be an explicit nil
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetAccountIdNil() {
	o.AccountId.Set(nil)
}

// UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) UnsetAccountId() {
	o.AccountId.Unset()
}

// GetHost returns the Host field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetHost() string {
	if o == nil || IsNil(o.Host) {
		var ret string
		return ret
	}
	return *o.Host
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetHostOk() (*string, bool) {
	if o == nil || IsNil(o.Host) {
		return nil, false
	}
	return o.Host, true
}

// HasHost returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasHost() bool {
	if o != nil && !IsNil(o.Host) {
		return true
	}

	return false
}

// SetHost gets a reference to the given string and assigns it to the Host field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetHost(v string) {
	o.Host = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPort() int32 {
	if o == nil || IsNil(o.Port) {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPortOk() (*int32, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPort(v int32) {
	o.Port = &v
}

func (o PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.PrivateEndpointConnectionId) {
		toSerialize["privateEndpointConnectionId"] = o.PrivateEndpointConnectionId
	}
	if !IsNil(o.ClusterId) {
		toSerialize["clusterId"] = o.ClusterId
	}
	if !IsNil(o.ClusterDisplayName) {
		toSerialize["clusterDisplayName"] = o.ClusterDisplayName
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	toSerialize["endpointId"] = o.EndpointId
	if o.PrivateIpAddress.IsSet() {
		toSerialize["privateIpAddress"] = o.PrivateIpAddress.Get()
	}
	if !IsNil(o.EndpointState) {
		toSerialize["endpointState"] = o.EndpointState
	}
	if !IsNil(o.Massage) {
		toSerialize["massage"] = o.Massage
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
	if !IsNil(o.PrivateLinkServiceName) {
		toSerialize["privateLinkServiceName"] = o.PrivateLinkServiceName
	}
	if !IsNil(o.PrivateLinkServiceState) {
		toSerialize["privateLinkServiceState"] = o.PrivateLinkServiceState
	}
	if !IsNil(o.TidbNodeGroupDisplayName) {
		toSerialize["tidbNodeGroupDisplayName"] = o.TidbNodeGroupDisplayName
	}
	if o.AccountId.IsSet() {
		toSerialize["accountId"] = o.AccountId.Get()
	}
	if !IsNil(o.Host) {
		toSerialize["host"] = o.Host
	}
	if !IsNil(o.Port) {
		toSerialize["port"] = o.Port
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"endpointId",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest := _PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{}

	err = json.Unmarshal(data, &varPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest)

	if err != nil {
		return err
	}

	*o = PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest(varPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "privateEndpointConnectionId")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "clusterDisplayName")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "endpointId")
		delete(additionalProperties, "privateIpAddress")
		delete(additionalProperties, "endpointState")
		delete(additionalProperties, "massage")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "privateLinkServiceName")
		delete(additionalProperties, "privateLinkServiceState")
		delete(additionalProperties, "tidbNodeGroupDisplayName")
		delete(additionalProperties, "accountId")
		delete(additionalProperties, "host")
		delete(additionalProperties, "port")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest struct {
	value *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest
	isSet bool
}

func (v NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) Get() *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	return v.value
}

func (v *NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) Set(val *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest(val *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) *NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	return &NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{value: val, isSet: true}
}

func (v NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
