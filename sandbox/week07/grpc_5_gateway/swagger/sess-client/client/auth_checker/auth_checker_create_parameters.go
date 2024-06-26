// Code generated by go-swagger; DO NOT EDIT.

package auth_checker

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

	"week07/grpc_5_gateway/swagger/sess-client/models"
)

// NewAuthCheckerCreateParams creates a new AuthCheckerCreateParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewAuthCheckerCreateParams() *AuthCheckerCreateParams {
	return &AuthCheckerCreateParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewAuthCheckerCreateParamsWithTimeout creates a new AuthCheckerCreateParams object
// with the ability to set a timeout on a request.
func NewAuthCheckerCreateParamsWithTimeout(timeout time.Duration) *AuthCheckerCreateParams {
	return &AuthCheckerCreateParams{
		timeout: timeout,
	}
}

// NewAuthCheckerCreateParamsWithContext creates a new AuthCheckerCreateParams object
// with the ability to set a context for a request.
func NewAuthCheckerCreateParamsWithContext(ctx context.Context) *AuthCheckerCreateParams {
	return &AuthCheckerCreateParams{
		Context: ctx,
	}
}

// NewAuthCheckerCreateParamsWithHTTPClient creates a new AuthCheckerCreateParams object
// with the ability to set a custom HTTPClient for a request.
func NewAuthCheckerCreateParamsWithHTTPClient(client *http.Client) *AuthCheckerCreateParams {
	return &AuthCheckerCreateParams{
		HTTPClient: client,
	}
}

/*
AuthCheckerCreateParams contains all the parameters to send to the API endpoint

	for the auth checker create operation.

	Typically these are written to a http.Request.
*/
type AuthCheckerCreateParams struct {

	// Body.
	Body *models.Grpc5GatewaySession

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the auth checker create params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AuthCheckerCreateParams) WithDefaults() *AuthCheckerCreateParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the auth checker create params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *AuthCheckerCreateParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the auth checker create params
func (o *AuthCheckerCreateParams) WithTimeout(timeout time.Duration) *AuthCheckerCreateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the auth checker create params
func (o *AuthCheckerCreateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the auth checker create params
func (o *AuthCheckerCreateParams) WithContext(ctx context.Context) *AuthCheckerCreateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the auth checker create params
func (o *AuthCheckerCreateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the auth checker create params
func (o *AuthCheckerCreateParams) WithHTTPClient(client *http.Client) *AuthCheckerCreateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the auth checker create params
func (o *AuthCheckerCreateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the auth checker create params
func (o *AuthCheckerCreateParams) WithBody(body *models.Grpc5GatewaySession) *AuthCheckerCreateParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the auth checker create params
func (o *AuthCheckerCreateParams) SetBody(body *models.Grpc5GatewaySession) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *AuthCheckerCreateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
