# \BranchServiceAPI

All URIs are relative to *https://serverless.tidbapi.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BranchServiceCreateBranch**](BranchServiceAPI.md#BranchServiceCreateBranch) | **Post** /v1beta1/clusters/{clusterId}/branches | Creates a branch.
[**BranchServiceDeleteBranch**](BranchServiceAPI.md#BranchServiceDeleteBranch) | **Delete** /v1beta1/clusters/{clusterId}/branches/{branchId} | Deletes a branch.
[**BranchServiceGetBranch**](BranchServiceAPI.md#BranchServiceGetBranch) | **Get** /v1beta1/clusters/{clusterId}/branches/{branchId} | Gets information about a branch.
[**BranchServiceListBranches**](BranchServiceAPI.md#BranchServiceListBranches) | **Get** /v1beta1/clusters/{clusterId}/branches | Lists information about branches.



## BranchServiceCreateBranch

> V1beta1Branch BranchServiceCreateBranch(ctx, clusterId).Branch(branch).Execute()

Creates a branch.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/branch"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The cluster ID of the branch
	branch := *openapiclient.NewV1beta1Branch("DisplayName_example") // V1beta1Branch | Required. The resource being created

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BranchServiceAPI.BranchServiceCreateBranch(context.Background(), clusterId).Branch(branch).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BranchServiceAPI.BranchServiceCreateBranch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BranchServiceCreateBranch`: V1beta1Branch
	fmt.Fprintf(os.Stdout, "Response from `BranchServiceAPI.BranchServiceCreateBranch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The cluster ID of the branch | 

### Other Parameters

Other parameters are passed through a pointer to a apiBranchServiceCreateBranchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **branch** | [**V1beta1Branch**](V1beta1Branch.md) | Required. The resource being created | 

### Return type

[**V1beta1Branch**](V1beta1Branch.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BranchServiceDeleteBranch

> V1beta1Branch BranchServiceDeleteBranch(ctx, clusterId, branchId).Execute()

Deletes a branch.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/branch"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The cluster ID of the branch
	branchId := "branchId_example" // string | Required. The branch ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BranchServiceAPI.BranchServiceDeleteBranch(context.Background(), clusterId, branchId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BranchServiceAPI.BranchServiceDeleteBranch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BranchServiceDeleteBranch`: V1beta1Branch
	fmt.Fprintf(os.Stdout, "Response from `BranchServiceAPI.BranchServiceDeleteBranch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The cluster ID of the branch | 
**branchId** | **string** | Required. The branch ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiBranchServiceDeleteBranchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**V1beta1Branch**](V1beta1Branch.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BranchServiceGetBranch

> V1beta1Branch BranchServiceGetBranch(ctx, clusterId, branchId).View(view).Execute()

Gets information about a branch.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/branch"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The cluster ID of the branch
	branchId := "branchId_example" // string | Required. The branch ID
	view := "view_example" // string | Optional. The view of the branch to return. Defaults to FULL   - BASIC: Basic response contains basic information for a branch.  - FULL: FULL response contains all detailed information for a branch. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BranchServiceAPI.BranchServiceGetBranch(context.Background(), clusterId, branchId).View(view).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BranchServiceAPI.BranchServiceGetBranch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BranchServiceGetBranch`: V1beta1Branch
	fmt.Fprintf(os.Stdout, "Response from `BranchServiceAPI.BranchServiceGetBranch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The cluster ID of the branch | 
**branchId** | **string** | Required. The branch ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiBranchServiceGetBranchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **view** | **string** | Optional. The view of the branch to return. Defaults to FULL   - BASIC: Basic response contains basic information for a branch.  - FULL: FULL response contains all detailed information for a branch. | 

### Return type

[**V1beta1Branch**](V1beta1Branch.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BranchServiceListBranches

> V1beta1ListBranchesResponse BranchServiceListBranches(ctx, clusterId).PageSize(pageSize).PageToken(pageToken).Execute()

Lists information about branches.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/branch"
)

func main() {
	clusterId := "clusterId_example" // string | Required. The ID of the project to which the clusters belong.
	pageSize := int32(56) // int32 | Optional. Requested page size. Server may return fewer items than requested. If unspecified, server will pick an appropriate default. (optional)
	pageToken := "pageToken_example" // string | Optional. A token identifying a page of results the server should return. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.BranchServiceAPI.BranchServiceListBranches(context.Background(), clusterId).PageSize(pageSize).PageToken(pageToken).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `BranchServiceAPI.BranchServiceListBranches``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `BranchServiceListBranches`: V1beta1ListBranchesResponse
	fmt.Fprintf(os.Stdout, "Response from `BranchServiceAPI.BranchServiceListBranches`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | Required. The ID of the project to which the clusters belong. | 

### Other Parameters

Other parameters are passed through a pointer to a apiBranchServiceListBranchesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **pageSize** | **int32** | Optional. Requested page size. Server may return fewer items than requested. If unspecified, server will pick an appropriate default. | 
 **pageToken** | **string** | Optional. A token identifying a page of results the server should return. | 

### Return type

[**V1beta1ListBranchesResponse**](V1beta1ListBranchesResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

