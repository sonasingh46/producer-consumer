package consumer

import (
	"errors"
	"fmt"
	"github.com/sonasingh46/producer-consumer/pkg/random"
	"github.com/sonasingh46/producer-consumer/producer"
	"github.com/sonasingh46/producer-consumer/widget"
	"sync"
	"time"
)

// Consumer is the interface that wraps the Consumer methods.
type Consumer interface {
	Consume(wg *sync.WaitGroup)
}

/*
WidgetConsumer consumes a random number widgets b/w 1 to 3 from a producer.
It stores the widgets consumed in its store of size x.
Once the consumer store is full, all the widgets are discarded from the
store and the consumer sleeps for a random period of time b/w 1 to 5
seconds.

If the number of widgets to be consumed is greater than the store size
of the WidgetConsumer then the WidgetConsumer does not consume any widget
and moves to next producer.

WidgetConsumer has reference to a list of producers and hence the queue of
the producer.
If no of widgets to be consumed is more than the widgets in the
producer queue than all the widgets are consumed.

WidgetConsumer tries to consume widgets from the producer and then advances
to the next producer to consume widgets.

If there are no available widgets to be consumed from a producer, the
consumer advances to the next producer.

When the consumer is done consuming widgets from last producer in the
list either successfully or unsuccessfully it starts from the first
producer in the list.

See consume method -- that implements the above story.
*/
type WidgetConsumer struct {
	// id is the id of the consumer.
	id int
	// currentProducerID is the id of the producer from which
	// the consumer will consume at a specific time.
	currentProducerID int
	// workers is the list of producers from which the consumer
	// will consume in a cyclic fashion.
	workers       []producer.Producer
	storeCapacity int
	// store is the store where the items consumed is stored.
	store []widget.Widget
	// CycleCount is the number of times the consumer has consumed
	// amount of widgets equal to the capacity of store. It should
	// be noted that once the consumer consumes amount of widgets
	// equal to the capacity of store -- the consumer store is reset.
	CycleCount int
	// CycleCountThreshold is the number of Cycle beyond which the
	// consumer won't go. If the consumer reached the CycleCountThreshold
	// it is done and exits.
	CycleCountThreshold int
}

// NewWidgetConsumer returns an empty instance of WidgetConsumer.
func NewWidgetConsumer() *WidgetConsumer {
	return &WidgetConsumer{}
}

// WithID sets the id of the widget consumer
func (wc *WidgetConsumer) WithID(id int) *WidgetConsumer {
	wc.id = id
	return wc
}

// WithWorkerPool sets the worker pool reference to the
// consumer. The consumer will consume from this worker
// pool.
func (wc *WidgetConsumer) WithWorkerPool(workerPool []producer.Producer) *WidgetConsumer {
	wc.workers = workerPool
	return wc
}

// WithStore sets the store for the consumer.
func (wc *WidgetConsumer) WithStoreCapacity(capacity int) *WidgetConsumer {
	wc.storeCapacity = capacity
	store := make([]widget.Widget, 0, capacity)
	wc.store = store
	return wc
}

// WithCycleCount sets the cycle count for the consumer.
func (wc *WidgetConsumer) WithCycleCountThreshold(limit int) *WidgetConsumer {
	wc.CycleCountThreshold = limit
	return wc
}

// Build returns the built widget consumer.
func (wc *WidgetConsumer) Build() (Consumer,error) {
	if wc.CycleCountThreshold>wc.storeCapacity{
		return nil,errors.New("Cycle count threshold cannot be greater than store capacity")
	}
	wc.currentProducerID = 1
	return wc,nil
}

// Consume consumes widget(s) from producer.
func (wc *WidgetConsumer) Consume(wg *sync.WaitGroup) {

	for {
		// If the consumer has consumed storeCapacity number of widgets
		// CycleCountThreshold times then exit it.
		if wc.CycleCount == wc.CycleCountThreshold {
			fmt.Printf("[Consumer %d]: Consumed %d widgets %d times. Exiting.\n",
				wc.id, wc.storeCapacity, wc.CycleCount)
			wc.CycleCount = 0
			break
		}

		// If consumer store is full -- discard all the widgets and
		// sleep for random amount of time.
		if len(wc.store) == cap(wc.store) {
			wc.CycleCount++
			fmt.Printf("[Consumer %d]: Discarded. Cycle Completion Count : %d\n",
				wc.id, wc.CycleCount)
			wc.store = make([]widget.Widget, 0, wc.storeCapacity)
			timeToSleep := random.GetRandomNumberInRange(1, 5)
			time.Sleep(time.Duration(timeToSleep) * time.Second)
			// ToDo: think about avoiding following condition and merge with above
			// same condition.
			if wc.CycleCount == wc.CycleCountThreshold{
				continue
			}
		}

		// Find how many widgets to consume
		widgetsToConsume := random.GetRandomNumberInRange(1, 3)
		// If widgets to consume is more than consumer store size then
		// advance to next producer.
		if widgetsToConsume > cap(wc.store)-len(wc.store) {
			wc.setToNextProducerID()
			continue
		}

		currentProducer := wc.getProducer()
		widgetList := currentProducer.Extract(widgetsToConsume)
		wc.store = append(wc.store, widgetList...)
		wc.setToNextProducerID()
		fmt.Printf("[Consumer %d]: %d Widgets consumed from producer %d\n",
			wc.id, len(widgetList), currentProducer.GetID())
	}

	defer wg.Done()
}

// setToNextProducerID sets the consumer to a new producer ID.
func (wc *WidgetConsumer) setToNextProducerID() {
	nextID := wc.currentProducerID + 1
	if nextID > len(wc.workers) {
		wc.currentProducerID = 1
	} else {
		wc.currentProducerID = nextID
	}
}

// getProducer returns the current producer for the consumer.
func (wc *WidgetConsumer) getProducer() producer.Producer {
	return wc.workers[wc.currentProducerID-1]
}
