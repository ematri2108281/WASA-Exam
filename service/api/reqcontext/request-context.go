/*
Package reqcontext contains the per-request context object.

An instance of RequestContext is created for each incoming HTTP request by
the middleware in api-context-wrapper.go and passed down to every handler.
Values inside should be considered valid only for the lifetime of that request.
*/
package reqcontext

import (
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// RequestContext carries request-scoped values for a single HTTP call.
type RequestContext struct {
	// ReqUUID is a unique identifier for this request, useful for tracing logs.
	ReqUUID uuid.UUID

	// Logger is a structured logger pre-populated with request fields.
	Logger logrus.FieldLogger

	// UserID is the authenticated user's identifier, extracted from the Bearer
	// token by the auth middleware. Empty string means unauthenticated.
	UserID string
}

// IsAuthenticated reports whether the request carries a valid user identity.
func (c RequestContext) IsAuthenticated() bool {
	return c.UserID != ""
}
