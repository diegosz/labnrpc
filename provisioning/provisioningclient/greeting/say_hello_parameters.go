// Code generated by go-swagger; DO NOT EDIT.

package greeting

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

	"labnrpc/provisioning/provisioningclient/models"
)

// NewSayHelloParams creates a new SayHelloParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSayHelloParams() *SayHelloParams {
	return &SayHelloParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSayHelloParamsWithTimeout creates a new SayHelloParams object
// with the ability to set a timeout on a request.
func NewSayHelloParamsWithTimeout(timeout time.Duration) *SayHelloParams {
	return &SayHelloParams{
		timeout: timeout,
	}
}

// NewSayHelloParamsWithContext creates a new SayHelloParams object
// with the ability to set a context for a request.
func NewSayHelloParamsWithContext(ctx context.Context) *SayHelloParams {
	return &SayHelloParams{
		Context: ctx,
	}
}

// NewSayHelloParamsWithHTTPClient creates a new SayHelloParams object
// with the ability to set a custom HTTPClient for a request.
func NewSayHelloParamsWithHTTPClient(client *http.Client) *SayHelloParams {
	return &SayHelloParams{
		HTTPClient: client,
	}
}

/*
SayHelloParams contains all the parameters to send to the API endpoint

	for the say hello operation.

	Typically these are written to a http.Request.
*/
type SayHelloParams struct {

	/* Body.

	   The request message containing the user's name.
	*/
	Body *models.ProvisioningpbSayHelloRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the say hello params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SayHelloParams) WithDefaults() *SayHelloParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the say hello params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SayHelloParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the say hello params
func (o *SayHelloParams) WithTimeout(timeout time.Duration) *SayHelloParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the say hello params
func (o *SayHelloParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the say hello params
func (o *SayHelloParams) WithContext(ctx context.Context) *SayHelloParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the say hello params
func (o *SayHelloParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the say hello params
func (o *SayHelloParams) WithHTTPClient(client *http.Client) *SayHelloParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the say hello params
func (o *SayHelloParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the say hello params
func (o *SayHelloParams) WithBody(body *models.ProvisioningpbSayHelloRequest) *SayHelloParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the say hello params
func (o *SayHelloParams) SetBody(body *models.ProvisioningpbSayHelloRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *SayHelloParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
