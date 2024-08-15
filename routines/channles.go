package main

import . "fmt"

func main()  {
	myChannel := make(chan string) // create a channel

	// create a go routine that will send a message to the channel
	go func() {
		Println("Sending message to channel...")
		myChannel <- "Data" // send a message to the channel
	}()

	Println("Channel => ", myChannel) // Output: Channel =>  0x14000120000

	msg := <- myChannel // receive a message from the channel

	Println("Message => ", msg)
}