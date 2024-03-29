// Code generated by go-swagger; DO NOT EDIT.

package serverless_service

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

// NewServerlessServicePartialUpdateClusterParams creates a new ServerlessServicePartialUpdateClusterParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewServerlessServicePartialUpdateClusterParams() *ServerlessServicePartialUpdateClusterParams {
	return &ServerlessServicePartialUpdateClusterParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewServerlessServicePartialUpdateClusterParamsWithTimeout creates a new ServerlessServicePartialUpdateClusterParams object
// with the ability to set a timeout on a request.
func NewServerlessServicePartialUpdateClusterParamsWithTimeout(timeout time.Duration) *ServerlessServicePartialUpdateClusterParams {
	return &ServerlessServicePartialUpdateClusterParams{
		timeout: timeout,
	}
}

// NewServerlessServicePartialUpdateClusterParamsWithContext creates a new ServerlessServicePartialUpdateClusterParams object
// with the ability to set a context for a request.
func NewServerlessServicePartialUpdateClusterParamsWithContext(ctx context.Context) *ServerlessServicePartialUpdateClusterParams {
	return &ServerlessServicePartialUpdateClusterParams{
		Context: ctx,
	}
}

// NewServerlessServicePartialUpdateClusterParamsWithHTTPClient creates a new ServerlessServicePartialUpdateClusterParams object
// with the ability to set a custom HTTPClient for a request.
func NewServerlessServicePartialUpdateClusterParamsWithHTTPClient(client *http.Client) *ServerlessServicePartialUpdateClusterParams {
	return &ServerlessServicePartialUpdateClusterParams{
		HTTPClient: client,
	}
}

/*
ServerlessServicePartialUpdateClusterParams contains all the parameters to send to the API endpoint

	for the serverless service partial update cluster operation.

	Typically these are written to a http.Request.
*/
type ServerlessServicePartialUpdateClusterParams struct {

	// Body.
	Body ServerlessServicePartialUpdateClusterBody

	/* ClusterClusterID.

	   Output_only. The unique ID of the cluster.
	*/
	ClusterClusterID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the serverless service partial update cluster params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ServerlessServicePartialUpdateClusterParams) WithDefaults() *ServerlessServicePartialUpdateClusterParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the serverless service partial update cluster params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ServerlessServicePartialUpdateClusterParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) WithTimeout(timeout time.Duration) *ServerlessServicePartialUpdateClusterParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) WithContext(ctx context.Context) *ServerlessServicePartialUpdateClusterParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) WithHTTPClient(client *http.Client) *ServerlessServicePartialUpdateClusterParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) WithBody(body ServerlessServicePartialUpdateClusterBody) *ServerlessServicePartialUpdateClusterParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) SetBody(body ServerlessServicePartialUpdateClusterBody) {
	o.Body = body
}

// WithClusterClusterID adds the clusterClusterID to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) WithClusterClusterID(clusterClusterID string) *ServerlessServicePartialUpdateClusterParams {
	o.SetClusterClusterID(clusterClusterID)
	return o
}

// SetClusterClusterID adds the clusterClusterId to the serverless service partial update cluster params
func (o *ServerlessServicePartialUpdateClusterParams) SetClusterClusterID(clusterClusterID string) {
	o.ClusterClusterID = clusterClusterID
}

// WriteToRequest writes these params to a swagger request
func (o *ServerlessServicePartialUpdateClusterParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	// path param cluster.clusterId
	if err := r.SetPathParam("cluster.clusterId", o.ClusterClusterID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
