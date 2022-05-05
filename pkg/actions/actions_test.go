package actions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client1Event = Event{
	Name: "client1",
	Meta: nil,
}
var client2Event = Event{
	Name: "client2",
	Meta: nil,
}

func Test_Marshal(t *testing.T) {
	a := New()
	a.AddEvent(client1Event)
	a.AddEvent(Event{
		Name: "one",
		Meta: "meta",
	})

	b, e := a.Marshal()
	assert.NoError(t, e)

	a2, e2 := UnMarshal(b)
	assert.NoError(t, e2)

	assert.Equal(t, a.GetEvents(), a2.GetEvents())
}

func TestNewActions(t *testing.T) {
	a := New()
	e1 := client1Event
	e2 := client2Event

	a.AddEvent(e1)
	ctxA := NewCtx(context.Background(), a)

	r := FromCtx(ctxA)
	assert.Equal(t, e1, r.GetEvents()[0])
	r.AddEvent(e2)

	ctxR := NewCtx(context.Background(), r)
	f := FromCtx(ctxR)

	assert.Equal(t, e1, f.GetEvents()[0])
	assert.Equal(t, e2, f.GetEvents()[1])

}
