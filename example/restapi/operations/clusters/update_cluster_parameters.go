// Code generated by gin-swagger; DO NOT EDIT.

package clusters

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/validate"
	"github.com/mikkeloscar/gin-swagger/api"
	"github.com/mikkeloscar/gin-swagger/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/mikkeloscar/gin-swagger/example/models"
)

// UpdateClusterEndpoint executes the core logic of the related
// route endpoint.
func UpdateClusterEndpoint(handler func(ctx *gin.Context, params *UpdateClusterParams) *api.Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := opentracing.SpanFromContext(tracing.Context(ctx))

		// attach tags to opentracing span
		if span != nil {
			ext.HTTPMethod.Set(span, ctx.Request.Method)
			ext.HTTPUrl.Set(span, ctx.Request.URL.String())
		}

		// generate params from request
		params := NewUpdateClusterParams()
		err := params.readRequest(ctx)
		if err != nil {
			errObj := err.(*errors.CompositeError)
			problem := api.Problem{
				Title:  "Unprocessable Entity.",
				Status: int(errObj.Code()),
				Detail: errObj.Error(),
			}

			// attach tags to opentracing span
			if span != nil {
				ext.HTTPStatusCode.Set(span, uint16(problem.Status))
			}

			ctx.Writer.Header().Set("Content-Type", "application/problem+json")
			ctx.JSON(problem.Status, problem)
			return
		}

		resp := handler(ctx, params)

		// attach tags to opentracing span
		if span != nil {
			ext.HTTPStatusCode.Set(span, uint16(resp.Code))
		}

		switch resp.Code {
		case http.StatusNoContent:
			ctx.AbortWithStatus(resp.Code)
		default:
			ctx.JSON(resp.Code, resp.Body)
		}
	}
}

// NewUpdateClusterParams creates a new UpdateClusterParams object
// with the default values initialized.
func NewUpdateClusterParams() *UpdateClusterParams {
	var ()
	return &UpdateClusterParams{}
}

// UpdateClusterParams contains all the bound params for the update cluster operation
// typically these are obtained from a http.Request
//
// swagger:parameters updateCluster
type UpdateClusterParams struct {

	/*Cluster that will be updated.
	  Required: true
	  In: body
	*/
	Cluster *models.ClusterUpdate
	/*ID of the cluster.
	  Required: true
	  Pattern: ^[a-z][a-z0-9-:]*[a-z0-9]$
	  In: path
	*/
	ClusterID string
}

// readRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *UpdateClusterParams) readRequest(ctx *gin.Context) error {
	var res []error
	formats := strfmt.NewFormats()

	if runtime.HasBody(ctx.Request) {
		var body models.ClusterUpdate
		if err := ctx.BindJSON(&body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("cluster", "body", ""))
			} else {
				res = append(res, errors.NewParseError("cluster", "body", "", err))
			}

		} else {
			if err := body.Validate(formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Cluster = &body
			}
		}

	} else {
		res = append(res, errors.Required("cluster", "body", ""))
	}

	rClusterID := []string{ctx.Param("cluster_id")}
	if err := o.bindClusterID(rClusterID, true, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateClusterParams) bindClusterID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.ClusterID = raw

	if err := o.validateClusterID(formats); err != nil {
		return err
	}

	return nil
}

func (o *UpdateClusterParams) validateClusterID(formats strfmt.Registry) error {

	if err := validate.Pattern("cluster_id", "path", o.ClusterID, `^[a-z][a-z0-9-:]*[a-z0-9]$`); err != nil {
		return err
	}

	return nil
}

// vim: ft=go
