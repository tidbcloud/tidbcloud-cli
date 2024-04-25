// Code generated by go-swagger; DO NOT EDIT.

package account

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new account API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for account API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	DeleteV1beta1ClustersClusterIDSQLUsersUserName(params *DeleteV1beta1ClustersClusterIDSQLUsersUserNameParams, opts ...ClientOption) (*DeleteV1beta1ClustersClusterIDSQLUsersUserNameOK, error)

	GetMspCustomers(params *GetMspCustomersParams, opts ...ClientOption) (*GetMspCustomersOK, error)

	GetMspCustomersCustomerOrgID(params *GetMspCustomersCustomerOrgIDParams, opts ...ClientOption) (*GetMspCustomersCustomerOrgIDOK, error)

	GetV1beta1ClustersClusterIDSQLUsers(params *GetV1beta1ClustersClusterIDSQLUsersParams, opts ...ClientOption) (*GetV1beta1ClustersClusterIDSQLUsersOK, error)

	GetV1beta1ClustersClusterIDSQLUsersUserName(params *GetV1beta1ClustersClusterIDSQLUsersUserNameParams, opts ...ClientOption) (*GetV1beta1ClustersClusterIDSQLUsersUserNameOK, error)

	GetV1beta1Projects(params *GetV1beta1ProjectsParams, opts ...ClientOption) (*GetV1beta1ProjectsOK, error)

	PatchV1beta1ClustersClusterIDSQLUsersUserName(params *PatchV1beta1ClustersClusterIDSQLUsersUserNameParams, opts ...ClientOption) (*PatchV1beta1ClustersClusterIDSQLUsersUserNameOK, error)

	PostCustomerSignupURL(params *PostCustomerSignupURLParams, opts ...ClientOption) (*PostCustomerSignupURLOK, error)

	PostV1beta1ClustersClusterIDSQLUsers(params *PostV1beta1ClustersClusterIDSQLUsersParams, opts ...ClientOption) (*PostV1beta1ClustersClusterIDSQLUsersOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
DeleteV1beta1ClustersClusterIDSQLUsersUserName deletes one sql user

This endpoint delete the sql user by user name.
*/
func (a *Client) DeleteV1beta1ClustersClusterIDSQLUsersUserName(params *DeleteV1beta1ClustersClusterIDSQLUsersUserNameParams, opts ...ClientOption) (*DeleteV1beta1ClustersClusterIDSQLUsersUserNameOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteV1beta1ClustersClusterIDSQLUsersUserNameParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "DeleteV1beta1ClustersClusterIDSQLUsersUserName",
		Method:             "DELETE",
		PathPattern:        "/v1beta1/clusters/{clusterId}/sqlUsers/{userName}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteV1beta1ClustersClusterIDSQLUsersUserNameReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DeleteV1beta1ClustersClusterIDSQLUsersUserNameOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for DeleteV1beta1ClustersClusterIDSQLUsersUserName: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetMspCustomers gets a list of m s p customers

This endpoint returns a list of MSP customers.
*/
func (a *Client) GetMspCustomers(params *GetMspCustomersParams, opts ...ClientOption) (*GetMspCustomersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMspCustomersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetMspCustomers",
		Method:             "GET",
		PathPattern:        "/mspCustomers",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetMspCustomersReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetMspCustomersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetMspCustomers: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetMspCustomersCustomerOrgID retrieves a single m s p customer

This endpoint retrieves a single MSP customer by their customer org ID.
*/
func (a *Client) GetMspCustomersCustomerOrgID(params *GetMspCustomersCustomerOrgIDParams, opts ...ClientOption) (*GetMspCustomersCustomerOrgIDOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetMspCustomersCustomerOrgIDParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetMspCustomersCustomerOrgID",
		Method:             "GET",
		PathPattern:        "/mspCustomers/{customerOrgId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetMspCustomersCustomerOrgIDReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetMspCustomersCustomerOrgIDOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetMspCustomersCustomerOrgID: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetV1beta1ClustersClusterIDSQLUsers gets all sql users

This endpoint retrieves all sql users  in the cluster.
*/
func (a *Client) GetV1beta1ClustersClusterIDSQLUsers(params *GetV1beta1ClustersClusterIDSQLUsersParams, opts ...ClientOption) (*GetV1beta1ClustersClusterIDSQLUsersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetV1beta1ClustersClusterIDSQLUsersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetV1beta1ClustersClusterIDSQLUsers",
		Method:             "GET",
		PathPattern:        "/v1beta1/clusters/{clusterId}/sqlUsers",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetV1beta1ClustersClusterIDSQLUsersReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetV1beta1ClustersClusterIDSQLUsersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetV1beta1ClustersClusterIDSQLUsers: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetV1beta1ClustersClusterIDSQLUsersUserName queries sql user

This endpoint retrieves a sql user by user name.
*/
func (a *Client) GetV1beta1ClustersClusterIDSQLUsersUserName(params *GetV1beta1ClustersClusterIDSQLUsersUserNameParams, opts ...ClientOption) (*GetV1beta1ClustersClusterIDSQLUsersUserNameOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetV1beta1ClustersClusterIDSQLUsersUserNameParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetV1beta1ClustersClusterIDSQLUsersUserName",
		Method:             "GET",
		PathPattern:        "/v1beta1/clusters/{clusterId}/sqlUsers/{userName}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetV1beta1ClustersClusterIDSQLUsersUserNameReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetV1beta1ClustersClusterIDSQLUsersUserNameOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetV1beta1ClustersClusterIDSQLUsersUserName: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetV1beta1Projects gets list of org projects

This endpoint returns a list of org projects.
*/
func (a *Client) GetV1beta1Projects(params *GetV1beta1ProjectsParams, opts ...ClientOption) (*GetV1beta1ProjectsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetV1beta1ProjectsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetV1beta1Projects",
		Method:             "GET",
		PathPattern:        "/v1beta1/projects",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetV1beta1ProjectsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetV1beta1ProjectsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetV1beta1Projects: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PatchV1beta1ClustersClusterIDSQLUsersUserName updates one sql user

This endpoint update one sql user in the cluster.
*/
func (a *Client) PatchV1beta1ClustersClusterIDSQLUsersUserName(params *PatchV1beta1ClustersClusterIDSQLUsersUserNameParams, opts ...ClientOption) (*PatchV1beta1ClustersClusterIDSQLUsersUserNameOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPatchV1beta1ClustersClusterIDSQLUsersUserNameParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PatchV1beta1ClustersClusterIDSQLUsersUserName",
		Method:             "PATCH",
		PathPattern:        "/v1beta1/clusters/{clusterId}/sqlUsers/{userName}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PatchV1beta1ClustersClusterIDSQLUsersUserNameReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PatchV1beta1ClustersClusterIDSQLUsersUserNameOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PatchV1beta1ClustersClusterIDSQLUsersUserName: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PostCustomerSignupURL creates a new signup URL for an m s p customer

This endpoint creates a new signup URL for an MSP customer.
*/
func (a *Client) PostCustomerSignupURL(params *PostCustomerSignupURLParams, opts ...ClientOption) (*PostCustomerSignupURLOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostCustomerSignupURLParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostCustomerSignupURL",
		Method:             "POST",
		PathPattern:        "/customerSignupUrl",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PostCustomerSignupURLReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostCustomerSignupURLOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostCustomerSignupURL: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
PostV1beta1ClustersClusterIDSQLUsers creates one sql user

This endpoint will create one sql user int the cluster.
*/
func (a *Client) PostV1beta1ClustersClusterIDSQLUsers(params *PostV1beta1ClustersClusterIDSQLUsersParams, opts ...ClientOption) (*PostV1beta1ClustersClusterIDSQLUsersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostV1beta1ClustersClusterIDSQLUsersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostV1beta1ClustersClusterIDSQLUsers",
		Method:             "POST",
		PathPattern:        "/v1beta1/clusters/{clusterId}/sqlUsers",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PostV1beta1ClustersClusterIDSQLUsersReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostV1beta1ClustersClusterIDSQLUsersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostV1beta1ClustersClusterIDSQLUsers: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
