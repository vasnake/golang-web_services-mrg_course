/*
	в golang sleep не блокирует системный тред
	этот пример надо смотреть перед cgo_sleep
	после запуска надо посомтреть сколько тредов запущено процессом
*/

package main

import (
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		go func() {
			// запускаем ГОшный sleep
			time.Sleep(time.Minute * 10)
			// don't lock current thread, event loop select another work, small amount of threads needed
		}()
	}
	time.Sleep(time.Minute * 11)
}
