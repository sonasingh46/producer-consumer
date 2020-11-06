package main

import (
	"fmt"
	"github.com/sonasingh46/producer-consumer/consumer"
	"github.com/sonasingh46/producer-consumer/producer"
	"sync"
)

const (
	ProducerCount = 4
	consumerCount = 4
)

var wg sync.WaitGroup

func main() {
	runSimulation(1)
}

func runSimulation(simulationCount int) {
	for i := 1; i <= simulationCount; i++ {
		fmt.Printf("*****************[%d]SIMULATION STARTED****************************\n",i)
		var workers []producer.Producer
		for i := 1; i <= ProducerCount; i++ {
			p := producer.NewWidgetProducer(i, 10)
			workers = append(workers, p)
			go p.Produce()
		}

		for i := 1; i <= consumerCount; i++ {
			c := consumer.NewWidgetConsumer(i, workers)
			wg.Add(1)
			go c.Consume(&wg)
		}
		wg.Wait()
		fmt.Println("*****************[%d]SIMULATION ENDED*****************************\n",i)
	}
}
