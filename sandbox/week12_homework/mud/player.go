package main

var _ IPlayer = &Player{} // type check

// NewPlayer: create new named player
func NewPlayer(name string) *Player {
	// GetOutput could be `GetOutput() map[string]any`
	// var foo = make(map[string]EmptyStruct, 16)
	// for bar := range foo {
	// 	var _ string = bar
	// }

	show("new Player ...")

	return &Player{
		cmdResponses:        make(chan string),
		currentLocationName: "unknown",
		collectedItems:      make(map[string]EmptyStruct, 16),
	}
}

// Player: IPlayer implementation
type Player struct {
	cmdResponses        chan string
	currentLocationName string
	collectedItems      map[string]EmptyStruct
}

func (p *Player) collectItem(item string) {
	p.collectedItems[item] = EmptyStruct{}
}

// hasItem implements IPlayer.
func (p *Player) hasItem(item string) bool {
	_, isExists := p.collectedItems[item]
	return isExists
}

// setLocation implements IPlayer.
func (p *Player) setLocation(targetLocation string) {
	p.currentLocationName = targetLocation
}

// getLocation implements IPlayer.
func (p *Player) getLocation() string {
	return p.currentLocationName
}

// commandReaction implements IPlayer.
func (p *Player) commandReaction(msg string) {
	p.cmdResponses <- msg
}

// HandleInput implements IPlayer.
// Process command and generate output message
func (p *Player) HandleInput(cmd string) {
	show("Player.HandleInput, cmd: ", cmd)

	typedCmd := parseCommand(cmd)
	if typedCmd == nil {
		panic("parseCommand failed, unknown command: " + cmd)
	}

	err := typedCmd.execute(game, p)
	panicOnError("command.execute failed", err)
}

// GetOutput implements IPlayer.
// Return reference to player output messages chan.
func (p *Player) GetOutput() chan string {
	show("Player.GetOutput ...")
	return p.cmdResponses
}
