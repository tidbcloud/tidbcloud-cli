/*
TiDB Cloud Dedicated API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Dedicated](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-dedicated) clusters. For TiDB Cloud Starter or TiDB Cloud Essential clusters, use the [TiDB Cloud Starter and Essential API](). For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/tidbcloud/api-overview/).  # Overview  The TiDB Cloud API is a [REST interface](https://en.wikipedia.org/wiki/Representational_state_transfer) that provides you with programmatic access to manage clusters and related resources within TiDB Cloud.  The API has the following features:  - **JSON entities.** All entities are expressed in JSON. - **HTTPS-only.** You can only access the API via HTTPS, ensuring all the data sent over the network is encrypted with TLS. - **Key-based access and digest authentication.** Before you access TiDB Cloud API, you must generate an API key. All requests are authenticated through [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication), ensuring the API key is never sent over the network.  # Get Started  This guide helps you make your first API call to TiDB Cloud API. You'll learn how to authenticate a request, build a request, and interpret the response.  ## Prerequisites  To complete this guide, you need to perform the following tasks:  - Create a [TiDB Cloud account](https://tidbcloud.com/free-trial) - Install [curl](https://curl.se/)  ## Step 1. Create an API key  To create an API key, log in to your TiDB Cloud console. Navigate to the [**API Keys**](https://tidbcloud.com/org-settings/api-keys) page of your organization, and create an API key.  An API key contains a public key and a private key. Copy and save them in a secure location. You will need to use the API key later in this guide.  For more details about creating API key, refer to [API Key Management](#section/Authentication/API-Key-Management).  ## Step 2. Make your first API call  ### Build an API call  TiDB Cloud API call consists of the following components:  - **A host**. The host for TiDB Cloud API is <https://dedicated.tidbapi.com>. - **An API Key**. The public key and the private key are required for authentication. - **A request**. When submitting data to a resource via `POST`, `PATCH`, or `PUT`, you must submit your payload in JSON.  In this guide, you call the [List clusters](#tag/Cluster/operation/ClusterService_ListClusters) endpoint. For the detailed description of the endpoint, see the [API reference](#tag/Cluster/operation/ClusterService_ListClusters).  ### Call an API endpoint  To get all clusters in your organization, run the following command in your terminal. Remember to change `YOUR_PUBLIC_KEY` to your public key and `YOUR_PRIVATE_KEY` to your private key.  ```shell curl --digest \\  --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\  --request GET \\  --url 'https://dedicated.tidbapi.com/v1beta1/clusters' ```  ## Step 3. Check the response  After making the API call, if the status code in response is `200` and you see details about all clusters in your organization, your request is successful.  # Authentication  The TiDB Cloud API uses [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication). It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, the TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. The role of the API key is always `Organization Owner` currently. 5. Click **Next**. Copy and save the public key and the private key. 6. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 7. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Edit**. 4. You can update the API key description. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # Rate Limiting  The TiDB Cloud API allows up to 100 requests per minute per API key. If you exceed the rate limit, the API returns a `429` error. For more quota, you can [submit a request](https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519) to contact our support team.  Each API request returns the following headers about the limit.  - `X-Ratelimit-Limit-Minute`: The number of requests allowed per minute. It is 100 currently. - `X-Ratelimit-Remaining-Minute`: The number of remaining requests in the current minute. When it reaches `0`, the API returns a `429` error and indicates that you exceed the rate limit. - `X-Ratelimit-Reset`: The time in seconds at which the current rate limit resets.  If you exceed the rate limit, an error response returns like this.  ``` > HTTP/2 429 > date: Fri, 22 Jul 2022 05:28:37 GMT > content-type: application/json > content-length: 66 > x-ratelimit-reset: 23 > x-ratelimit-remaining-minute: 0 > x-ratelimit-limit-minute: 100 > x-kong-response-latency: 2 > server: kong/2.8.1  > {\"details\":[],\"code\":49900007,\"message\":\"The request exceeded the limit of 100 times per apikey per minute. For more quota, please contact us: https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519\"} ```  # API Changelog  This changelog lists all changes to the TiDB Cloud API.  <!-- In reverse chronological order -->  ## 20250812  - Initial release of the TiDB Cloud Dedicated API, including the following resources and endpoints:  * Cluster    * [List clusters](#tag/Cluster/operation/ClusterService_ListClusters)    * [Create a cluster](#tag/Cluster/operation/ClusterService_CreateCluster)    * [Get a cluster](#tag/Cluster/operation/ClusterService_GetCluster)    * [Delete a cluster](#tag/Cluster/operation/ClusterService_DeleteCluster)    * [Update a cluster](#tag/Cluster/operation/ClusterService_UpdateCluster)    * [Pause a cluster](#tag/Cluster/operation/ClusterService_PauseCluster)    * [Resume a cluster](#tag/Cluster/operation/ClusterService_ResumeCluster)    * [Reset the root password of a cluster](#tag/Cluster/operation/ClusterService_ResetRootPassword)    * [List node quotas for your organization](#tag/Cluster/operation/ClusterService_ShowNodeQuota)    * [Get log redaction policy](#tag/Cluster/operation/ClusterService_GetLogRedactionPolicy)   * Region    * [List regions](#tag/Region/operation/RegionService_ListRegions)    * [Get a region](#tag/Region/operation/RegionService_GetRegion)    * [List cloud providers](#tag/Region/operation/RegionService_ShowCloudProviders)    * [List node specs](#tag/Region/operation/RegionService_ListNodeSpecs)    * [Get a node spec](#tag/Region/operation/RegionService_GetNodeSpec)   * Private Endpoint Connection    * [Get private link service for a TiDB node group](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateLinkService)    * [Create a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_CreatePrivateEndpointConnection)    * [List private endpoint connections](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_ListPrivateEndpointConnections)    * [Get a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateEndpointConnection)    * [Delete a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_DeletePrivateEndpointConnection)   * Import    * [List import tasks](#tag/Import/operation/ListImports)    * [Create an import task](#tag/Import/operation/CreateImport)    * [Get an import task](#tag/Import/operation/GetImport)    * [Cancel an import task](#tag/Import/operation/CancelImport)

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"fmt"
)

// checks if the Dedicatedv1beta1LogRedactionPolicy type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1LogRedactionPolicy{}

// Dedicatedv1beta1LogRedactionPolicy struct for Dedicatedv1beta1LogRedactionPolicy
type Dedicatedv1beta1LogRedactionPolicy struct {
	// The ID of the cluster to which the log redaction policy applies.
	ClusterId string `json:"clusterId"`
	// Indicates whether to redact sensitive information in logs.  The default policy is `true`, which means sensitive data is redacted in logs.  **Note**: Disabling log redaction might expose sensitive information and increase the risk of data leakage. Ensure that you understand and acknowledge this risk before proceeding. Remember to re-enable it as soon as you complete your diagnostic or maintenance task.
	RedactLog            NullableBool `json:"redactLog,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1LogRedactionPolicy Dedicatedv1beta1LogRedactionPolicy

// NewDedicatedv1beta1LogRedactionPolicy instantiates a new Dedicatedv1beta1LogRedactionPolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1LogRedactionPolicy(clusterId string) *Dedicatedv1beta1LogRedactionPolicy {
	this := Dedicatedv1beta1LogRedactionPolicy{}
	this.ClusterId = clusterId
	var redactLog bool = true
	this.RedactLog = *NewNullableBool(&redactLog)
	return &this
}

// NewDedicatedv1beta1LogRedactionPolicyWithDefaults instantiates a new Dedicatedv1beta1LogRedactionPolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1LogRedactionPolicyWithDefaults() *Dedicatedv1beta1LogRedactionPolicy {
	this := Dedicatedv1beta1LogRedactionPolicy{}
	var redactLog bool = true
	this.RedactLog = *NewNullableBool(&redactLog)
	return &this
}

// GetClusterId returns the ClusterId field value
func (o *Dedicatedv1beta1LogRedactionPolicy) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1LogRedactionPolicy) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *Dedicatedv1beta1LogRedactionPolicy) SetClusterId(v string) {
	o.ClusterId = v
}

// GetRedactLog returns the RedactLog field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Dedicatedv1beta1LogRedactionPolicy) GetRedactLog() bool {
	if o == nil || IsNil(o.RedactLog.Get()) {
		var ret bool
		return ret
	}
	return *o.RedactLog.Get()
}

// GetRedactLogOk returns a tuple with the RedactLog field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Dedicatedv1beta1LogRedactionPolicy) GetRedactLogOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.RedactLog.Get(), o.RedactLog.IsSet()
}

// HasRedactLog returns a boolean if a field has been set.
func (o *Dedicatedv1beta1LogRedactionPolicy) HasRedactLog() bool {
	if o != nil && o.RedactLog.IsSet() {
		return true
	}

	return false
}

// SetRedactLog gets a reference to the given NullableBool and assigns it to the RedactLog field.
func (o *Dedicatedv1beta1LogRedactionPolicy) SetRedactLog(v bool) {
	o.RedactLog.Set(&v)
}

// SetRedactLogNil sets the value for RedactLog to be an explicit nil
func (o *Dedicatedv1beta1LogRedactionPolicy) SetRedactLogNil() {
	o.RedactLog.Set(nil)
}

// UnsetRedactLog ensures that no value is present for RedactLog, not even an explicit nil
func (o *Dedicatedv1beta1LogRedactionPolicy) UnsetRedactLog() {
	o.RedactLog.Unset()
}

func (o Dedicatedv1beta1LogRedactionPolicy) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1LogRedactionPolicy) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["clusterId"] = o.ClusterId
	if o.RedactLog.IsSet() {
		toSerialize["redactLog"] = o.RedactLog.Get()
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1LogRedactionPolicy) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"clusterId",
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

	varDedicatedv1beta1LogRedactionPolicy := _Dedicatedv1beta1LogRedactionPolicy{}

	err = json.Unmarshal(data, &varDedicatedv1beta1LogRedactionPolicy)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1LogRedactionPolicy(varDedicatedv1beta1LogRedactionPolicy)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "redactLog")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1LogRedactionPolicy struct {
	value *Dedicatedv1beta1LogRedactionPolicy
	isSet bool
}

func (v NullableDedicatedv1beta1LogRedactionPolicy) Get() *Dedicatedv1beta1LogRedactionPolicy {
	return v.value
}

func (v *NullableDedicatedv1beta1LogRedactionPolicy) Set(val *Dedicatedv1beta1LogRedactionPolicy) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1LogRedactionPolicy) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1LogRedactionPolicy) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1LogRedactionPolicy(val *Dedicatedv1beta1LogRedactionPolicy) *NullableDedicatedv1beta1LogRedactionPolicy {
	return &NullableDedicatedv1beta1LogRedactionPolicy{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1LogRedactionPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1LogRedactionPolicy) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
