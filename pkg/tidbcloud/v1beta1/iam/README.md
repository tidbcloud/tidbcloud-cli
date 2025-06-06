# Go API client for iam

This is account open api.

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
import iam "github.com/GIT_USER_ID/GIT_REPO_ID"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```go
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `iam.ContextServerIndex` of type `int`.

```go
ctx := context.WithValue(context.Background(), iam.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `iam.ContextServerVariables` of type `map[string]string`.

```go
ctx := context.WithValue(context.Background(), iam.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `iam.ContextOperationServerIndices` and `iam.ContextOperationServerVariables` context maps.

```go
ctx := context.WithValue(context.Background(), iam.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), iam.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *https://iam.tidbapi.com*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AccountAPI* | [**CustomerSignupUrlPost**](docs/AccountAPI.md#customersignupurlpost) | **Post** /customerSignupUrl | Create a new signup URL for an MSP customer
*AccountAPI* | [**MspCustomersCustomerOrgIdGet**](docs/AccountAPI.md#mspcustomerscustomerorgidget) | **Get** /mspCustomers/{customerOrgId} | Retrieve a single MSP customer
*AccountAPI* | [**MspCustomersGet**](docs/AccountAPI.md#mspcustomersget) | **Get** /mspCustomers | Get a list of MSP customers
*AccountAPI* | [**V1beta1ClustersClusterIdDbuserGet**](docs/AccountAPI.md#v1beta1clustersclusteriddbuserget) | **Get** /v1beta1/clusters/{clusterId}/dbuser | get one dbuser
*AccountAPI* | [**V1beta1ClustersClusterIdSqlUsersGet**](docs/AccountAPI.md#v1beta1clustersclusteridsqlusersget) | **Get** /v1beta1/clusters/{clusterId}/sqlUsers | Get all sql users
*AccountAPI* | [**V1beta1ClustersClusterIdSqlUsersPost**](docs/AccountAPI.md#v1beta1clustersclusteridsqluserspost) | **Post** /v1beta1/clusters/{clusterId}/sqlUsers | Create one sql user
*AccountAPI* | [**V1beta1ClustersClusterIdSqlUsersUserNameDelete**](docs/AccountAPI.md#v1beta1clustersclusteridsqlusersusernamedelete) | **Delete** /v1beta1/clusters/{clusterId}/sqlUsers/{userName} | Delete one sql user
*AccountAPI* | [**V1beta1ClustersClusterIdSqlUsersUserNameGet**](docs/AccountAPI.md#v1beta1clustersclusteridsqlusersusernameget) | **Get** /v1beta1/clusters/{clusterId}/sqlUsers/{userName} | Query sql user
*AccountAPI* | [**V1beta1ClustersClusterIdSqlUsersUserNamePatch**](docs/AccountAPI.md#v1beta1clustersclusteridsqlusersusernamepatch) | **Patch** /v1beta1/clusters/{clusterId}/sqlUsers/{userName} | Update one sql user
*AccountAPI* | [**V1beta1ProjectsGet**](docs/AccountAPI.md#v1beta1projectsget) | **Get** /v1beta1/projects | Get  list of org projects


## Documentation For Models

 - [ApiBasicResp](docs/ApiBasicResp.md)
 - [ApiCreateSqlUserReq](docs/ApiCreateSqlUserReq.md)
 - [ApiGetDbuserRsp](docs/ApiGetDbuserRsp.md)
 - [ApiListProjectsRsp](docs/ApiListProjectsRsp.md)
 - [ApiListSqlUsersRsp](docs/ApiListSqlUsersRsp.md)
 - [ApiOpenApiCreateMspCustomerSignupUrlReq](docs/ApiOpenApiCreateMspCustomerSignupUrlReq.md)
 - [ApiOpenApiError](docs/ApiOpenApiError.md)
 - [ApiOpenApiListMspCustomerRsp](docs/ApiOpenApiListMspCustomerRsp.md)
 - [ApiOpenApiMspCustomer](docs/ApiOpenApiMspCustomer.md)
 - [ApiOpenApiMspCustomerSignupUrl](docs/ApiOpenApiMspCustomerSignupUrl.md)
 - [ApiProject](docs/ApiProject.md)
 - [ApiSqlUser](docs/ApiSqlUser.md)
 - [ApiUpdateSqlUserReq](docs/ApiUpdateSqlUserReq.md)


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



