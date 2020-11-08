# Problem Statement (producer-consumer)

The producer/consumer system works as follows:

 

1. There are four goroutines producing "widgets". Each goroutine can produce
   between 1 and 3 widgets at a time, and will hold up to 10 widgets before
   waiting for widgets to be consumed. In this case, no more widgets are
   produced before some widgets are consumed.

2. There are two consumer goroutines. These goroutines consume between 1 and 3
   widgets from each of the producer goroutines at a time. Once a consumer
   goroutine acquires 10 widgets, the widgets are discarded, and that consumer
   goroutine sleeps for a random amount of time between 1 to 5 seconds.

3. If a producer is out of widgets at the time the consumer requests widgets
   from that producer, the producer is awakened and produces a random number
   of widgets between 1 and 3 widgets. If the consumer requests more widgets
   than the producer currently has, the consumer acquires all that producer's
   widgets and advances to the next producer.

4. Each producer is only "visited" once by each consumer before advancing to the
   next producer.

5. When a consumer acquires widgets from producer 4, it "wraps" back to producer
   1 next.

Create 6 goroutines that model the above behaviour. Run the simulation until both consumers
have consumed 10 widgets, 10 times each. Each time a producer produces widgets,
print the number of widgets produced. Each time a consumer consumes widgets,
print the number consumed. When a consumer acquires 10 widgets, print "discarded".

Ensure there are no deadlocks, and that neither consumer is
"starved" for widgets.

# Solution 

$ go run main.go 
