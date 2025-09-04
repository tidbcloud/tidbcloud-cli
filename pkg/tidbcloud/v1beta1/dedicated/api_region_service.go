/*
TiDB Cloud Dedicated API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Dedicated](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-dedicated) clusters. For TiDB Cloud Starter or TiDB Cloud Essential clusters, use the [TiDB Cloud Starter and Essential API](). For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/tidbcloud/api-overview/).  # Overview  The TiDB Cloud API is a [REST interface](https://en.wikipedia.org/wiki/Representational_state_transfer) that provides you with programmatic access to manage clusters and related resources within TiDB Cloud.  The API has the following features:  - **JSON entities.** All entities are expressed in JSON. - **HTTPS-only.** You can only access the API via HTTPS, ensuring all the data sent over the network is encrypted with TLS. - **Key-based access and digest authentication.** Before you access TiDB Cloud API, you must generate an API key. All requests are authenticated through [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication), ensuring the API key is never sent over the network.  # Get Started  This guide helps you make your first API call to TiDB Cloud API. You'll learn how to authenticate a request, build a request, and interpret the response.  ## Prerequisites  To complete this guide, you need to perform the following tasks:  - Create a [TiDB Cloud account](https://tidbcloud.com/free-trial) - Install [curl](https://curl.se/)  ## Step 1. Create an API key  To create an API key, log in to your TiDB Cloud console. Navigate to the [**API Keys**](https://tidbcloud.com/org-settings/api-keys) page of your organization, and create an API key.  An API key contains a public key and a private key. Copy and save them in a secure location. You will need to use the API key later in this guide.  For more details about creating API key, refer to [API Key Management](#section/Authentication/API-Key-Management).  ## Step 2. Make your first API call  ### Build an API call  TiDB Cloud API call consists of the following components:  - **A host**. The host for TiDB Cloud API is <https://dedicated.tidbapi.com>. - **An API Key**. The public key and the private key are required for authentication. - **A request**. When submitting data to a resource via `POST`, `PATCH`, or `PUT`, you must submit your payload in JSON.  In this guide, you call the [List clusters](#tag/Cluster/operation/ClusterService_ListClusters) endpoint. For the detailed description of the endpoint, see the [API reference](#tag/Cluster/operation/ClusterService_ListClusters).  ### Call an API endpoint  To get all clusters in your organization, run the following command in your terminal. Remember to change `YOUR_PUBLIC_KEY` to your public key and `YOUR_PRIVATE_KEY` to your private key.  ```shell curl --digest \\  --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\  --request GET \\  --url 'https://dedicated.tidbapi.com/v1beta1/clusters' ```  ## Step 3. Check the response  After making the API call, if the status code in response is `200` and you see details about all clusters in your organization, your request is successful.  # Authentication  The TiDB Cloud API uses [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication). It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, the TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. The role of the API key is always `Organization Owner` currently. 5. Click **Next**. Copy and save the public key and the private key. 6. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 7. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Edit**. 4. You can update the API key description. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # Rate Limiting  The TiDB Cloud API allows up to 100 requests per minute per API key. If you exceed the rate limit, the API returns a `429` error. For more quota, you can [submit a request](https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519) to contact our support team.  Each API request returns the following headers about the limit.  - `X-Ratelimit-Limit-Minute`: The number of requests allowed per minute. It is 100 currently. - `X-Ratelimit-Remaining-Minute`: The number of remaining requests in the current minute. When it reaches `0`, the API returns a `429` error and indicates that you exceed the rate limit. - `X-Ratelimit-Reset`: The time in seconds at which the current rate limit resets.  If you exceed the rate limit, an error response returns like this.  ``` > HTTP/2 429 > date: Fri, 22 Jul 2022 05:28:37 GMT > content-type: application/json > content-length: 66 > x-ratelimit-reset: 23 > x-ratelimit-remaining-minute: 0 > x-ratelimit-limit-minute: 100 > x-kong-response-latency: 2 > server: kong/2.8.1  > {\"details\":[],\"code\":49900007,\"message\":\"The request exceeded the limit of 100 times per apikey per minute. For more quota, please contact us: https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519\"} ```  # API Changelog  This changelog lists all changes to the TiDB Cloud API.  <!-- In reverse chronological order -->  ## 20250812  - Initial release of the TiDB Cloud Dedicated API, including the following resources and endpoints:  * Cluster    * [List clusters](#tag/Cluster/operation/ClusterService_ListClusters)    * [Create a cluster](#tag/Cluster/operation/ClusterService_CreateCluster)    * [Get a cluster](#tag/Cluster/operation/ClusterService_GetCluster)    * [Delete a cluster](#tag/Cluster/operation/ClusterService_DeleteCluster)    * [Update a cluster](#tag/Cluster/operation/ClusterService_UpdateCluster)    * [Pause a cluster](#tag/Cluster/operation/ClusterService_PauseCluster)    * [Resume a cluster](#tag/Cluster/operation/ClusterService_ResumeCluster)    * [Reset the root password of a cluster](#tag/Cluster/operation/ClusterService_ResetRootPassword)    * [List node quotas for your organization](#tag/Cluster/operation/ClusterService_ShowNodeQuota)    * [Get log redaction policy](#tag/Cluster/operation/ClusterService_GetLogRedactionPolicy)   * Region    * [List regions](#tag/Region/operation/RegionService_ListRegions)    * [Get a region](#tag/Region/operation/RegionService_GetRegion)    * [List cloud providers](#tag/Region/operation/RegionService_ShowCloudProviders)    * [List node specs](#tag/Region/operation/RegionService_ListNodeSpecs)    * [Get a node spec](#tag/Region/operation/RegionService_GetNodeSpec)   * Private Endpoint Connection    * [Get private link service for a TiDB node group](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateLinkService)    * [Create a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_CreatePrivateEndpointConnection)    * [List private endpoint connections](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_ListPrivateEndpointConnections)    * [Get a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateEndpointConnection)    * [Delete a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_DeletePrivateEndpointConnection)   * Import    * [List import tasks](#tag/Import/operation/ListImports)    * [Create an import task](#tag/Import/operation/CreateImport)    * [Get an import task](#tag/Import/operation/GetImport)    * [Cancel an import task](#tag/Import/operation/CancelImport)

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RegionServiceAPIService RegionServiceAPI service
type RegionServiceAPIService service

type ApiRegionServiceGetNodeSpecRequest struct {
	ctx           context.Context
	ApiService    *RegionServiceAPIService
	regionId      string
	componentType ClusterServiceListNodeInstancesComponentTypeParameter
	nodeSpecKey   string
	projectId     *string
	clusterId     *string
}

// The ID of the project for which to retrieve the node spec. If not specified, the project ID of the default project is used.
func (r ApiRegionServiceGetNodeSpecRequest) ProjectId(projectId string) ApiRegionServiceGetNodeSpecRequest {
	r.projectId = &projectId
	return r
}

// The ID of the cluster. If specified, only node specs that are available to the specified cluster are returned. If not specified, all available node specs are returned.
func (r ApiRegionServiceGetNodeSpecRequest) ClusterId(clusterId string) ApiRegionServiceGetNodeSpecRequest {
	r.clusterId = &clusterId
	return r
}

func (r ApiRegionServiceGetNodeSpecRequest) Execute() (*V1beta1NodeSpec, *http.Response, error) {
	return r.ApiService.RegionServiceGetNodeSpecExecute(r)
}

/*
RegionServiceGetNodeSpec Get a node spec

Retrieves details of a node spec used for creating or scaling a cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param regionId The region ID of the node spec to retrieve, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-east-1`. If not specified, the default region ID is used.
	@param componentType Filters the results to include only node specs for the specified component type.
	@param nodeSpecKey The key of the node spec to retrieve. For example, `8C32G`.
	@return ApiRegionServiceGetNodeSpecRequest
*/
func (a *RegionServiceAPIService) RegionServiceGetNodeSpec(ctx context.Context, regionId string, componentType ClusterServiceListNodeInstancesComponentTypeParameter, nodeSpecKey string) ApiRegionServiceGetNodeSpecRequest {
	return ApiRegionServiceGetNodeSpecRequest{
		ApiService:    a,
		ctx:           ctx,
		regionId:      regionId,
		componentType: componentType,
		nodeSpecKey:   nodeSpecKey,
	}
}

// Execute executes the request
//
//	@return V1beta1NodeSpec
func (a *RegionServiceAPIService) RegionServiceGetNodeSpecExecute(r ApiRegionServiceGetNodeSpecRequest) (*V1beta1NodeSpec, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *V1beta1NodeSpec
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RegionServiceAPIService.RegionServiceGetNodeSpec")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/regions/{regionId}/componentTypes/{componentType}/nodeSpecs/{nodeSpecKey}"
	localVarPath = strings.Replace(localVarPath, "{"+"regionId"+"}", url.PathEscape(parameterValueToString(r.regionId, "regionId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"componentType"+"}", url.PathEscape(parameterValueToString(r.componentType, "componentType")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"nodeSpecKey"+"}", url.PathEscape(parameterValueToString(r.nodeSpecKey, "nodeSpecKey")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.projectId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "projectId", r.projectId, "", "")
	}
	if r.clusterId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "clusterId", r.clusterId, "", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		var v GooglerpcStatus
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiRegionServiceGetRegionRequest struct {
	ctx        context.Context
	ApiService *RegionServiceAPIService
	regionId   string
	projectId  *string
}

// The ID of the project for which to retrieve the region. If not specified, the project ID of the default project is used.
func (r ApiRegionServiceGetRegionRequest) ProjectId(projectId string) ApiRegionServiceGetRegionRequest {
	r.projectId = &projectId
	return r
}

func (r ApiRegionServiceGetRegionRequest) Execute() (*Commonv1beta1Region, *http.Response, error) {
	return r.ApiService.RegionServiceGetRegionExecute(r)
}

/*
RegionServiceGetRegion Get a region

Retrieves details of a specific region by region ID.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param regionId The ID of the region to retrieve, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-east-1`.
	@return ApiRegionServiceGetRegionRequest
*/
func (a *RegionServiceAPIService) RegionServiceGetRegion(ctx context.Context, regionId string) ApiRegionServiceGetRegionRequest {
	return ApiRegionServiceGetRegionRequest{
		ApiService: a,
		ctx:        ctx,
		regionId:   regionId,
	}
}

// Execute executes the request
//
//	@return Commonv1beta1Region
func (a *RegionServiceAPIService) RegionServiceGetRegionExecute(r ApiRegionServiceGetRegionRequest) (*Commonv1beta1Region, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Commonv1beta1Region
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RegionServiceAPIService.RegionServiceGetRegion")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/regions/{regionId}"
	localVarPath = strings.Replace(localVarPath, "{"+"regionId"+"}", url.PathEscape(parameterValueToString(r.regionId, "regionId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.projectId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "projectId", r.projectId, "", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		var v GooglerpcStatus
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiRegionServiceListNodeSpecsRequest struct {
	ctx           context.Context
	ApiService    *RegionServiceAPIService
	regionId      string
	componentType *ClusterServiceListNodeInstancesComponentTypeParameter
	projectId     *string
	clusterId     *string
	pageSize      *int32
	pageToken     *string
	skip          *int32
}

// Filters the results to include only node specs for the specified component type.
func (r ApiRegionServiceListNodeSpecsRequest) ComponentType(componentType ClusterServiceListNodeInstancesComponentTypeParameter) ApiRegionServiceListNodeSpecsRequest {
	r.componentType = &componentType
	return r
}

// Filters the results to include only node specs available to the specified project. If not specified, the project ID of the default project is used.
func (r ApiRegionServiceListNodeSpecsRequest) ProjectId(projectId string) ApiRegionServiceListNodeSpecsRequest {
	r.projectId = &projectId
	return r
}

// Filters the results to include only node specs available to the specified cluster. If not specified, all available node specs are returned.
func (r ApiRegionServiceListNodeSpecsRequest) ClusterId(clusterId string) ApiRegionServiceListNodeSpecsRequest {
	r.clusterId = &clusterId
	return r
}

// The maximum number of node specs to return. If not specified, at most 10 node specs will be returned. The maximum value is &#x60;100&#x60;. Values greater than &#x60;100&#x60; are set to &#x60;100&#x60;.
func (r ApiRegionServiceListNodeSpecsRequest) PageSize(pageSize int32) ApiRegionServiceListNodeSpecsRequest {
	r.pageSize = &pageSize
	return r
}

// The pagination token received from a previous [List node specs](#tag/Region/operation/RegionService_ListNodeSpecs) request. Use this token to retrieve the next page of results.  **Note**: When paginating, all other parameters must match the original request.
func (r ApiRegionServiceListNodeSpecsRequest) PageToken(pageToken string) ApiRegionServiceListNodeSpecsRequest {
	r.pageToken = &pageToken
	return r
}

// The number of node specs to skip before returning results. If the value exceeds the total number of clusters, the response is &#x60;200&#x60; with an empty list and no &#x60;nextPageToken&#x60;.
func (r ApiRegionServiceListNodeSpecsRequest) Skip(skip int32) ApiRegionServiceListNodeSpecsRequest {
	r.skip = &skip
	return r
}

func (r ApiRegionServiceListNodeSpecsRequest) Execute() (*V1beta1ListNodeSpecsResponse, *http.Response, error) {
	return r.ApiService.RegionServiceListNodeSpecsExecute(r)
}

/*
RegionServiceListNodeSpecs List node specs

Retrieves a paginated list of node specifications (specs) available for creating or scaling a cluster in the specified region. You can filter the results using the `componentType`, `projectId`, or `clusterId` parameter.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param regionId The ID of the region, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-east-1`.
	@return ApiRegionServiceListNodeSpecsRequest
*/
func (a *RegionServiceAPIService) RegionServiceListNodeSpecs(ctx context.Context, regionId string) ApiRegionServiceListNodeSpecsRequest {
	return ApiRegionServiceListNodeSpecsRequest{
		ApiService: a,
		ctx:        ctx,
		regionId:   regionId,
	}
}

// Execute executes the request
//
//	@return V1beta1ListNodeSpecsResponse
func (a *RegionServiceAPIService) RegionServiceListNodeSpecsExecute(r ApiRegionServiceListNodeSpecsRequest) (*V1beta1ListNodeSpecsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *V1beta1ListNodeSpecsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RegionServiceAPIService.RegionServiceListNodeSpecs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/regions/{regionId}/nodeSpecs"
	localVarPath = strings.Replace(localVarPath, "{"+"regionId"+"}", url.PathEscape(parameterValueToString(r.regionId, "regionId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.componentType != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "componentType", r.componentType, "", "")
	}
	if r.projectId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "projectId", r.projectId, "", "")
	}
	if r.clusterId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "clusterId", r.clusterId, "", "")
	}
	if r.pageSize != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageSize", r.pageSize, "", "")
	} else {
		var defaultValue int32 = 10
		r.pageSize = &defaultValue
	}
	if r.pageToken != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageToken", r.pageToken, "", "")
	}
	if r.skip != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "skip", r.skip, "", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		var v GooglerpcStatus
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiRegionServiceListRegionsRequest struct {
	ctx           context.Context
	ApiService    *RegionServiceAPIService
	cloudProvider *string
	projectId     *string
	pageSize      *int32
	pageToken     *string
	skip          *int32
}

// The cloud provider where the region is located. If specified, only regions of the specified cloud provider are returned.
func (r ApiRegionServiceListRegionsRequest) CloudProvider(cloudProvider string) ApiRegionServiceListRegionsRequest {
	r.cloudProvider = &cloudProvider
	return r
}

// The ID of the project to list regions for. If not specified, the project ID of the default project is used.
func (r ApiRegionServiceListRegionsRequest) ProjectId(projectId string) ApiRegionServiceListRegionsRequest {
	r.projectId = &projectId
	return r
}

// The maximum number of regions to return. If not specified, at most 20 regions will be returned. The maximum value is &#x60;100&#x60;. Values greater than &#x60;100&#x60; are set to &#x60;100&#x60;.
func (r ApiRegionServiceListRegionsRequest) PageSize(pageSize int32) ApiRegionServiceListRegionsRequest {
	r.pageSize = &pageSize
	return r
}

// The pagination token received from a previous [List regions](#tag/Region/operation/RegionService_ListRegions) request. Use this token to retrieve the next page of results.  **Note**: When paginating, all other parameters must match the original request.
func (r ApiRegionServiceListRegionsRequest) PageToken(pageToken string) ApiRegionServiceListRegionsRequest {
	r.pageToken = &pageToken
	return r
}

// The number of regions to skip before returning results. If the value exceeds the total number of clusters, the response is &#x60;200&#x60; with an empty list and no &#x60;nextPageToken&#x60;.
func (r ApiRegionServiceListRegionsRequest) Skip(skip int32) ApiRegionServiceListRegionsRequest {
	r.skip = &skip
	return r
}

func (r ApiRegionServiceListRegionsRequest) Execute() (*TidbCloudOpenApidedicatedv1beta1ListRegionsResponse, *http.Response, error) {
	return r.ApiService.RegionServiceListRegionsExecute(r)
}

/*
RegionServiceListRegions List regions

Lists the regions where you can create a cluster, across all supported cloud providers.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiRegionServiceListRegionsRequest
*/
func (a *RegionServiceAPIService) RegionServiceListRegions(ctx context.Context) ApiRegionServiceListRegionsRequest {
	return ApiRegionServiceListRegionsRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return TidbCloudOpenApidedicatedv1beta1ListRegionsResponse
func (a *RegionServiceAPIService) RegionServiceListRegionsExecute(r ApiRegionServiceListRegionsRequest) (*TidbCloudOpenApidedicatedv1beta1ListRegionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *TidbCloudOpenApidedicatedv1beta1ListRegionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RegionServiceAPIService.RegionServiceListRegions")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/regions"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.cloudProvider != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "cloudProvider", r.cloudProvider, "", "")
	}
	if r.projectId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "projectId", r.projectId, "", "")
	}
	if r.pageSize != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageSize", r.pageSize, "", "")
	} else {
		var defaultValue int32 = 20
		r.pageSize = &defaultValue
	}
	if r.pageToken != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageToken", r.pageToken, "", "")
	}
	if r.skip != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "skip", r.skip, "", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		var v GooglerpcStatus
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiRegionServiceShowCloudProvidersRequest struct {
	ctx        context.Context
	ApiService *RegionServiceAPIService
	projectId  *string
}

// The ID of the project. If not specified, the project ID of the default project is used.
func (r ApiRegionServiceShowCloudProvidersRequest) ProjectId(projectId string) ApiRegionServiceShowCloudProvidersRequest {
	r.projectId = &projectId
	return r
}

func (r ApiRegionServiceShowCloudProvidersRequest) Execute() (*V1beta1ShowCloudProvidersResponse, *http.Response, error) {
	return r.ApiService.RegionServiceShowCloudProvidersExecute(r)
}

/*
RegionServiceShowCloudProviders List cloud providers

Lists the cloud providers available for creating a cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiRegionServiceShowCloudProvidersRequest
*/
func (a *RegionServiceAPIService) RegionServiceShowCloudProviders(ctx context.Context) ApiRegionServiceShowCloudProvidersRequest {
	return ApiRegionServiceShowCloudProvidersRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return V1beta1ShowCloudProvidersResponse
func (a *RegionServiceAPIService) RegionServiceShowCloudProvidersExecute(r ApiRegionServiceShowCloudProvidersRequest) (*V1beta1ShowCloudProvidersResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *V1beta1ShowCloudProvidersResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RegionServiceAPIService.RegionServiceShowCloudProviders")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/regions:showCloudProviders"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.projectId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "projectId", r.projectId, "", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 400 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 429 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v GooglerpcStatus
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		var v GooglerpcStatus
		err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
		if err != nil {
			newErr.error = err.Error()
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
		newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
