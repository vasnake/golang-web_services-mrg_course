package main

import (
	"fmt"
)

func main() {
	cancelCh := make(chan bool)
	dataCh := make(chan int)

	go func(cancelCh chan bool, dataCh chan int) {
		val := 0
		for {
			select {

			case <-cancelCh: // if got cancel signal
				println("closing dataCh")
				close(dataCh)
				return

			case dataCh <- val: // write to data
				fmt.Println("data sent", val)
				val++

			}
		}
	}(cancelCh, dataCh)

	for curVal := range dataCh {
		fmt.Println("got from dataCh", curVal)
		if curVal > 3 {
			fmt.Println("sending cancel ...")
			cancelCh <- true
			// break
		}
	}

}
