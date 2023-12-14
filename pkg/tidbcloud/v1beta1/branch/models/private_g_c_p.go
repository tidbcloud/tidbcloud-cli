// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PrivateGCP Message for GCP Private Link Service.
//
// swagger:model PrivateGCP
type PrivateGCP struct {

	// Output Only. Target Service Account for Private Link Service.
	// Read Only: true
	TargetServiceAccount string `json:"targetServiceAccount,omitempty"`
}

// Validate validates this private g c p
func (m *PrivateGCP) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validate this private g c p based on the context it is used
func (m *PrivateGCP) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateTargetServiceAccount(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PrivateGCP) contextValidateTargetServiceAccount(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "targetServiceAccount", "body", string(m.TargetServiceAccount)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PrivateGCP) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PrivateGCP) UnmarshalBinary(b []byte) error {
	var res PrivateGCP
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}