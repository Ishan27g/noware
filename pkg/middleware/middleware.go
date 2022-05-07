package middleware

import (
	"context"
	"net/http"

	"github.com/Ishan27g/noware/pkg/actions"
	"github.com/Ishan27g/noware/pkg/noop"
	"github.com/gin-gonic/gin"
)

// HttpRequest returns a clone of the passed request with `noop`
// into the request's header if the context is `noop`. If `noop`, then `actions` are
// injected is context has `actions`
func HttpRequest(req *http.Request) *http.Request {
	if noop.ContainsNoop(req.Context()) {
		return httpReqInjectNoop(httpReqInjectActions(noop.AddHeader(req)))
	}
	return req
}

// Http middleware extracts `noop` and `actions` from the request's header and adds to the request's context
func Http(n http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if noop.CheckHeader(req) {
			if actions.CheckHeader(req) {
				req = httpReqExtractActions(req)
			}
			req = httpReqInjectNoop(req)
		}
		n.ServeHTTP(w, req)
	})
}

// Gin middleware extracts `noop` and `actions` from the request's header and adds to the request's context
func Gin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if noop.CheckHeader(c.Request) {
			if actions.CheckHeader(c.Request) {
				c.Request = httpReqExtractActions(c.Request)
			}
			c.Request = httpReqInjectNoop(c.Request)
		}
	}
}

// adds noop ctx to request
func httpReqInjectNoop(req *http.Request) *http.Request {
	return req.Clone(noop.NewCtxWithNoop(req.Context(), true))
}

// adds actions to request's header present in request's context
func httpReqInjectActions(req *http.Request) *http.Request {
	var a *actions.Actions
	if a = actions.FromCtx(req.Context()); a == nil {
		return req
	}
	actionJson, _ := a.Marshal()
	req = actions.AddHeader(req, string(actionJson))
	return req
}

// adds actions to request's context if present in request's header
func httpReqExtractActions(req *http.Request) *http.Request {
	a := actions.FromHeader(req.Header)
	if a == nil {
		return req
	}
	return req.Clone(actions.NewCtx(context.Background(), a))
}
