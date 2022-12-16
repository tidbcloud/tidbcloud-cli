// Code generated by go-swagger; DO NOT EDIT.

package operations

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

// NewGetImportParams creates a new GetImportParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetImportParams() *GetImportParams {
	return &GetImportParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetImportParamsWithTimeout creates a new GetImportParams object
// with the ability to set a timeout on a request.
func NewGetImportParamsWithTimeout(timeout time.Duration) *GetImportParams {
	return &GetImportParams{
		timeout: timeout,
	}
}

// NewGetImportParamsWithContext creates a new GetImportParams object
// with the ability to set a context for a request.
func NewGetImportParamsWithContext(ctx context.Context) *GetImportParams {
	return &GetImportParams{
		Context: ctx,
	}
}

// NewGetImportParamsWithHTTPClient creates a new GetImportParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetImportParamsWithHTTPClient(client *http.Client) *GetImportParams {
	return &GetImportParams{
		HTTPClient: client,
	}
}

/*
GetImportParams contains all the parameters to send to the API endpoint

	for the get import operation.

	Typically these are written to a http.Request.
*/
type GetImportParams struct {

	// ClusterID.
	//
	// Format: uint64
	ClusterID string

	// ID.
	//
	// Format: uint64
	ID string

	// ProjectID.
	//
	// Format: uint64
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get import params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetImportParams) WithDefaults() *GetImportParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get import params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetImportParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get import params
func (o *GetImportParams) WithTimeout(timeout time.Duration) *GetImportParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get import params
func (o *GetImportParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get import params
func (o *GetImportParams) WithContext(ctx context.Context) *GetImportParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get import params
func (o *GetImportParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get import params
func (o *GetImportParams) WithHTTPClient(client *http.Client) *GetImportParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get import params
func (o *GetImportParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the get import params
func (o *GetImportParams) WithClusterID(clusterID string) *GetImportParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the get import params
func (o *GetImportParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithID adds the id to the get import params
func (o *GetImportParams) WithID(id string) *GetImportParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get import params
func (o *GetImportParams) SetID(id string) {
	o.ID = id
}

// WithProjectID adds the projectID to the get import params
func (o *GetImportParams) WithProjectID(projectID string) *GetImportParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the get import params
func (o *GetImportParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *GetImportParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param cluster_id
	if err := r.SetPathParam("cluster_id", o.ClusterID); err != nil {
		return err
	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	// path param project_id
	if err := r.SetPathParam("project_id", o.ProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}