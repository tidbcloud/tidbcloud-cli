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

// NewBranchServiceDeleteBranchParams creates a new BranchServiceDeleteBranchParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewBranchServiceDeleteBranchParams() *BranchServiceDeleteBranchParams {
	return &BranchServiceDeleteBranchParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewBranchServiceDeleteBranchParamsWithTimeout creates a new BranchServiceDeleteBranchParams object
// with the ability to set a timeout on a request.
func NewBranchServiceDeleteBranchParamsWithTimeout(timeout time.Duration) *BranchServiceDeleteBranchParams {
	return &BranchServiceDeleteBranchParams{
		timeout: timeout,
	}
}

// NewBranchServiceDeleteBranchParamsWithContext creates a new BranchServiceDeleteBranchParams object
// with the ability to set a context for a request.
func NewBranchServiceDeleteBranchParamsWithContext(ctx context.Context) *BranchServiceDeleteBranchParams {
	return &BranchServiceDeleteBranchParams{
		Context: ctx,
	}
}

// NewBranchServiceDeleteBranchParamsWithHTTPClient creates a new BranchServiceDeleteBranchParams object
// with the ability to set a custom HTTPClient for a request.
func NewBranchServiceDeleteBranchParamsWithHTTPClient(client *http.Client) *BranchServiceDeleteBranchParams {
	return &BranchServiceDeleteBranchParams{
		HTTPClient: client,
	}
}

/*
BranchServiceDeleteBranchParams contains all the parameters to send to the API endpoint

	for the branch service delete branch operation.

	Typically these are written to a http.Request.
*/
type BranchServiceDeleteBranchParams struct {

	/* BranchID.

	   Required. The branch ID
	*/
	BranchID string

	/* ClusterID.

	   Required. The cluster ID of the branch
	*/
	ClusterID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the branch service delete branch params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *BranchServiceDeleteBranchParams) WithDefaults() *BranchServiceDeleteBranchParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the branch service delete branch params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *BranchServiceDeleteBranchParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) WithTimeout(timeout time.Duration) *BranchServiceDeleteBranchParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) WithContext(ctx context.Context) *BranchServiceDeleteBranchParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) WithHTTPClient(client *http.Client) *BranchServiceDeleteBranchParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBranchID adds the branchID to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) WithBranchID(branchID string) *BranchServiceDeleteBranchParams {
	o.SetBranchID(branchID)
	return o
}

// SetBranchID adds the branchId to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) SetBranchID(branchID string) {
	o.BranchID = branchID
}

// WithClusterID adds the clusterID to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) WithClusterID(clusterID string) *BranchServiceDeleteBranchParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the branch service delete branch params
func (o *BranchServiceDeleteBranchParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WriteToRequest writes these params to a swagger request
func (o *BranchServiceDeleteBranchParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param branchId
	if err := r.SetPathParam("branchId", o.BranchID); err != nil {
		return err
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