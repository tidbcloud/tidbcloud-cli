# \ExportServiceAPI

All URIs are relative to *https://serverless.tidbapi.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ExportServiceCancelExport**](ExportServiceAPI.md#ExportServiceCancelExport) | **Post** /v1beta1/clusters/{clusterId}/exports/{exportId}:cancel | Cancel a specific export job.
[**ExportServiceCreateExport**](ExportServiceAPI.md#ExportServiceCreateExport) | **Post** /v1beta1/clusters/{clusterId}/exports | Create an export job
[**ExportServiceDeleteExport**](ExportServiceAPI.md#ExportServiceDeleteExport) | **Delete** /v1beta1/clusters/{clusterId}/exports/{exportId} | Delete an export job
[**ExportServiceDownloadExport**](ExportServiceAPI.md#ExportServiceDownloadExport) | **Post** /v1beta1/clusters/{clusterId}/exports/{exportId}:download | Generate download url
[**ExportServiceGetExport**](ExportServiceAPI.md#ExportServiceGetExport) | **Get** /v1beta1/clusters/{clusterId}/exports/{exportId} | Retrieves details of an export job.
[**ExportServiceListExports**](ExportServiceAPI.md#ExportServiceListExports) | **Get** /v1beta1/clusters/{clusterId}/exports | Provides a list of export jobs.



## ExportServiceCancelExport

> Export ExportServiceCancelExport(ctx, clusterId, exportId).Body(body).Execute()

Cancel a specific export job.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The ID of the cluster.
	exportId := "exportId_example" // string | Required. The ID of the export to be retrieved.
	body := map[string]interface{}{ ... } // map[string]interface{} | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExportServiceAPI.ExportServiceCancelExport(context.Background(), clusterId, exportId).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExportServiceAPI.ExportServiceCancelExport``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExportServiceCancelExport`: Export
	fmt.Fprintf(os.Stdout, "Response from `ExportServiceAPI.ExportServiceCancelExport`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The ID of the cluster. | 
**exportId** | **string** | Required. The ID of the export to be retrieved. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExportServiceCancelExportRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | **map[string]interface{}** |  | 

### Return type

[**Export**](Export.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExportServiceCreateExport

> Export ExportServiceCreateExport(ctx, clusterId).Body(body).Execute()

Create an export job

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The ID of the cluster.
	body := *openapiclient.NewExportServiceCreateExportBody() // ExportServiceCreateExportBody | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExportServiceAPI.ExportServiceCreateExport(context.Background(), clusterId).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExportServiceAPI.ExportServiceCreateExport``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExportServiceCreateExport`: Export
	fmt.Fprintf(os.Stdout, "Response from `ExportServiceAPI.ExportServiceCreateExport`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The ID of the cluster. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExportServiceCreateExportRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**ExportServiceCreateExportBody**](ExportServiceCreateExportBody.md) |  | 

### Return type

[**Export**](Export.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExportServiceDeleteExport

> Export ExportServiceDeleteExport(ctx, clusterId, exportId).Execute()

Delete an export job

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The ID of the cluster.
	exportId := "exportId_example" // string | Required. The ID of the export to be retrieved.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExportServiceAPI.ExportServiceDeleteExport(context.Background(), clusterId, exportId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExportServiceAPI.ExportServiceDeleteExport``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExportServiceDeleteExport`: Export
	fmt.Fprintf(os.Stdout, "Response from `ExportServiceAPI.ExportServiceDeleteExport`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The ID of the cluster. | 
**exportId** | **string** | Required. The ID of the export to be retrieved. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExportServiceDeleteExportRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Export**](Export.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExportServiceDownloadExport

> DownloadExportsResponse ExportServiceDownloadExport(ctx, clusterId, exportId).Body(body).Execute()

Generate download url

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The ID of the cluster.
	exportId := "exportId_example" // string | Required. The ID of the export to be retrieved.
	body := map[string]interface{}{ ... } // map[string]interface{} | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExportServiceAPI.ExportServiceDownloadExport(context.Background(), clusterId, exportId).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExportServiceAPI.ExportServiceDownloadExport``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExportServiceDownloadExport`: DownloadExportsResponse
	fmt.Fprintf(os.Stdout, "Response from `ExportServiceAPI.ExportServiceDownloadExport`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The ID of the cluster. | 
**exportId** | **string** | Required. The ID of the export to be retrieved. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExportServiceDownloadExportRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | **map[string]interface{}** |  | 

### Return type

[**DownloadExportsResponse**](DownloadExportsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExportServiceGetExport

> Export ExportServiceGetExport(ctx, clusterId, exportId).Execute()

Retrieves details of an export job.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The ID of the cluster.
	exportId := "exportId_example" // string | Required. The ID of the export to be retrieved.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExportServiceAPI.ExportServiceGetExport(context.Background(), clusterId, exportId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExportServiceAPI.ExportServiceGetExport``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExportServiceGetExport`: Export
	fmt.Fprintf(os.Stdout, "Response from `ExportServiceAPI.ExportServiceGetExport`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The ID of the cluster. | 
**exportId** | **string** | Required. The ID of the export to be retrieved. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExportServiceGetExportRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**Export**](Export.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExportServiceListExports

> ListExportsResponse ExportServiceListExports(ctx, clusterId).PageSize(pageSize).PageToken(pageToken).OrderBy(orderBy).Execute()

Provides a list of export jobs.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The cluster ID to list exports for.
	pageSize := int32(56) // int32 | Optional. The maximum number of clusters to return. Default is 10. (optional)
	pageToken := "pageToken_example" // string | Optional. The page token from the previous response for pagination. (optional)
	orderBy := "orderBy_example" // string | Optional. List exports order by, separated by comma, default is ascending. Example: \"foo, bar desc\". Supported field: create_time (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ExportServiceAPI.ExportServiceListExports(context.Background(), clusterId).PageSize(pageSize).PageToken(pageToken).OrderBy(orderBy).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ExportServiceAPI.ExportServiceListExports``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ExportServiceListExports`: ListExportsResponse
	fmt.Fprintf(os.Stdout, "Response from `ExportServiceAPI.ExportServiceListExports`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The cluster ID to list exports for. | 

### Other Parameters

Other parameters are passed through a pointer to a apiExportServiceListExportsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **pageSize** | **int32** | Optional. The maximum number of clusters to return. Default is 10. | 
 **pageToken** | **string** | Optional. The page token from the previous response for pagination. | 
 **orderBy** | **string** | Optional. List exports order by, separated by comma, default is ascending. Example: \&quot;foo, bar desc\&quot;. Supported field: create_time | 

### Return type

[**ListExportsResponse**](ListExportsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

