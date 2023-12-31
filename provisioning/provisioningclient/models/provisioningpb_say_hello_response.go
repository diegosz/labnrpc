// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ProvisioningpbSayHelloResponse The response message containing the greetings
//
// swagger:model provisioningpbSayHelloResponse
type ProvisioningpbSayHelloResponse struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this provisioningpb say hello response
func (m *ProvisioningpbSayHelloResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this provisioningpb say hello response based on context it is used
func (m *ProvisioningpbSayHelloResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ProvisioningpbSayHelloResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProvisioningpbSayHelloResponse) UnmarshalBinary(b []byte) error {
	var res ProvisioningpbSayHelloResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
