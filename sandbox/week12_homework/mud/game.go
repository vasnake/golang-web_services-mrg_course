package main

// сюда писать код
// на сервер грузить только этот файл

type IPlayer interface {
	GetOutput() chan string // or map[string]any
	HandleInput(inp string)
}

var _ IPlayer = Player{} // type check

type Player struct {
}

func NewPlayer(name string) *Player {
	// GetOutput could be `GetOutput() map[string]any`
	var foo = make(map[string]EmptyStruct, 16)
	for bar := range foo {
		var _ string = bar
	}

	show("new Player ...")
	return &Player{}
}

// HandleInput implements IPlayer.
func (p Player) HandleInput(cmd string) {
	show("Player.HandleInput: ", cmd)
}

// GetOutput implements IPlayer.
func (p Player) GetOutput() chan string {
	show("Player.GetOutput ...")
	ch := make(chan string)
	close(ch)
	return ch
}

func initGame() {
	show("initGame ...")
}

func addPlayer(p IPlayer) {
	show("addPlayer: ", p)
}

type EmptyStruct struct{}
