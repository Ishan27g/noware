package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func CheckHeader(req *http.Request) bool {
	hasActions := req.Header.Get(actionKey)
	return hasActions != ""
}
func AddHeader(req *http.Request, val string) *http.Request {
	val = strings.TrimPrefix(val, "[")
	val = strings.TrimSuffix(val, "]")
	req.Header.Set(actionKey, val)
	return req
}

func FromHeader(header http.Header) *Actions {
	var a *Actions = New()
	for _, ev := range header.Values(actionKey) {
		var e = Event{}
		err := json.Unmarshal([]byte(ev), &e)
		if err != nil {
			fmt.Println(err.Error())
		}
		a.AddEvent(e)
	}
	//s := header.Get(actionKey)
	//if s == "" {
	//	fmt.Println("no action bytes")
	//	return FromCtx(nil)
	//}
	//fmt.Println("FROM HEADER _ ", s)
	//a, err := UnMarshal([]byte(s))
	//fmt.Println("actions from headers - ", a.GetEvents())
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return FromCtx(nil)
	//}
	return a
}
