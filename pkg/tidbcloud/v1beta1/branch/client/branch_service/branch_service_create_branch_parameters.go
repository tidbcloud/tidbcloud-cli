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

	"tidbcloud-cli/pkg/tidbcloud/v1beta1/branch/models"
)

// NewBranchServiceCreateBranchParams creates a new BranchServiceCreateBranchParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewBranchServiceCreateBranchParams() *BranchServiceCreateBranchParams {
	return &BranchServiceCreateBranchParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewBranchServiceCreateBranchParamsWithTimeout creates a new BranchServiceCreateBranchParams object
// with the ability to set a timeout on a request.
func NewBranchServiceCreateBranchParamsWithTimeout(timeout time.Duration) *BranchServiceCreateBranchParams {
	return &BranchServiceCreateBranchParams{
		timeout: timeout,
	}
}

// NewBranchServiceCreateBranchParamsWithContext creates a new BranchServiceCreateBranchParams object
// with the ability to set a context for a request.
func NewBranchServiceCreateBranchParamsWithContext(ctx context.Context) *BranchServiceCreateBranchParams {
	return &BranchServiceCreateBranchParams{
		Context: ctx,
	}
}

// NewBranchServiceCreateBranchParamsWithHTTPClient creates a new BranchServiceCreateBranchParams object
// with the ability to set a custom HTTPClient for a request.
func NewBranchServiceCreateBranchParamsWithHTTPClient(client *http.Client) *BranchServiceCreateBranchParams {
	return &BranchServiceCreateBranchParams{
		HTTPClient: client,
	}
}

/*
BranchServiceCreateBranchParams contains all the parameters to send to the API endpoint

	for the branch service create branch operation.

	Typically these are written to a http.Request.
*/
type BranchServiceCreateBranchParams struct {

	/* Branch.

	   Required. The resource being created
	*/
	Branch *models.V1beta1Branch

	/* ClusterID.

	   Required. The cluster ID of the branch
	*/
	ClusterID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the branch service create branch params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *BranchServiceCreateBranchParams) WithDefaults() *BranchServiceCreateBranchParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the branch service create branch params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *BranchServiceCreateBranchParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the branch service create branch params
func (o *BranchServiceCreateBranchParams) WithTimeout(timeout time.Duration) *BranchServiceCreateBranchParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the branch service create branch params
func (o *BranchServiceCreateBranchParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the branch service create branch params
func (o *BranchServiceCreateBranchParams) WithContext(ctx context.Context) *BranchServiceCreateBranchParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the branch service create branch params
func (o *BranchServiceCreateBranchParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the branch service create branch params
func (o *BranchServiceCreateBranchParams) WithHTTPClient(client *http.Client) *BranchServiceCreateBranchParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the branch service create branch params
func (o *BranchServiceCreateBranchParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBranch adds the branch to the branch service create branch params
func (o *BranchServiceCreateBranchParams) WithBranch(branch *models.V1beta1Branch) *BranchServiceCreateBranchParams {
	o.SetBranch(branch)
	return o
}

// SetBranch adds the branch to the branch service create branch params
func (o *BranchServiceCreateBranchParams) SetBranch(branch *models.V1beta1Branch) {
	o.Branch = branch
}

// WithClusterID adds the clusterID to the branch service create branch params
func (o *BranchServiceCreateBranchParams) WithClusterID(clusterID string) *BranchServiceCreateBranchParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the branch service create branch params
func (o *BranchServiceCreateBranchParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WriteToRequest writes these params to a swagger request
func (o *BranchServiceCreateBranchParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Branch != nil {
		if err := r.SetBodyParam(o.Branch); err != nil {
			return err
		}
	}

	// path param clusterId
	if err := r.SetPathParam("clusterId", o.ClusterID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
