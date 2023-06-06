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

// NewDeleteUserParams creates a new DeleteUserParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDeleteUserParams() *DeleteUserParams {
	return &DeleteUserParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteUserParamsWithTimeout creates a new DeleteUserParams object
// with the ability to set a timeout on a request.
func NewDeleteUserParamsWithTimeout(timeout time.Duration) *DeleteUserParams {
	return &DeleteUserParams{
		timeout: timeout,
	}
}

// NewDeleteUserParamsWithContext creates a new DeleteUserParams object
// with the ability to set a context for a request.
func NewDeleteUserParamsWithContext(ctx context.Context) *DeleteUserParams {
	return &DeleteUserParams{
		Context: ctx,
	}
}

// NewDeleteUserParamsWithHTTPClient creates a new DeleteUserParams object
// with the ability to set a custom HTTPClient for a request.
func NewDeleteUserParamsWithHTTPClient(client *http.Client) *DeleteUserParams {
	return &DeleteUserParams{
		HTTPClient: client,
	}
}

/*
DeleteUserParams contains all the parameters to send to the API endpoint

	for the delete user operation.

	Typically these are written to a http.Request.
*/
type DeleteUserParams struct {

	/* BranchName.

	   The name of creating branch.
	*/
	BranchName string

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

	// Username.
	Username string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the delete user params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteUserParams) WithDefaults() *DeleteUserParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the delete user params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DeleteUserParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the delete user params
func (o *DeleteUserParams) WithTimeout(timeout time.Duration) *DeleteUserParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete user params
func (o *DeleteUserParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete user params
func (o *DeleteUserParams) WithContext(ctx context.Context) *DeleteUserParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete user params
func (o *DeleteUserParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete user params
func (o *DeleteUserParams) WithHTTPClient(client *http.Client) *DeleteUserParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete user params
func (o *DeleteUserParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBranchName adds the branchName to the delete user params
func (o *DeleteUserParams) WithBranchName(branchName string) *DeleteUserParams {
	o.SetBranchName(branchName)
	return o
}

// SetBranchName adds the branchName to the delete user params
func (o *DeleteUserParams) SetBranchName(branchName string) {
	o.BranchName = branchName
}

// WithClusterID adds the clusterID to the delete user params
func (o *DeleteUserParams) WithClusterID(clusterID string) *DeleteUserParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the delete user params
func (o *DeleteUserParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithProjectID adds the projectID to the delete user params
func (o *DeleteUserParams) WithProjectID(projectID string) *DeleteUserParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the delete user params
func (o *DeleteUserParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WithUsername adds the username to the delete user params
func (o *DeleteUserParams) WithUsername(username string) *DeleteUserParams {
	o.SetUsername(username)
	return o
}

// SetUsername adds the username to the delete user params
func (o *DeleteUserParams) SetUsername(username string) {
	o.Username = username
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteUserParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param branch_name
	if err := r.SetPathParam("branch_name", o.BranchName); err != nil {
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

	// path param username
	if err := r.SetPathParam("username", o.Username); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
