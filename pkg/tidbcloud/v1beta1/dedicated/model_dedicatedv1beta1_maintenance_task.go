/*
TiDB Cloud Dedicated API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Dedicated](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-dedicated) clusters. For TiDB Cloud Starter or TiDB Cloud Essential clusters, use the [TiDB Cloud Starter and Essential API](). For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/tidbcloud/api-overview/).  # Overview  The TiDB Cloud API is a [REST interface](https://en.wikipedia.org/wiki/Representational_state_transfer) that provides you with programmatic access to manage clusters and related resources within TiDB Cloud.  The API has the following features:  - **JSON entities.** All entities are expressed in JSON. - **HTTPS-only.** You can only access the API via HTTPS, ensuring all the data sent over the network is encrypted with TLS. - **Key-based access and digest authentication.** Before you access TiDB Cloud API, you must generate an API key. All requests are authenticated through [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication), ensuring the API key is never sent over the network.  # Get Started  This guide helps you make your first API call to TiDB Cloud API. You'll learn how to authenticate a request, build a request, and interpret the response.  ## Prerequisites  To complete this guide, you need to perform the following tasks:  - Create a [TiDB Cloud account](https://tidbcloud.com/free-trial) - Install [curl](https://curl.se/)  ## Step 1. Create an API key  To create an API key, log in to your TiDB Cloud console. Navigate to the [**API Keys**](https://tidbcloud.com/org-settings/api-keys) page of your organization, and create an API key.  An API key contains a public key and a private key. Copy and save them in a secure location. You will need to use the API key later in this guide.  For more details about creating API key, refer to [API Key Management](#section/Authentication/API-Key-Management).  ## Step 2. Make your first API call  ### Build an API call  TiDB Cloud API call consists of the following components:  - **A host**. The host for TiDB Cloud API is <https://dedicated.tidbapi.com>. - **An API Key**. The public key and the private key are required for authentication. - **A request**. When submitting data to a resource via `POST`, `PATCH`, or `PUT`, you must submit your payload in JSON.  In this guide, you call the [List clusters](#tag/Cluster/operation/ClusterService_ListClusters) endpoint. For the detailed description of the endpoint, see the [API reference](#tag/Cluster/operation/ClusterService_ListClusters).  ### Call an API endpoint  To get all clusters in your organization, run the following command in your terminal. Remember to change `YOUR_PUBLIC_KEY` to your public key and `YOUR_PRIVATE_KEY` to your private key.  ```shell curl --digest \\  --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\  --request GET \\  --url 'https://dedicated.tidbapi.com/v1beta1/clusters' ```  ## Step 3. Check the response  After making the API call, if the status code in response is `200` and you see details about all clusters in your organization, your request is successful.  # Authentication  The TiDB Cloud API uses [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication). It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, the TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. The role of the API key is always `Organization Owner` currently. 5. Click **Next**. Copy and save the public key and the private key. 6. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 7. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Edit**. 4. You can update the API key description. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # Rate Limiting  The TiDB Cloud API allows up to 100 requests per minute per API key. If you exceed the rate limit, the API returns a `429` error. For more quota, you can [submit a request](https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519) to contact our support team.  Each API request returns the following headers about the limit.  - `X-Ratelimit-Limit-Minute`: The number of requests allowed per minute. It is 100 currently. - `X-Ratelimit-Remaining-Minute`: The number of remaining requests in the current minute. When it reaches `0`, the API returns a `429` error and indicates that you exceed the rate limit. - `X-Ratelimit-Reset`: The time in seconds at which the current rate limit resets.  If you exceed the rate limit, an error response returns like this.  ``` > HTTP/2 429 > date: Fri, 22 Jul 2022 05:28:37 GMT > content-type: application/json > content-length: 66 > x-ratelimit-reset: 23 > x-ratelimit-remaining-minute: 0 > x-ratelimit-limit-minute: 100 > x-kong-response-latency: 2 > server: kong/2.8.1  > {\"details\":[],\"code\":49900007,\"message\":\"The request exceeded the limit of 100 times per apikey per minute. For more quota, please contact us: https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519\"} ```  # API Changelog  This changelog lists all changes to the TiDB Cloud API.  <!-- In reverse chronological order -->  ## 20250812  - Initial release of the TiDB Cloud Dedicated API, including the following resources and endpoints:  * Cluster    * [List clusters](#tag/Cluster/operation/ClusterService_ListClusters)    * [Create a cluster](#tag/Cluster/operation/ClusterService_CreateCluster)    * [Get a cluster](#tag/Cluster/operation/ClusterService_GetCluster)    * [Delete a cluster](#tag/Cluster/operation/ClusterService_DeleteCluster)    * [Update a cluster](#tag/Cluster/operation/ClusterService_UpdateCluster)    * [Pause a cluster](#tag/Cluster/operation/ClusterService_PauseCluster)    * [Resume a cluster](#tag/Cluster/operation/ClusterService_ResumeCluster)    * [Reset the root password of a cluster](#tag/Cluster/operation/ClusterService_ResetRootPassword)    * [List node quotas for your organization](#tag/Cluster/operation/ClusterService_ShowNodeQuota)    * [Get log redaction policy](#tag/Cluster/operation/ClusterService_GetLogRedactionPolicy)   * Region    * [List regions](#tag/Region/operation/RegionService_ListRegions)    * [Get a region](#tag/Region/operation/RegionService_GetRegion)    * [List cloud providers](#tag/Region/operation/RegionService_ShowCloudProviders)    * [List node specs](#tag/Region/operation/RegionService_ListNodeSpecs)    * [Get a node spec](#tag/Region/operation/RegionService_GetNodeSpec)   * Private Endpoint Connection    * [Get private link service for a TiDB node group](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateLinkService)    * [Create a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_CreatePrivateEndpointConnection)    * [List private endpoint connections](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_ListPrivateEndpointConnections)    * [Get a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateEndpointConnection)    * [Delete a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_DeletePrivateEndpointConnection)   * Import    * [List import tasks](#tag/Import/operation/ListImports)    * [Create an import task](#tag/Import/operation/CreateImport)    * [Get an import task](#tag/Import/operation/GetImport)    * [Cancel an import task](#tag/Import/operation/CancelImport)

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"time"
)

// checks if the Dedicatedv1beta1MaintenanceTask type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1MaintenanceTask{}

// Dedicatedv1beta1MaintenanceTask All fields are output only.
type Dedicatedv1beta1MaintenanceTask struct {
	Name              *string                               `json:"name,omitempty"`
	MaintenanceTaskId *string                               `json:"maintenanceTaskId,omitempty"`
	ProjectId         *string                               `json:"projectId,omitempty"`
	Description       *string                               `json:"description,omitempty"`
	State             *Dedicatedv1beta1MaintenanceTaskState `json:"state,omitempty"`
	// Timestamp when the task was created.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Timestamp when the task run.
	ScheduledApplyTime *time.Time `json:"scheduledApplyTime,omitempty"`
	// Timestamp when the task will be expired.
	Deadline             *time.Time `json:"deadline,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1MaintenanceTask Dedicatedv1beta1MaintenanceTask

// NewDedicatedv1beta1MaintenanceTask instantiates a new Dedicatedv1beta1MaintenanceTask object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1MaintenanceTask() *Dedicatedv1beta1MaintenanceTask {
	this := Dedicatedv1beta1MaintenanceTask{}
	return &this
}

// NewDedicatedv1beta1MaintenanceTaskWithDefaults instantiates a new Dedicatedv1beta1MaintenanceTask object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1MaintenanceTaskWithDefaults() *Dedicatedv1beta1MaintenanceTask {
	this := Dedicatedv1beta1MaintenanceTask{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1MaintenanceTask) SetName(v string) {
	o.Name = &v
}

// GetMaintenanceTaskId returns the MaintenanceTaskId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetMaintenanceTaskId() string {
	if o == nil || IsNil(o.MaintenanceTaskId) {
		var ret string
		return ret
	}
	return *o.MaintenanceTaskId
}

// GetMaintenanceTaskIdOk returns a tuple with the MaintenanceTaskId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetMaintenanceTaskIdOk() (*string, bool) {
	if o == nil || IsNil(o.MaintenanceTaskId) {
		return nil, false
	}
	return o.MaintenanceTaskId, true
}

// HasMaintenanceTaskId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasMaintenanceTaskId() bool {
	if o != nil && !IsNil(o.MaintenanceTaskId) {
		return true
	}

	return false
}

// SetMaintenanceTaskId gets a reference to the given string and assigns it to the MaintenanceTaskId field.
func (o *Dedicatedv1beta1MaintenanceTask) SetMaintenanceTaskId(v string) {
	o.MaintenanceTaskId = &v
}

// GetProjectId returns the ProjectId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetProjectId() string {
	if o == nil || IsNil(o.ProjectId) {
		var ret string
		return ret
	}
	return *o.ProjectId
}

// GetProjectIdOk returns a tuple with the ProjectId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetProjectIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProjectId) {
		return nil, false
	}
	return o.ProjectId, true
}

// HasProjectId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasProjectId() bool {
	if o != nil && !IsNil(o.ProjectId) {
		return true
	}

	return false
}

// SetProjectId gets a reference to the given string and assigns it to the ProjectId field.
func (o *Dedicatedv1beta1MaintenanceTask) SetProjectId(v string) {
	o.ProjectId = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Dedicatedv1beta1MaintenanceTask) SetDescription(v string) {
	o.Description = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetState() Dedicatedv1beta1MaintenanceTaskState {
	if o == nil || IsNil(o.State) {
		var ret Dedicatedv1beta1MaintenanceTaskState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetStateOk() (*Dedicatedv1beta1MaintenanceTaskState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Dedicatedv1beta1MaintenanceTaskState and assigns it to the State field.
func (o *Dedicatedv1beta1MaintenanceTask) SetState(v Dedicatedv1beta1MaintenanceTaskState) {
	o.State = &v
}

// GetCreateTime returns the CreateTime field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetCreateTime() time.Time {
	if o == nil || IsNil(o.CreateTime) {
		var ret time.Time
		return ret
	}
	return *o.CreateTime
}

// GetCreateTimeOk returns a tuple with the CreateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetCreateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateTime) {
		return nil, false
	}
	return o.CreateTime, true
}

// HasCreateTime returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasCreateTime() bool {
	if o != nil && !IsNil(o.CreateTime) {
		return true
	}

	return false
}

// SetCreateTime gets a reference to the given time.Time and assigns it to the CreateTime field.
func (o *Dedicatedv1beta1MaintenanceTask) SetCreateTime(v time.Time) {
	o.CreateTime = &v
}

// GetScheduledApplyTime returns the ScheduledApplyTime field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetScheduledApplyTime() time.Time {
	if o == nil || IsNil(o.ScheduledApplyTime) {
		var ret time.Time
		return ret
	}
	return *o.ScheduledApplyTime
}

// GetScheduledApplyTimeOk returns a tuple with the ScheduledApplyTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetScheduledApplyTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ScheduledApplyTime) {
		return nil, false
	}
	return o.ScheduledApplyTime, true
}

// HasScheduledApplyTime returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasScheduledApplyTime() bool {
	if o != nil && !IsNil(o.ScheduledApplyTime) {
		return true
	}

	return false
}

// SetScheduledApplyTime gets a reference to the given time.Time and assigns it to the ScheduledApplyTime field.
func (o *Dedicatedv1beta1MaintenanceTask) SetScheduledApplyTime(v time.Time) {
	o.ScheduledApplyTime = &v
}

// GetDeadline returns the Deadline field value if set, zero value otherwise.
func (o *Dedicatedv1beta1MaintenanceTask) GetDeadline() time.Time {
	if o == nil || IsNil(o.Deadline) {
		var ret time.Time
		return ret
	}
	return *o.Deadline
}

// GetDeadlineOk returns a tuple with the Deadline field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1MaintenanceTask) GetDeadlineOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Deadline) {
		return nil, false
	}
	return o.Deadline, true
}

// HasDeadline returns a boolean if a field has been set.
func (o *Dedicatedv1beta1MaintenanceTask) HasDeadline() bool {
	if o != nil && !IsNil(o.Deadline) {
		return true
	}

	return false
}

// SetDeadline gets a reference to the given time.Time and assigns it to the Deadline field.
func (o *Dedicatedv1beta1MaintenanceTask) SetDeadline(v time.Time) {
	o.Deadline = &v
}

func (o Dedicatedv1beta1MaintenanceTask) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1MaintenanceTask) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.MaintenanceTaskId) {
		toSerialize["maintenanceTaskId"] = o.MaintenanceTaskId
	}
	if !IsNil(o.ProjectId) {
		toSerialize["projectId"] = o.ProjectId
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.CreateTime) {
		toSerialize["createTime"] = o.CreateTime
	}
	if !IsNil(o.ScheduledApplyTime) {
		toSerialize["scheduledApplyTime"] = o.ScheduledApplyTime
	}
	if !IsNil(o.Deadline) {
		toSerialize["deadline"] = o.Deadline
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1MaintenanceTask) UnmarshalJSON(data []byte) (err error) {
	varDedicatedv1beta1MaintenanceTask := _Dedicatedv1beta1MaintenanceTask{}

	err = json.Unmarshal(data, &varDedicatedv1beta1MaintenanceTask)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1MaintenanceTask(varDedicatedv1beta1MaintenanceTask)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "maintenanceTaskId")
		delete(additionalProperties, "projectId")
		delete(additionalProperties, "description")
		delete(additionalProperties, "state")
		delete(additionalProperties, "createTime")
		delete(additionalProperties, "scheduledApplyTime")
		delete(additionalProperties, "deadline")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1MaintenanceTask struct {
	value *Dedicatedv1beta1MaintenanceTask
	isSet bool
}

func (v NullableDedicatedv1beta1MaintenanceTask) Get() *Dedicatedv1beta1MaintenanceTask {
	return v.value
}

func (v *NullableDedicatedv1beta1MaintenanceTask) Set(val *Dedicatedv1beta1MaintenanceTask) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1MaintenanceTask) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1MaintenanceTask) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1MaintenanceTask(val *Dedicatedv1beta1MaintenanceTask) *NullableDedicatedv1beta1MaintenanceTask {
	return &NullableDedicatedv1beta1MaintenanceTask{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1MaintenanceTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1MaintenanceTask) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
