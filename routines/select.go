package main

import "fmt"
import "time"

func main()  {
	myChannel1 := make(chan string) // create a channel
	myChannel2 := make(chan int) // create a channel

	// create a go routine that will send a message to the channel
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("Sending message to channel1...")
		myChannel1 <- "Data" // send a message to the channel
	}()

	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("Sending message to channel2...")
		myChannel2 <- 1 // send a message to the channel
	}()

	for i := 0; i < 2; i++ {
		// The select statement in Go is used to wait on multiple channel operations
		select {
		case msg1 := <- myChannel1:
			fmt.Println("Message1 => ", msg1)
		case msg2 := <- myChannel2:
			fmt.Println("Message2 => ", msg2)
		}
	}

}