package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)

	i := 0
	for tick := range ticker.C { // signal in chan Time each second
		i++
		fmt.Println("step", i, "time", tick)

		if i >= 5 {
			// надо останавливать, иначе потечет
			ticker.Stop()
			break
		}
	}
	fmt.Println("total ticks", i)

	// return

	// не может быть остановлен и собран сборщиком мусора
	// используйте только если должен работать вечно
	c := time.Tick(time.Second)
	i = 0
	for tickTime := range c { // signal every second
		i++
		fmt.Println("step", i, "time", tickTime)

		if i >= 5 {
			break
		}
	}

}
