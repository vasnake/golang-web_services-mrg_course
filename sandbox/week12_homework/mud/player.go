package main

import "fmt"

var _ IPlayer = &Player{} // type check

// NewPlayer: create new named player
func NewPlayer(name string) *Player {
	show("new Player ", name)

	return &Player{
		name:                name,
		cmdResponses:        make(chan string),
		currentLocationName: "unknown",
		isBagReady:          false,
		collectedItems:      make(map[string]EmptyStruct, 16),
		locationsStates:     make(map[string]EmptyStruct, 16),
	}
}

// Player: IPlayer implementation
type Player struct {
	name                string
	cmdResponses        chan string
	currentLocationName string
	isBagReady          bool // inventory
	collectedItems      map[string]EmptyStruct
	locationsStates     map[string]EmptyStruct
}

// getName implements IPlayer.
func (p *Player) getName() string {
	return p.name
}

// isStateInLocationState implements IPlayer.
func (p *Player) isStateInLocationState(locName string, state string) bool {
	// if player.isStateInLocationState("коридор", "дверь открыта")
	_, isIn := p.locationsStates[fmt.Sprintf("%s : %s", locName, state)] // коридор : дверь открыта
	return isIn
}

func (p *Player) setStateInLocationState(locName string, state string) {
	// player.setStateInLocationState("коридор", "дверь открыта")
	p.locationsStates[fmt.Sprintf("%s : %s", locName, state)] = EmptyStruct{} // коридор : дверь открыта
}

// hasBag implements IPlayer.
func (p *Player) hasBag() bool {
	return p.isBagReady
}

func (p *Player) collectItem(item string) {
	p.collectedItems[item] = EmptyStruct{}
	if item == "рюкзак" {
		p.isBagReady = true
	}
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
	show("Player.HandleInput: ", p.getName(), p.getLocation(), cmd)

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
