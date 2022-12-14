// Code generated by go-swagger; DO NOT EDIT.

package import_service

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

// NewGenerateUploadURLParams creates a new GenerateUploadURLParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGenerateUploadURLParams() *GenerateUploadURLParams {
	return &GenerateUploadURLParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGenerateUploadURLParamsWithTimeout creates a new GenerateUploadURLParams object
// with the ability to set a timeout on a request.
func NewGenerateUploadURLParamsWithTimeout(timeout time.Duration) *GenerateUploadURLParams {
	return &GenerateUploadURLParams{
		timeout: timeout,
	}
}

// NewGenerateUploadURLParamsWithContext creates a new GenerateUploadURLParams object
// with the ability to set a context for a request.
func NewGenerateUploadURLParamsWithContext(ctx context.Context) *GenerateUploadURLParams {
	return &GenerateUploadURLParams{
		Context: ctx,
	}
}

// NewGenerateUploadURLParamsWithHTTPClient creates a new GenerateUploadURLParams object
// with the ability to set a custom HTTPClient for a request.
func NewGenerateUploadURLParamsWithHTTPClient(client *http.Client) *GenerateUploadURLParams {
	return &GenerateUploadURLParams{
		HTTPClient: client,
	}
}

/*
GenerateUploadURLParams contains all the parameters to send to the API endpoint

	for the generate upload URL operation.

	Typically these are written to a http.Request.
*/
type GenerateUploadURLParams struct {

	// Body.
	Body GenerateUploadURLBody

	/* ClusterID.

	   The ID of the cluster.

	   Format: uint64
	*/
	ClusterID string

	/* ProjectID.

	   The ID of the project.

	   Format: uint64
	*/
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the generate upload URL params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GenerateUploadURLParams) WithDefaults() *GenerateUploadURLParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the generate upload URL params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GenerateUploadURLParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the generate upload URL params
func (o *GenerateUploadURLParams) WithTimeout(timeout time.Duration) *GenerateUploadURLParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the generate upload URL params
func (o *GenerateUploadURLParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the generate upload URL params
func (o *GenerateUploadURLParams) WithContext(ctx context.Context) *GenerateUploadURLParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the generate upload URL params
func (o *GenerateUploadURLParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the generate upload URL params
func (o *GenerateUploadURLParams) WithHTTPClient(client *http.Client) *GenerateUploadURLParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the generate upload URL params
func (o *GenerateUploadURLParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the generate upload URL params
func (o *GenerateUploadURLParams) WithBody(body GenerateUploadURLBody) *GenerateUploadURLParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the generate upload URL params
func (o *GenerateUploadURLParams) SetBody(body GenerateUploadURLBody) {
	o.Body = body
}

// WithClusterID adds the clusterID to the generate upload URL params
func (o *GenerateUploadURLParams) WithClusterID(clusterID string) *GenerateUploadURLParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the generate upload URL params
func (o *GenerateUploadURLParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithProjectID adds the projectID to the generate upload URL params
func (o *GenerateUploadURLParams) WithProjectID(projectID string) *GenerateUploadURLParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the generate upload URL params
func (o *GenerateUploadURLParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *GenerateUploadURLParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param project_id
	if err := r.SetPathParam("project_id", o.ProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
