package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Ishan27g/noware/pkg/actions"
	nats2 "github.com/Ishan27g/noware/pkg/examples/async/nats"
	"github.com/Ishan27g/noware/pkg/examples/async/stage"
	"github.com/Ishan27g/noware/pkg/examples/async/types"
	"github.com/Ishan27g/noware/pkg/noop"
)

const (
	host = "http://localhost"
)

func buildUrl(service string) string {
	return host + service
}

func sendHttpReq(ctx context.Context, to string, data string) []actions.Event {
	payload, err := json.Marshal(types.Data{Data: data})
	if err != nil {
		return nil
	}
	request, err := http.NewRequestWithContext(ctx, "POST", to, bytes.NewReader(payload))
	if err != nil {
		return nil
	}
	client := http.Client{
		Timeout: 6 * time.Second,
	}
	resp, err := client.Do(noop.HttpRequest(request))
	if err != nil {
		fmt.Println(err.Error() + " ok")
		return nil
	}
	rsp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	action, _ := actions.UnMarshal(rsp)
	return action.GetEvents()
}
func sendRequestWithActions(action *actions.Actions, service, data string) bool {
	ctx := actions.NewCtx(noop.NewCtxWithNoop(context.Background()), action)
	rsp := sendHttpReq(ctx, buildUrl(service), data)
	if rsp != nil {
		action.AddEvent(rsp...)
		return true
	}
	return false
}

// SendHttpNoop sends http to individually trigger each service in order
func SendHttpNoop(s1, s2, s3 string, triggerAll bool) {

	data := "pipeline-data"
	endpoint := "/endpoint"

	actions := actions.New()
	if !sendRequestWithActions(actions, s1+endpoint, data) {
		fmt.Println("error at " + s1)
		return
	}
	if !sendRequestWithActions(actions, s2+endpoint, data) {
		fmt.Println("error at " + s2)
		return
	}

	if !sendRequestWithActions(actions, s3+endpoint, data) {
		fmt.Println("error at " + s3)
		return
	}

	fmt.Println("Final response actions →")
	for _, event := range actions.GetEvents() {
		fmt.Println(fmt.Sprintf("%+v", event))
	}

}

func main() {
	ctx, can := context.WithCancel(context.Background())
	defer can()

	// async-services have a http server that wraps their respective async-subscribe method
	svc1 := stage.Start(":9091", ctx, "stage-1", "stage-2")
	_ = stage.Start(":9092", ctx, "stage-2", "stage-3")
	_ = stage.Start(":9093", ctx, "stage-3", "stage-4")

	<-time.After(2 * time.Second)
	// example - trigger normal operation i.e. async msg received by stage 1
	nats2.Publish(svc1.PublishTopic(), "pipeline-data", nil)

	// NOOP/EVENTS - trigger http request with noop & events to stage 1,2,3
	<-time.After(2 * time.Second)
	SendHttpNoop(":9091", ":9092", ":9093", true)
}
