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
	"time"
)

// checks if the TidbCloudOpenApiserverlessv1beta1Cluster type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbCloudOpenApiserverlessv1beta1Cluster{}

// TidbCloudOpenApiserverlessv1beta1Cluster Message for a TiDB Cloud Starter or Essential cluster.
type TidbCloudOpenApiserverlessv1beta1Cluster struct {
	// The unique identifier for the TiDB cluster, which is generated by the API and follows the format `clusters/{clusterId}`.
	Name *string `json:"name,omitempty"`
	// The ID of the cluster.
	ClusterId *string `json:"clusterId,omitempty"`
	// The user-defined name of the cluster.
	DisplayName string `json:"displayName" validate:"regexp=^[A-Za-z0-9][-A-Za-z0-9]{2,62}[A-Za-z0-9]$"`
	// The region where the cluster is deployed.
	Region Commonv1beta1Region `json:"region"`
	// The [spending limit](https://docs.pingcap.com/tidbcloud/tidb-cloud-glossary/#spending-limit) for the cluster.
	SpendingLimit *ClusterSpendingLimit `json:"spendingLimit,omitempty"`
	// The schedule and retention rules for automated database backups.
	AutomatedBackupPolicy *V1beta1ClusterAutomatedBackupPolicy `json:"automatedBackupPolicy,omitempty"`
	// The connection endpoints for accessing the cluster.
	Endpoints *V1beta1ClusterEndpoints `json:"endpoints,omitempty"`
	// The root password of the cluster. It must be between 8 and 64 characters long and can contain letters, numbers, and special characters.
	RootPassword *string `json:"rootPassword,omitempty" validate:"regexp=^.{8,64}$"`
	// The data encryption settings for the cluster.
	EncryptionConfig     *V1beta1ClusterEncryptionConfig               `json:"encryptionConfig,omitempty"`
	HighAvailabilityType *Serverlessv1beta1ClusterHighAvailabilityType `json:"highAvailabilityType,omitempty"`
	// The TiDB version of the cluster.
	Version *string `json:"version,omitempty"`
	// The email address of the user who create the cluster.
	CreatedBy *string `json:"createdBy,omitempty"`
	// The unique prefix automatically generated for SQL usernames on this cluster. TiDB Cloud uses this prefix to distinguish between clusters. For more information, see [User name prefix](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#user-name-prefix).
	UserPrefix *string `json:"userPrefix,omitempty"`
	// The current state of the cluster.
	State *Commonv1beta1ClusterState `json:"state,omitempty"`
	// The labels for the cluster.  - `tidb.cloud/organization`: the ID of the organization where the cluster belongs. - `tidb.cloud/project`: the ID of the project where the cluster belongs.
	Labels *map[string]string `json:"labels,omitempty"`
	// The annotations for the cluster. The following lists some predefined annotations.  - `tidb.cloud/has-set-password`: indicates whether the cluster has a root password set. - `tidb.cloud/available-features`: lists available features of the cluster.
	Annotations *map[string]string `json:"annotations,omitempty"`
	// The timestamp when the cluster was created, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// The timestamp when the cluster was last updated, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The audit log configuration for the cluster.
	AuditLogConfig *V1beta1ClusterAuditLogConfig `json:"auditLogConfig,omitempty"`
	ClusterPlan    *ClusterClusterPlan           `json:"clusterPlan,omitempty"`
	// The auto-scaling configuration for the cluster.
	AutoScaling          *V1beta1ClusterAutoScaling `json:"autoScaling,omitempty"`
	ServicePlan          *Commonv1beta1ServicePlan  `json:"servicePlan,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _TidbCloudOpenApiserverlessv1beta1Cluster TidbCloudOpenApiserverlessv1beta1Cluster

// NewTidbCloudOpenApiserverlessv1beta1Cluster instantiates a new TidbCloudOpenApiserverlessv1beta1Cluster object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbCloudOpenApiserverlessv1beta1Cluster(displayName string, region Commonv1beta1Region) *TidbCloudOpenApiserverlessv1beta1Cluster {
	this := TidbCloudOpenApiserverlessv1beta1Cluster{}
	this.DisplayName = displayName
	this.Region = region
	return &this
}

// NewTidbCloudOpenApiserverlessv1beta1ClusterWithDefaults instantiates a new TidbCloudOpenApiserverlessv1beta1Cluster object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbCloudOpenApiserverlessv1beta1ClusterWithDefaults() *TidbCloudOpenApiserverlessv1beta1Cluster {
	this := TidbCloudOpenApiserverlessv1beta1Cluster{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetName(v string) {
	o.Name = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}
	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetDisplayName returns the DisplayName field value
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetDisplayName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DisplayName, true
}

// SetDisplayName sets field value
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetDisplayName(v string) {
	o.DisplayName = v
}

// GetRegion returns the Region field value
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetRegion() Commonv1beta1Region {
	if o == nil {
		var ret Commonv1beta1Region
		return ret
	}

	return o.Region
}

// GetRegionOk returns a tuple with the Region field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetRegionOk() (*Commonv1beta1Region, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Region, true
}

// SetRegion sets field value
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetRegion(v Commonv1beta1Region) {
	o.Region = v
}

// GetSpendingLimit returns the SpendingLimit field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetSpendingLimit() ClusterSpendingLimit {
	if o == nil || IsNil(o.SpendingLimit) {
		var ret ClusterSpendingLimit
		return ret
	}
	return *o.SpendingLimit
}

// GetSpendingLimitOk returns a tuple with the SpendingLimit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetSpendingLimitOk() (*ClusterSpendingLimit, bool) {
	if o == nil || IsNil(o.SpendingLimit) {
		return nil, false
	}
	return o.SpendingLimit, true
}

// HasSpendingLimit returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasSpendingLimit() bool {
	if o != nil && !IsNil(o.SpendingLimit) {
		return true
	}

	return false
}

// SetSpendingLimit gets a reference to the given ClusterSpendingLimit and assigns it to the SpendingLimit field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetSpendingLimit(v ClusterSpendingLimit) {
	o.SpendingLimit = &v
}

// GetAutomatedBackupPolicy returns the AutomatedBackupPolicy field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAutomatedBackupPolicy() V1beta1ClusterAutomatedBackupPolicy {
	if o == nil || IsNil(o.AutomatedBackupPolicy) {
		var ret V1beta1ClusterAutomatedBackupPolicy
		return ret
	}
	return *o.AutomatedBackupPolicy
}

// GetAutomatedBackupPolicyOk returns a tuple with the AutomatedBackupPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAutomatedBackupPolicyOk() (*V1beta1ClusterAutomatedBackupPolicy, bool) {
	if o == nil || IsNil(o.AutomatedBackupPolicy) {
		return nil, false
	}
	return o.AutomatedBackupPolicy, true
}

// HasAutomatedBackupPolicy returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasAutomatedBackupPolicy() bool {
	if o != nil && !IsNil(o.AutomatedBackupPolicy) {
		return true
	}

	return false
}

// SetAutomatedBackupPolicy gets a reference to the given V1beta1ClusterAutomatedBackupPolicy and assigns it to the AutomatedBackupPolicy field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetAutomatedBackupPolicy(v V1beta1ClusterAutomatedBackupPolicy) {
	o.AutomatedBackupPolicy = &v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetEndpoints() V1beta1ClusterEndpoints {
	if o == nil || IsNil(o.Endpoints) {
		var ret V1beta1ClusterEndpoints
		return ret
	}
	return *o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetEndpointsOk() (*V1beta1ClusterEndpoints, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}
	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given V1beta1ClusterEndpoints and assigns it to the Endpoints field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetEndpoints(v V1beta1ClusterEndpoints) {
	o.Endpoints = &v
}

// GetRootPassword returns the RootPassword field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetRootPassword() string {
	if o == nil || IsNil(o.RootPassword) {
		var ret string
		return ret
	}
	return *o.RootPassword
}

// GetRootPasswordOk returns a tuple with the RootPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetRootPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.RootPassword) {
		return nil, false
	}
	return o.RootPassword, true
}

// HasRootPassword returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasRootPassword() bool {
	if o != nil && !IsNil(o.RootPassword) {
		return true
	}

	return false
}

// SetRootPassword gets a reference to the given string and assigns it to the RootPassword field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetRootPassword(v string) {
	o.RootPassword = &v
}

// GetEncryptionConfig returns the EncryptionConfig field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetEncryptionConfig() V1beta1ClusterEncryptionConfig {
	if o == nil || IsNil(o.EncryptionConfig) {
		var ret V1beta1ClusterEncryptionConfig
		return ret
	}
	return *o.EncryptionConfig
}

// GetEncryptionConfigOk returns a tuple with the EncryptionConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetEncryptionConfigOk() (*V1beta1ClusterEncryptionConfig, bool) {
	if o == nil || IsNil(o.EncryptionConfig) {
		return nil, false
	}
	return o.EncryptionConfig, true
}

// HasEncryptionConfig returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasEncryptionConfig() bool {
	if o != nil && !IsNil(o.EncryptionConfig) {
		return true
	}

	return false
}

// SetEncryptionConfig gets a reference to the given V1beta1ClusterEncryptionConfig and assigns it to the EncryptionConfig field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetEncryptionConfig(v V1beta1ClusterEncryptionConfig) {
	o.EncryptionConfig = &v
}

// GetHighAvailabilityType returns the HighAvailabilityType field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetHighAvailabilityType() Serverlessv1beta1ClusterHighAvailabilityType {
	if o == nil || IsNil(o.HighAvailabilityType) {
		var ret Serverlessv1beta1ClusterHighAvailabilityType
		return ret
	}
	return *o.HighAvailabilityType
}

// GetHighAvailabilityTypeOk returns a tuple with the HighAvailabilityType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetHighAvailabilityTypeOk() (*Serverlessv1beta1ClusterHighAvailabilityType, bool) {
	if o == nil || IsNil(o.HighAvailabilityType) {
		return nil, false
	}
	return o.HighAvailabilityType, true
}

// HasHighAvailabilityType returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasHighAvailabilityType() bool {
	if o != nil && !IsNil(o.HighAvailabilityType) {
		return true
	}

	return false
}

// SetHighAvailabilityType gets a reference to the given Serverlessv1beta1ClusterHighAvailabilityType and assigns it to the HighAvailabilityType field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetHighAvailabilityType(v Serverlessv1beta1ClusterHighAvailabilityType) {
	o.HighAvailabilityType = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetVersion(v string) {
	o.Version = &v
}

// GetCreatedBy returns the CreatedBy field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetCreatedBy() string {
	if o == nil || IsNil(o.CreatedBy) {
		var ret string
		return ret
	}
	return *o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetCreatedByOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedBy) {
		return nil, false
	}
	return o.CreatedBy, true
}

// HasCreatedBy returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasCreatedBy() bool {
	if o != nil && !IsNil(o.CreatedBy) {
		return true
	}

	return false
}

// SetCreatedBy gets a reference to the given string and assigns it to the CreatedBy field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetCreatedBy(v string) {
	o.CreatedBy = &v
}

// GetUserPrefix returns the UserPrefix field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUserPrefix() string {
	if o == nil || IsNil(o.UserPrefix) {
		var ret string
		return ret
	}
	return *o.UserPrefix
}

// GetUserPrefixOk returns a tuple with the UserPrefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUserPrefixOk() (*string, bool) {
	if o == nil || IsNil(o.UserPrefix) {
		return nil, false
	}
	return o.UserPrefix, true
}

// HasUserPrefix returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasUserPrefix() bool {
	if o != nil && !IsNil(o.UserPrefix) {
		return true
	}

	return false
}

// SetUserPrefix gets a reference to the given string and assigns it to the UserPrefix field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetUserPrefix(v string) {
	o.UserPrefix = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetState() Commonv1beta1ClusterState {
	if o == nil || IsNil(o.State) {
		var ret Commonv1beta1ClusterState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetStateOk() (*Commonv1beta1ClusterState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Commonv1beta1ClusterState and assigns it to the State field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetState(v Commonv1beta1ClusterState) {
	o.State = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetAnnotations returns the Annotations field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAnnotations() map[string]string {
	if o == nil || IsNil(o.Annotations) {
		var ret map[string]string
		return ret
	}
	return *o.Annotations
}

// GetAnnotationsOk returns a tuple with the Annotations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAnnotationsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Annotations) {
		return nil, false
	}
	return o.Annotations, true
}

// HasAnnotations returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasAnnotations() bool {
	if o != nil && !IsNil(o.Annotations) {
		return true
	}

	return false
}

// SetAnnotations gets a reference to the given map[string]string and assigns it to the Annotations field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetAnnotations(v map[string]string) {
	o.Annotations = &v
}

// GetCreateTime returns the CreateTime field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetCreateTime() time.Time {
	if o == nil || IsNil(o.CreateTime) {
		var ret time.Time
		return ret
	}
	return *o.CreateTime
}

// GetCreateTimeOk returns a tuple with the CreateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetCreateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateTime) {
		return nil, false
	}
	return o.CreateTime, true
}

// HasCreateTime returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasCreateTime() bool {
	if o != nil && !IsNil(o.CreateTime) {
		return true
	}

	return false
}

// SetCreateTime gets a reference to the given time.Time and assigns it to the CreateTime field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetCreateTime(v time.Time) {
	o.CreateTime = &v
}

// GetUpdateTime returns the UpdateTime field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUpdateTime() time.Time {
	if o == nil || IsNil(o.UpdateTime) {
		var ret time.Time
		return ret
	}
	return *o.UpdateTime
}

// GetUpdateTimeOk returns a tuple with the UpdateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUpdateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.UpdateTime) {
		return nil, false
	}
	return o.UpdateTime, true
}

// HasUpdateTime returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasUpdateTime() bool {
	if o != nil && !IsNil(o.UpdateTime) {
		return true
	}

	return false
}

// SetUpdateTime gets a reference to the given time.Time and assigns it to the UpdateTime field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetUpdateTime(v time.Time) {
	o.UpdateTime = &v
}

// GetAuditLogConfig returns the AuditLogConfig field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAuditLogConfig() V1beta1ClusterAuditLogConfig {
	if o == nil || IsNil(o.AuditLogConfig) {
		var ret V1beta1ClusterAuditLogConfig
		return ret
	}
	return *o.AuditLogConfig
}

// GetAuditLogConfigOk returns a tuple with the AuditLogConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAuditLogConfigOk() (*V1beta1ClusterAuditLogConfig, bool) {
	if o == nil || IsNil(o.AuditLogConfig) {
		return nil, false
	}
	return o.AuditLogConfig, true
}

// HasAuditLogConfig returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasAuditLogConfig() bool {
	if o != nil && !IsNil(o.AuditLogConfig) {
		return true
	}

	return false
}

// SetAuditLogConfig gets a reference to the given V1beta1ClusterAuditLogConfig and assigns it to the AuditLogConfig field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetAuditLogConfig(v V1beta1ClusterAuditLogConfig) {
	o.AuditLogConfig = &v
}

// GetClusterPlan returns the ClusterPlan field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetClusterPlan() ClusterClusterPlan {
	if o == nil || IsNil(o.ClusterPlan) {
		var ret ClusterClusterPlan
		return ret
	}
	return *o.ClusterPlan
}

// GetClusterPlanOk returns a tuple with the ClusterPlan field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetClusterPlanOk() (*ClusterClusterPlan, bool) {
	if o == nil || IsNil(o.ClusterPlan) {
		return nil, false
	}
	return o.ClusterPlan, true
}

// HasClusterPlan returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasClusterPlan() bool {
	if o != nil && !IsNil(o.ClusterPlan) {
		return true
	}

	return false
}

// SetClusterPlan gets a reference to the given ClusterClusterPlan and assigns it to the ClusterPlan field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetClusterPlan(v ClusterClusterPlan) {
	o.ClusterPlan = &v
}

// GetAutoScaling returns the AutoScaling field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAutoScaling() V1beta1ClusterAutoScaling {
	if o == nil || IsNil(o.AutoScaling) {
		var ret V1beta1ClusterAutoScaling
		return ret
	}
	return *o.AutoScaling
}

// GetAutoScalingOk returns a tuple with the AutoScaling field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAutoScalingOk() (*V1beta1ClusterAutoScaling, bool) {
	if o == nil || IsNil(o.AutoScaling) {
		return nil, false
	}
	return o.AutoScaling, true
}

// HasAutoScaling returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasAutoScaling() bool {
	if o != nil && !IsNil(o.AutoScaling) {
		return true
	}

	return false
}

// SetAutoScaling gets a reference to the given V1beta1ClusterAutoScaling and assigns it to the AutoScaling field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetAutoScaling(v V1beta1ClusterAutoScaling) {
	o.AutoScaling = &v
}

// GetServicePlan returns the ServicePlan field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetServicePlan() Commonv1beta1ServicePlan {
	if o == nil || IsNil(o.ServicePlan) {
		var ret Commonv1beta1ServicePlan
		return ret
	}
	return *o.ServicePlan
}

// GetServicePlanOk returns a tuple with the ServicePlan field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetServicePlanOk() (*Commonv1beta1ServicePlan, bool) {
	if o == nil || IsNil(o.ServicePlan) {
		return nil, false
	}
	return o.ServicePlan, true
}

// HasServicePlan returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasServicePlan() bool {
	if o != nil && !IsNil(o.ServicePlan) {
		return true
	}

	return false
}

// SetServicePlan gets a reference to the given Commonv1beta1ServicePlan and assigns it to the ServicePlan field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetServicePlan(v Commonv1beta1ServicePlan) {
	o.ServicePlan = &v
}

func (o TidbCloudOpenApiserverlessv1beta1Cluster) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbCloudOpenApiserverlessv1beta1Cluster) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.ClusterId) {
		toSerialize["clusterId"] = o.ClusterId
	}
	toSerialize["displayName"] = o.DisplayName
	toSerialize["region"] = o.Region
	if !IsNil(o.SpendingLimit) {
		toSerialize["spendingLimit"] = o.SpendingLimit
	}
	if !IsNil(o.AutomatedBackupPolicy) {
		toSerialize["automatedBackupPolicy"] = o.AutomatedBackupPolicy
	}
	if !IsNil(o.Endpoints) {
		toSerialize["endpoints"] = o.Endpoints
	}
	if !IsNil(o.RootPassword) {
		toSerialize["rootPassword"] = o.RootPassword
	}
	if !IsNil(o.EncryptionConfig) {
		toSerialize["encryptionConfig"] = o.EncryptionConfig
	}
	if !IsNil(o.HighAvailabilityType) {
		toSerialize["highAvailabilityType"] = o.HighAvailabilityType
	}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !IsNil(o.CreatedBy) {
		toSerialize["createdBy"] = o.CreatedBy
	}
	if !IsNil(o.UserPrefix) {
		toSerialize["userPrefix"] = o.UserPrefix
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	if !IsNil(o.Annotations) {
		toSerialize["annotations"] = o.Annotations
	}
	if !IsNil(o.CreateTime) {
		toSerialize["createTime"] = o.CreateTime
	}
	if !IsNil(o.UpdateTime) {
		toSerialize["updateTime"] = o.UpdateTime
	}
	if !IsNil(o.AuditLogConfig) {
		toSerialize["auditLogConfig"] = o.AuditLogConfig
	}
	if !IsNil(o.ClusterPlan) {
		toSerialize["clusterPlan"] = o.ClusterPlan
	}
	if !IsNil(o.AutoScaling) {
		toSerialize["autoScaling"] = o.AutoScaling
	}
	if !IsNil(o.ServicePlan) {
		toSerialize["servicePlan"] = o.ServicePlan
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TidbCloudOpenApiserverlessv1beta1Cluster) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"displayName",
		"region",
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

	varTidbCloudOpenApiserverlessv1beta1Cluster := _TidbCloudOpenApiserverlessv1beta1Cluster{}

	err = json.Unmarshal(data, &varTidbCloudOpenApiserverlessv1beta1Cluster)

	if err != nil {
		return err
	}

	*o = TidbCloudOpenApiserverlessv1beta1Cluster(varTidbCloudOpenApiserverlessv1beta1Cluster)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "region")
		delete(additionalProperties, "spendingLimit")
		delete(additionalProperties, "automatedBackupPolicy")
		delete(additionalProperties, "endpoints")
		delete(additionalProperties, "rootPassword")
		delete(additionalProperties, "encryptionConfig")
		delete(additionalProperties, "highAvailabilityType")
		delete(additionalProperties, "version")
		delete(additionalProperties, "createdBy")
		delete(additionalProperties, "userPrefix")
		delete(additionalProperties, "state")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "annotations")
		delete(additionalProperties, "createTime")
		delete(additionalProperties, "updateTime")
		delete(additionalProperties, "auditLogConfig")
		delete(additionalProperties, "clusterPlan")
		delete(additionalProperties, "autoScaling")
		delete(additionalProperties, "servicePlan")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbCloudOpenApiserverlessv1beta1Cluster struct {
	value *TidbCloudOpenApiserverlessv1beta1Cluster
	isSet bool
}

func (v NullableTidbCloudOpenApiserverlessv1beta1Cluster) Get() *TidbCloudOpenApiserverlessv1beta1Cluster {
	return v.value
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1Cluster) Set(val *TidbCloudOpenApiserverlessv1beta1Cluster) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbCloudOpenApiserverlessv1beta1Cluster) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1Cluster) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbCloudOpenApiserverlessv1beta1Cluster(val *TidbCloudOpenApiserverlessv1beta1Cluster) *NullableTidbCloudOpenApiserverlessv1beta1Cluster {
	return &NullableTidbCloudOpenApiserverlessv1beta1Cluster{value: val, isSet: true}
}

func (v NullableTidbCloudOpenApiserverlessv1beta1Cluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1Cluster) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
