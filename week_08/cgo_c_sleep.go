/*
	cgo блокирует весь системный тред
	golang-рантайм не может больше в нём запускать никакие другие горутины
	если в СИ операция была блокирующей, например, sleep, то она заблокирует весь тред
	после запуска надо посомтреть сколько тредов запущено процессом
	будет больше чем в cgo_go_sleep
*/

package main

//#include <unistd.h>
import "C"
import (
	"time"
)

func main() {
	for i := 0; i < 100; i++ {
		go func() {
			// запускаем СИшный sleep
			C.sleep(60 * 10) // 10 минут
			// lock current thread, for other work Go start another thread
		}()
	}
	time.Sleep(11 * time.Minute)
}
