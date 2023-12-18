// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// FormsBillsOtherName forms bills other name
//
// swagger:model forms.BillsOtherName
type FormsBillsOtherName struct {

	// charge name
	// Example: Support Plan
	ChargeName string `json:"chargeName,omitempty"`

	// The total credits held by the organization in this bill. The value of this field is expressed in cents (100ths of one US Dollar).
	// Example: 0.00
	Credits string `json:"credits,omitempty"`

	// Total amount of discounts applied to this bill. The value of this field is expressed in cents (100ths of one US Dollar).
	// Example: 0.00
	Discounts string `json:"discounts,omitempty"`

	// The sum of services that the specified organization consumed in the period during this bill period. The value of this field is expressed in cents (100ths of one US Dollar).
	// Example: 0.00
	RunningTotal string `json:"runningTotal,omitempty"`

	// The total amount that the specified organization should pay toward this bill. The value of this field is expressed in cents (100ths of one US Dollar).
	// `total_cost` = `running_total` - `discounts` - `credits`.
	// Example: 0.00
	TotalCost string `json:"totalCost,omitempty"`
}

// Validate validates this forms bills other name
func (m *FormsBillsOtherName) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this forms bills other name based on context it is used
func (m *FormsBillsOtherName) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *FormsBillsOtherName) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FormsBillsOtherName) UnmarshalBinary(b []byte) error {
	var res FormsBillsOtherName
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
