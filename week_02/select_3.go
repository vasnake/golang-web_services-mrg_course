package main

import (
	"fmt"
)

func main() {
	cancelCh := make(chan bool) // commands channel
	dataCh := make(chan int)

	go func(cancelCh chan bool, dataCh chan int) { // call anon. func in async mode
		val := 0
		for {
			select { // N.B. no default branch, you have to be sure that someone reading data channel
			// Think about protocol here ...

			case <-cancelCh: // if got cancel signal
				println("closing dataCh")
				close(dataCh) // help reader to stop iterations in `range dataCh`
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
