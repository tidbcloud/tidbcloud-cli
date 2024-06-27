// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1beta1ImportSourceViewAzBlobSource v1beta1 import source view az blob source
//
// swagger:model v1beta1ImportSourceViewAzBlobSource
type V1beta1ImportSourceViewAzBlobSource struct {

	// blob Uri
	BlobURI string `json:"blobUri,omitempty"`

	// sas token
	SasToken V1beta1ImportSourceViewSASToken `json:"sasToken,omitempty"`

	// type
	Type V1beta1ImportSourceViewAzBlobSourceAuthType `json:"type,omitempty"`
}

// Validate validates this v1beta1 import source view az blob source
func (m *V1beta1ImportSourceViewAzBlobSource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1ImportSourceViewAzBlobSource) validateType(formats strfmt.Registry) error {
	if swag.IsZero(m.Type) { // not required
		return nil
	}

	if err := m.Type.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("type")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("type")
		}
		return err
	}

	return nil
}

// ContextValidate validate this v1beta1 import source view az blob source based on the context it is used
func (m *V1beta1ImportSourceViewAzBlobSource) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1ImportSourceViewAzBlobSource) contextValidateType(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	if err := m.Type.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("type")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("type")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1beta1ImportSourceViewAzBlobSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1beta1ImportSourceViewAzBlobSource) UnmarshalBinary(b []byte) error {
	var res V1beta1ImportSourceViewAzBlobSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
