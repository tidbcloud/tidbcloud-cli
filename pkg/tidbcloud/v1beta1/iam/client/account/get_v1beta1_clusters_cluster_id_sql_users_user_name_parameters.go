// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetV1beta1ClustersClusterIDSQLUsersUserNameParams creates a new GetV1beta1ClustersClusterIDSQLUsersUserNameParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetV1beta1ClustersClusterIDSQLUsersUserNameParams() *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	return &GetV1beta1ClustersClusterIDSQLUsersUserNameParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetV1beta1ClustersClusterIDSQLUsersUserNameParamsWithTimeout creates a new GetV1beta1ClustersClusterIDSQLUsersUserNameParams object
// with the ability to set a timeout on a request.
func NewGetV1beta1ClustersClusterIDSQLUsersUserNameParamsWithTimeout(timeout time.Duration) *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	return &GetV1beta1ClustersClusterIDSQLUsersUserNameParams{
		timeout: timeout,
	}
}

// NewGetV1beta1ClustersClusterIDSQLUsersUserNameParamsWithContext creates a new GetV1beta1ClustersClusterIDSQLUsersUserNameParams object
// with the ability to set a context for a request.
func NewGetV1beta1ClustersClusterIDSQLUsersUserNameParamsWithContext(ctx context.Context) *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	return &GetV1beta1ClustersClusterIDSQLUsersUserNameParams{
		Context: ctx,
	}
}

// NewGetV1beta1ClustersClusterIDSQLUsersUserNameParamsWithHTTPClient creates a new GetV1beta1ClustersClusterIDSQLUsersUserNameParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetV1beta1ClustersClusterIDSQLUsersUserNameParamsWithHTTPClient(client *http.Client) *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	return &GetV1beta1ClustersClusterIDSQLUsersUserNameParams{
		HTTPClient: client,
	}
}

/*
GetV1beta1ClustersClusterIDSQLUsersUserNameParams contains all the parameters to send to the API endpoint

	for the get v1beta1 clusters cluster ID SQL users user name operation.

	Typically these are written to a http.Request.
*/
type GetV1beta1ClustersClusterIDSQLUsersUserNameParams struct {

	/* ClusterID.

	   The id of the cluster.
	*/
	ClusterID string

	/* UserName.

	   The name of the sql user.
	*/
	UserName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get v1beta1 clusters cluster ID SQL users user name params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) WithDefaults() *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get v1beta1 clusters cluster ID SQL users user name params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) WithTimeout(timeout time.Duration) *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) WithContext(ctx context.Context) *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) WithHTTPClient(client *http.Client) *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) WithClusterID(clusterID string) *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithUserName adds the userName to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) WithUserName(userName string) *GetV1beta1ClustersClusterIDSQLUsersUserNameParams {
	o.SetUserName(userName)
	return o
}

// SetUserName adds the userName to the get v1beta1 clusters cluster ID SQL users user name params
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) SetUserName(userName string) {
	o.UserName = userName
}

// WriteToRequest writes these params to a swagger request
func (o *GetV1beta1ClustersClusterIDSQLUsersUserNameParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param clusterId
	if err := r.SetPathParam("clusterId", o.ClusterID); err != nil {
		return err
	}

	// path param userName
	if err := r.SetPathParam("userName", o.UserName); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
