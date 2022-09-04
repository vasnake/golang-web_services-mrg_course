package main

import (
	"fmt"
)

// SM as interface to some implementation
var sessManager SessionManagerI

func main() {
	// use interface (and fabric) to get an object with implementation

	sessManager = NewSessManager()

	// создаем сессию
	sessId, err := sessManager.Create(
		&Session{
			Login:     "rvasily",
			Useragent: "chrome",
		})
	fmt.Println("sessId", sessId, err)

	// проеряем сессию
	sess := sessManager.Check(
		&SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess)

	// удаляем сессию
	sessManager.Delete(
		&SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = sessManager.Check(
		&SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess)

}
