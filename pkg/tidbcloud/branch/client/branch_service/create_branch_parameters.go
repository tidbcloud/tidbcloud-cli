// Code generated by go-swagger; DO NOT EDIT.

package branch_service

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

// NewCreateBranchParams creates a new CreateBranchParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateBranchParams() *CreateBranchParams {
	return &CreateBranchParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateBranchParamsWithTimeout creates a new CreateBranchParams object
// with the ability to set a timeout on a request.
func NewCreateBranchParamsWithTimeout(timeout time.Duration) *CreateBranchParams {
	return &CreateBranchParams{
		timeout: timeout,
	}
}

// NewCreateBranchParamsWithContext creates a new CreateBranchParams object
// with the ability to set a context for a request.
func NewCreateBranchParamsWithContext(ctx context.Context) *CreateBranchParams {
	return &CreateBranchParams{
		Context: ctx,
	}
}

// NewCreateBranchParamsWithHTTPClient creates a new CreateBranchParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateBranchParamsWithHTTPClient(client *http.Client) *CreateBranchParams {
	return &CreateBranchParams{
		HTTPClient: client,
	}
}

/*
CreateBranchParams contains all the parameters to send to the API endpoint

	for the create branch operation.

	Typically these are written to a http.Request.
*/
type CreateBranchParams struct {

	// Body.
	Body CreateBranchBody

	/* ClusterID.

	   The ID of the cluster.

	   Format: uint64
	*/
	ClusterID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create branch params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateBranchParams) WithDefaults() *CreateBranchParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create branch params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateBranchParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create branch params
func (o *CreateBranchParams) WithTimeout(timeout time.Duration) *CreateBranchParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create branch params
func (o *CreateBranchParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create branch params
func (o *CreateBranchParams) WithContext(ctx context.Context) *CreateBranchParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create branch params
func (o *CreateBranchParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create branch params
func (o *CreateBranchParams) WithHTTPClient(client *http.Client) *CreateBranchParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create branch params
func (o *CreateBranchParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create branch params
func (o *CreateBranchParams) WithBody(body CreateBranchBody) *CreateBranchParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create branch params
func (o *CreateBranchParams) SetBody(body CreateBranchBody) {
	o.Body = body
}

// WithClusterID adds the clusterID to the create branch params
func (o *CreateBranchParams) WithClusterID(clusterID string) *CreateBranchParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the create branch params
func (o *CreateBranchParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WriteToRequest writes these params to a swagger request
func (o *CreateBranchParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	// path param cluster_id
	if err := r.SetPathParam("cluster_id", o.ClusterID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
