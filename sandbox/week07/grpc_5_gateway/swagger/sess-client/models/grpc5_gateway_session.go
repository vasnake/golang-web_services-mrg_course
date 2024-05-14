// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Grpc5GatewaySession grpc 5 gateway session
//
// swagger:model grpc_5_gatewaySession
type Grpc5GatewaySession struct {

	// login
	Login string `json:"login,omitempty"`

	// useragent
	Useragent string `json:"useragent,omitempty"`
}

// Validate validates this grpc 5 gateway session
func (m *Grpc5GatewaySession) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this grpc 5 gateway session based on context it is used
func (m *Grpc5GatewaySession) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Grpc5GatewaySession) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Grpc5GatewaySession) UnmarshalBinary(b []byte) error {
	var res Grpc5GatewaySession
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}