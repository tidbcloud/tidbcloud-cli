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

// checks if the Dedicatedv1beta1AuditLogFilterRule type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1AuditLogFilterRule{}

// Dedicatedv1beta1AuditLogFilterRule struct for Dedicatedv1beta1AuditLogFilterRule
type Dedicatedv1beta1AuditLogFilterRule struct {
	Name *string `json:"name,omitempty"`
	// audit_log_filter_rule_id is a unique identifier for the rule. Currently, the value is the same as user_expr. For legacy, the audit_log_filter_rule_id is a number.
	AuditLogFilterRuleId *string `json:"auditLogFilterRuleId,omitempty"`
	ClusterId            string  `json:"clusterId"`
	// Format: <user>@<host>. The audit user <user>@<host> consists of the username and hostname with @ as separator, where @ and <host> are optional. Both username and hostname can be identifiers with wildcards:   - % for matching any username/hostname.   - _ for matching any character. For multiple filter rules, if their user_expr fields match the same username, that matches the longest username takes effect. In general, this field is required. Only when `AuditLogConfig.legacy` is true, this field is optional. For legacy, the express is a regex.
	UserExpr *string `json:"userExpr,omitempty"`
	// Default is enabled.
	Disabled *bool `json:"disabled,omitempty"`
	// If empty, it means all events.
	EventTypes []Dedicatedv1beta1AuditLogFilterRuleEventType `json:"eventTypes,omitempty"`
	// If not set, default to the value of `user_expr`.
	DisplayName *string `json:"displayName,omitempty"`
	// Deprecated. Only available when `AuditLogConfig.legacy` is true. The express is a glob pattern. Refer to https://github.com/pingcap/tidb/blob/master/pkg/util/table-filter/README.md#wildcards.
	DbExpr *string `json:"dbExpr,omitempty"`
	// Deprecated. Only available when `AuditLogConfig.legacy` is true. The express is a glob pattern. Refer to https://github.com/pingcap/tidb/blob/master/pkg/util/table-filter/README.md#wildcards.
	TableExpr *string `json:"tableExpr,omitempty"`
	// Deprecated. Only available when `AuditLogConfig.legacy` is true.
	AccessTypeList       []string `json:"accessTypeList,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1AuditLogFilterRule Dedicatedv1beta1AuditLogFilterRule

// NewDedicatedv1beta1AuditLogFilterRule instantiates a new Dedicatedv1beta1AuditLogFilterRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1AuditLogFilterRule(clusterId string) *Dedicatedv1beta1AuditLogFilterRule {
	this := Dedicatedv1beta1AuditLogFilterRule{}
	this.ClusterId = clusterId
	return &this
}

// NewDedicatedv1beta1AuditLogFilterRuleWithDefaults instantiates a new Dedicatedv1beta1AuditLogFilterRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1AuditLogFilterRuleWithDefaults() *Dedicatedv1beta1AuditLogFilterRule {
	this := Dedicatedv1beta1AuditLogFilterRule{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetName(v string) {
	o.Name = &v
}

// GetAuditLogFilterRuleId returns the AuditLogFilterRuleId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetAuditLogFilterRuleId() string {
	if o == nil || IsNil(o.AuditLogFilterRuleId) {
		var ret string
		return ret
	}
	return *o.AuditLogFilterRuleId
}

// GetAuditLogFilterRuleIdOk returns a tuple with the AuditLogFilterRuleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetAuditLogFilterRuleIdOk() (*string, bool) {
	if o == nil || IsNil(o.AuditLogFilterRuleId) {
		return nil, false
	}
	return o.AuditLogFilterRuleId, true
}

// HasAuditLogFilterRuleId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasAuditLogFilterRuleId() bool {
	if o != nil && !IsNil(o.AuditLogFilterRuleId) {
		return true
	}

	return false
}

// SetAuditLogFilterRuleId gets a reference to the given string and assigns it to the AuditLogFilterRuleId field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetAuditLogFilterRuleId(v string) {
	o.AuditLogFilterRuleId = &v
}

// GetClusterId returns the ClusterId field value
func (o *Dedicatedv1beta1AuditLogFilterRule) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *Dedicatedv1beta1AuditLogFilterRule) SetClusterId(v string) {
	o.ClusterId = v
}

// GetUserExpr returns the UserExpr field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetUserExpr() string {
	if o == nil || IsNil(o.UserExpr) {
		var ret string
		return ret
	}
	return *o.UserExpr
}

// GetUserExprOk returns a tuple with the UserExpr field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetUserExprOk() (*string, bool) {
	if o == nil || IsNil(o.UserExpr) {
		return nil, false
	}
	return o.UserExpr, true
}

// HasUserExpr returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasUserExpr() bool {
	if o != nil && !IsNil(o.UserExpr) {
		return true
	}

	return false
}

// SetUserExpr gets a reference to the given string and assigns it to the UserExpr field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetUserExpr(v string) {
	o.UserExpr = &v
}

// GetDisabled returns the Disabled field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetDisabled() bool {
	if o == nil || IsNil(o.Disabled) {
		var ret bool
		return ret
	}
	return *o.Disabled
}

// GetDisabledOk returns a tuple with the Disabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetDisabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Disabled) {
		return nil, false
	}
	return o.Disabled, true
}

// HasDisabled returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasDisabled() bool {
	if o != nil && !IsNil(o.Disabled) {
		return true
	}

	return false
}

// SetDisabled gets a reference to the given bool and assigns it to the Disabled field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetDisabled(v bool) {
	o.Disabled = &v
}

// GetEventTypes returns the EventTypes field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetEventTypes() []Dedicatedv1beta1AuditLogFilterRuleEventType {
	if o == nil || IsNil(o.EventTypes) {
		var ret []Dedicatedv1beta1AuditLogFilterRuleEventType
		return ret
	}
	return o.EventTypes
}

// GetEventTypesOk returns a tuple with the EventTypes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetEventTypesOk() ([]Dedicatedv1beta1AuditLogFilterRuleEventType, bool) {
	if o == nil || IsNil(o.EventTypes) {
		return nil, false
	}
	return o.EventTypes, true
}

// HasEventTypes returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasEventTypes() bool {
	if o != nil && !IsNil(o.EventTypes) {
		return true
	}

	return false
}

// SetEventTypes gets a reference to the given []Dedicatedv1beta1AuditLogFilterRuleEventType and assigns it to the EventTypes field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetEventTypes(v []Dedicatedv1beta1AuditLogFilterRuleEventType) {
	o.EventTypes = v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetDbExpr returns the DbExpr field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetDbExpr() string {
	if o == nil || IsNil(o.DbExpr) {
		var ret string
		return ret
	}
	return *o.DbExpr
}

// GetDbExprOk returns a tuple with the DbExpr field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetDbExprOk() (*string, bool) {
	if o == nil || IsNil(o.DbExpr) {
		return nil, false
	}
	return o.DbExpr, true
}

// HasDbExpr returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasDbExpr() bool {
	if o != nil && !IsNil(o.DbExpr) {
		return true
	}

	return false
}

// SetDbExpr gets a reference to the given string and assigns it to the DbExpr field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetDbExpr(v string) {
	o.DbExpr = &v
}

// GetTableExpr returns the TableExpr field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetTableExpr() string {
	if o == nil || IsNil(o.TableExpr) {
		var ret string
		return ret
	}
	return *o.TableExpr
}

// GetTableExprOk returns a tuple with the TableExpr field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetTableExprOk() (*string, bool) {
	if o == nil || IsNil(o.TableExpr) {
		return nil, false
	}
	return o.TableExpr, true
}

// HasTableExpr returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasTableExpr() bool {
	if o != nil && !IsNil(o.TableExpr) {
		return true
	}

	return false
}

// SetTableExpr gets a reference to the given string and assigns it to the TableExpr field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetTableExpr(v string) {
	o.TableExpr = &v
}

// GetAccessTypeList returns the AccessTypeList field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetAccessTypeList() []string {
	if o == nil || IsNil(o.AccessTypeList) {
		var ret []string
		return ret
	}
	return o.AccessTypeList
}

// GetAccessTypeListOk returns a tuple with the AccessTypeList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) GetAccessTypeListOk() ([]string, bool) {
	if o == nil || IsNil(o.AccessTypeList) {
		return nil, false
	}
	return o.AccessTypeList, true
}

// HasAccessTypeList returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogFilterRule) HasAccessTypeList() bool {
	if o != nil && !IsNil(o.AccessTypeList) {
		return true
	}

	return false
}

// SetAccessTypeList gets a reference to the given []string and assigns it to the AccessTypeList field.
func (o *Dedicatedv1beta1AuditLogFilterRule) SetAccessTypeList(v []string) {
	o.AccessTypeList = v
}

func (o Dedicatedv1beta1AuditLogFilterRule) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1AuditLogFilterRule) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.AuditLogFilterRuleId) {
		toSerialize["auditLogFilterRuleId"] = o.AuditLogFilterRuleId
	}
	toSerialize["clusterId"] = o.ClusterId
	if !IsNil(o.UserExpr) {
		toSerialize["userExpr"] = o.UserExpr
	}
	if !IsNil(o.Disabled) {
		toSerialize["disabled"] = o.Disabled
	}
	if !IsNil(o.EventTypes) {
		toSerialize["eventTypes"] = o.EventTypes
	}
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}
	if !IsNil(o.DbExpr) {
		toSerialize["dbExpr"] = o.DbExpr
	}
	if !IsNil(o.TableExpr) {
		toSerialize["tableExpr"] = o.TableExpr
	}
	if !IsNil(o.AccessTypeList) {
		toSerialize["accessTypeList"] = o.AccessTypeList
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1AuditLogFilterRule) UnmarshalJSON(data []byte) (err error) {
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

	varDedicatedv1beta1AuditLogFilterRule := _Dedicatedv1beta1AuditLogFilterRule{}

	err = json.Unmarshal(data, &varDedicatedv1beta1AuditLogFilterRule)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1AuditLogFilterRule(varDedicatedv1beta1AuditLogFilterRule)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "auditLogFilterRuleId")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "userExpr")
		delete(additionalProperties, "disabled")
		delete(additionalProperties, "eventTypes")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "dbExpr")
		delete(additionalProperties, "tableExpr")
		delete(additionalProperties, "accessTypeList")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1AuditLogFilterRule struct {
	value *Dedicatedv1beta1AuditLogFilterRule
	isSet bool
}

func (v NullableDedicatedv1beta1AuditLogFilterRule) Get() *Dedicatedv1beta1AuditLogFilterRule {
	return v.value
}

func (v *NullableDedicatedv1beta1AuditLogFilterRule) Set(val *Dedicatedv1beta1AuditLogFilterRule) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1AuditLogFilterRule) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1AuditLogFilterRule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1AuditLogFilterRule(val *Dedicatedv1beta1AuditLogFilterRule) *NullableDedicatedv1beta1AuditLogFilterRule {
	return &NullableDedicatedv1beta1AuditLogFilterRule{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1AuditLogFilterRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1AuditLogFilterRule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
