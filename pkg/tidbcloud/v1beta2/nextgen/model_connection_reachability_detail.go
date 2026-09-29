/*
TiDB Cloud Premium API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Premium](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-premium) instances. For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/api/tidb-cloud-api-overview/).  The public connection setting endpoints also support TiDB Cloud Essential V2 instances.  # Overview  The TiDB Cloud Premium API (v1beta2) provides [REST](https://en.wikipedia.org/wiki/REST) endpoints to manage TiDB Cloud Premium instances and related resources.  You can use this API to manage the following resources:  - **TiDB Cloud Premium instance**: manage the lifecycle and configuration of TiDB Cloud Premium instances, including passwords, CA certificates, and cloud provider information. - **Public connection setting**: enable or disable public connections and manage IP access lists for Premium and Essential V2 instances. - **Customer-managed encryption key**: retrieve the IAM principal and verify KMS access before creating a Premium instance on AWS or Alibaba Cloud. - **Backup**: manage backups for TiDB Cloud Premium instances, including backup-based restore. - **Region**: retrieve available regions for deploying TiDB Cloud Premium instances.  # Get Started  This guide helps you make your first API call to the TiDB Cloud Premium API. You will learn how to authenticate a request, build a request, and interpret the response.  1. Create a [TiDB Cloud account](https://tidbcloud.com/signup) if you do not already have one. 2. In the [TiDB Cloud console](https://tidbcloud.com/), go to **Organization** > **API Keys** and create an API key. For more information, see [API key management](#section/Authentication/API-key-management). 3. Make your first API call.   To get all TiDB Cloud Premium instances in your organization, run the following command in your terminal. Replace `YOUR_PUBLIC_KEY` and `YOUR_PRIVATE_KEY` with your own key values.   ```bash  curl --digest \\    --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\    --request GET \\    --url 'https://cloud.tidbapi.com/v1beta2/tidbs' \\    --header 'Accept: application/json'  ```  4. The API returns a JSON list of your TiDB Cloud Premium instances. If none exist, the response contains an empty list.  # Authentication  The TiDB Cloud API supports [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication) with API keys and OAuth Bearer tokens. OAuth clients send the token in the `Authorization: Bearer <token>` header. It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. 5. Configure the role and scope for the API key. For more information about the permissions of a role, see [User roles](https://docs.pingcap.com/tidbcloud/manage-user-access/#user-roles). 6. Click **Generate API Key**. Copy and save the public key and the private key. 7. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 8. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Update Role**. 4. You can update the description and role of the API key. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # API Changelog  This changelog lists all changes to the TiDB Cloud Premium API (v1beta2).  <!-- In reverse chronological order -->  ## 20260924  - Published CMEK IAM principal retrieval and verification APIs for Premium instances on AWS and Alibaba Cloud.  ## 20260923  - Added APIs to get and update public connection settings for Premium and Essential V2 instances.  ## 20260428  - Initial release of the TiDB Cloud Premium API (v1beta2), including the following resources and endpoints:  * TiDB Cloud Premium instance   * [List TiDB Cloud Premium instances](#tag/TiDB-Instance/operation/TidbService_ListTidbs)   * [Create a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_CreateTidb)   * [Get a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetTidb)   * [Delete a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_DeleteTidb)   * [Update a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_UpdateTidb)   * [Reset the root password of a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_ResetRootPassword)   * [Get cloud provider information for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCloudProviderInfo)   * [Get the CA certificate download URL for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCaCertificateDownloadUrl)  * Backup   * [List backups for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_ListTidbBackups)   * [Delete a backup for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_DeleteTidbBackup)   * [Restore a TiDB Cloud Premium instance from a backup](#tag/Backup/operation/TidbService_RestoreTidb)   * [Get restore status for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_GetRestoreStatus)  * Region   * [List regions](#tag/Region/operation/RegionService_ListRegions)

API version: v1beta2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package nextgen

import (
	"encoding/json"
)

// checks if the ConnectionReachabilityDetail type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConnectionReachabilityDetail{}

// ConnectionReachabilityDetail struct for ConnectionReachabilityDetail
type ConnectionReachabilityDetail struct {
	// Indicates whether the service is active.
	ServiceActive *bool `json:"serviceActive,omitempty"`
	// Indicates whether the endpoint is active.
	EndpointActive *bool `json:"endpointActive,omitempty"`
	// Indicates whether the DNS is reachable.
	DnsReachable *bool `json:"dnsReachable,omitempty"`
	// A human-readable message that describes the reachability status.
	Message              *string `json:"message,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ConnectionReachabilityDetail ConnectionReachabilityDetail

// NewConnectionReachabilityDetail instantiates a new ConnectionReachabilityDetail object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConnectionReachabilityDetail() *ConnectionReachabilityDetail {
	this := ConnectionReachabilityDetail{}
	return &this
}

// NewConnectionReachabilityDetailWithDefaults instantiates a new ConnectionReachabilityDetail object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConnectionReachabilityDetailWithDefaults() *ConnectionReachabilityDetail {
	this := ConnectionReachabilityDetail{}
	return &this
}

// GetServiceActive returns the ServiceActive field value if set, zero value otherwise.
func (o *ConnectionReachabilityDetail) GetServiceActive() bool {
	if o == nil || IsNil(o.ServiceActive) {
		var ret bool
		return ret
	}
	return *o.ServiceActive
}

// GetServiceActiveOk returns a tuple with the ServiceActive field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionReachabilityDetail) GetServiceActiveOk() (*bool, bool) {
	if o == nil || IsNil(o.ServiceActive) {
		return nil, false
	}
	return o.ServiceActive, true
}

// HasServiceActive returns a boolean if a field has been set.
func (o *ConnectionReachabilityDetail) HasServiceActive() bool {
	if o != nil && !IsNil(o.ServiceActive) {
		return true
	}

	return false
}

// SetServiceActive gets a reference to the given bool and assigns it to the ServiceActive field.
func (o *ConnectionReachabilityDetail) SetServiceActive(v bool) {
	o.ServiceActive = &v
}

// GetEndpointActive returns the EndpointActive field value if set, zero value otherwise.
func (o *ConnectionReachabilityDetail) GetEndpointActive() bool {
	if o == nil || IsNil(o.EndpointActive) {
		var ret bool
		return ret
	}
	return *o.EndpointActive
}

// GetEndpointActiveOk returns a tuple with the EndpointActive field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionReachabilityDetail) GetEndpointActiveOk() (*bool, bool) {
	if o == nil || IsNil(o.EndpointActive) {
		return nil, false
	}
	return o.EndpointActive, true
}

// HasEndpointActive returns a boolean if a field has been set.
func (o *ConnectionReachabilityDetail) HasEndpointActive() bool {
	if o != nil && !IsNil(o.EndpointActive) {
		return true
	}

	return false
}

// SetEndpointActive gets a reference to the given bool and assigns it to the EndpointActive field.
func (o *ConnectionReachabilityDetail) SetEndpointActive(v bool) {
	o.EndpointActive = &v
}

// GetDnsReachable returns the DnsReachable field value if set, zero value otherwise.
func (o *ConnectionReachabilityDetail) GetDnsReachable() bool {
	if o == nil || IsNil(o.DnsReachable) {
		var ret bool
		return ret
	}
	return *o.DnsReachable
}

// GetDnsReachableOk returns a tuple with the DnsReachable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionReachabilityDetail) GetDnsReachableOk() (*bool, bool) {
	if o == nil || IsNil(o.DnsReachable) {
		return nil, false
	}
	return o.DnsReachable, true
}

// HasDnsReachable returns a boolean if a field has been set.
func (o *ConnectionReachabilityDetail) HasDnsReachable() bool {
	if o != nil && !IsNil(o.DnsReachable) {
		return true
	}

	return false
}

// SetDnsReachable gets a reference to the given bool and assigns it to the DnsReachable field.
func (o *ConnectionReachabilityDetail) SetDnsReachable(v bool) {
	o.DnsReachable = &v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *ConnectionReachabilityDetail) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConnectionReachabilityDetail) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *ConnectionReachabilityDetail) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *ConnectionReachabilityDetail) SetMessage(v string) {
	o.Message = &v
}

func (o ConnectionReachabilityDetail) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConnectionReachabilityDetail) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ServiceActive) {
		toSerialize["serviceActive"] = o.ServiceActive
	}
	if !IsNil(o.EndpointActive) {
		toSerialize["endpointActive"] = o.EndpointActive
	}
	if !IsNil(o.DnsReachable) {
		toSerialize["dnsReachable"] = o.DnsReachable
	}
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ConnectionReachabilityDetail) UnmarshalJSON(data []byte) (err error) {
	varConnectionReachabilityDetail := _ConnectionReachabilityDetail{}

	err = json.Unmarshal(data, &varConnectionReachabilityDetail)

	if err != nil {
		return err
	}

	*o = ConnectionReachabilityDetail(varConnectionReachabilityDetail)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "serviceActive")
		delete(additionalProperties, "endpointActive")
		delete(additionalProperties, "dnsReachable")
		delete(additionalProperties, "message")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableConnectionReachabilityDetail struct {
	value *ConnectionReachabilityDetail
	isSet bool
}

func (v NullableConnectionReachabilityDetail) Get() *ConnectionReachabilityDetail {
	return v.value
}

func (v *NullableConnectionReachabilityDetail) Set(val *ConnectionReachabilityDetail) {
	v.value = val
	v.isSet = true
}

func (v NullableConnectionReachabilityDetail) IsSet() bool {
	return v.isSet
}

func (v *NullableConnectionReachabilityDetail) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConnectionReachabilityDetail(val *ConnectionReachabilityDetail) *NullableConnectionReachabilityDetail {
	return &NullableConnectionReachabilityDetail{value: val, isSet: true}
}

func (v NullableConnectionReachabilityDetail) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConnectionReachabilityDetail) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
