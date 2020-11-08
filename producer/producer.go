package producer

import (
	"fmt"
	"github.com/sonasingh46/producer-consumer/pkg/random"
	"github.com/sonasingh46/producer-consumer/widget"
)

// Producer is the interface that wraps the Producer methods.
type Producer interface {
	Produce()
	Extract(widgetCount int) []widget.Widget
	GetID() int
}

/*

WidgetProducer keeps producing widgets to whole of its lifetime.
It can produce random number of widget b/w 1 to 3.

If number of widgets to be produced is greater than the available
queue space then the producer blocks until it can produce new widgets.
Producer will again start producing widget once it has enough queue
size available.
Producer will produce random number of widgets at a time and
then only move to a new random number of widgets to be produced.
For example, if a producer has to produce 2 widgets, and it has
only 1 space available then it will push 1 widget in the queue
and block until the remaining 1 can be pushed. Once the required
number of widgets ( in this example it is 2) are pushed to the
queue then only a new random number is generated to tell the
number of widgets to be produced.
See -- Produce method implementation that does this above story.
*/

// WidgetProducer is the concrete type that implements Producer
// interface
type WidgetProducer struct {
	// id is the id the producer.
	id int
	// queue is the queue where widgets are stored.
	queue chan widget.Widget
}

// NewWidgetProducer returns a new widget producer instance
func NewWidgetProducer(id, bufferSize int) Producer {
	return &WidgetProducer{
		id:    id,
		queue: make(chan widget.Widget, bufferSize),
	}
}

// Produce produces widgets.
func (wp *WidgetProducer) Produce() {
	for {
		widgetsToProduce := random.GetRandomNumberInRange(1, 3)
		temp := widgetsToProduce
		for widgetsToProduce != 0 {
			select {
			case wp.queue <- widget.NewWidget():
				widgetsToProduce--
			}
		}
		fmt.Printf("[Producer %d]: %d Widgets prodcued\n", wp.id, temp)
	}
}

// Extract extracts widget from a producer.
// This method should be called by a consumer that intends to
// consume form this producer.
func (wp *WidgetProducer) Extract(widgetCount int) []widget.Widget {
	widget := make([]widget.Widget, 0)
	for widgetCount != 0 {
		select {
		case wg := <-wp.queue:
			widget = append(widget, wg)
			widgetCount--
		}
	}
	return widget
}

// Extract extracts widget from a producer.
// This method should be called by a consumer that intends to
// consume form this producer.
func (wp *WidgetProducer) GetID() int {
	return wp.id
}
