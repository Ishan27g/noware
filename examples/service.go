package example

//
//import (
//	"github.com/nats-io/nats.go"
//)
//
//type Service struct {
//}
//
//func New() *Service {
//	s := Service{}
//	return &s
//}
//
//// GenericMethod is any generic async method in a request pipeline
//// It would normally be triggered by nats/kafka. After processing
//// this method would ideally publish to nats/kafka to trigger the next
//// service/stage of the request pipeline
///* Noop-Testing: Before the final publish operation ->
//- add publish `msg` to the context
//- return if ctx is noop, otherwise publish `msg` to queue
//*/
//
//func (s *Service) GenericMethod(msg *nats.Msg) bool {
//	return NatsPublish("s.subjForNext()", string(msg.Data), nil)
//}
