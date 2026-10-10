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

// Nextgenv1beta2PrivateLinkServiceState the model 'Nextgenv1beta2PrivateLinkServiceState'
type Nextgenv1beta2PrivateLinkServiceState string

// List of nextgenv1beta2PrivateLinkServiceState
const (
	NEXTGENV1BETA2PRIVATELINKSERVICESTATE_CREATING Nextgenv1beta2PrivateLinkServiceState = "CREATING"
	NEXTGENV1BETA2PRIVATELINKSERVICESTATE_ACTIVE   Nextgenv1beta2PrivateLinkServiceState = "ACTIVE"
	NEXTGENV1BETA2PRIVATELINKSERVICESTATE_DELETING Nextgenv1beta2PrivateLinkServiceState = "DELETING"
)

// All allowed values of Nextgenv1beta2PrivateLinkServiceState enum
var AllowedNextgenv1beta2PrivateLinkServiceStateEnumValues = []Nextgenv1beta2PrivateLinkServiceState{
	"CREATING",
	"ACTIVE",
	"DELETING",
}

func (v *Nextgenv1beta2PrivateLinkServiceState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Nextgenv1beta2PrivateLinkServiceState(value)
	for _, existing := range AllowedNextgenv1beta2PrivateLinkServiceStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	*v = Nextgenv1beta2PrivateLinkServiceState(value)
	return nil
}

// NewNextgenv1beta2PrivateLinkServiceStateFromValue returns a pointer to a valid Nextgenv1beta2PrivateLinkServiceState for the value passed as argument
func NewNextgenv1beta2PrivateLinkServiceStateFromValue(v string) *Nextgenv1beta2PrivateLinkServiceState {
	ev := Nextgenv1beta2PrivateLinkServiceState(v)
	return &ev
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Nextgenv1beta2PrivateLinkServiceState) IsValid() bool {
	for _, existing := range AllowedNextgenv1beta2PrivateLinkServiceStateEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to nextgenv1beta2PrivateLinkServiceState value
func (v Nextgenv1beta2PrivateLinkServiceState) Ptr() *Nextgenv1beta2PrivateLinkServiceState {
	return &v
}

type NullableNextgenv1beta2PrivateLinkServiceState struct {
	value *Nextgenv1beta2PrivateLinkServiceState
	isSet bool
}

func (v NullableNextgenv1beta2PrivateLinkServiceState) Get() *Nextgenv1beta2PrivateLinkServiceState {
	return v.value
}

func (v *NullableNextgenv1beta2PrivateLinkServiceState) Set(val *Nextgenv1beta2PrivateLinkServiceState) {
	v.value = val
	v.isSet = true
}

func (v NullableNextgenv1beta2PrivateLinkServiceState) IsSet() bool {
	return v.isSet
}

func (v *NullableNextgenv1beta2PrivateLinkServiceState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNextgenv1beta2PrivateLinkServiceState(val *Nextgenv1beta2PrivateLinkServiceState) *NullableNextgenv1beta2PrivateLinkServiceState {
	return &NullableNextgenv1beta2PrivateLinkServiceState{value: val, isSet: true}
}

func (v NullableNextgenv1beta2PrivateLinkServiceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNextgenv1beta2PrivateLinkServiceState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
