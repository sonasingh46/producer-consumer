package consumer

import (
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
It stores the widgets consumed in its store of size 10.
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
	workers []producer.Producer
	// store is the store where the items consumed is stored.
	store []widget.Widget
	// TenCompletion is the number of time the consumer has consumed
	// 10 items. It should be noted that once the consumer consumes
	// 10 items -- the consumer store is reset.
	TenCompletion int
}

// NewWidgetConsumer returns a new widget consumer.
func NewWidgetConsumer(id int, workerPool []producer.Producer)Consumer{
	return &WidgetConsumer{
		id:id,
		currentProducerID:1,
		workers:workerPool,
		store:make([]widget.Widget,0),
	}
}

// Consume consumes widget(s) from producer.
func (wc *WidgetConsumer)Consume(wg *sync.WaitGroup){

	for {
		// If the consumer has consumed 10 widgets 10 times
		// each then exit it.
		if wc.TenCompletion == 10 {
			fmt.Printf("[Consumer %d]: Consumed 10 widgets 10 times. Exiting.\n",wc.id)
			break
		}

		// If consumer store is full -- discard all the widgets and
		// sleep for random amount of time.
		if len(wc.store)==10{
			wc.TenCompletion++
			fmt.Printf("[Consumer %d]: Discarded. Ten Completion Count : %d \n",
				wc.id,wc.TenCompletion)
			wc.store=make([]widget.Widget,0)
			timeToSleep:=random.GetRandomNumberInRange(1,5)
			time.Sleep(time.Duration(timeToSleep)*time.Second)
		}
		// Find how many widgets to consume
		widgetsToConsume:=random.GetRandomNumberInRange(1,3)
		// If widgets to consume is more than consumer store size then
		// advance to next producer.
		if widgetsToConsume > 10 - len(wc.store){
			wc.setToNextProducerID()
			continue
		}

		currentProducer:=wc.getProducer()
		widgetList:=currentProducer.Extract(widgetsToConsume)
		wc.store = append(wc.store,widgetList...)
		wc.setToNextProducerID()
		fmt.Printf("[Consumer %d]: %d Widgets consumed from %d producer\n",
			wc.id,len(widgetList),currentProducer.GetID())
	}

	defer wg.Done()
}

// setToNextProducerID sets the consumer to a new producer ID.
func (wc* WidgetConsumer)setToNextProducerID()  {
	nextID:=wc.currentProducerID+1
	if nextID > len(wc.workers){
		wc.currentProducerID = 1
	}else {
		wc.currentProducerID = nextID
	}
}

// getProducer returns the current producer for the consumer.
func (wc* WidgetConsumer)getProducer() producer.Producer {
	return wc.workers[wc.currentProducerID-1]
}