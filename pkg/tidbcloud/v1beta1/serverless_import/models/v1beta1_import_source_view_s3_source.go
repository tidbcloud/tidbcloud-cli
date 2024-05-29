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

// V1beta1ImportSourceViewS3Source v1beta1 import source view s3 source
//
// swagger:model v1beta1ImportSourceViewS3Source
type V1beta1ImportSourceViewS3Source struct {

	// access key
	AccessKey *V1beta1ImportSourceViewAccessKey `json:"accessKey,omitempty"`

	// role arn
	RoleArn *V1beta1ImportSourceViewRoleArn `json:"roleArn,omitempty"`

	// s3 Uri
	S3URI string `json:"s3Uri,omitempty"`

	// target database
	TargetDatabase *string `json:"targetDatabase,omitempty"`

	// type
	Type V1beta1AuthType `json:"type,omitempty"`
}

// Validate validates this v1beta1 import source view s3 source
func (m *V1beta1ImportSourceViewS3Source) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccessKey(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRoleArn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1ImportSourceViewS3Source) validateAccessKey(formats strfmt.Registry) error {
	if swag.IsZero(m.AccessKey) { // not required
		return nil
	}

	if m.AccessKey != nil {
		if err := m.AccessKey.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("accessKey")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("accessKey")
			}
			return err
		}
	}

	return nil
}

func (m *V1beta1ImportSourceViewS3Source) validateRoleArn(formats strfmt.Registry) error {
	if swag.IsZero(m.RoleArn) { // not required
		return nil
	}

	if m.RoleArn != nil {
		if err := m.RoleArn.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("roleArn")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("roleArn")
			}
			return err
		}
	}

	return nil
}

func (m *V1beta1ImportSourceViewS3Source) validateType(formats strfmt.Registry) error {
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

// ContextValidate validate this v1beta1 import source view s3 source based on the context it is used
func (m *V1beta1ImportSourceViewS3Source) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAccessKey(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRoleArn(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1ImportSourceViewS3Source) contextValidateAccessKey(ctx context.Context, formats strfmt.Registry) error {

	if m.AccessKey != nil {

		if swag.IsZero(m.AccessKey) { // not required
			return nil
		}

		if err := m.AccessKey.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("accessKey")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("accessKey")
			}
			return err
		}
	}

	return nil
}

func (m *V1beta1ImportSourceViewS3Source) contextValidateRoleArn(ctx context.Context, formats strfmt.Registry) error {

	if m.RoleArn != nil {

		if swag.IsZero(m.RoleArn) { // not required
			return nil
		}

		if err := m.RoleArn.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("roleArn")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("roleArn")
			}
			return err
		}
	}

	return nil
}

func (m *V1beta1ImportSourceViewS3Source) contextValidateType(ctx context.Context, formats strfmt.Registry) error {

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
func (m *V1beta1ImportSourceViewS3Source) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1beta1ImportSourceViewS3Source) UnmarshalBinary(b []byte) error {
	var res V1beta1ImportSourceViewS3Source
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
