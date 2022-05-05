package actions

import (
	"fmt"
	"net/http"
)

func CheckHeader(req *http.Request) bool {
	hasActions := req.Header.Get(actionKey)
	return hasActions != ""
}
func AddHeader(req *http.Request, val string) *http.Request {
	req.Header.Add(actionKey, val)
	return req
}

func FromHeader(header http.Header) *Actions {
	s := header.Get(actionKey)
	if s == "" {
		fmt.Println("no action bytes")
		return FromCtx(nil)
	}
	a, err := UnMarshal([]byte(s))
	if err != nil {
		return FromCtx(nil)
	}
	return &a
}
