package main

//
//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"log"
//	"net/http"
//	"time"
//
//	"github.com/Ishan27g/noware/pkg/actions"
//	"github.com/Ishan27g/noware/pkg/middleware"
//	"github.com/Ishan27g/noware/pkg/noop"
//)
//
//const (
//	urlGolang = "http://localhost:8081/new/user"
//	urlNode   = "http://localhost:8082/delete/user"
//)
//
//type Request struct {
//	Id   string `json:"data,omitempty"`
//	Name string `json:"name,omitempty"`
//}
//type Response Request
//
//func timeIt(from time.Time) {
//	log.Println("took", time.Since(from).String())
//}
//func requestWithNoopCtx(ctx context.Context, url string, payload Request) bool {
//	defer timeIt(time.Now())
//
//	p, _ := json.Marshal(payload)
//	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(p))
//	if err != nil {
//		log.Fatalf("%v", err)
//	}
//
//	// add noop ctx
//	r := middleware.HttpRequest(req)
//
//	// fmt.Println("Sending - ", req.Header)
//	client := http.DefaultClient
//	res, err := client.Do(r)
//	if err != nil {
//		log.Fatalf("%v", err)
//	}
//
//	//rsp, err := ioutil.ReadAll(res.Body)
//	//if err != nil {
//	//	return false
//	//}
//	//res.Body.Close()
//	//a, _ := actions.UnMarshal(rsp)
//	//fmt.Println("response is ", a.GetEvents())
//
//	return res.StatusCode == http.StatusOK
//}
//func addAction(ctx context.Context) context.Context {
//	a := actions.New()
//	a.AddEvent(actions.Event{
//		Name: "one",
//		Meta: []int{1, 2, 3, 4},
//	})
//
//	ctx = actions.NewCtx(ctx, a)
//	return ctx
//}
//func main() {
//
//	ctx := noop.NewCtxWithNoop(context.Background(), true)
//	ctx = addAction(ctx)
//	requestWithNoopCtx(ctx, "http://localhost:8081/new/user", Request{Name: "someone"})
//	requestWithNoopCtx(ctx, "http://localhost:8082/1", Request{Name: "someone"})
//
//	requestWithNoopCtx(context.Background(), "http://localhost:8081/new/user", Request{Name: "someone"})
//	requestWithNoopCtx(context.Background(), "http://localhost:8082/1", Request{Name: "someone"})
//
//}
