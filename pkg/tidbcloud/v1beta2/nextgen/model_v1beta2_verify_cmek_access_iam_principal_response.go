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

// checks if the V1beta2VerifyCmekAccessIamPrincipalResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta2VerifyCmekAccessIamPrincipalResponse{}

// V1beta2VerifyCmekAccessIamPrincipalResponse struct for V1beta2VerifyCmekAccessIamPrincipalResponse
type V1beta2VerifyCmekAccessIamPrincipalResponse struct {
	Valid                *bool   `json:"valid,omitempty"`
	FailReason           *string `json:"failReason,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta2VerifyCmekAccessIamPrincipalResponse V1beta2VerifyCmekAccessIamPrincipalResponse

// NewV1beta2VerifyCmekAccessIamPrincipalResponse instantiates a new V1beta2VerifyCmekAccessIamPrincipalResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta2VerifyCmekAccessIamPrincipalResponse() *V1beta2VerifyCmekAccessIamPrincipalResponse {
	this := V1beta2VerifyCmekAccessIamPrincipalResponse{}
	return &this
}

// NewV1beta2VerifyCmekAccessIamPrincipalResponseWithDefaults instantiates a new V1beta2VerifyCmekAccessIamPrincipalResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta2VerifyCmekAccessIamPrincipalResponseWithDefaults() *V1beta2VerifyCmekAccessIamPrincipalResponse {
	this := V1beta2VerifyCmekAccessIamPrincipalResponse{}
	return &this
}

// GetValid returns the Valid field value if set, zero value otherwise.
func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) GetValid() bool {
	if o == nil || IsNil(o.Valid) {
		var ret bool
		return ret
	}
	return *o.Valid
}

// GetValidOk returns a tuple with the Valid field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) GetValidOk() (*bool, bool) {
	if o == nil || IsNil(o.Valid) {
		return nil, false
	}
	return o.Valid, true
}

// HasValid returns a boolean if a field has been set.
func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) HasValid() bool {
	if o != nil && !IsNil(o.Valid) {
		return true
	}

	return false
}

// SetValid gets a reference to the given bool and assigns it to the Valid field.
func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) SetValid(v bool) {
	o.Valid = &v
}

// GetFailReason returns the FailReason field value if set, zero value otherwise.
func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) GetFailReason() string {
	if o == nil || IsNil(o.FailReason) {
		var ret string
		return ret
	}
	return *o.FailReason
}

// GetFailReasonOk returns a tuple with the FailReason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) GetFailReasonOk() (*string, bool) {
	if o == nil || IsNil(o.FailReason) {
		return nil, false
	}
	return o.FailReason, true
}

// HasFailReason returns a boolean if a field has been set.
func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) HasFailReason() bool {
	if o != nil && !IsNil(o.FailReason) {
		return true
	}

	return false
}

// SetFailReason gets a reference to the given string and assigns it to the FailReason field.
func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) SetFailReason(v string) {
	o.FailReason = &v
}

func (o V1beta2VerifyCmekAccessIamPrincipalResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta2VerifyCmekAccessIamPrincipalResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Valid) {
		toSerialize["valid"] = o.Valid
	}
	if !IsNil(o.FailReason) {
		toSerialize["failReason"] = o.FailReason
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta2VerifyCmekAccessIamPrincipalResponse) UnmarshalJSON(data []byte) (err error) {
	varV1beta2VerifyCmekAccessIamPrincipalResponse := _V1beta2VerifyCmekAccessIamPrincipalResponse{}

	err = json.Unmarshal(data, &varV1beta2VerifyCmekAccessIamPrincipalResponse)

	if err != nil {
		return err
	}

	*o = V1beta2VerifyCmekAccessIamPrincipalResponse(varV1beta2VerifyCmekAccessIamPrincipalResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "valid")
		delete(additionalProperties, "failReason")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta2VerifyCmekAccessIamPrincipalResponse struct {
	value *V1beta2VerifyCmekAccessIamPrincipalResponse
	isSet bool
}

func (v NullableV1beta2VerifyCmekAccessIamPrincipalResponse) Get() *V1beta2VerifyCmekAccessIamPrincipalResponse {
	return v.value
}

func (v *NullableV1beta2VerifyCmekAccessIamPrincipalResponse) Set(val *V1beta2VerifyCmekAccessIamPrincipalResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta2VerifyCmekAccessIamPrincipalResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta2VerifyCmekAccessIamPrincipalResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta2VerifyCmekAccessIamPrincipalResponse(val *V1beta2VerifyCmekAccessIamPrincipalResponse) *NullableV1beta2VerifyCmekAccessIamPrincipalResponse {
	return &NullableV1beta2VerifyCmekAccessIamPrincipalResponse{value: val, isSet: true}
}

func (v NullableV1beta2VerifyCmekAccessIamPrincipalResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta2VerifyCmekAccessIamPrincipalResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
