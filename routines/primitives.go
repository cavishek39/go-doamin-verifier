package main
import . "fmt"
import "time"

func goRoutine(
	a int,
)  {
	Println("Value => ", a)
}

// Start here
func main()  {

	go goRoutine(1) // fork -> child process... goRoutine(1) will be executed in a separate go routine
	go goRoutine(2) // fork -> child process... goRoutine(2) will be executed in a separate go routine
	go goRoutine(3) // fork -> child process... goRoutine(3) will be executed in a separate go routine

	time.Sleep(time.Second * 2)
	// join the go routines

	Println("Hello, World!")
}

// End here