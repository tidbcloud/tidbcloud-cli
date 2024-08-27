# \AccountAPI

All URIs are relative to *http://iam.tidbapi.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CustomerSignupUrlPost**](AccountAPI.md#CustomerSignupUrlPost) | **Post** /customerSignupUrl | Create a new signup URL for an MSP customer
[**MspCustomersCustomerOrgIdGet**](AccountAPI.md#MspCustomersCustomerOrgIdGet) | **Get** /mspCustomers/{customerOrgId} | Retrieve a single MSP customer
[**MspCustomersGet**](AccountAPI.md#MspCustomersGet) | **Get** /mspCustomers | Get a list of MSP customers
[**V1beta1ClustersClusterIdDbuserGet**](AccountAPI.md#V1beta1ClustersClusterIdDbuserGet) | **Get** /v1beta1/clusters/{clusterId}/dbuser | get one dbuser
[**V1beta1ClustersClusterIdSqlUsersGet**](AccountAPI.md#V1beta1ClustersClusterIdSqlUsersGet) | **Get** /v1beta1/clusters/{clusterId}/sqlUsers | Get all sql users
[**V1beta1ClustersClusterIdSqlUsersPost**](AccountAPI.md#V1beta1ClustersClusterIdSqlUsersPost) | **Post** /v1beta1/clusters/{clusterId}/sqlUsers | Create one sql user
[**V1beta1ClustersClusterIdSqlUsersUserNameDelete**](AccountAPI.md#V1beta1ClustersClusterIdSqlUsersUserNameDelete) | **Delete** /v1beta1/clusters/{clusterId}/sqlUsers/{userName} | Delete one sql user
[**V1beta1ClustersClusterIdSqlUsersUserNameGet**](AccountAPI.md#V1beta1ClustersClusterIdSqlUsersUserNameGet) | **Get** /v1beta1/clusters/{clusterId}/sqlUsers/{userName} | Query sql user
[**V1beta1ClustersClusterIdSqlUsersUserNamePatch**](AccountAPI.md#V1beta1ClustersClusterIdSqlUsersUserNamePatch) | **Patch** /v1beta1/clusters/{clusterId}/sqlUsers/{userName} | Update one sql user
[**V1beta1ProjectsGet**](AccountAPI.md#V1beta1ProjectsGet) | **Get** /v1beta1/projects | Get  list of org projects



## CustomerSignupUrlPost

> ApiOpenApiMspCustomerSignupUrl CustomerSignupUrlPost(ctx).MspCustomerOrgId(mspCustomerOrgId).Execute()

Create a new signup URL for an MSP customer



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
	mspCustomerOrgId := *openapiclient.NewApiOpenApiCreateMspCustomerSignupUrlReq() // ApiOpenApiCreateMspCustomerSignupUrlReq | The MSP customer org ID.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.CustomerSignupUrlPost(context.Background()).MspCustomerOrgId(mspCustomerOrgId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.CustomerSignupUrlPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CustomerSignupUrlPost`: ApiOpenApiMspCustomerSignupUrl
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.CustomerSignupUrlPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCustomerSignupUrlPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mspCustomerOrgId** | [**ApiOpenApiCreateMspCustomerSignupUrlReq**](ApiOpenApiCreateMspCustomerSignupUrlReq.md) | The MSP customer org ID. | 

### Return type

[**ApiOpenApiMspCustomerSignupUrl**](ApiOpenApiMspCustomerSignupUrl.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MspCustomersCustomerOrgIdGet

> ApiOpenApiMspCustomer MspCustomersCustomerOrgIdGet(ctx, customerOrgId).Execute()

Retrieve a single MSP customer



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
	customerOrgId := "customerOrgId_example" // string | The MSP customer org ID.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.MspCustomersCustomerOrgIdGet(context.Background(), customerOrgId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.MspCustomersCustomerOrgIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `MspCustomersCustomerOrgIdGet`: ApiOpenApiMspCustomer
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.MspCustomersCustomerOrgIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**customerOrgId** | **string** | The MSP customer org ID. | 

### Other Parameters

Other parameters are passed through a pointer to a apiMspCustomersCustomerOrgIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ApiOpenApiMspCustomer**](ApiOpenApiMspCustomer.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MspCustomersGet

> ApiOpenApiListMspCustomerRsp MspCustomersGet(ctx).PageToken(pageToken).PageSize(pageSize).Execute()

Get a list of MSP customers



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
	pageToken := "pageToken_example" // string | The page token of the next page. (optional)
	pageSize := int32(56) // int32 | The page size of the next page. If `pageSize` is set to 0, it returns all MSP customers in one page. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.MspCustomersGet(context.Background()).PageToken(pageToken).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.MspCustomersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `MspCustomersGet`: ApiOpenApiListMspCustomerRsp
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.MspCustomersGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMspCustomersGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageToken** | **string** | The page token of the next page. | 
 **pageSize** | **int32** | The page size of the next page. If &#x60;pageSize&#x60; is set to 0, it returns all MSP customers in one page. | 

### Return type

[**ApiOpenApiListMspCustomerRsp**](ApiOpenApiListMspCustomerRsp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1beta1ClustersClusterIdDbuserGet

> ApiGetDbuserRsp V1beta1ClustersClusterIdDbuserGet(ctx, clusterId).Execute()

get one dbuser



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
	clusterId := "clusterId_example" // string | The id of the cluster.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.V1beta1ClustersClusterIdDbuserGet(context.Background(), clusterId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.V1beta1ClustersClusterIdDbuserGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1beta1ClustersClusterIdDbuserGet`: ApiGetDbuserRsp
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.V1beta1ClustersClusterIdDbuserGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | The id of the cluster. | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1beta1ClustersClusterIdDbuserGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ApiGetDbuserRsp**](ApiGetDbuserRsp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1beta1ClustersClusterIdSqlUsersGet

> ApiListSqlUsersRsp V1beta1ClustersClusterIdSqlUsersGet(ctx, clusterId).PageToken(pageToken).PageSize(pageSize).Execute()

Get all sql users



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
	clusterId := "clusterId_example" // string | The id of the cluster.
	pageToken := "pageToken_example" // string | The page token of the next page. (optional)
	pageSize := int32(56) // int32 | The page size of the next page. If `pageSize` is set to 0, it returns 100 records in one page. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.V1beta1ClustersClusterIdSqlUsersGet(context.Background(), clusterId).PageToken(pageToken).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.V1beta1ClustersClusterIdSqlUsersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1beta1ClustersClusterIdSqlUsersGet`: ApiListSqlUsersRsp
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.V1beta1ClustersClusterIdSqlUsersGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | The id of the cluster. | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1beta1ClustersClusterIdSqlUsersGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **pageToken** | **string** | The page token of the next page. | 
 **pageSize** | **int32** | The page size of the next page. If &#x60;pageSize&#x60; is set to 0, it returns 100 records in one page. | 

### Return type

[**ApiListSqlUsersRsp**](ApiListSqlUsersRsp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1beta1ClustersClusterIdSqlUsersPost

> ApiSqlUser V1beta1ClustersClusterIdSqlUsersPost(ctx, clusterId).SqlUser(sqlUser).Execute()

Create one sql user



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
	clusterId := "clusterId_example" // string | The id of the cluster.
	sqlUser := *openapiclient.NewApiCreateSqlUserReq() // ApiCreateSqlUserReq | create sql user request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.V1beta1ClustersClusterIdSqlUsersPost(context.Background(), clusterId).SqlUser(sqlUser).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.V1beta1ClustersClusterIdSqlUsersPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1beta1ClustersClusterIdSqlUsersPost`: ApiSqlUser
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.V1beta1ClustersClusterIdSqlUsersPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | The id of the cluster. | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1beta1ClustersClusterIdSqlUsersPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **sqlUser** | [**ApiCreateSqlUserReq**](ApiCreateSqlUserReq.md) | create sql user request | 

### Return type

[**ApiSqlUser**](ApiSqlUser.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1beta1ClustersClusterIdSqlUsersUserNameDelete

> ApiBasicResp V1beta1ClustersClusterIdSqlUsersUserNameDelete(ctx, clusterId, userName).Execute()

Delete one sql user



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
	clusterId := "clusterId_example" // string | The id of the cluster.
	userName := "userName_example" // string | The name of the sql user.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNameDelete(context.Background(), clusterId, userName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNameDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1beta1ClustersClusterIdSqlUsersUserNameDelete`: ApiBasicResp
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNameDelete`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | The id of the cluster. | 
**userName** | **string** | The name of the sql user. | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1beta1ClustersClusterIdSqlUsersUserNameDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ApiBasicResp**](ApiBasicResp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1beta1ClustersClusterIdSqlUsersUserNameGet

> ApiSqlUser V1beta1ClustersClusterIdSqlUsersUserNameGet(ctx, clusterId, userName).Execute()

Query sql user



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
	clusterId := "clusterId_example" // string | The id of the cluster.
	userName := "userName_example" // string | The name of the sql user.

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNameGet(context.Background(), clusterId, userName).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNameGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1beta1ClustersClusterIdSqlUsersUserNameGet`: ApiSqlUser
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNameGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | The id of the cluster. | 
**userName** | **string** | The name of the sql user. | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1beta1ClustersClusterIdSqlUsersUserNameGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ApiSqlUser**](ApiSqlUser.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1beta1ClustersClusterIdSqlUsersUserNamePatch

> ApiSqlUser V1beta1ClustersClusterIdSqlUsersUserNamePatch(ctx, clusterId, userName).SqlUser(sqlUser).Execute()

Update one sql user



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
	clusterId := "clusterId_example" // string | The id of the cluster.
	userName := "userName_example" // string | The name of the sql user.
	sqlUser := *openapiclient.NewApiUpdateSqlUserReq() // ApiUpdateSqlUserReq | update sql user request

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNamePatch(context.Background(), clusterId, userName).SqlUser(sqlUser).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNamePatch``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1beta1ClustersClusterIdSqlUsersUserNamePatch`: ApiSqlUser
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.V1beta1ClustersClusterIdSqlUsersUserNamePatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**clusterId** | **string** | The id of the cluster. | 
**userName** | **string** | The name of the sql user. | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1beta1ClustersClusterIdSqlUsersUserNamePatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **sqlUser** | [**ApiUpdateSqlUserReq**](ApiUpdateSqlUserReq.md) | update sql user request | 

### Return type

[**ApiSqlUser**](ApiSqlUser.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1beta1ProjectsGet

> ApiListProjectsRsp V1beta1ProjectsGet(ctx).PageToken(pageToken).PageSize(pageSize).Execute()

Get  list of org projects



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
	pageToken := "pageToken_example" // string | The page token of the next page. (optional)
	pageSize := int32(56) // int32 | The page size of the next page. If `pageSize` is set to 0, it returns 100 records in one page. (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AccountAPI.V1beta1ProjectsGet(context.Background()).PageToken(pageToken).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AccountAPI.V1beta1ProjectsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1beta1ProjectsGet`: ApiListProjectsRsp
	fmt.Fprintf(os.Stdout, "Response from `AccountAPI.V1beta1ProjectsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiV1beta1ProjectsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageToken** | **string** | The page token of the next page. | 
 **pageSize** | **int32** | The page size of the next page. If &#x60;pageSize&#x60; is set to 0, it returns 100 records in one page. | 

### Return type

[**ApiListProjectsRsp**](ApiListProjectsRsp.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

