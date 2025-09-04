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
	"time"
)

// DatabaseAuditLogServiceAPIService DatabaseAuditLogServiceAPI service
type DatabaseAuditLogServiceAPIService service

type ApiDatabaseAuditLogServiceCreateAuditLogFilterRuleRequest struct {
	ctx                         context.Context
	ApiService                  *DatabaseAuditLogServiceAPIService
	auditLogFilterRuleClusterId string
	auditLogFilterRule          *DatabaseAuditLogServiceCreateAuditLogFilterRuleRequest
}

func (r ApiDatabaseAuditLogServiceCreateAuditLogFilterRuleRequest) AuditLogFilterRule(auditLogFilterRule DatabaseAuditLogServiceCreateAuditLogFilterRuleRequest) ApiDatabaseAuditLogServiceCreateAuditLogFilterRuleRequest {
	r.auditLogFilterRule = &auditLogFilterRule
	return r
}

func (r ApiDatabaseAuditLogServiceCreateAuditLogFilterRuleRequest) Execute() (*Dedicatedv1beta1AuditLogFilterRule, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceCreateAuditLogFilterRuleExecute(r)
}

/*
DatabaseAuditLogServiceCreateAuditLogFilterRule Create an audit log filter rule

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param auditLogFilterRuleClusterId
	@return ApiDatabaseAuditLogServiceCreateAuditLogFilterRuleRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceCreateAuditLogFilterRule(ctx context.Context, auditLogFilterRuleClusterId string) ApiDatabaseAuditLogServiceCreateAuditLogFilterRuleRequest {
	return ApiDatabaseAuditLogServiceCreateAuditLogFilterRuleRequest{
		ApiService:                  a,
		ctx:                         ctx,
		auditLogFilterRuleClusterId: auditLogFilterRuleClusterId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1AuditLogFilterRule
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceCreateAuditLogFilterRuleExecute(r ApiDatabaseAuditLogServiceCreateAuditLogFilterRuleRequest) (*Dedicatedv1beta1AuditLogFilterRule, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1AuditLogFilterRule
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceCreateAuditLogFilterRule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{auditLogFilterRule.clusterId}/auditLogFilterRules"
	localVarPath = strings.Replace(localVarPath, "{"+"auditLogFilterRule.clusterId"+"}", url.PathEscape(parameterValueToString(r.auditLogFilterRuleClusterId, "auditLogFilterRuleClusterId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.auditLogFilterRule == nil {
		return localVarReturnValue, nil, reportError("auditLogFilterRule is required and must be specified")
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
	localVarPostBody = r.auditLogFilterRule
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

type ApiDatabaseAuditLogServiceDeleteAuditLogFilterRuleRequest struct {
	ctx                  context.Context
	ApiService           *DatabaseAuditLogServiceAPIService
	clusterId            string
	auditLogFilterRuleId string
}

func (r ApiDatabaseAuditLogServiceDeleteAuditLogFilterRuleRequest) Execute() (map[string]interface{}, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceDeleteAuditLogFilterRuleExecute(r)
}

/*
DatabaseAuditLogServiceDeleteAuditLogFilterRule Delete an audit log filter rule

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId
	@param auditLogFilterRuleId
	@return ApiDatabaseAuditLogServiceDeleteAuditLogFilterRuleRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceDeleteAuditLogFilterRule(ctx context.Context, clusterId string, auditLogFilterRuleId string) ApiDatabaseAuditLogServiceDeleteAuditLogFilterRuleRequest {
	return ApiDatabaseAuditLogServiceDeleteAuditLogFilterRuleRequest{
		ApiService:           a,
		ctx:                  ctx,
		clusterId:            clusterId,
		auditLogFilterRuleId: auditLogFilterRuleId,
	}
}

// Execute executes the request
//
//	@return map[string]interface{}
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceDeleteAuditLogFilterRuleExecute(r ApiDatabaseAuditLogServiceDeleteAuditLogFilterRuleRequest) (map[string]interface{}, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue map[string]interface{}
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceDeleteAuditLogFilterRule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/auditLogFilterRules/{auditLogFilterRuleId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"auditLogFilterRuleId"+"}", url.PathEscape(parameterValueToString(r.auditLogFilterRuleId, "auditLogFilterRuleId")), -1)

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

type ApiDatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressRequest struct {
	ctx             context.Context
	ApiService      *DatabaseAuditLogServiceAPIService
	clusterId       string
	auditLogFileUri *string
}

func (r ApiDatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressRequest) AuditLogFileUri(auditLogFileUri string) ApiDatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressRequest {
	r.auditLogFileUri = &auditLogFileUri
	return r
}

func (r ApiDatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressRequest) Execute() (*Dedicatedv1beta1GenerateAuditLogFileDownloadAddressResponse, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressExecute(r)
}

/*
DatabaseAuditLogServiceGenerateAuditLogFileDownloadAddress Generate the download address for an audit log file, the address have an 15 minutes expiration time

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId
	@return ApiDatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceGenerateAuditLogFileDownloadAddress(ctx context.Context, clusterId string) ApiDatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressRequest {
	return ApiDatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1GenerateAuditLogFileDownloadAddressResponse
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressExecute(r ApiDatabaseAuditLogServiceGenerateAuditLogFileDownloadAddressRequest) (*Dedicatedv1beta1GenerateAuditLogFileDownloadAddressResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1GenerateAuditLogFileDownloadAddressResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceGenerateAuditLogFileDownloadAddress")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/auditLogConfig:generateAuditLogFileDownloadAddress"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.auditLogFileUri == nil {
		return localVarReturnValue, nil, reportError("auditLogFileUri is required and must be specified")
	}

	parameterAddToHeaderOrQuery(localVarQueryParams, "auditLogFileUri", r.auditLogFileUri, "", "")
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

type ApiDatabaseAuditLogServiceGetAuditLogConfigRequest struct {
	ctx        context.Context
	ApiService *DatabaseAuditLogServiceAPIService
	clusterId  string
}

func (r ApiDatabaseAuditLogServiceGetAuditLogConfigRequest) Execute() (*Dedicatedv1beta1AuditLogConfig, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceGetAuditLogConfigExecute(r)
}

/*
DatabaseAuditLogServiceGetAuditLogConfig Get the audit log config of a cluster

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId
	@return ApiDatabaseAuditLogServiceGetAuditLogConfigRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceGetAuditLogConfig(ctx context.Context, clusterId string) ApiDatabaseAuditLogServiceGetAuditLogConfigRequest {
	return ApiDatabaseAuditLogServiceGetAuditLogConfigRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1AuditLogConfig
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceGetAuditLogConfigExecute(r ApiDatabaseAuditLogServiceGetAuditLogConfigRequest) (*Dedicatedv1beta1AuditLogConfig, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1AuditLogConfig
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceGetAuditLogConfig")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/auditLogConfig"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)

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

type ApiDatabaseAuditLogServiceGetAuditLogFilterRuleRequest struct {
	ctx                  context.Context
	ApiService           *DatabaseAuditLogServiceAPIService
	clusterId            string
	auditLogFilterRuleId string
}

func (r ApiDatabaseAuditLogServiceGetAuditLogFilterRuleRequest) Execute() (*Dedicatedv1beta1AuditLogFilterRule, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceGetAuditLogFilterRuleExecute(r)
}

/*
DatabaseAuditLogServiceGetAuditLogFilterRule Get an audit log filter rule

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId
	@param auditLogFilterRuleId
	@return ApiDatabaseAuditLogServiceGetAuditLogFilterRuleRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceGetAuditLogFilterRule(ctx context.Context, clusterId string, auditLogFilterRuleId string) ApiDatabaseAuditLogServiceGetAuditLogFilterRuleRequest {
	return ApiDatabaseAuditLogServiceGetAuditLogFilterRuleRequest{
		ApiService:           a,
		ctx:                  ctx,
		clusterId:            clusterId,
		auditLogFilterRuleId: auditLogFilterRuleId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1AuditLogFilterRule
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceGetAuditLogFilterRuleExecute(r ApiDatabaseAuditLogServiceGetAuditLogFilterRuleRequest) (*Dedicatedv1beta1AuditLogFilterRule, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1AuditLogFilterRule
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceGetAuditLogFilterRule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/auditLogFilterRules/{auditLogFilterRuleId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"auditLogFilterRuleId"+"}", url.PathEscape(parameterValueToString(r.auditLogFilterRuleId, "auditLogFilterRuleId")), -1)

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

type ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest struct {
	ctx        context.Context
	ApiService *DatabaseAuditLogServiceAPIService
	clusterId  string
	pageSize   *int32
	pageToken  *string
	skip       *int32
}

func (r ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest) PageSize(pageSize int32) ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest {
	r.pageSize = &pageSize
	return r
}

func (r ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest) PageToken(pageToken string) ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest {
	r.pageToken = &pageToken
	return r
}

// The number of individual resources to skip before starting to return results. If the skip value causes the cursor to move past the end of the collection, the response will be 200 OK with an empty result set and no next_page_token.
func (r ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest) Skip(skip int32) ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest {
	r.skip = &skip
	return r
}

func (r ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest) Execute() (*Dedicatedv1beta1ListAuditLogFilterRulesResponse, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceListAuditLogFilterRulesExecute(r)
}

/*
DatabaseAuditLogServiceListAuditLogFilterRules List audit log filter rules

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId
	@return ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceListAuditLogFilterRules(ctx context.Context, clusterId string) ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest {
	return ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1ListAuditLogFilterRulesResponse
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceListAuditLogFilterRulesExecute(r ApiDatabaseAuditLogServiceListAuditLogFilterRulesRequest) (*Dedicatedv1beta1ListAuditLogFilterRulesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1ListAuditLogFilterRulesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceListAuditLogFilterRules")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/auditLogFilterRules"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

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

type ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest struct {
	ctx        context.Context
	ApiService *DatabaseAuditLogServiceAPIService
	clusterId  string
	startTime  *time.Time
	endTime    *time.Time
	pageSize   *int32
	pageToken  *string
	skip       *int32
}

// the start time of the audit logs created time range.
func (r ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest) StartTime(startTime time.Time) ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest {
	r.startTime = &startTime
	return r
}

// the end time of the audit logs created time range.
func (r ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest) EndTime(endTime time.Time) ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest {
	r.endTime = &endTime
	return r
}

func (r ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest) PageSize(pageSize int32) ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest {
	r.pageSize = &pageSize
	return r
}

func (r ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest) PageToken(pageToken string) ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest {
	r.pageToken = &pageToken
	return r
}

// The number of individual resources to skip before starting to return results. If the skip value causes the cursor to move past the end of the collection, the response will be 200 OK with an empty result set and no next_page_token.
func (r ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest) Skip(skip int32) ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest {
	r.skip = &skip
	return r
}

func (r ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest) Execute() (*Dedicatedv1beta1QueryAuditLogFilesResponse, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceQueryAuditLogFilesExecute(r)
}

/*
DatabaseAuditLogServiceQueryAuditLogFiles Query audit log files

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId
	@return ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceQueryAuditLogFiles(ctx context.Context, clusterId string) ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest {
	return ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1QueryAuditLogFilesResponse
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceQueryAuditLogFilesExecute(r ApiDatabaseAuditLogServiceQueryAuditLogFilesRequest) (*Dedicatedv1beta1QueryAuditLogFilesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1QueryAuditLogFilesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceQueryAuditLogFiles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/auditLogConfig:queryAuditLogFiles"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.startTime != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "startTime", r.startTime, "", "")
	}
	if r.endTime != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "endTime", r.endTime, "", "")
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

type ApiDatabaseAuditLogServiceReplaceAuditLogFilterRuleRequest struct {
	ctx                                    context.Context
	ApiService                             *DatabaseAuditLogServiceAPIService
	auditLogFilterRuleClusterId            string
	auditLogFilterRuleAuditLogFilterRuleId string
	auditLogFilterRule                     *RequiredTheAuditLogFilterRuleToUpdate
}

// Required. The audit log filter rule to update.
func (r ApiDatabaseAuditLogServiceReplaceAuditLogFilterRuleRequest) AuditLogFilterRule(auditLogFilterRule RequiredTheAuditLogFilterRuleToUpdate) ApiDatabaseAuditLogServiceReplaceAuditLogFilterRuleRequest {
	r.auditLogFilterRule = &auditLogFilterRule
	return r
}

func (r ApiDatabaseAuditLogServiceReplaceAuditLogFilterRuleRequest) Execute() (*V1beta1ReplaceAuditLogFilterRuleResponse, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceReplaceAuditLogFilterRuleExecute(r)
}

/*
DatabaseAuditLogServiceReplaceAuditLogFilterRule Replace an audit log filter rule. All fields of the rule will be replaced with the provided values.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param auditLogFilterRuleClusterId
	@param auditLogFilterRuleAuditLogFilterRuleId
	@return ApiDatabaseAuditLogServiceReplaceAuditLogFilterRuleRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceReplaceAuditLogFilterRule(ctx context.Context, auditLogFilterRuleClusterId string, auditLogFilterRuleAuditLogFilterRuleId string) ApiDatabaseAuditLogServiceReplaceAuditLogFilterRuleRequest {
	return ApiDatabaseAuditLogServiceReplaceAuditLogFilterRuleRequest{
		ApiService:                             a,
		ctx:                                    ctx,
		auditLogFilterRuleClusterId:            auditLogFilterRuleClusterId,
		auditLogFilterRuleAuditLogFilterRuleId: auditLogFilterRuleAuditLogFilterRuleId,
	}
}

// Execute executes the request
//
//	@return V1beta1ReplaceAuditLogFilterRuleResponse
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceReplaceAuditLogFilterRuleExecute(r ApiDatabaseAuditLogServiceReplaceAuditLogFilterRuleRequest) (*V1beta1ReplaceAuditLogFilterRuleResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *V1beta1ReplaceAuditLogFilterRuleResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceReplaceAuditLogFilterRule")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{auditLogFilterRule.clusterId}/auditLogFilterRules/{auditLogFilterRule.auditLogFilterRuleId}:replaceAuditLogFilterRule"
	localVarPath = strings.Replace(localVarPath, "{"+"auditLogFilterRule.clusterId"+"}", url.PathEscape(parameterValueToString(r.auditLogFilterRuleClusterId, "auditLogFilterRuleClusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"auditLogFilterRule.auditLogFilterRuleId"+"}", url.PathEscape(parameterValueToString(r.auditLogFilterRuleAuditLogFilterRuleId, "auditLogFilterRuleAuditLogFilterRuleId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.auditLogFilterRule == nil {
		return localVarReturnValue, nil, reportError("auditLogFilterRule is required and must be specified")
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
	localVarPostBody = r.auditLogFilterRule
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

type ApiDatabaseAuditLogServiceShowObjectStorageAccessIamPrincipalRequest struct {
	ctx        context.Context
	ApiService *DatabaseAuditLogServiceAPIService
	clusterId  string
}

func (r ApiDatabaseAuditLogServiceShowObjectStorageAccessIamPrincipalRequest) Execute() (*V1beta1ShowObjectStorageAccessIamPrincipalResponse, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceShowObjectStorageAccessIamPrincipalExecute(r)
}

/*
DatabaseAuditLogServiceShowObjectStorageAccessIamPrincipal Show IAM principal of TiDB Cloud for accessing customer's object storage

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId
	@return ApiDatabaseAuditLogServiceShowObjectStorageAccessIamPrincipalRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceShowObjectStorageAccessIamPrincipal(ctx context.Context, clusterId string) ApiDatabaseAuditLogServiceShowObjectStorageAccessIamPrincipalRequest {
	return ApiDatabaseAuditLogServiceShowObjectStorageAccessIamPrincipalRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
	}
}

// Execute executes the request
//
//	@return V1beta1ShowObjectStorageAccessIamPrincipalResponse
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceShowObjectStorageAccessIamPrincipalExecute(r ApiDatabaseAuditLogServiceShowObjectStorageAccessIamPrincipalRequest) (*V1beta1ShowObjectStorageAccessIamPrincipalResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *V1beta1ShowObjectStorageAccessIamPrincipalResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceShowObjectStorageAccessIamPrincipal")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{clusterId}/auditLogConfig:showObjectStorageAccessIamPrincipal"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)

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

type ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest struct {
	ctx                     context.Context
	ApiService              *DatabaseAuditLogServiceAPIService
	auditLogConfigClusterId string
	auditLogConfig          *DatabaseAuditLogServiceUpdateAuditLogConfigRequest
	validateOnly            *bool
}

func (r ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest) AuditLogConfig(auditLogConfig DatabaseAuditLogServiceUpdateAuditLogConfigRequest) ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest {
	r.auditLogConfig = &auditLogConfig
	return r
}

func (r ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest) ValidateOnly(validateOnly bool) ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest {
	r.validateOnly = &validateOnly
	return r
}

func (r ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest) Execute() (*Dedicatedv1beta1AuditLogConfig, *http.Response, error) {
	return r.ApiService.DatabaseAuditLogServiceUpdateAuditLogConfigExecute(r)
}

/*
DatabaseAuditLogServiceUpdateAuditLogConfig Update the audit log config of a cluster

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param auditLogConfigClusterId
	@return ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest
*/
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceUpdateAuditLogConfig(ctx context.Context, auditLogConfigClusterId string) ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest {
	return ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest{
		ApiService:              a,
		ctx:                     ctx,
		auditLogConfigClusterId: auditLogConfigClusterId,
	}
}

// Execute executes the request
//
//	@return Dedicatedv1beta1AuditLogConfig
func (a *DatabaseAuditLogServiceAPIService) DatabaseAuditLogServiceUpdateAuditLogConfigExecute(r ApiDatabaseAuditLogServiceUpdateAuditLogConfigRequest) (*Dedicatedv1beta1AuditLogConfig, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Dedicatedv1beta1AuditLogConfig
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "DatabaseAuditLogServiceAPIService.DatabaseAuditLogServiceUpdateAuditLogConfig")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/clusters/{auditLogConfig.clusterId}/auditLogConfig"
	localVarPath = strings.Replace(localVarPath, "{"+"auditLogConfig.clusterId"+"}", url.PathEscape(parameterValueToString(r.auditLogConfigClusterId, "auditLogConfigClusterId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.auditLogConfig == nil {
		return localVarReturnValue, nil, reportError("auditLogConfig is required and must be specified")
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
	localVarPostBody = r.auditLogConfig
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
