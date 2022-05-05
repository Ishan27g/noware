package main

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"log"
//	"net/http"
//	"time"
//
//	"github.com/Ishan27g/noware/pkg/actions"
//	"github.com/Ishan27g/noware/pkg/middleware"
//	"github.com/Ishan27g/noware/pkg/noop"
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//)
//
//const (
//	port        = ":8081"
//	endpointUrl = "/new/user"
//)
//
//type Request struct {
//	Id   string `json:"data,omitempty"`
//	Name string `json:"name,omitempty"`
//}
//type Response Request
//
//func endpointHandler(c *gin.Context) {
//	var data Request
//	var response Response
//	var rsp []byte
//	var err error
//
//	if err = c.ShouldBindJSON(&data); err != nil {
//		c.JSON(http.StatusExpectationFailed, "bad payload")
//		return
//	}
//	fmt.Println("context has noop ?- ", noop.ContainsNoop(c.Request.Context()))
//	fmt.Println(" creating id for ", data.Name)
//
//	ctx := c.Request.Context() // noop.NewCtxWithNoop(context.Background(), isNoop) // todo new ctx or gin-request context?
//
//	uid := uuid.New().String()
//	response = Response{Id: uid}
//	// create actions for this request
//	// add an event to it
//	// set this as response
//	a := actions.FromCtx(ctx)
//	if a != nil {
//		a.AddEvent(actions.Event{Name: data.Name + c.Request.RequestURI, Meta: response})
//		if r, err := a.Marshal(); err == nil {
//			rsp = r
//		}
//	}
//
//	// finally, check if noop operation to trigger or skip the external service call
//	if !noop.ContainsNoop(ctx) {
//		// do db update
//		log.Println("updating database")
//
//		<-time.After(2 * time.Second)
//
//		log.Println("updated database")
//
//		rsp, _ = json.Marshal(response)
//	}
//
//	fmt.Println(string(rsp))
//	// respond
//	c.Writer.WriteHeader(http.StatusOK)
//	_, _ = c.Writer.Write(rsp)
//}
//
//func ginServer(port string) *http.Server {
//	gin.SetMode(gin.ReleaseMode)
//	r := gin.New()
//	r.Use(gin.Logger())
//
//	// use noop middleware
//	r.Use(middleware.Gin())
//
//	r.POST(endpointUrl, endpointHandler)
//	return &http.Server{Addr: port, Handler: r}
//}
//
//func main() {
//	ctx, can := context.WithCancel(context.Background())
//	defer can()
//
//	server := ginServer(port)
//	go func() {
//		log.Println("starting on", server.Addr)
//		if err := server.ListenAndServe(); err != nil {
//			if err != http.ErrServerClosed {
//				log.Fatal(err)
//			}
//		}
//	}()
//	<-ctx.Done()
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	if err := server.Shutdown(ctx); err != nil {
//		log.Fatal(err)
//	}
//	log.Println("shutdown on - ", server.Addr)
//}
