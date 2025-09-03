# Go API client for dedicated

*TiDB Cloud API is in beta.*

This API manages [TiDB Cloud Dedicated](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-dedicated) clusters. For TiDB Cloud Starter or TiDB Cloud Essential clusters, use the [TiDB Cloud Starter and Essential API](). For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/tidbcloud/api-overview/).

# Overview

The TiDB Cloud API is a [REST interface](https://en.wikipedia.org/wiki/Representational_state_transfer) that provides you with programmatic access to manage clusters and related resources within TiDB Cloud.

The API has the following features:

- **JSON entities.** All entities are expressed in JSON.
- **HTTPS-only.** You can only access the API via HTTPS, ensuring all the data sent over the network is encrypted with TLS.
- **Key-based access and digest authentication.** Before you access TiDB Cloud API, you must generate an API key. All requests are authenticated through [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication), ensuring the API key is never sent over the network.

# Get Started

This guide helps you make your first API call to TiDB Cloud API. You'll learn how to authenticate a request, build a request, and interpret the response.

## Prerequisites

To complete this guide, you need to perform the following tasks:

- Create a [TiDB Cloud account](https://tidbcloud.com/free-trial)
- Install [curl](https://curl.se/)

## Step 1. Create an API key

To create an API key, log in to your TiDB Cloud console. Navigate to the [**API Keys**](https://tidbcloud.com/org-settings/api-keys) page of your organization, and create an API key.

An API key contains a public key and a private key. Copy and save them in a secure location. You will need to use the API key later in this guide.

For more details about creating API key, refer to [API Key Management](#section/Authentication/API-Key-Management).

## Step 2. Make your first API call

### Build an API call

TiDB Cloud API call consists of the following components:

- **A host**. The host for TiDB Cloud API is <https://dedicated.tidbapi.com>.
- **An API Key**. The public key and the private key are required for authentication.
- **A request**. When submitting data to a resource via `POST`, `PATCH`, or `PUT`, you must submit your payload in JSON.

In this guide, you call the [List clusters](#tag/Cluster/operation/ClusterService_ListClusters) endpoint. For the detailed description of the endpoint, see the [API reference](#tag/Cluster/operation/ClusterService_ListClusters).

### Call an API endpoint

To get all clusters in your organization, run the following command in your terminal. Remember to change `YOUR_PUBLIC_KEY` to your public key and `YOUR_PRIVATE_KEY` to your private key.

```shell
curl --digest \\
 --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\
 --request GET \\
 --url 'https://dedicated.tidbapi.com/v1beta1/clusters'
```

## Step 3. Check the response

After making the API call, if the status code in response is `200` and you see details about all clusters in your organization, your request is successful.

# Authentication

The TiDB Cloud API uses [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication). It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).

## API key overview

- The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation.
- The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles).
- You must provide the correct API key in every request. Otherwise, the TiDB Cloud responds with a `401` error.

## API key management

### Create an API key

Only the **owner** of an organization can create an API key.

To create an API key in an organization, perform the following steps:

1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner.
2. In the left navigation pane, click **Organization Settings** > **API Keys**.
3. On the **API Keys** page, click **Create API Key**.
4. Enter a description for your API key. The role of the API key is always `Organization Owner` currently.
5. Click **Next**. Copy and save the public key and the private key.
6. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again.
7. Click **Done**.

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
3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Edit**.
4. You can update the API key description.
5. Click **Update**.

### Delete an API key

Only the **owner** of an organization can delete an API key.

To delete an API key in an organization, perform the following steps:

1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner.
2. In the left navigation pane, click **Organization Settings** > **API Keys**.
3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**.
4. Click **I understand, delete it.**

# Rate Limiting

The TiDB Cloud API allows up to 100 requests per minute per API key. If you exceed the rate limit, the API returns a `429` error. For more quota, you can [submit a request](https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519) to contact our support team.

Each API request returns the following headers about the limit.

- `X-Ratelimit-Limit-Minute`: The number of requests allowed per minute. It is 100 currently.
- `X-Ratelimit-Remaining-Minute`: The number of remaining requests in the current minute. When it reaches `0`, the API returns a `429` error and indicates that you exceed the rate limit.
- `X-Ratelimit-Reset`: The time in seconds at which the current rate limit resets.

If you exceed the rate limit, an error response returns like this.

```
> HTTP/2 429
> date: Fri, 22 Jul 2022 05:28:37 GMT
> content-type: application/json
> content-length: 66
> x-ratelimit-reset: 23
> x-ratelimit-remaining-minute: 0
> x-ratelimit-limit-minute: 100
> x-kong-response-latency: 2
> server: kong/2.8.1

> {\"details\":[],\"code\":49900007,\"message\":\"The request exceeded the limit of 100 times per apikey per minute. For more quota, please contact us: https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519\"}
```

# API Changelog

This changelog lists all changes to the TiDB Cloud API.

<!-- In reverse chronological order -->

## 20250812

- Initial release of the TiDB Cloud Dedicated API, including the following resources and endpoints:
 * Cluster
  * [List clusters](#tag/Cluster/operation/ClusterService_ListClusters)
  * [Create a cluster](#tag/Cluster/operation/ClusterService_CreateCluster)
  * [Get a cluster](#tag/Cluster/operation/ClusterService_GetCluster)
  * [Delete a cluster](#tag/Cluster/operation/ClusterService_DeleteCluster)
  * [Update a cluster](#tag/Cluster/operation/ClusterService_UpdateCluster)
  * [Pause a cluster](#tag/Cluster/operation/ClusterService_PauseCluster)
  * [Resume a cluster](#tag/Cluster/operation/ClusterService_ResumeCluster)
  * [Reset the root password of a cluster](#tag/Cluster/operation/ClusterService_ResetRootPassword)
  * [List node quotas for your organization](#tag/Cluster/operation/ClusterService_ShowNodeQuota)
  * [Get log redaction policy](#tag/Cluster/operation/ClusterService_GetLogRedactionPolicy)
 * Region
  * [List regions](#tag/Region/operation/RegionService_ListRegions)
  * [Get a region](#tag/Region/operation/RegionService_GetRegion)
  * [List cloud providers](#tag/Region/operation/RegionService_ShowCloudProviders)
  * [List node specs](#tag/Region/operation/RegionService_ListNodeSpecs)
  * [Get a node spec](#tag/Region/operation/RegionService_GetNodeSpec)
 * Private Endpoint Connection
  * [Get private link service for a TiDB node group](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateLinkService)
  * [Create a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_CreatePrivateEndpointConnection)
  * [List private endpoint connections](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_ListPrivateEndpointConnections)
  * [Get a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateEndpointConnection)
  * [Delete a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_DeletePrivateEndpointConnection)
 * Import
  * [List import tasks](#tag/Import/operation/ListImports)
  * [Create an import task](#tag/Import/operation/CreateImport)
  * [Get an import task](#tag/Import/operation/GetImport)
  * [Cancel an import task](#tag/Import/operation/CancelImport)


## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: v1beta1
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
import dedicated "github.com/GIT_USER_ID/GIT_REPO_ID"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```go
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `dedicated.ContextServerIndex` of type `int`.

```go
ctx := context.WithValue(context.Background(), dedicated.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `dedicated.ContextServerVariables` of type `map[string]string`.

```go
ctx := context.WithValue(context.Background(), dedicated.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `dedicated.ContextOperationServerIndices` and `dedicated.ContextOperationServerVariables` context maps.

```go
ctx := context.WithValue(context.Background(), dedicated.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), dedicated.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *https://dedicated.tidbapi.com/v1beta1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*ClusterServiceAPI* | [**ClusterServiceCreateCluster**](docs/ClusterServiceAPI.md#clusterservicecreatecluster) | **Post** /clusters | Create a cluster
*ClusterServiceAPI* | [**ClusterServiceDeleteCluster**](docs/ClusterServiceAPI.md#clusterservicedeletecluster) | **Delete** /clusters/{clusterId} | Delete a cluster
*ClusterServiceAPI* | [**ClusterServiceGetCluster**](docs/ClusterServiceAPI.md#clusterservicegetcluster) | **Get** /clusters/{clusterId} | Get a cluster
*ClusterServiceAPI* | [**ClusterServiceGetLogRedactionPolicy**](docs/ClusterServiceAPI.md#clusterservicegetlogredactionpolicy) | **Get** /clusters/{clusterId}/logRedactionPolicy | Get log redaction policy
*ClusterServiceAPI* | [**ClusterServiceGetNodeInstance**](docs/ClusterServiceAPI.md#clusterservicegetnodeinstance) | **Get** /clusters/{clusterId}/nodeInstances/{instanceId} | Get a node instance
*ClusterServiceAPI* | [**ClusterServiceListClusters**](docs/ClusterServiceAPI.md#clusterservicelistclusters) | **Get** /clusters | List clusters
*ClusterServiceAPI* | [**ClusterServiceListNodeInstances**](docs/ClusterServiceAPI.md#clusterservicelistnodeinstances) | **Get** /clusters/{clusterId}/nodeInstances | List node instances in a cluster
*ClusterServiceAPI* | [**ClusterServicePauseCluster**](docs/ClusterServiceAPI.md#clusterservicepausecluster) | **Post** /clusters/{clusterId}:pauseCluster | Pause a cluster
*ClusterServiceAPI* | [**ClusterServiceResetRootPassword**](docs/ClusterServiceAPI.md#clusterserviceresetrootpassword) | **Post** /clusters/{clusterId}:resetRootPassword | Reset the root password of a cluster
*ClusterServiceAPI* | [**ClusterServiceResumeCluster**](docs/ClusterServiceAPI.md#clusterserviceresumecluster) | **Post** /clusters/{clusterId}:resumeCluster | Resume a cluster
*ClusterServiceAPI* | [**ClusterServiceShowNodeQuota**](docs/ClusterServiceAPI.md#clusterserviceshownodequota) | **Get** /clusters:showNodeQuota | List node quotas for your organization
*ClusterServiceAPI* | [**ClusterServiceUpdateCluster**](docs/ClusterServiceAPI.md#clusterserviceupdatecluster) | **Patch** /clusters/{cluster.clusterId} | Update a cluster
*ClusterServiceAPI* | [**ClusterServiceUpdateLogRedactionPolicy**](docs/ClusterServiceAPI.md#clusterserviceupdatelogredactionpolicy) | **Patch** /clusters/{logRedactionPolicy.clusterId}/logRedactionPolicy | Update log redaction policy
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceCreateAuditLogFilterRule**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogservicecreateauditlogfilterrule) | **Post** /clusters/{auditLogFilterRule.clusterId}/auditLogFilterRules | Create an audit log filter rule
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceDeleteAuditLogFilterRule**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogservicedeleteauditlogfilterrule) | **Delete** /clusters/{clusterId}/auditLogFilterRules/{auditLogFilterRuleId} | Delete an audit log filter rule
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceGenerateAuditLogFileDownloadAddress**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogservicegenerateauditlogfiledownloadaddress) | **Post** /clusters/{clusterId}/auditLogConfig:generateAuditLogFileDownloadAddress | Generate the download address for an audit log file, the address have an 15 minutes expiration time
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceGetAuditLogConfig**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogservicegetauditlogconfig) | **Get** /clusters/{clusterId}/auditLogConfig | Get the audit log config of a cluster
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceGetAuditLogFilterRule**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogservicegetauditlogfilterrule) | **Get** /clusters/{clusterId}/auditLogFilterRules/{auditLogFilterRuleId} | Get an audit log filter rule
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceListAuditLogFilterRules**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogservicelistauditlogfilterrules) | **Get** /clusters/{clusterId}/auditLogFilterRules | List audit log filter rules
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceQueryAuditLogFiles**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogservicequeryauditlogfiles) | **Get** /clusters/{clusterId}/auditLogConfig:queryAuditLogFiles | Query audit log files
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceReplaceAuditLogFilterRule**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogservicereplaceauditlogfilterrule) | **Post** /clusters/{auditLogFilterRule.clusterId}/auditLogFilterRules/{auditLogFilterRule.auditLogFilterRuleId}:replaceAuditLogFilterRule | Replace an audit log filter rule. All fields of the rule will be replaced with the provided values.
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceShowObjectStorageAccessIamPrincipal**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogserviceshowobjectstorageaccessiamprincipal) | **Get** /clusters/{clusterId}/auditLogConfig:showObjectStorageAccessIamPrincipal | Show IAM principal of TiDB Cloud for accessing customer&#39;s object storage
*DatabaseAuditLogServiceAPI* | [**DatabaseAuditLogServiceUpdateAuditLogConfig**](docs/DatabaseAuditLogServiceAPI.md#databaseauditlogserviceupdateauditlogconfig) | **Patch** /clusters/{auditLogConfig.clusterId}/auditLogConfig | Update the audit log config of a cluster
*MaintenanceServiceAPI* | [**MaintenanceServiceGetMaintenanceTask**](docs/MaintenanceServiceAPI.md#maintenanceservicegetmaintenancetask) | **Get** /maintenanceTasks/{maintenanceTaskId} | Get a maintenance task
*MaintenanceServiceAPI* | [**MaintenanceServiceGetMaintenanceWindow**](docs/MaintenanceServiceAPI.md#maintenanceservicegetmaintenancewindow) | **Get** /maintenanceWindows/{maintenanceWindowId} | Get a maintenance window
*MaintenanceServiceAPI* | [**MaintenanceServiceListMaintenanceTasks**](docs/MaintenanceServiceAPI.md#maintenanceservicelistmaintenancetasks) | **Get** /maintenanceTasks | List maintenance tasks
*MaintenanceServiceAPI* | [**MaintenanceServiceListMaintenanceWindows**](docs/MaintenanceServiceAPI.md#maintenanceservicelistmaintenancewindows) | **Get** /maintenanceWindows | List maintenance windows
*MaintenanceServiceAPI* | [**MaintenanceServiceUpdateMaintenanceTask**](docs/MaintenanceServiceAPI.md#maintenanceserviceupdatemaintenancetask) | **Patch** /maintenanceTasks/{maintenanceTask.maintenanceTaskId} | Update a maintenance task
*MaintenanceServiceAPI* | [**MaintenanceServiceUpdateMaintenanceWindow**](docs/MaintenanceServiceAPI.md#maintenanceserviceupdatemaintenancewindow) | **Patch** /maintenanceWindows/{maintenanceWindow.maintenanceWindowId} | Update a maintenance window
*NetworkContainerServiceAPI* | [**NetworkContainerServiceCreateNetworkContainer**](docs/NetworkContainerServiceAPI.md#networkcontainerservicecreatenetworkcontainer) | **Post** /networkContainers | Create a network container
*NetworkContainerServiceAPI* | [**NetworkContainerServiceCreateVpcPeering**](docs/NetworkContainerServiceAPI.md#networkcontainerservicecreatevpcpeering) | **Post** /vpcPeerings | Create a VPC peering
*NetworkContainerServiceAPI* | [**NetworkContainerServiceDeleteNetworkContainer**](docs/NetworkContainerServiceAPI.md#networkcontainerservicedeletenetworkcontainer) | **Delete** /networkContainers/{networkContainerId} | Delete a network container
*NetworkContainerServiceAPI* | [**NetworkContainerServiceDeleteVpcPeering**](docs/NetworkContainerServiceAPI.md#networkcontainerservicedeletevpcpeering) | **Delete** /vpcPeerings/{vpcPeeringId} | Delete a VPC peering
*NetworkContainerServiceAPI* | [**NetworkContainerServiceGetNetworkContainer**](docs/NetworkContainerServiceAPI.md#networkcontainerservicegetnetworkcontainer) | **Get** /networkContainers/{networkContainerId} | Get a network container
*NetworkContainerServiceAPI* | [**NetworkContainerServiceGetVpcPeering**](docs/NetworkContainerServiceAPI.md#networkcontainerservicegetvpcpeering) | **Get** /vpcPeerings/{vpcPeeringId} | Get a VPC peering
*NetworkContainerServiceAPI* | [**NetworkContainerServiceListNetworkContainers**](docs/NetworkContainerServiceAPI.md#networkcontainerservicelistnetworkcontainers) | **Get** /networkContainers | List network containers
*NetworkContainerServiceAPI* | [**NetworkContainerServiceListVpcPeerings**](docs/NetworkContainerServiceAPI.md#networkcontainerservicelistvpcpeerings) | **Get** /vpcPeerings | List VPC peerings
*PrivateEndpointConnectionServiceAPI* | [**PrivateEndpointConnectionServiceCreatePrivateEndpointConnection**](docs/PrivateEndpointConnectionServiceAPI.md#privateendpointconnectionservicecreateprivateendpointconnection) | **Post** /clusters/{clusterId}/tidbNodeGroups/{privateEndpointConnection.tidbNodeGroupId}/privateEndpointConnections | Create a private endpoint connection
*PrivateEndpointConnectionServiceAPI* | [**PrivateEndpointConnectionServiceDeletePrivateEndpointConnection**](docs/PrivateEndpointConnectionServiceAPI.md#privateendpointconnectionservicedeleteprivateendpointconnection) | **Delete** /clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/privateEndpointConnections/{privateEndpointConnectionId} | Delete a private endpoint connection
*PrivateEndpointConnectionServiceAPI* | [**PrivateEndpointConnectionServiceGetPrivateEndpointConnection**](docs/PrivateEndpointConnectionServiceAPI.md#privateendpointconnectionservicegetprivateendpointconnection) | **Get** /clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/privateEndpointConnections/{privateEndpointConnectionId} | Get a private endpoint connection
*PrivateEndpointConnectionServiceAPI* | [**PrivateEndpointConnectionServiceGetPrivateLinkService**](docs/PrivateEndpointConnectionServiceAPI.md#privateendpointconnectionservicegetprivatelinkservice) | **Get** /clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/privateLinkService | Get private link service for a TiDB node group
*PrivateEndpointConnectionServiceAPI* | [**PrivateEndpointConnectionServiceListPrivateEndpointConnections**](docs/PrivateEndpointConnectionServiceAPI.md#privateendpointconnectionservicelistprivateendpointconnections) | **Get** /clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/privateEndpointConnections | List private endpoint connections
*RegionServiceAPI* | [**RegionServiceGetNodeSpec**](docs/RegionServiceAPI.md#regionservicegetnodespec) | **Get** /regions/{regionId}/componentTypes/{componentType}/nodeSpecs/{nodeSpecKey} | Get a node spec
*RegionServiceAPI* | [**RegionServiceGetRegion**](docs/RegionServiceAPI.md#regionservicegetregion) | **Get** /regions/{regionId} | Get a region
*RegionServiceAPI* | [**RegionServiceListNodeSpecs**](docs/RegionServiceAPI.md#regionservicelistnodespecs) | **Get** /regions/{regionId}/nodeSpecs | List node specs
*RegionServiceAPI* | [**RegionServiceListRegions**](docs/RegionServiceAPI.md#regionservicelistregions) | **Get** /regions | List regions
*RegionServiceAPI* | [**RegionServiceShowCloudProviders**](docs/RegionServiceAPI.md#regionserviceshowcloudproviders) | **Get** /regions:showCloudProviders | List cloud providers
*TidbNodeGroupServiceAPI* | [**TidbNodeGroupServiceCreateTidbNodeGroup**](docs/TidbNodeGroupServiceAPI.md#tidbnodegroupservicecreatetidbnodegroup) | **Post** /clusters/{tidbNodeGroup.clusterId}/tidbNodeGroups | Create a TiDB Node Group
*TidbNodeGroupServiceAPI* | [**TidbNodeGroupServiceDeleteTidbNodeGroup**](docs/TidbNodeGroupServiceAPI.md#tidbnodegroupservicedeletetidbnodegroup) | **Delete** /clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId} | Delete a TiDB Node Group
*TidbNodeGroupServiceAPI* | [**TidbNodeGroupServiceGetPublicEndpointSetting**](docs/TidbNodeGroupServiceAPI.md#tidbnodegroupservicegetpublicendpointsetting) | **Get** /clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/publicEndpointSetting | Get the public endpoint setting of a TiDB Node Group
*TidbNodeGroupServiceAPI* | [**TidbNodeGroupServiceGetTidbNodeGroup**](docs/TidbNodeGroupServiceAPI.md#tidbnodegroupservicegettidbnodegroup) | **Get** /clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId} | Get a TiDB Node Group
*TidbNodeGroupServiceAPI* | [**TidbNodeGroupServiceListTidbNodeGroups**](docs/TidbNodeGroupServiceAPI.md#tidbnodegroupservicelisttidbnodegroups) | **Get** /clusters/{clusterId}/tidbNodeGroups | List TiDB Node Groups
*TidbNodeGroupServiceAPI* | [**TidbNodeGroupServiceUpdatePublicEndpointSetting**](docs/TidbNodeGroupServiceAPI.md#tidbnodegroupserviceupdatepublicendpointsetting) | **Patch** /clusters/{clusterId}/tidbNodeGroups/{publicEndpointSetting.tidbNodeGroupId}/publicEndpointSetting | Update the public endpoint setting of a TiDB Node Group
*TidbNodeGroupServiceAPI* | [**TidbNodeGroupServiceUpdateTidbNodeGroup**](docs/TidbNodeGroupServiceAPI.md#tidbnodegroupserviceupdatetidbnodegroup) | **Patch** /clusters/{tidbNodeGroup.clusterId}/tidbNodeGroups/{tidbNodeGroup.tidbNodeGroupId} | Update a TiDB Node Group


## Documentation For Models

 - [AuditLogConfigBucketManager](docs/AuditLogConfigBucketManager.md)
 - [ClusterNodeChangingProgress](docs/ClusterNodeChangingProgress.md)
 - [ClusterServiceListClustersClusterStatesParameterInner](docs/ClusterServiceListClustersClusterStatesParameterInner.md)
 - [ClusterServiceListNodeInstancesComponentTypeParameter](docs/ClusterServiceListNodeInstancesComponentTypeParameter.md)
 - [ClusterServiceUpdateLogRedactionPolicyRequest](docs/ClusterServiceUpdateLogRedactionPolicyRequest.md)
 - [ClusterStorageNodeSettingStorageType](docs/ClusterStorageNodeSettingStorageType.md)
 - [Commonv1beta1ClusterState](docs/Commonv1beta1ClusterState.md)
 - [Commonv1beta1Region](docs/Commonv1beta1Region.md)
 - [Commonv1beta1ServicePlan](docs/Commonv1beta1ServicePlan.md)
 - [DatabaseAuditLogServiceCreateAuditLogFilterRuleRequest](docs/DatabaseAuditLogServiceCreateAuditLogFilterRuleRequest.md)
 - [DatabaseAuditLogServiceUpdateAuditLogConfigRequest](docs/DatabaseAuditLogServiceUpdateAuditLogConfigRequest.md)
 - [Dedicatedv1beta1AuditLogConfig](docs/Dedicatedv1beta1AuditLogConfig.md)
 - [Dedicatedv1beta1AuditLogFilterRule](docs/Dedicatedv1beta1AuditLogFilterRule.md)
 - [Dedicatedv1beta1AuditLogFilterRuleEventType](docs/Dedicatedv1beta1AuditLogFilterRuleEventType.md)
 - [Dedicatedv1beta1ClusterPausePlan](docs/Dedicatedv1beta1ClusterPausePlan.md)
 - [Dedicatedv1beta1ClusterPausePlanType](docs/Dedicatedv1beta1ClusterPausePlanType.md)
 - [Dedicatedv1beta1ComponentType](docs/Dedicatedv1beta1ComponentType.md)
 - [Dedicatedv1beta1GenerateAuditLogFileDownloadAddressResponse](docs/Dedicatedv1beta1GenerateAuditLogFileDownloadAddressResponse.md)
 - [Dedicatedv1beta1ListAuditLogFilterRulesResponse](docs/Dedicatedv1beta1ListAuditLogFilterRulesResponse.md)
 - [Dedicatedv1beta1ListMaintenanceTasksResponse](docs/Dedicatedv1beta1ListMaintenanceTasksResponse.md)
 - [Dedicatedv1beta1ListMaintenanceWindowsResponse](docs/Dedicatedv1beta1ListMaintenanceWindowsResponse.md)
 - [Dedicatedv1beta1ListPrivateEndpointConnectionsResponse](docs/Dedicatedv1beta1ListPrivateEndpointConnectionsResponse.md)
 - [Dedicatedv1beta1ListTidbNodeGroupsResponse](docs/Dedicatedv1beta1ListTidbNodeGroupsResponse.md)
 - [Dedicatedv1beta1ListVpcPeeringsResponse](docs/Dedicatedv1beta1ListVpcPeeringsResponse.md)
 - [Dedicatedv1beta1LogRedactionPolicy](docs/Dedicatedv1beta1LogRedactionPolicy.md)
 - [Dedicatedv1beta1MaintenanceTask](docs/Dedicatedv1beta1MaintenanceTask.md)
 - [Dedicatedv1beta1MaintenanceTaskState](docs/Dedicatedv1beta1MaintenanceTaskState.md)
 - [Dedicatedv1beta1MaintenanceWindow](docs/Dedicatedv1beta1MaintenanceWindow.md)
 - [Dedicatedv1beta1PrivateEndpointConnection](docs/Dedicatedv1beta1PrivateEndpointConnection.md)
 - [Dedicatedv1beta1PrivateEndpointConnectionEndpointState](docs/Dedicatedv1beta1PrivateEndpointConnectionEndpointState.md)
 - [Dedicatedv1beta1PrivateLinkService](docs/Dedicatedv1beta1PrivateLinkService.md)
 - [Dedicatedv1beta1PrivateLinkServiceState](docs/Dedicatedv1beta1PrivateLinkServiceState.md)
 - [Dedicatedv1beta1QueryAuditLogFilesResponse](docs/Dedicatedv1beta1QueryAuditLogFilesResponse.md)
 - [Dedicatedv1beta1QueryAuditLogFilesResponseAuditLogFile](docs/Dedicatedv1beta1QueryAuditLogFilesResponseAuditLogFile.md)
 - [Dedicatedv1beta1TidbNodeGroup](docs/Dedicatedv1beta1TidbNodeGroup.md)
 - [Dedicatedv1beta1TidbNodeGroupEndpoint](docs/Dedicatedv1beta1TidbNodeGroupEndpoint.md)
 - [Dedicatedv1beta1TidbNodeGroupState](docs/Dedicatedv1beta1TidbNodeGroupState.md)
 - [Dedicatedv1beta1TidbNodeGroupTiProxySetting](docs/Dedicatedv1beta1TidbNodeGroupTiProxySetting.md)
 - [Dedicatedv1beta1VpcPeering](docs/Dedicatedv1beta1VpcPeering.md)
 - [Dedicatedv1beta1VpcPeeringState](docs/Dedicatedv1beta1VpcPeeringState.md)
 - [GooglerpcStatus](docs/GooglerpcStatus.md)
 - [MaintenanceServiceUpdateMaintenanceTaskRequest](docs/MaintenanceServiceUpdateMaintenanceTaskRequest.md)
 - [MaintenanceServiceUpdateMaintenanceWindowRequest](docs/MaintenanceServiceUpdateMaintenanceWindowRequest.md)
 - [PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest](docs/PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest.md)
 - [PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter](docs/PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter.md)
 - [ProtobufAny](docs/ProtobufAny.md)
 - [RequiredTheAuditLogFilterRuleToUpdate](docs/RequiredTheAuditLogFilterRuleToUpdate.md)
 - [ShowNodeQuotaResponseComponentQuota](docs/ShowNodeQuotaResponseComponentQuota.md)
 - [TheUpdatedClusterConfiguration](docs/TheUpdatedClusterConfiguration.md)
 - [TidbCloudOpenApidedicatedv1beta1Cluster](docs/TidbCloudOpenApidedicatedv1beta1Cluster.md)
 - [TidbCloudOpenApidedicatedv1beta1ListClustersResponse](docs/TidbCloudOpenApidedicatedv1beta1ListClustersResponse.md)
 - [TidbCloudOpenApidedicatedv1beta1ListRegionsResponse](docs/TidbCloudOpenApidedicatedv1beta1ListRegionsResponse.md)
 - [TidbNodeGroupServiceCreateTidbNodeGroupRequest](docs/TidbNodeGroupServiceCreateTidbNodeGroupRequest.md)
 - [TidbNodeGroupServiceUpdatePublicEndpointSettingRequest](docs/TidbNodeGroupServiceUpdatePublicEndpointSettingRequest.md)
 - [TidbNodeGroupServiceUpdateTidbNodeGroupRequest](docs/TidbNodeGroupServiceUpdateTidbNodeGroupRequest.md)
 - [UpdateClusterRequestTidbNodeSettingTidbNodeGroup](docs/UpdateClusterRequestTidbNodeSettingTidbNodeGroup.md)
 - [V1beta1AuditLogConfigBucketWriteCheck](docs/V1beta1AuditLogConfigBucketWriteCheck.md)
 - [V1beta1AuditLogConfigFormat](docs/V1beta1AuditLogConfigFormat.md)
 - [V1beta1AuditLogConfigRotationPolicy](docs/V1beta1AuditLogConfigRotationPolicy.md)
 - [V1beta1ClusterServiceResetRootPasswordBody](docs/V1beta1ClusterServiceResetRootPasswordBody.md)
 - [V1beta1ClusterStorageNodeSetting](docs/V1beta1ClusterStorageNodeSetting.md)
 - [V1beta1ClusterTidbNodeSetting](docs/V1beta1ClusterTidbNodeSetting.md)
 - [V1beta1ListNetworkContainersResponse](docs/V1beta1ListNetworkContainersResponse.md)
 - [V1beta1ListNodeInstancesResponse](docs/V1beta1ListNodeInstancesResponse.md)
 - [V1beta1ListNodeSpecsResponse](docs/V1beta1ListNodeSpecsResponse.md)
 - [V1beta1NetworkContainer](docs/V1beta1NetworkContainer.md)
 - [V1beta1NetworkContainerState](docs/V1beta1NetworkContainerState.md)
 - [V1beta1NodeInstance](docs/V1beta1NodeInstance.md)
 - [V1beta1NodeInstanceState](docs/V1beta1NodeInstanceState.md)
 - [V1beta1NodeSpec](docs/V1beta1NodeSpec.md)
 - [V1beta1PauseClusterResponse](docs/V1beta1PauseClusterResponse.md)
 - [V1beta1PublicEndpointSetting](docs/V1beta1PublicEndpointSetting.md)
 - [V1beta1PublicEndpointSettingIpAccessList](docs/V1beta1PublicEndpointSettingIpAccessList.md)
 - [V1beta1RegionCloudProvider](docs/V1beta1RegionCloudProvider.md)
 - [V1beta1ReplaceAuditLogFilterRuleRequestAuditLogFilterRule](docs/V1beta1ReplaceAuditLogFilterRuleRequestAuditLogFilterRule.md)
 - [V1beta1ReplaceAuditLogFilterRuleResponse](docs/V1beta1ReplaceAuditLogFilterRuleResponse.md)
 - [V1beta1ResumeClusterResponse](docs/V1beta1ResumeClusterResponse.md)
 - [V1beta1ShowCloudProvidersResponse](docs/V1beta1ShowCloudProvidersResponse.md)
 - [V1beta1ShowNodeQuotaResponse](docs/V1beta1ShowNodeQuotaResponse.md)
 - [V1beta1ShowObjectStorageAccessIamPrincipalResponse](docs/V1beta1ShowObjectStorageAccessIamPrincipalResponse.md)
 - [V1beta1TidbNodeGroupEndpointConnectionType](docs/V1beta1TidbNodeGroupEndpointConnectionType.md)
 - [V1beta1UpdateAuditLogConfigRequestAuditLogConfig](docs/V1beta1UpdateAuditLogConfigRequestAuditLogConfig.md)
 - [V1beta1UpdateClusterRequestCluster](docs/V1beta1UpdateClusterRequestCluster.md)
 - [V1beta1UpdateClusterRequestStorageNodeSetting](docs/V1beta1UpdateClusterRequestStorageNodeSetting.md)
 - [V1beta1UpdateClusterRequestTidbNodeSetting](docs/V1beta1UpdateClusterRequestTidbNodeSetting.md)
 - [V1beta1UpdateMaintenanceTaskRequestMaintenanceTask](docs/V1beta1UpdateMaintenanceTaskRequestMaintenanceTask.md)
 - [V1beta1UpdateMaintenanceWindowRequestMaintenanceWindow](docs/V1beta1UpdateMaintenanceWindowRequestMaintenanceWindow.md)
 - [V1beta1UpdateTidbNodeGroupRequestTidbNodeGroup](docs/V1beta1UpdateTidbNodeGroupRequestTidbNodeGroup.md)


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



