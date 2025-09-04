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

// PrivateEndpointConnectionServiceAPIService PrivateEndpointConnectionServiceAPI service
type PrivateEndpointConnectionServiceAPIService service

type ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest struct {
	ctx                                      context.Context
	ApiService                               *PrivateEndpointConnectionServiceAPIService
	clusterId                                string
	privateEndpointConnectionTidbNodeGroupId string
	privateEndpointConnection                *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest
	validateOnly                             *bool
}

func (r ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) PrivateEndpointConnection(privateEndpointConnection PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	r.privateEndpointConnection = &privateEndpointConnection
	return r
}

// (Azure only) If set to &#x60;true&#x60;, the request is validated but not executed. Defaults to &#x60;false&#x60;.
func (r ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) ValidateOnly(validateOnly bool) ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	r.validateOnly = &validateOnly
	return r
}

func (r ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) Execute() (*Dedicatedv1beta1PrivateEndpointConnection, *http.Response, error) {
	return r.ApiService.PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionExecute(r)
}

/*
PrivateEndpointConnectionServiceCreatePrivateEndpointConnection Create a private endpoint connection

Creates a private endpoint connection for a TiDB node group. For more information, see the following documents:

- [Connect via Private Endpoint with AWS](https://docs.pingcap.com/tidbcloud/set-up-private-endpoint-connections/)
- [Connect via Private Endpoint with Azure](https://docs.pingcap.com/tidbcloud/set-up-private-endpoint-connections-on-azure/)
- [Connect via Private Endpoint with Google Cloud](https://docs.pingcap.com/tidbcloud/set-up-private-endpoint-connections-on-google-cloud/)

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId
	@param privateEndpointConnectionTidbNodeGroupId The ID of the TiDB group to which the private endpoint connection belongs.
	@return ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest
*/
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceCreatePrivateEndpointConnection(ctx context.Context, clusterId string, privateEndpointConnectionTidbNodeGroupId string) ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	return ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{
		ApiService:                               a,
		ctx:                                      ctx,
		clusterId:                                clusterId,
		privateEndpointConnectionTidbNodeGroupId: privateEndpointConnectionTidbNodeGroupId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1PrivateEndpointConnection
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionExecute(r ApiPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) (*Dedicatedv1beta1PrivateEndpointConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1PrivateEndpointConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PrivateEndpointConnectionServiceAPIService.PrivateEndpointConnectionServiceCreatePrivateEndpointConnection")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/tidbNodeGroups/{privateEndpointConnection.tidbNodeGroupId}/privateEndpointConnections"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"privateEndpointConnection.tidbNodeGroupId"+"}", url.PathEscape(parameterValueToString(r.privateEndpointConnectionTidbNodeGroupId, "privateEndpointConnectionTidbNodeGroupId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.privateEndpointConnection == nil {
		return localVarReturnValue, nil, reportError("privateEndpointConnection is required and must be specified")
	}

	if r.validateOnly != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "validateOnly", r.validateOnly, "", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

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
	// body params
	localVarPostBody = r.privateEndpointConnection
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

type ApiPrivateEndpointConnectionServiceDeletePrivateEndpointConnectionRequest struct {
	ctx                         context.Context
	ApiService                  *PrivateEndpointConnectionServiceAPIService
	clusterId                   string
	tidbNodeGroupId             string
	privateEndpointConnectionId string
}

func (r ApiPrivateEndpointConnectionServiceDeletePrivateEndpointConnectionRequest) Execute() (map[string]interface{}, *http.Response, error) {
	return r.ApiService.PrivateEndpointConnectionServiceDeletePrivateEndpointConnectionExecute(r)
}

/*
PrivateEndpointConnectionServiceDeletePrivateEndpointConnection Delete a private endpoint connection

Deletes a private endpoint connection for a TiDB node group.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId The ID of the cluster from which to delete the private endpoint connection.
	@param tidbNodeGroupId The ID of the TiDB node group from which to delete the private endpoint connection.
	@param privateEndpointConnectionId The ID of the private endpoint connection to delete.
	@return ApiPrivateEndpointConnectionServiceDeletePrivateEndpointConnectionRequest
*/
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceDeletePrivateEndpointConnection(ctx context.Context, clusterId string, tidbNodeGroupId string, privateEndpointConnectionId string) ApiPrivateEndpointConnectionServiceDeletePrivateEndpointConnectionRequest {
	return ApiPrivateEndpointConnectionServiceDeletePrivateEndpointConnectionRequest{
		ApiService:                  a,
		ctx:                         ctx,
		clusterId:                   clusterId,
		tidbNodeGroupId:             tidbNodeGroupId,
		privateEndpointConnectionId: privateEndpointConnectionId,
	}
}

// Execute executes the request
//
//	@return map[string]interface{}
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceDeletePrivateEndpointConnectionExecute(r ApiPrivateEndpointConnectionServiceDeletePrivateEndpointConnectionRequest) (map[string]interface{}, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue map[string]interface{}
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PrivateEndpointConnectionServiceAPIService.PrivateEndpointConnectionServiceDeletePrivateEndpointConnection")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/privateEndpointConnections/{privateEndpointConnectionId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"tidbNodeGroupId"+"}", url.PathEscape(parameterValueToString(r.tidbNodeGroupId, "tidbNodeGroupId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"privateEndpointConnectionId"+"}", url.PathEscape(parameterValueToString(r.privateEndpointConnectionId, "privateEndpointConnectionId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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

type ApiPrivateEndpointConnectionServiceGetPrivateEndpointConnectionRequest struct {
	ctx                         context.Context
	ApiService                  *PrivateEndpointConnectionServiceAPIService
	clusterId                   string
	tidbNodeGroupId             string
	privateEndpointConnectionId string
}

func (r ApiPrivateEndpointConnectionServiceGetPrivateEndpointConnectionRequest) Execute() (*Dedicatedv1beta1PrivateEndpointConnection, *http.Response, error) {
	return r.ApiService.PrivateEndpointConnectionServiceGetPrivateEndpointConnectionExecute(r)
}

/*
PrivateEndpointConnectionServiceGetPrivateEndpointConnection Get a private endpoint connection

Retrieves a private endpoint connection for a TiDB node group.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId The ID of the cluster for which to get the private endpoint connection.
	@param tidbNodeGroupId The ID of the TiDB node group for which to get the private endpoint connection.
	@param privateEndpointConnectionId The ID of the private endpoint connection to get.
	@return ApiPrivateEndpointConnectionServiceGetPrivateEndpointConnectionRequest
*/
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceGetPrivateEndpointConnection(ctx context.Context, clusterId string, tidbNodeGroupId string, privateEndpointConnectionId string) ApiPrivateEndpointConnectionServiceGetPrivateEndpointConnectionRequest {
	return ApiPrivateEndpointConnectionServiceGetPrivateEndpointConnectionRequest{
		ApiService:                  a,
		ctx:                         ctx,
		clusterId:                   clusterId,
		tidbNodeGroupId:             tidbNodeGroupId,
		privateEndpointConnectionId: privateEndpointConnectionId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1PrivateEndpointConnection
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceGetPrivateEndpointConnectionExecute(r ApiPrivateEndpointConnectionServiceGetPrivateEndpointConnectionRequest) (*Dedicatedv1beta1PrivateEndpointConnection, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1PrivateEndpointConnection
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PrivateEndpointConnectionServiceAPIService.PrivateEndpointConnectionServiceGetPrivateEndpointConnection")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/privateEndpointConnections/{privateEndpointConnectionId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"tidbNodeGroupId"+"}", url.PathEscape(parameterValueToString(r.tidbNodeGroupId, "tidbNodeGroupId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"privateEndpointConnectionId"+"}", url.PathEscape(parameterValueToString(r.privateEndpointConnectionId, "privateEndpointConnectionId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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

type ApiPrivateEndpointConnectionServiceGetPrivateLinkServiceRequest struct {
	ctx             context.Context
	ApiService      *PrivateEndpointConnectionServiceAPIService
	clusterId       string
	tidbNodeGroupId string
}

func (r ApiPrivateEndpointConnectionServiceGetPrivateLinkServiceRequest) Execute() (*Dedicatedv1beta1PrivateLinkService, *http.Response, error) {
	return r.ApiService.PrivateEndpointConnectionServiceGetPrivateLinkServiceExecute(r)
}

/*
PrivateEndpointConnectionServiceGetPrivateLinkService Get private link service for a TiDB node group

Retrieves the private link service configuration for a [TiDB node group](https://docs.pingcap.com/tidbcloud/tidb-node-group-overview/). You can the private link service to establish private endpoint connections to your TiDB cluster.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId The ID of the cluster for which to get the private link service.
	@param tidbNodeGroupId The ID of the TiDB node group for which to get the private link service.
	@return ApiPrivateEndpointConnectionServiceGetPrivateLinkServiceRequest
*/
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceGetPrivateLinkService(ctx context.Context, clusterId string, tidbNodeGroupId string) ApiPrivateEndpointConnectionServiceGetPrivateLinkServiceRequest {
	return ApiPrivateEndpointConnectionServiceGetPrivateLinkServiceRequest{
		ApiService:      a,
		ctx:             ctx,
		clusterId:       clusterId,
		tidbNodeGroupId: tidbNodeGroupId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1PrivateLinkService
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceGetPrivateLinkServiceExecute(r ApiPrivateEndpointConnectionServiceGetPrivateLinkServiceRequest) (*Dedicatedv1beta1PrivateLinkService, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1PrivateLinkService
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PrivateEndpointConnectionServiceAPIService.PrivateEndpointConnectionServiceGetPrivateLinkService")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/privateLinkService"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"tidbNodeGroupId"+"}", url.PathEscape(parameterValueToString(r.tidbNodeGroupId, "tidbNodeGroupId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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

type ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest struct {
	ctx             context.Context
	ApiService      *PrivateEndpointConnectionServiceAPIService
	clusterId       string
	tidbNodeGroupId string
	projectId       *string
	cloudProvider   *PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter
	pageSize        *int32
	pageToken       *string
	skip            *int32
}

// Pass this project_id only when listing under a project across different clusters and TiDB groups. Both cluster_id and tidb_node_group_id should be set to \&quot;-\&quot; in this case. If unspecified, the project ID of the cluster is used.
func (r ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest) ProjectId(projectId string) ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest {
	r.projectId = &projectId
	return r
}

// If specified, only private link connections of the specified cloud provider will be returned.   - aws: Amazon Web Services.  - gcp: Google Cloud Platform.  - azure: Microsoft Azure.  - alicloud: Alibaba Cloud
func (r ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest) CloudProvider(cloudProvider PrivateEndpointConnectionServiceListPrivateEndpointConnectionsCloudProviderParameter) ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest {
	r.cloudProvider = &cloudProvider
	return r
}

// The maximum number of private link connections to return. The service may return fewer than this value. If unspecified, at most 10 private link connections will be returned. The maximum value is 100; values above 100 will be coerced to 100.
func (r ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest) PageSize(pageSize int32) ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest {
	r.pageSize = &pageSize
	return r
}

// A page token, received from a previous &#x60;ListPrivateLinkConnections&#x60; call. Provide this to retrieve the subsequent page.
func (r ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest) PageToken(pageToken string) ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest {
	r.pageToken = &pageToken
	return r
}

// The number of individual resources to skip before starting to return results. If the skip value causes the cursor to move past the end of the collection, the response will be 200 OK with an empty result set and no next_page_token.
func (r ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest) Skip(skip int32) ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest {
	r.skip = &skip
	return r
}

func (r ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest) Execute() (*Dedicatedv1beta1ListPrivateEndpointConnectionsResponse, *http.Response, error) {
	return r.ApiService.PrivateEndpointConnectionServiceListPrivateEndpointConnectionsExecute(r)
}

/*
PrivateEndpointConnectionServiceListPrivateEndpointConnections List private endpoint connections

Lists private endpoint connections for a cluster and TiDB node group. To list private link connections across different clusters and TiDB groups, set cluster_id and tidb_node_group_id to "-". If you want to list private link connections under a project across different clusters and TiDB groups, set both cluster_id and tidb_node_group_id to "-", and pass the project_id in the query parameter.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId The id of the cluster for which to list private link connections. To list private link connections across different clusters and TiDB groups, set as `-``.
	@param tidbNodeGroupId To list private link connections across different clusters and TiDB groups, set tidb_node_group_id to \"-\".
	@return ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest
*/
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceListPrivateEndpointConnections(ctx context.Context, clusterId string, tidbNodeGroupId string) ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest {
	return ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest{
		ApiService:      a,
		ctx:             ctx,
		clusterId:       clusterId,
		tidbNodeGroupId: tidbNodeGroupId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1ListPrivateEndpointConnectionsResponse
func (a *PrivateEndpointConnectionServiceAPIService) PrivateEndpointConnectionServiceListPrivateEndpointConnectionsExecute(r ApiPrivateEndpointConnectionServiceListPrivateEndpointConnectionsRequest) (*Dedicatedv1beta1ListPrivateEndpointConnectionsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1ListPrivateEndpointConnectionsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "PrivateEndpointConnectionServiceAPIService.PrivateEndpointConnectionServiceListPrivateEndpointConnections")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/tidbNodeGroups/{tidbNodeGroupId}/privateEndpointConnections"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"tidbNodeGroupId"+"}", url.PathEscape(parameterValueToString(r.tidbNodeGroupId, "tidbNodeGroupId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.projectId != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "projectId", r.projectId, "", "")
	}
	if r.cloudProvider != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "cloudProvider", r.cloudProvider, "", "")
	}
	if r.pageSize != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageSize", r.pageSize, "", "")
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
