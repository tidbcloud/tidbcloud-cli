# Go API client for auditlog

TiDB Cloud Serverless Database Audit Logging Open API

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
import auditlog "github.com/GIT_USER_ID/GIT_REPO_ID"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```go
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `auditlog.ContextServerIndex` of type `int`.

```go
ctx := context.WithValue(context.Background(), auditlog.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `auditlog.ContextServerVariables` of type `map[string]string`.

```go
ctx := context.WithValue(context.Background(), auditlog.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `auditlog.ContextOperationServerIndices` and `auditlog.ContextOperationServerVariables` context maps.

```go
ctx := context.WithValue(context.Background(), auditlog.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), auditlog.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *https://serverless.tidbapi.com*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AuditLogServiceAPI* | [**AuditLogServiceCreateAuditLogFilterRule**](docs/AuditLogServiceAPI.md#auditlogservicecreateauditlogfilterrule) | **Post** /v1beta1/clusters/{clusterId}/auditlogs/filterRules | Create audit log filter rule.
*AuditLogServiceAPI* | [**AuditLogServiceDeleteAuditLogFilterRule**](docs/AuditLogServiceAPI.md#auditlogservicedeleteauditlogfilterrule) | **Delete** /v1beta1/clusters/{clusterId}/auditlogs/filterRules/{name} | Delete audit log filter rule.
*AuditLogServiceAPI* | [**AuditLogServiceDownloadAuditLogs**](docs/AuditLogServiceAPI.md#auditlogservicedownloadauditlogs) | **Post** /v1beta1/clusters/{clusterId}/auditlogs:download | Generate audit logs download url
*AuditLogServiceAPI* | [**AuditLogServiceGetAuditLogFilterRule**](docs/AuditLogServiceAPI.md#auditlogservicegetauditlogfilterrule) | **Get** /v1beta1/clusters/{clusterId}/auditlogs/filterRules/{name} | Get audit log filter rule.
*AuditLogServiceAPI* | [**AuditLogServiceListAuditLogFilterRules**](docs/AuditLogServiceAPI.md#auditlogservicelistauditlogfilterrules) | **Get** /v1beta1/clusters/{clusterId}/auditlogs/filterRules | List audit log filter rules.
*AuditLogServiceAPI* | [**AuditLogServiceListAuditLogs**](docs/AuditLogServiceAPI.md#auditlogservicelistauditlogs) | **Get** /v1beta1/clusters/{clusterId}/auditlogs | List database audit logs.
*AuditLogServiceAPI* | [**AuditLogServiceUpdateAuditLogFilterRule**](docs/AuditLogServiceAPI.md#auditlogserviceupdateauditlogfilterrule) | **Patch** /v1beta1/clusters/{clusterId}/auditlogs/filterRules/{name} | Update audit log filter rule.


## Documentation For Models

 - [Any](docs/Any.md)
 - [AuditLog](docs/AuditLog.md)
 - [AuditLogFilter](docs/AuditLogFilter.md)
 - [AuditLogFilterRule](docs/AuditLogFilterRule.md)
 - [AuditLogServiceCreateAuditLogFilterRuleBody](docs/AuditLogServiceCreateAuditLogFilterRuleBody.md)
 - [AuditLogServiceDownloadAuditLogsBody](docs/AuditLogServiceDownloadAuditLogsBody.md)
 - [AuditLogServiceUpdateAuditLogFilterRuleBody](docs/AuditLogServiceUpdateAuditLogFilterRuleBody.md)
 - [DownloadAuditLogsResponse](docs/DownloadAuditLogsResponse.md)
 - [ListAuditLogFilterRulesResponse](docs/ListAuditLogFilterRulesResponse.md)
 - [ListAuditLogsResponse](docs/ListAuditLogsResponse.md)
 - [Status](docs/Status.md)


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



