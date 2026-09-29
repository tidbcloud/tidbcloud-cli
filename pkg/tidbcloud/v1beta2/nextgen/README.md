# Go API client for nextgen

*TiDB Cloud API is in beta.*

This API manages [TiDB Cloud Premium](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-premium) instances. For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/api/tidb-cloud-api-overview/).

The public connection setting endpoints also support TiDB Cloud Essential V2 instances.

# Overview

The TiDB Cloud Premium API (v1beta2) provides [REST](https://en.wikipedia.org/wiki/REST) endpoints to manage TiDB Cloud Premium instances and related resources.

You can use this API to manage the following resources:

- **TiDB Cloud Premium instance**: manage the lifecycle and configuration of TiDB Cloud Premium instances, including passwords, CA certificates, and cloud provider information.
- **Public connection setting**: enable or disable public connections and manage IP access lists for Premium and Essential V2 instances.
- **Customer-managed encryption key**: retrieve the IAM principal and verify KMS access before creating a Premium instance on AWS or Alibaba Cloud.
- **Backup**: manage backups for TiDB Cloud Premium instances, including backup-based restore.
- **Region**: retrieve available regions for deploying TiDB Cloud Premium instances.

# Get Started

This guide helps you make your first API call to the TiDB Cloud Premium API. You will learn how to authenticate a request, build a request, and interpret the response.

1. Create a [TiDB Cloud account](https://tidbcloud.com/signup) if you do not already have one.
2. In the [TiDB Cloud console](https://tidbcloud.com/), go to **Organization** > **API Keys** and create an API key. For more information, see [API key management](#section/Authentication/API-key-management).
3. Make your first API call.

 To get all TiDB Cloud Premium instances in your organization, run the following command in your terminal. Replace `YOUR_PUBLIC_KEY` and `YOUR_PRIVATE_KEY` with your own key values.

 ```bash
 curl --digest \\
   --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\
   --request GET \\
   --url 'https://cloud.tidbapi.com/v1beta2/tidbs' \\
   --header 'Accept: application/json'
 ```

4. The API returns a JSON list of your TiDB Cloud Premium instances. If none exist, the response contains an empty list.

# Authentication

The TiDB Cloud API supports [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication) with API keys and OAuth Bearer tokens. OAuth clients send the token in the `Authorization: Bearer <token>` header. It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).

## API key overview

- The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation.
- The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles).
- You must provide the correct API key in every request. Otherwise, TiDB Cloud responds with a `401` error.

## API key management

### Create an API key

Only the **owner** of an organization can create an API key.

To create an API key in an organization, perform the following steps:

1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner.
2. In the left navigation pane, click **Organization Settings** > **API Keys**.
3. On the **API Keys** page, click **Create API Key**.
4. Enter a description for your API key.
5. Configure the role and scope for the API key. For more information about the permissions of a role, see [User roles](https://docs.pingcap.com/tidbcloud/manage-user-access/#user-roles).
6. Click **Generate API Key**. Copy and save the public key and the private key.
7. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again.
8. Click **Done**.

### View details of an API key

To view details of an API key, perform the following steps:

1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner.
2. In the left navigation pane, click **Organization Settings** > **API Keys**.
3. You can view the details of the API keys on the page.

### Edit an API key

Only the **owner** of an organization can modify an API key.

To edit an API key in an organization, perform the following steps:

1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner.
2. In the left navigation pane, click **Organization Settings** > **API Keys**.
3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Update Role**.
4. You can update the description and role of the API key.
5. Click **Update**.

### Delete an API key

Only the **owner** of an organization can delete an API key.

To delete an API key in an organization, perform the following steps:

1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner.
2. In the left navigation pane, click **Organization Settings** > **API Keys**.
3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**.
4. Click **I understand, delete it.**

# API Changelog

This changelog lists all changes to the TiDB Cloud Premium API (v1beta2).

<!-- In reverse chronological order -->

## 20260924

- Published CMEK IAM principal retrieval and verification APIs for Premium instances on AWS and Alibaba Cloud.

## 20260923

- Added APIs to get and update public connection settings for Premium and Essential V2 instances.

## 20260428

- Initial release of the TiDB Cloud Premium API (v1beta2), including the following resources and endpoints:
 * TiDB Cloud Premium instance
  * [List TiDB Cloud Premium instances](#tag/TiDB-Instance/operation/TidbService_ListTidbs)
  * [Create a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_CreateTidb)
  * [Get a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetTidb)
  * [Delete a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_DeleteTidb)
  * [Update a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_UpdateTidb)
  * [Reset the root password of a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_ResetRootPassword)
  * [Get cloud provider information for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCloudProviderInfo)
  * [Get the CA certificate download URL for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCaCertificateDownloadUrl)
 * Backup
  * [List backups for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_ListTidbBackups)
  * [Delete a backup for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_DeleteTidbBackup)
  * [Restore a TiDB Cloud Premium instance from a backup](#tag/Backup/operation/TidbService_RestoreTidb)
  * [Get restore status for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_GetRestoreStatus)
 * Region
  * [List regions](#tag/Region/operation/RegionService_ListRegions)


## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: v1beta2
- Package version: 1.0.0
- Generator version: 7.12.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen

## Installation

Install the following dependencies:

```sh
go get github.com/stretchr/testify/assert
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```go
import nextgen "github.com/GIT_USER_ID/GIT_REPO_ID"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```go
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `nextgen.ContextServerIndex` of type `int`.

```go
ctx := context.WithValue(context.Background(), nextgen.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `nextgen.ContextServerVariables` of type `map[string]string`.

```go
ctx := context.WithValue(context.Background(), nextgen.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `nextgen.ContextOperationServerIndices` and `nextgen.ContextOperationServerVariables` context maps.

```go
ctx := context.WithValue(context.Background(), nextgen.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), nextgen.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *https://cloud.tidbapi.com/v1beta2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*BackupAPI* | [**BackupServiceDeleteBackup**](docs/BackupAPI.md#backupservicedeletebackup) | **Delete** /backups/{backupId} | Delete a backup
*BackupAPI* | [**BackupServiceListBackups**](docs/BackupAPI.md#backupservicelistbackups) | **Get** /backups | List backups
*BackupAPI* | [**TidbServiceDeleteTidbBackup**](docs/BackupAPI.md#tidbservicedeletetidbbackup) | **Delete** /tidbs/{tidbId}/backups/{backupId} | Delete a backup for a TiDB Cloud Premium instance
*BackupAPI* | [**TidbServiceGetRestoreStatus**](docs/BackupAPI.md#tidbservicegetrestorestatus) | **Get** /tidbs/{tidbId}:getRestoreStatus | Get the restore status for a TiDB Cloud Premium instance
*BackupAPI* | [**TidbServiceListTidbBackups**](docs/BackupAPI.md#tidbservicelisttidbbackups) | **Get** /tidbs/{tidbId}/backups | List backups for a TiDB Cloud Premium instance
*BackupAPI* | [**TidbServiceRestoreTidb**](docs/BackupAPI.md#tidbservicerestoretidb) | **Post** /tidbs:restore | Restore a TiDB Cloud Premium instance from a backup
*CustomerManagedEncryptionKeyServiceAPI* | [**CustomerManagedEncryptionKeyServiceGetCmekAccessIamPrincipal**](docs/CustomerManagedEncryptionKeyServiceAPI.md#customermanagedencryptionkeyservicegetcmekaccessiamprincipal) | **Get** /cmeks:principal | Get CMEK IAM principal
*CustomerManagedEncryptionKeyServiceAPI* | [**CustomerManagedEncryptionKeyServiceVerifyCmekAccessIamPrincipal**](docs/CustomerManagedEncryptionKeyServiceAPI.md#customermanagedencryptionkeyserviceverifycmekaccessiamprincipal) | **Post** /cmeks:verifyPrincipal | Verify CMEK IAM principal
*PublicConnectionSettingServiceAPI* | [**PublicConnectionSettingServiceGetPublicConnectionSetting**](docs/PublicConnectionSettingServiceAPI.md#publicconnectionsettingservicegetpublicconnectionsetting) | **Get** /tidbs/{tidbId}/publicConnectionSetting | Get public connection setting
*PublicConnectionSettingServiceAPI* | [**PublicConnectionSettingServiceUpdatePublicConnectionSetting**](docs/PublicConnectionSettingServiceAPI.md#publicconnectionsettingserviceupdatepublicconnectionsetting) | **Patch** /tidbs/{tidbId}/publicConnectionSetting | Update public connection setting
*RegionAPI* | [**RegionServiceListRegions**](docs/RegionAPI.md#regionservicelistregions) | **Get** /regions | List regions
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceCreateTidb**](docs/TiDBCloudPremiumInstanceAPI.md#tidbservicecreatetidb) | **Post** /tidbs | Create a TiDB Cloud Premium instance
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceDeleteTidb**](docs/TiDBCloudPremiumInstanceAPI.md#tidbservicedeletetidb) | **Delete** /tidbs/{tidbId} | Delete a TiDB Cloud Premium instance
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceGetCaCertificateDownloadUrl**](docs/TiDBCloudPremiumInstanceAPI.md#tidbservicegetcacertificatedownloadurl) | **Get** /tidbs/{tidbId}/caCertificateUrl | Get the CA certificate download URL for a TiDB Cloud Premium instance
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceGetCloudProviderInfo**](docs/TiDBCloudPremiumInstanceAPI.md#tidbservicegetcloudproviderinfo) | **Get** /tidbs/{tidbId}/cloudProviderInfo | Get cloud provider information for a TiDB Cloud Premium instance
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceGetTidb**](docs/TiDBCloudPremiumInstanceAPI.md#tidbservicegettidb) | **Get** /tidbs/{tidbId} | Get a TiDB Cloud Premium instance
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceListTidbs**](docs/TiDBCloudPremiumInstanceAPI.md#tidbservicelisttidbs) | **Get** /tidbs | List TiDB Cloud Premium instances
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceResetRootPassword**](docs/TiDBCloudPremiumInstanceAPI.md#tidbserviceresetrootpassword) | **Post** /tidbs/{tidbId}:resetRootPassword | Reset the root password of a TiDB Cloud Premium instance
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceResumeTidb**](docs/TiDBCloudPremiumInstanceAPI.md#tidbserviceresumetidb) | **Post** /tidbs/{tidbId}:resumeTidb | Resume a TiDB Cloud Premium instance
*TiDBCloudPremiumInstanceAPI* | [**TidbServiceUpdateTidb**](docs/TiDBCloudPremiumInstanceAPI.md#tidbserviceupdatetidb) | **Patch** /tidbs/{tidb.tidbId} | Update a TiDB Cloud Premium instance


## Documentation For Models

 - [BackupServiceListBackupsTidbStatesParameterInner](docs/BackupServiceListBackupsTidbStatesParameterInner.md)
 - [BackupTriggerType](docs/BackupTriggerType.md)
 - [CmekAccessIamPrincipalAliyunCmekPrincipal](docs/CmekAccessIamPrincipalAliyunCmekPrincipal.md)
 - [CmekAccessIamPrincipalAwsCmekPrincipal](docs/CmekAccessIamPrincipalAwsCmekPrincipal.md)
 - [ConnectionReachabilityDetail](docs/ConnectionReachabilityDetail.md)
 - [CustomerManagedEncryptionKeyAliyunKms](docs/CustomerManagedEncryptionKeyAliyunKms.md)
 - [CustomerManagedEncryptionKeyAwsKms](docs/CustomerManagedEncryptionKeyAwsKms.md)
 - [EndpointConnectionType](docs/EndpointConnectionType.md)
 - [Nextgenv1beta2Tidb](docs/Nextgenv1beta2Tidb.md)
 - [ProtobufAny](docs/ProtobufAny.md)
 - [RegionCloudProvider](docs/RegionCloudProvider.md)
 - [ReservedCapacityScope](docs/ReservedCapacityScope.md)
 - [RpcStatus](docs/RpcStatus.md)
 - [TheTiDBCloudPremiumInstanceToUpdate](docs/TheTiDBCloudPremiumInstanceToUpdate.md)
 - [TidbDualLayerDataEncryption](docs/TidbDualLayerDataEncryption.md)
 - [TidbEndpoint](docs/TidbEndpoint.md)
 - [TidbServiceListTidbBackupsStatesParameterInner](docs/TidbServiceListTidbBackupsStatesParameterInner.md)
 - [TidbServiceListTidbBackupsTriggerTypesParameterInner](docs/TidbServiceListTidbBackupsTriggerTypesParameterInner.md)
 - [TidbServiceListTidbsServicePlanParameter](docs/TidbServiceListTidbsServicePlanParameter.md)
 - [TidbServiceResetRootPasswordBody](docs/TidbServiceResetRootPasswordBody.md)
 - [V1beta1ClusterState](docs/V1beta1ClusterState.md)
 - [V1beta1EngineType](docs/V1beta1EngineType.md)
 - [V1beta1Region](docs/V1beta1Region.md)
 - [V1beta1ServicePlan](docs/V1beta1ServicePlan.md)
 - [V1beta1ServicePlanInfo](docs/V1beta1ServicePlanInfo.md)
 - [V1beta2Backup](docs/V1beta2Backup.md)
 - [V1beta2BackupState](docs/V1beta2BackupState.md)
 - [V1beta2BackupStorageAuth](docs/V1beta2BackupStorageAuth.md)
 - [V1beta2BackupType](docs/V1beta2BackupType.md)
 - [V1beta2CaCertificateDownloadUrl](docs/V1beta2CaCertificateDownloadUrl.md)
 - [V1beta2ClassicBackupStorage](docs/V1beta2ClassicBackupStorage.md)
 - [V1beta2CloudProviderInfo](docs/V1beta2CloudProviderInfo.md)
 - [V1beta2CmekAccessIamPrincipal](docs/V1beta2CmekAccessIamPrincipal.md)
 - [V1beta2ConnectionReachability](docs/V1beta2ConnectionReachability.md)
 - [V1beta2CustomerManagedEncryptionKey](docs/V1beta2CustomerManagedEncryptionKey.md)
 - [V1beta2GetRestoreStatusResponse](docs/V1beta2GetRestoreStatusResponse.md)
 - [V1beta2ListBackupsResponse](docs/V1beta2ListBackupsResponse.md)
 - [V1beta2ListRegionsResponse](docs/V1beta2ListRegionsResponse.md)
 - [V1beta2ListTidbBackupsResponse](docs/V1beta2ListTidbBackupsResponse.md)
 - [V1beta2ListTidbsResponse](docs/V1beta2ListTidbsResponse.md)
 - [V1beta2PublicConnectionSetting](docs/V1beta2PublicConnectionSetting.md)
 - [V1beta2PublicConnectionSettingIpAccessList](docs/V1beta2PublicConnectionSettingIpAccessList.md)
 - [V1beta2ReservedCapacity](docs/V1beta2ReservedCapacity.md)
 - [V1beta2RestoreMode](docs/V1beta2RestoreMode.md)
 - [V1beta2RestoreState](docs/V1beta2RestoreState.md)
 - [V1beta2RestoreTidbRequest](docs/V1beta2RestoreTidbRequest.md)
 - [V1beta2RestoreTidbResponse](docs/V1beta2RestoreTidbResponse.md)
 - [V1beta2TidbHighAvailabilityType](docs/V1beta2TidbHighAvailabilityType.md)
 - [V1beta2VerifyCmekAccessIamPrincipalResponse](docs/V1beta2VerifyCmekAccessIamPrincipalResponse.md)


## Documentation For Authorization

Endpoints do not require authorization.


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author



