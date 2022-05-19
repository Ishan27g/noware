package service

import (
	"context"

	nats2 "github.com/Ishan27g/noware/examples/async/nats"
	"github.com/nats-io/nats.go"
)

type Service struct {
	listenForTopic string
	publishToTopic string
}

func (s *Service) PublishTopic() string {
	return s.publishToTopic
}

func (s *Service) SetPublishToTopic(publishToTopic string) {
	s.publishToTopic = publishToTopic
}
func New(listenForTopic string, publishToTopic string) *Service {
	s := Service{listenForTopic: listenForTopic, publishToTopic: publishToTopic}
	go nats2.Subscribe(s.listenForTopic, func(ctx context.Context, msg *nats.Msg) bool {
		return s.GenericMethod(msg)
	})
	return &s
}

// GenericMethod is any generic async method in a request pipeline
// It would normally be triggered by nats/kafka. After processing
// this method would ideally publish to nats/kafka to trigger the next
// service/stage of the request pipeline
/* Noop-Testing: Before the final publish operation ->
- add publish `msg` to the context
- return if ctx is noop, otherwise publish `msg` to queue
*/
func (s *Service) GenericMethod(msg *nats.Msg) bool {
	return s.Publish(msg)
}

// Publish the actual  async method in a request pipeline
func (s *Service) Publish(msg *nats.Msg) bool {
	return nats2.Publish(s.publishToTopic, string(msg.Data), nil)
}
