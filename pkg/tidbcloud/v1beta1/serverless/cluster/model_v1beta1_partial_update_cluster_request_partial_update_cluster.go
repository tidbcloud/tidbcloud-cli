/*
TiDB Cloud Starter and Essential API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Starter](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-serverless) and [TiDB Cloud Essential](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#essential) clusters. For [TiDB Cloud Dedicated](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-dedicated) clusters, use the [TiDB Cloud Dedicated API](https://docs.pingcap.com/tidbcloud/api/v1beta1/dedicated/). For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/tidbcloud/api-overview/).  # Overview  The TiDB Cloud API is a [REST interface](https://en.wikipedia.org/wiki/Representational_state_transfer) that provides you with programmatic access to manage clusters and related resources within TiDB Cloud.  The API has the following features:  - **JSON entities.** All entities are expressed in JSON. - **HTTPS-only.** You can only access the API via HTTPS, ensuring all the data sent over the network is encrypted with TLS. - **Key-based access and digest authentication.** Before you access TiDB Cloud API, you must generate an API key. All requests are authenticated through [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication), ensuring the API key is never sent over the network.  # Get Started  This guide helps you make your first API call to TiDB Cloud API. You'll learn how to authenticate a request, build a request, and interpret the response.  ## Prerequisites  To complete this guide, you need to perform the following tasks:  - Create a [TiDB Cloud account](https://tidbcloud.com/free-trial) - Install [curl](https://curl.se/)  ## Step 1. Create an API key  To create an API key, log in to your TiDB Cloud console. Navigate to the [**API Keys**](https://tidbcloud.com/org-settings/api-keys) page of your organization, and create an API key.  An API key contains a public key and a private key. Copy and save them in a secure location. You will need to use the API key later in this guide.  For more details about creating API key, refer to [API Key Management](#section/Authentication/API-Key-Management).  ## Step 2. Make your first API call  ### Build an API call  TiDB Cloud API call consists of the following components:  - **A host**. The host for TiDB Cloud API is <https://serverless.tidbapi.com>. - **An API Key**. The public key and the private key are required for authentication. - **A request**. When submitting data to a resource via `POST`, `PATCH`, or `PUT`, you must submit your payload in JSON.  In this guide, you call the [List all clusters](#tag/Cluster/operation/ClusterService_ListClusters) endpoint. For the detailed description of the endpoint, see the [API reference](#tag/Cluster/operation/ClusterService_ListClusters).  ### Call an API endpoint  To get all clusters in your organization, run the following command in your terminal. Remember to change `YOUR_PUBLIC_KEY` to your public key and `YOUR_PRIVATE_KEY` to your private key.  ```shell curl --digest \\  --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\  --request GET \\  --url 'https://serverless.tidbapi.com/v1beta1/clusters' ```  ## Step 3. Check the response  After making the API call, if the status code in response is `200` and you see details about all clusters in your organization, your request is successful.  # Authentication  The TiDB Cloud API uses [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication). It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, the TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. The role of the API key is always `Organization Owner` currently. 5. Click **Next**. Copy and save the public key and the private key. 6. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 7. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Edit**. 4. You can update the API key description. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # Rate Limiting  The TiDB Cloud API allows up to 100 requests per minute per API key. If you exceed the rate limit, the API returns a `429` error. For more quota, you can [submit a request](https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519) to contact our support team.  Each API request returns the following headers about the limit.  - `X-Ratelimit-Limit-Minute`: The number of requests allowed per minute. It is 100 currently. - `X-Ratelimit-Remaining-Minute`: The number of remaining requests in the current minute. When it reaches `0`, the API returns a `429` error and indicates that you exceed the rate limit. - `X-Ratelimit-Reset`: The time in seconds at which the current rate limit resets.  If you exceed the rate limit, an error response returns like this.  ``` > HTTP/2 429 > date: Fri, 22 Jul 2022 05:28:37 GMT > content-type: application/json > content-length: 66 > x-ratelimit-reset: 23 > x-ratelimit-remaining-minute: 0 > x-ratelimit-limit-minute: 100 > x-kong-response-latency: 2 > server: kong/2.8.1  > {\"details\":[],\"code\":49900007,\"message\":\"The request exceeded the limit of 100 times per apikey per minute. For more quota, please contact us: https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519\"} ```  # API Changelog  This changelog lists all changes to the TiDB Cloud API.  <!-- In reverse chronological order -->  ## 20250812  - Initial release of the TiDB Cloud Starter and Essential API, including the following resources and endpoints:   - Cluster:  - [List TiDB Cloud Starter and Essential clusters](#tag/Cluster/operation/ClusterService_ListClusters)  - [Create a new TiDB Cloud Starter or Essential cluster](#tag/Cluster/operation/ClusterService_CreateCluster)  - [Get details of a TiDB Cloud Starter or Essential cluster](#tag/Cluster/operation/ClusterService_GetCluster)  - [Delete a TiDB Cloud Starter or Essential cluster](#tag/Cluster/operation/ClusterService_DeleteCluster)  - [Update a TiDB Cloud Starter or Essential cluster](#tag/Cluster/operation/ClusterService_PartialUpdateCluster)  - [List available regions for an organization](#tag/Cluster/operation/ClusterService_ListRegions)  - Branch:  - [List branches](#tag/Branch/operation/BranchService_ListBranches)  - [Create a branch](#tag/Branch/operation/BranchService_CreateBranch)  - [Get details of a branch](#tag/Branch/operation/BranchService_GetBranch)  - [Delete a branch](#tag/Branch/operation/BranchService_DeleteBranch)  - [Reset a branch](#tag/Branch/operation/BranchService_ResetBranch)  - Export:  - [List export tasks for a cluster](#tag/Export/operation/ExportService_ListExports)  - [Create an export task](#tag/Export/operation/ExportService_CreateExport)  - [Get details of an export task](#tag/Export/operation/ExportService_GetExport)  - [Delete an export task](#tag/Export/operation/ExportService_DeleteExport)  - [Cancel an export task](#tag/Export/operation/ExportService_CancelExport)  - Import:  - [List import tasks for a cluster](#tag/Import/operation/ImportService_ListImports)  - [Create an import task](#tag/Import/operation/ImportService_CreateImport)  - [Get an import task](#tag/Import/operation/ImportService_GetImport)  - [Cancel an import task](#tag/Import/operation/ImportService_CancelImport)

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cluster

import (
	"encoding/json"
	"fmt"
)

// checks if the V1beta1PartialUpdateClusterRequestPartialUpdateCluster type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1PartialUpdateClusterRequestPartialUpdateCluster{}

// V1beta1PartialUpdateClusterRequestPartialUpdateCluster The updated cluster configuration.
type V1beta1PartialUpdateClusterRequestPartialUpdateCluster struct {
	// The ID of the cluster to update.
	ClusterId string `json:"clusterId"`
	// The user-defined name of the cluster.
	DisplayName *string `json:"displayName,omitempty"`
	// The [spending limit](https://docs.pingcap.com/tidbcloud/tidb-cloud-glossary/#spending-limit) for the cluster.
	SpendingLimit *ClusterSpendingLimit `json:"spendingLimit,omitempty"`
	// The schedule and retention rules for automated database backups.
	AutomatedBackupPolicy *V1beta1ClusterAutomatedBackupPolicy `json:"automatedBackupPolicy,omitempty"`
	// The connection endpoints for accessing the cluster.
	Endpoints *V1beta1ClusterEndpoints `json:"endpoints,omitempty"`
	// The labels for the cluster.  - `tidb.cloud/organization`: the ID of the organization where the cluster belongs. - `tidb.cloud/project`: the ID of the project where the cluster belongs.
	Labels *map[string]string `json:"labels,omitempty"`
	// The audit log configuration for the cluster.
	AuditLogConfig *V1beta1ClusterAuditLogConfig `json:"auditLogConfig,omitempty"`
	// The auto-scaling configuration for the cluster.
	AutoScaling          *V1beta1ClusterAutoScaling `json:"autoScaling,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1PartialUpdateClusterRequestPartialUpdateCluster V1beta1PartialUpdateClusterRequestPartialUpdateCluster

// NewV1beta1PartialUpdateClusterRequestPartialUpdateCluster instantiates a new V1beta1PartialUpdateClusterRequestPartialUpdateCluster object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1PartialUpdateClusterRequestPartialUpdateCluster(clusterId string) *V1beta1PartialUpdateClusterRequestPartialUpdateCluster {
	this := V1beta1PartialUpdateClusterRequestPartialUpdateCluster{}
	this.ClusterId = clusterId
	return &this
}

// NewV1beta1PartialUpdateClusterRequestPartialUpdateClusterWithDefaults instantiates a new V1beta1PartialUpdateClusterRequestPartialUpdateCluster object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1PartialUpdateClusterRequestPartialUpdateClusterWithDefaults() *V1beta1PartialUpdateClusterRequestPartialUpdateCluster {
	this := V1beta1PartialUpdateClusterRequestPartialUpdateCluster{}
	return &this
}

// GetClusterId returns the ClusterId field value
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) SetClusterId(v string) {
	o.ClusterId = v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetSpendingLimit returns the SpendingLimit field value if set, zero value otherwise.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetSpendingLimit() ClusterSpendingLimit {
	if o == nil || IsNil(o.SpendingLimit) {
		var ret ClusterSpendingLimit
		return ret
	}
	return *o.SpendingLimit
}

// GetSpendingLimitOk returns a tuple with the SpendingLimit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetSpendingLimitOk() (*ClusterSpendingLimit, bool) {
	if o == nil || IsNil(o.SpendingLimit) {
		return nil, false
	}
	return o.SpendingLimit, true
}

// HasSpendingLimit returns a boolean if a field has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) HasSpendingLimit() bool {
	if o != nil && !IsNil(o.SpendingLimit) {
		return true
	}

	return false
}

// SetSpendingLimit gets a reference to the given ClusterSpendingLimit and assigns it to the SpendingLimit field.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) SetSpendingLimit(v ClusterSpendingLimit) {
	o.SpendingLimit = &v
}

// GetAutomatedBackupPolicy returns the AutomatedBackupPolicy field value if set, zero value otherwise.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetAutomatedBackupPolicy() V1beta1ClusterAutomatedBackupPolicy {
	if o == nil || IsNil(o.AutomatedBackupPolicy) {
		var ret V1beta1ClusterAutomatedBackupPolicy
		return ret
	}
	return *o.AutomatedBackupPolicy
}

// GetAutomatedBackupPolicyOk returns a tuple with the AutomatedBackupPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetAutomatedBackupPolicyOk() (*V1beta1ClusterAutomatedBackupPolicy, bool) {
	if o == nil || IsNil(o.AutomatedBackupPolicy) {
		return nil, false
	}
	return o.AutomatedBackupPolicy, true
}

// HasAutomatedBackupPolicy returns a boolean if a field has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) HasAutomatedBackupPolicy() bool {
	if o != nil && !IsNil(o.AutomatedBackupPolicy) {
		return true
	}

	return false
}

// SetAutomatedBackupPolicy gets a reference to the given V1beta1ClusterAutomatedBackupPolicy and assigns it to the AutomatedBackupPolicy field.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) SetAutomatedBackupPolicy(v V1beta1ClusterAutomatedBackupPolicy) {
	o.AutomatedBackupPolicy = &v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetEndpoints() V1beta1ClusterEndpoints {
	if o == nil || IsNil(o.Endpoints) {
		var ret V1beta1ClusterEndpoints
		return ret
	}
	return *o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetEndpointsOk() (*V1beta1ClusterEndpoints, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}
	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given V1beta1ClusterEndpoints and assigns it to the Endpoints field.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) SetEndpoints(v V1beta1ClusterEndpoints) {
	o.Endpoints = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetAuditLogConfig returns the AuditLogConfig field value if set, zero value otherwise.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetAuditLogConfig() V1beta1ClusterAuditLogConfig {
	if o == nil || IsNil(o.AuditLogConfig) {
		var ret V1beta1ClusterAuditLogConfig
		return ret
	}
	return *o.AuditLogConfig
}

// GetAuditLogConfigOk returns a tuple with the AuditLogConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetAuditLogConfigOk() (*V1beta1ClusterAuditLogConfig, bool) {
	if o == nil || IsNil(o.AuditLogConfig) {
		return nil, false
	}
	return o.AuditLogConfig, true
}

// HasAuditLogConfig returns a boolean if a field has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) HasAuditLogConfig() bool {
	if o != nil && !IsNil(o.AuditLogConfig) {
		return true
	}

	return false
}

// SetAuditLogConfig gets a reference to the given V1beta1ClusterAuditLogConfig and assigns it to the AuditLogConfig field.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) SetAuditLogConfig(v V1beta1ClusterAuditLogConfig) {
	o.AuditLogConfig = &v
}

// GetAutoScaling returns the AutoScaling field value if set, zero value otherwise.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetAutoScaling() V1beta1ClusterAutoScaling {
	if o == nil || IsNil(o.AutoScaling) {
		var ret V1beta1ClusterAutoScaling
		return ret
	}
	return *o.AutoScaling
}

// GetAutoScalingOk returns a tuple with the AutoScaling field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) GetAutoScalingOk() (*V1beta1ClusterAutoScaling, bool) {
	if o == nil || IsNil(o.AutoScaling) {
		return nil, false
	}
	return o.AutoScaling, true
}

// HasAutoScaling returns a boolean if a field has been set.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) HasAutoScaling() bool {
	if o != nil && !IsNil(o.AutoScaling) {
		return true
	}

	return false
}

// SetAutoScaling gets a reference to the given V1beta1ClusterAutoScaling and assigns it to the AutoScaling field.
func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) SetAutoScaling(v V1beta1ClusterAutoScaling) {
	o.AutoScaling = &v
}

func (o V1beta1PartialUpdateClusterRequestPartialUpdateCluster) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1PartialUpdateClusterRequestPartialUpdateCluster) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["clusterId"] = o.ClusterId
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}
	if !IsNil(o.SpendingLimit) {
		toSerialize["spendingLimit"] = o.SpendingLimit
	}
	if !IsNil(o.AutomatedBackupPolicy) {
		toSerialize["automatedBackupPolicy"] = o.AutomatedBackupPolicy
	}
	if !IsNil(o.Endpoints) {
		toSerialize["endpoints"] = o.Endpoints
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	if !IsNil(o.AuditLogConfig) {
		toSerialize["auditLogConfig"] = o.AuditLogConfig
	}
	if !IsNil(o.AutoScaling) {
		toSerialize["autoScaling"] = o.AutoScaling
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) UnmarshalJSON(data []byte) (err error) {
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

	varV1beta1PartialUpdateClusterRequestPartialUpdateCluster := _V1beta1PartialUpdateClusterRequestPartialUpdateCluster{}

	err = json.Unmarshal(data, &varV1beta1PartialUpdateClusterRequestPartialUpdateCluster)

	if err != nil {
		return err
	}

	*o = V1beta1PartialUpdateClusterRequestPartialUpdateCluster(varV1beta1PartialUpdateClusterRequestPartialUpdateCluster)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "spendingLimit")
		delete(additionalProperties, "automatedBackupPolicy")
		delete(additionalProperties, "endpoints")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "auditLogConfig")
		delete(additionalProperties, "autoScaling")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster struct {
	value *V1beta1PartialUpdateClusterRequestPartialUpdateCluster
	isSet bool
}

func (v NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster) Get() *V1beta1PartialUpdateClusterRequestPartialUpdateCluster {
	return v.value
}

func (v *NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster) Set(val *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster(val *V1beta1PartialUpdateClusterRequestPartialUpdateCluster) *NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster {
	return &NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster{value: val, isSet: true}
}

func (v NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1PartialUpdateClusterRequestPartialUpdateCluster) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
