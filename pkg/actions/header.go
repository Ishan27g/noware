package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func CheckHeader(req *http.Request) bool {
	hasActions := req.Header.Get(actionKey)
	return hasActions != ""
}

// Inject adds actions to request's header present in request's context
func Inject(req *http.Request) *http.Request {
	var a *Actions
	if a = FromCtx(req.Context()); a == nil {
		return req
	}
	actionJson, _ := a.Marshal()
	req = addHeader(req, string(actionJson))
	return req
}

// Extract actions to request's context if present in request's header
func Extract(req *http.Request) *http.Request {
	a := fromHeader(req.Header)
	if a == nil {
		return req
	}
	return req.Clone(NewCtx(context.Background(), a))
}

func addHeader(req *http.Request, val string) *http.Request {
	val = strings.TrimPrefix(val, "[")
	val = strings.TrimSuffix(val, "]")
	req.Header.Set(actionKey, val)
	return req
}

func fromHeader(header http.Header) *Actions {
	var a = New()
	for _, ev := range header.Values(actionKey) {
		var e = Event{}
		err := json.Unmarshal([]byte(ev), &e)
		if err != nil {
			fmt.Println(err.Error())
		}
		a.AddEvent(e)
	}
	return a
}
