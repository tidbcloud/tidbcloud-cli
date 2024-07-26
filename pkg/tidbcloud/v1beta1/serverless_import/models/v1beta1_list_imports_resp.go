// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// V1beta1ListImportsResp v1beta1 list imports resp
//
// swagger:model v1beta1ListImportsResp
type V1beta1ListImportsResp struct {

	// The imports.
	// Read Only: true
	Imports []*V1beta1Import `json:"imports"`

	// The next page token.
	// Read Only: true
	NextPageToken string `json:"nextPageToken,omitempty"`

	// The total size of the imports.
	// Read Only: true
	TotalSize int64 `json:"totalSize,omitempty"`
}

// Validate validates this v1beta1 list imports resp
func (m *V1beta1ListImportsResp) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateImports(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1ListImportsResp) validateImports(formats strfmt.Registry) error {
	if swag.IsZero(m.Imports) { // not required
		return nil
	}

	for i := 0; i < len(m.Imports); i++ {
		if swag.IsZero(m.Imports[i]) { // not required
			continue
		}

		if m.Imports[i] != nil {
			if err := m.Imports[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("imports" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("imports" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1beta1 list imports resp based on the context it is used
func (m *V1beta1ListImportsResp) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateImports(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateNextPageToken(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTotalSize(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1ListImportsResp) contextValidateImports(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "imports", "body", []*V1beta1Import(m.Imports)); err != nil {
		return err
	}

	for i := 0; i < len(m.Imports); i++ {

		if m.Imports[i] != nil {

			if swag.IsZero(m.Imports[i]) { // not required
				return nil
			}

			if err := m.Imports[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("imports" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("imports" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *V1beta1ListImportsResp) contextValidateNextPageToken(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "nextPageToken", "body", string(m.NextPageToken)); err != nil {
		return err
	}

	return nil
}

func (m *V1beta1ListImportsResp) contextValidateTotalSize(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "totalSize", "body", int64(m.TotalSize)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1beta1ListImportsResp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1beta1ListImportsResp) UnmarshalBinary(b []byte) error {
	var res V1beta1ListImportsResp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
