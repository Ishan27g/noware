package middleware

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/Ishan27g/noware/pkg/actions"
	"github.com/Ishan27g/noware/pkg/noop"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	endpoint    = "http://localhost"
	port1       = ":8081"
	port2       = ":8082"
	endpointUrl = "/url"

	headerKey1 = "some"
	headerKey2 = "another"
	headerVal1 = "header"
	headerVal2 = "header"
)

var (
	RequestWithoutNoop = "Request does not have noop ctx"
	RequestHasNoopCtx  = "Request has noop ctx"
	NoopCtxHasActions  = "Ctx has action events"
)
var metas = []interface{}{nil, "metadata", []int{1, 2, 3, 4}, map[string]string{"1": "1",
	"2": "2",
	"3": "3",
}, map[string]int{
	"1": 1,
	"2": 2,
	"3": 3,
}}

func createEvents(count int) []actions.Event {
	var events []actions.Event
	if count > len(metas) {
		count = len(metas)
	}
	for i := 0; i < count; i++ {
		events = append(events, actions.Event{
			Name: strconv.Itoa(i),
			Meta: metas[i],
		})
	}
	return events
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(headerKey1) != headerVal1 || r.Header.Get(headerKey2) != headerVal2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing headers"))
		return
	}
	w.WriteHeader(http.StatusAccepted)

	if noop.ContainsNoop(r.Context()) {
		a := actions.FromCtx(r.Context())
		if a != nil {
			if len(a.GetEvents()) > 0 {
				w.Write([]byte(NoopCtxHasActions))
				return
			}
		}
		w.Write([]byte(RequestHasNoopCtx))
		return
	}

	w.Write([]byte(RequestWithoutNoop))
	return
}
func runServer(server *http.Server, ctx context.Context) {
	go func() {
		log.Println("starting on", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("shutdown on - ", server.Addr)
}
func sendRequest(addr string, requestIsNoop bool, events []actions.Event) string {
	ctx := context.Background()
	if requestIsNoop {
		ctx = noop.NewCtxWithNoop(ctx)
	}
	if len(events) > 0 {
		a := actions.New()
		a.AddEvent(events...)
		ctx = actions.NewCtx(ctx, a)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+addr+endpointUrl, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	req.Header.Add(headerKey1, headerVal1)
	req.Header.Add(headerKey2, headerVal2)

	res, err := http.DefaultClient.Do(noop.HttpRequest(req))
	if err != nil {
		log.Fatalf("%v", err)
	}
	rsp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	res.Body.Close()
	return string(rsp)
	// return strings.EqualFold(string(rsp), expectedResponse)
}

func ginServer(port string) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(Gin())
	r.GET(endpointUrl, gin.WrapF(endpointHandler))
	return &http.Server{Addr: port, Handler: r}

}
func httpServer(port string) *http.Server {
	var handler = http.NewServeMux()
	handler.Handle(endpointUrl, Http(endpointHandler))
	return &http.Server{Addr: port, Handler: handler}
}

func TestMiddleware(t *testing.T) {
	t.Parallel()

	servers := []struct {
		name string
		*http.Server
	}{
		{name: "gin - Middleware on port " + port1, Server: ginServer(port1)},
		{name: "http - Middleware on port" + port2, Server: httpServer(port2)},
	}

	tests := []struct {
		name             string
		requestIsNoop    bool
		events           []actions.Event
		expectedResponse string
	}{
		{name: "", requestIsNoop: false, events: nil, expectedResponse: RequestWithoutNoop},
		{name: "", requestIsNoop: true, events: nil, expectedResponse: RequestHasNoopCtx},
		{name: "", requestIsNoop: true, events: createEvents(0), expectedResponse: RequestHasNoopCtx},
		{name: "", requestIsNoop: true, events: createEvents(1), expectedResponse: NoopCtxHasActions},
		{name: "", requestIsNoop: true, events: createEvents(10), expectedResponse: NoopCtxHasActions},
	}

	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	for _, server := range servers {

		go runServer(server.Server, ctx)
		<-time.After(2 * time.Second)

		for _, tt := range tests {
			t.Run(server.name+"\t"+tt.expectedResponse, func(t *testing.T) {
				assert.Equal(t, tt.expectedResponse, sendRequest(server.Addr, tt.requestIsNoop, tt.events), server.name+"\t"+tt.expectedResponse)
			})
		}
	}

}
