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

// checks if the TidbEndpoint type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbEndpoint{}

// TidbEndpoint struct for TidbEndpoint
type TidbEndpoint struct {
	// The hostname of the endpoint.  When `connectionType` is `\"PRIVATE_ENDPOINT\"`, the behavior of the `host` field varies by cloud provider: - For AWS, this field is available after the Private Link service is ready. - For Google Cloud and Azure, this field is available after you create a private endpoint connection.
	Host *string `json:"host,omitempty"`
	// The port number of the endpoint.
	Port *int32 `json:"port,omitempty"`
	// The connection type of the endpoint.  - `PUBLIC`: a public endpoint. For more information, see [Network](https://docs.pingcap.com/tidbcloud/connect-to-tidb-instance/#network).  - `PRIVATE_ENDPOINT`: a private endpoint. For more information, see [Network](https://docs.pingcap.com/tidbcloud/connect-to-tidb-instance/#network).  - `VPC_PEERING`: a VPC peering endpoint (not yet supported, planned for a future release).
	ConnectionType *EndpointConnectionType `json:"connectionType,omitempty"`
	// The reachability status of the endpoint connection.
	ConnectionReachability *V1beta2ConnectionReachability `json:"connectionReachability,omitempty"`
	// The public egress IP addresses used by the TiDB Cloud Premium instance.  This field is returned only for public endpoints. You can use these IP addresses to configure allowlists on upstream services when connecting through a public endpoint.
	EgressIps            []string `json:"egressIps,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _TidbEndpoint TidbEndpoint

// NewTidbEndpoint instantiates a new TidbEndpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbEndpoint() *TidbEndpoint {
	this := TidbEndpoint{}
	return &this
}

// NewTidbEndpointWithDefaults instantiates a new TidbEndpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbEndpointWithDefaults() *TidbEndpoint {
	this := TidbEndpoint{}
	return &this
}

// GetHost returns the Host field value if set, zero value otherwise.
func (o *TidbEndpoint) GetHost() string {
	if o == nil || IsNil(o.Host) {
		var ret string
		return ret
	}
	return *o.Host
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbEndpoint) GetHostOk() (*string, bool) {
	if o == nil || IsNil(o.Host) {
		return nil, false
	}
	return o.Host, true
}

// HasHost returns a boolean if a field has been set.
func (o *TidbEndpoint) HasHost() bool {
	if o != nil && !IsNil(o.Host) {
		return true
	}

	return false
}

// SetHost gets a reference to the given string and assigns it to the Host field.
func (o *TidbEndpoint) SetHost(v string) {
	o.Host = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *TidbEndpoint) GetPort() int32 {
	if o == nil || IsNil(o.Port) {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbEndpoint) GetPortOk() (*int32, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *TidbEndpoint) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *TidbEndpoint) SetPort(v int32) {
	o.Port = &v
}

// GetConnectionType returns the ConnectionType field value if set, zero value otherwise.
func (o *TidbEndpoint) GetConnectionType() EndpointConnectionType {
	if o == nil || IsNil(o.ConnectionType) {
		var ret EndpointConnectionType
		return ret
	}
	return *o.ConnectionType
}

// GetConnectionTypeOk returns a tuple with the ConnectionType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbEndpoint) GetConnectionTypeOk() (*EndpointConnectionType, bool) {
	if o == nil || IsNil(o.ConnectionType) {
		return nil, false
	}
	return o.ConnectionType, true
}

// HasConnectionType returns a boolean if a field has been set.
func (o *TidbEndpoint) HasConnectionType() bool {
	if o != nil && !IsNil(o.ConnectionType) {
		return true
	}

	return false
}

// SetConnectionType gets a reference to the given EndpointConnectionType and assigns it to the ConnectionType field.
func (o *TidbEndpoint) SetConnectionType(v EndpointConnectionType) {
	o.ConnectionType = &v
}

// GetConnectionReachability returns the ConnectionReachability field value if set, zero value otherwise.
func (o *TidbEndpoint) GetConnectionReachability() V1beta2ConnectionReachability {
	if o == nil || IsNil(o.ConnectionReachability) {
		var ret V1beta2ConnectionReachability
		return ret
	}
	return *o.ConnectionReachability
}

// GetConnectionReachabilityOk returns a tuple with the ConnectionReachability field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbEndpoint) GetConnectionReachabilityOk() (*V1beta2ConnectionReachability, bool) {
	if o == nil || IsNil(o.ConnectionReachability) {
		return nil, false
	}
	return o.ConnectionReachability, true
}

// HasConnectionReachability returns a boolean if a field has been set.
func (o *TidbEndpoint) HasConnectionReachability() bool {
	if o != nil && !IsNil(o.ConnectionReachability) {
		return true
	}

	return false
}

// SetConnectionReachability gets a reference to the given V1beta2ConnectionReachability and assigns it to the ConnectionReachability field.
func (o *TidbEndpoint) SetConnectionReachability(v V1beta2ConnectionReachability) {
	o.ConnectionReachability = &v
}

// GetEgressIps returns the EgressIps field value if set, zero value otherwise.
func (o *TidbEndpoint) GetEgressIps() []string {
	if o == nil || IsNil(o.EgressIps) {
		var ret []string
		return ret
	}
	return o.EgressIps
}

// GetEgressIpsOk returns a tuple with the EgressIps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbEndpoint) GetEgressIpsOk() ([]string, bool) {
	if o == nil || IsNil(o.EgressIps) {
		return nil, false
	}
	return o.EgressIps, true
}

// HasEgressIps returns a boolean if a field has been set.
func (o *TidbEndpoint) HasEgressIps() bool {
	if o != nil && !IsNil(o.EgressIps) {
		return true
	}

	return false
}

// SetEgressIps gets a reference to the given []string and assigns it to the EgressIps field.
func (o *TidbEndpoint) SetEgressIps(v []string) {
	o.EgressIps = v
}

func (o TidbEndpoint) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbEndpoint) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Host) {
		toSerialize["host"] = o.Host
	}
	if !IsNil(o.Port) {
		toSerialize["port"] = o.Port
	}
	if !IsNil(o.ConnectionType) {
		toSerialize["connectionType"] = o.ConnectionType
	}
	if !IsNil(o.ConnectionReachability) {
		toSerialize["connectionReachability"] = o.ConnectionReachability
	}
	if !IsNil(o.EgressIps) {
		toSerialize["egressIps"] = o.EgressIps
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TidbEndpoint) UnmarshalJSON(data []byte) (err error) {
	varTidbEndpoint := _TidbEndpoint{}

	err = json.Unmarshal(data, &varTidbEndpoint)

	if err != nil {
		return err
	}

	*o = TidbEndpoint(varTidbEndpoint)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "host")
		delete(additionalProperties, "port")
		delete(additionalProperties, "connectionType")
		delete(additionalProperties, "connectionReachability")
		delete(additionalProperties, "egressIps")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbEndpoint struct {
	value *TidbEndpoint
	isSet bool
}

func (v NullableTidbEndpoint) Get() *TidbEndpoint {
	return v.value
}

func (v *NullableTidbEndpoint) Set(val *TidbEndpoint) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbEndpoint) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbEndpoint) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbEndpoint(val *TidbEndpoint) *NullableTidbEndpoint {
	return &NullableTidbEndpoint{value: val, isSet: true}
}

func (v NullableTidbEndpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbEndpoint) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
