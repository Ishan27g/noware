package stage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ishan27g/noware/examples/async/service"
	"github.com/Ishan27g/noware/examples/async/types"
	"github.com/Ishan27g/noware/pkg/actions"
	"github.com/Ishan27g/noware/pkg/middleware"
	"github.com/Ishan27g/noware/pkg/noop"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

// existingPubMethod is any async method that subscribes to nats and publishes some data at the end
type existingPubMethod func(msg *nats.Msg) bool

type wrapper struct {
	methods map[string]gin.HandlerFunc
}

/* Noop-Testing:
Wrap async method as its own http handler
- extracts the noop & action from request's context
- triggers underlying async method
- responds with actions taken by this method
*/
func wrapAsync(endpoint string, async existingPubMethod) func(c *gin.Context) {
	return func(c *gin.Context) {
		// dataFromSub in the request should be the dataFromSub that would normally be received by the async method
		var dataFromSub types.Data
		var dataToPub nats.Msg
		if e := c.ShouldBindJSON(&dataFromSub); e != nil {
			c.JSON(http.StatusExpectationFailed, nil)
			return
		}

		fmt.Println("Context has actions? - ", noop.ContainsNoop(c.Request.Context()))

		ctx := c.Request.Context() // noop.NewCtxWithNoop(context.Background(), isNoop) // todo new ctx or gin-request context?

		dataToPub.Data = []byte("some data for the next subscriber")
		// create actions for this request
		a := actions.New()

		// add an event to represent the dataFromSub that would be sent by the async method
		a.AddEvent(actions.Event{Name: "Endpoint hit at " + endpoint, Meta: string(dataToPub.Data)})
		ctx = actions.NewCtx(ctx, a)

		// finally, check if noop operation to trigger or skip the async method call
		if !noop.ContainsNoop(ctx) {
			if !async(&nats.Msg{Data: []byte(dataFromSub.Data)}) {
				c.JSON(http.StatusExpectationFailed, nil)
				return
			}
		}

		// respond with the
		if rsp, err := a.Marshal(); err == nil {
			c.Writer.WriteHeader(http.StatusOK)
			_, _ = c.Writer.Write(rsp)
			return
		}
		c.JSON(http.StatusExpectationFailed, nil)
	}
}

func (w *wrapper) wrap(endpoint string, async existingPubMethod) {
	w.methods[endpoint] = wrapAsync(endpoint, async)
}
func (w *wrapper) runServer(port string, ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Gin())
	for endpoint, method := range w.methods {
		r.POST(endpoint, method)
	}
	server := http.Server{
		Addr:    port,
		Handler: r,
	}
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
func Start(port string, ctx context.Context, listenForTopic string, publishToTopic string) *service.Service {
	svc := service.New(listenForTopic, publishToTopic)
	w := wrapper{map[string]gin.HandlerFunc{}}

	w.wrap("/endpoint", svc.GenericMethod)

	go w.runServer(port, ctx)
	return svc
}
