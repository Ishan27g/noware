package middleware

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/Ishan27g/noware/pkg/actions"
	"github.com/Ishan27g/noware/pkg/noop"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	endpoint    = "http://localhost"
	port        = ":8081"
	endpointUrl = "/url"
	urlE        = endpoint + port + endpointUrl

	headerKey1 = "some"
	headerKey2 = "another"
	headerVal1 = "header"
	headerVal2 = "header"
)

var clientEvent = actions.Event{
	Name: "client",
	Meta: map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	},
}

func endpointHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(headerKey1) != headerVal1 || r.Header.Get(headerKey2) != headerVal2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing headers"))
		return
	}

	fmt.Println("Events - ", actions.FromCtx(r.Context()).GetEvents())

	if noop.ContainsNoop(r.Context()) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Ctx with Noop"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ctx without Noop"))
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

func request(ctx context.Context, expect string) bool {

	a := actions.New()
	a.AddEvent(clientEvent)
	a.AddEvent(actions.Event{
		Name: "one",
		Meta: []int{1, 2, 3, 4},
	})

	ctx = actions.NewCtx(ctx, a)
	req, err := http.NewRequestWithContext(ctx, "GET", urlE, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	req.Header.Add(headerKey1, headerVal1)
	req.Header.Add(headerKey2, headerVal2)

	r := HttpRequest(req)

	//	fmt.Println("Sending - ", req.Header)
	client := http.DefaultClient
	res, err := client.Do(r)
	if err != nil {
		log.Fatalf("%v", err)
	}
	rsp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false
	}
	res.Body.Close()
	return expect == string(rsp)
}
func testWith(t *testing.T, ctx context.Context, expect string) {
	assert.Equal(t, true, request(ctx, expect))
}
func testServer(t *testing.T, h *http.Server) {
	ctx, can := context.WithCancel(context.Background())
	defer can()
	go runServer(h, ctx)
	<-time.After(2 * time.Second)

	testWith(t, context.Background(), "Ctx without Noop")
	testWith(t, noop.NewCtxWithNoop(context.Background(), true), "Ctx with Noop")
}

func TestServer(t *testing.T) {
	testServer(t, httpServer(port))
	testServer(t, ginServer(port))
}
