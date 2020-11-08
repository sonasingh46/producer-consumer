package main

import (
	"fmt"
	"github.com/sonasingh46/producer-consumer/consumer"
	"github.com/sonasingh46/producer-consumer/producer"
	"k8s.io/kubernetes/pkg/kubelet/kubeletconfig/util/log"
	"sync"
)

const (
	ProducerCount = 4
	consumerCount = 2
)

var wg sync.WaitGroup

func main() {
	runSimulation(1)
}

func runSimulation(simulationCount int) {
	for i := 1; i <= simulationCount; i++ {
		fmt.Printf("*****************[%d]SIMULATION STARTED****************************\n", i)
		var workers []producer.Producer
		for i := 1; i <= ProducerCount; i++ {
			// bufferSize id the number of widgets the producer
			// can hold up before waiting for widgets to be consumed.
			bufferSize := 10
			p := producer.NewWidgetProducer(i, bufferSize)
			workers = append(workers, p)
			go p.Produce()
		}

		for i := 1; i <= consumerCount; i++ {

			newWidgetConsumer,err := consumer.NewWidgetConsumer().
				WithID(i).
				// The consumer can store 10 widgets before it discards
				WithStoreCapacity(10).
				// The producers pool from which the consumer will
				// acquire widgets in a cyclic manner.
				WithWorkerPool(workers).
				// The consumer will discard 10 times and then exit.
				WithCycleCountThreshold(10).
				Build()
			if err!=nil{
				log.Errorf("Could not launch %d consumer:{%s}",i,err.Error())
			}
			wg.Add(1)
			go newWidgetConsumer.Consume(&wg)
		}
		wg.Wait()
		fmt.Printf("*****************[%d]SIMULATION ENDED*****************************\n", i)
	}
}
