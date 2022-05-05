package example

//
//import (
//	"fmt"
//	"log"
//	"net/http"
//
//	"github.com/Ishan27g/noware/pkg/actions"
//	"github.com/Ishan27g/noware/pkg/middleware"
//	"github.com/Ishan27g/noware/pkg/noop"
//	"github.com/gin-gonic/gin"
//	"github.com/nats-io/nats.go"
//)
//
//// AsyncPubMethod is any async method that subscribes to nats and publishes some data at the end
//type AsyncPubMethod func(msg *nats.Msg) bool
//
//type wrapper struct {
//	methods map[string]gin.HandlerFunc
//}
//type d struct {
//	Data string `json:"Data"`
//}
//
///* Noop-Testing:
//   Wrap async method as its own http handler
//   - extracts the noop & action from request's context
//   - triggers underlying async method
//   - responds with actions taken by this method
//*/
//func wrapAsync(endpoint string, async AsyncPubMethod) func(c *gin.Context) {
//	return func(c *gin.Context) {
//		// dataFromSub in the request should be the dataFromSub that would normally be received by the async method
//		var dataFromSub d
//		var dataToPub nats.Msg
//		if e := c.ShouldBindJSON(&dataFromSub); e != nil {
//			c.JSON(http.StatusExpectationFailed, nil)
//			return
//		}
//
//		fmt.Println("from context - ", noop.ContainsNoop(c.Request.Context()))
//
//		ctx := c.Request.Context() // noop.NewCtxWithNoop(context.Background(), isNoop) // todo new ctx or gin-request context?
//
//		//
//		dataToPub.Data = []byte("data to next subscriber")
//		// create actions for this request
//		a := actions.New()
//
//		// add an event to represent the dataFromSub that would be sent by the async method
//		a.AddEvent(actions.Event{Name: "Endpoint hit at " + endpoint, Meta: string(dataToPub.Data)})
//		ctx = actions.NewCtx(ctx, a)
//
//		// finally, check if noop operation to trigger or skip the async method call
//		if !noop.ContainsNoop(ctx) {
//			if !async(&nats.Msg{Data: []byte(dataFromSub.Data)}) {
//				c.JSON(http.StatusExpectationFailed, nil)
//				return
//			}
//		}
//
//		// respond with the
//		if rsp, err := a.Marshal(); err == nil {
//			c.Writer.WriteHeader(http.StatusOK)
//			_, _ = c.Writer.Write(rsp)
//			return
//		}
//		c.JSON(http.StatusExpectationFailed, nil)
//	}
//}
//
//func (w *wrapper) wrap(endpoint string, async AsyncPubMethod) {
//	w.methods[endpoint] = wrapAsync(endpoint, async)
//}
//func (w *wrapper) start() {
//	gin.SetMode(gin.ReleaseMode)
//	r := gin.New()
//	r.Use(gin.Recovery())
//	r.Use(middleware.Gin())
//	for endpoint, method := range w.methods {
//		r.POST(endpoint, method)
//	}
//	port := ":9091"
//	log.Println("starting on", port)
//	go func() {
//		err := r.Run(port)
//		if err != nil {
//			log.Fatalf(err.Error(), err)
//		}
//	}()
//}
//func setup() {
//
//	w := wrapper{map[string]gin.HandlerFunc{}}
//
//	w.wrap("/endpoint", New().GenericMethod)
//
//	w.start()
//
//	<-make(chan bool)
//}
