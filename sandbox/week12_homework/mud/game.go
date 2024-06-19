package main

import "fmt"

// сюда писать код
// на сервер грузить только этот файл

// Commands dictionary
var _ = `
- осмотреться => описание локации (что есть, куда можно, что нужно, етц)
- идти $location => описание локации (состояние, куда можно пройти дальше)
	locations: коридор, комната, улица, кухня
- одеть $item => что было сделано
	items: рюкзак
- взять $item => что было сделано
	items: ключи, конспекты, телефон
- применить $item $object => что было сделано
	items: ключи, конспекты, телефон
	objects: дверь, шкаф
	examples:
		ключи дверь
		телефон шкаф
		ключи шкаф
- сказать $phrase => широковещание всем игрокам в локации, цитирование сказанного
	examples: Пора топать в универ
- сказать_игроку $player $phrase_option => адресное обращение к игроку в локации, цитирование
	examples:
		Izolda Может ещё по чаю лучше?
		Tristan
LookAroundCmd
GotoCmd
PutOnCmd
TakeCmd
ApplyCmd
ShoutCmd
SayToPlayerCmd
`

// Game logic, notes
var _ = `
- стейт игры это набор локаций
- игрок находится в текущей локации
- локации "связаны": игроку можно переходить из одной в другую
- локации содержат обьекты и айтемы, некоторое описание (или шаблоны описания)
- айтемы игрок может класть в "инвентарь", это меняет стейт локации
- айтемы из инвентаря игрок может применять к обьектам, это меняет стейт локации
- игрок может посылать другим игрокам (в своей локации) сообщения, бродкаст или персональное
`

// addPlayer: add player to game
func addPlayer(p IPlayer) {
	show("addPlayer: ", p)
	p.setLocation("кухня")
}

// initGame: create initial game state (locations, ...)
func initGame() {
	show("initGame ...")

	g := &Game{
		locations: make(map[string]*Location, 16),
	}

	// {18, "идти кухня", "кухня, ничего интересного. можно пройти - коридор"},
	// {1, "осмотреться", "ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"},
	loc := &Location{
		name:              "кухня",
		gotoDescription:   "кухня, ничего интересного. можно пройти - коридор",
		lookupDescription: "ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор",
		objects:           make([]*ObjectInLocation, 0, 16),
	}
	g.addLocation(loc)

	// {2, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"}
	loc = &Location{
		name:              "коридор",
		gotoDescription:   "ничего интересного. можно пройти - кухня, комната, улица",
		lookupDescription: "ты находишься в коридоре. можно пройти - кухня, комната, улица",
		objects:           make([]*ObjectInLocation, 0, 16),
	}
	door := ObjectInLocation{
		name:           "дверь",
		compatibleWith: []string{"ключи"},
		currentState:   "дверь закрыта",
		activatedState: "дверь открыта", // {9, "применить ключи дверь", "дверь открыта"}
	}
	loc.objects = append(loc.objects, &door)
	g.addLocation(loc)

	// {3, "идти комната", "ты в своей комнате. можно пройти - коридор"},
	// {4, "осмотреться", "на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор"},
	loc = &Location{
		name:              "комната",
		gotoDescription:   "ты в своей комнате. можно пройти - коридор",
		lookupDescription: "на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор",
		objects:           make([]*ObjectInLocation, 0, 16),
	}
	g.addLocation(loc)

	var I_AM_HERE = `
024-06-19T17:55:44.232Z: Player.HandleInput, cmd: "идти улица";
2024-06-19T17:55:44.232Z: Game.getLocationObj, search for loc: "улица";
--- FAIL: TestGameSingleplayer (0.01s)
panic: game getLocationObj failed: g.locations doesn't have location with name 'улица'

{11, "идти улица", "на улице весна. можно пройти - домой"}
	`

	game = g
}

type Game struct {
	locations map[string]*Location
}

// getLocation implements IGame.
func (g *Game) getLocation(name string) *Location {
	loc, isExists := g.locations[name]
	if isExists {
		return loc
	}
	panic(fmt.Sprintf("Game.getLocation, location '%s' not exists", name))
}

func (g *Game) addLocation(loc *Location) *Game {
	g.locations[loc.name] = loc
	return g
}

// getLookupLocationDescription implements IGame.
func (g *Game) getLookupLocationDescription(location string) string {
	locObj, err := g.getLocationObj(location)
	panicOnError("game getLocationObj failed", err)
	return locObj.getLookupDescription()
}

// getGoToLocationDescription implements IGame.
func (g *Game) getGoToLocationDescription(location string) string {
	locObj, err := g.getLocationObj(location)
	panicOnError("game getLocationObj failed", err)
	return locObj.getGoToDescription()
}

func (g *Game) getLocationObj(locName string) (*Location, error) {
	show("Game.getLocationObj, search for loc: ", locName)
	loc, isExists := g.locations[locName]
	if isExists {
		return loc, nil
	}

	return nil, fmt.Errorf("g.locations doesn't have location with name '%s'", locName)
}

// isLocationsConnected implements IGame.
func (g *Game) isLocationsConnected(currLoc string, targetLoc string) bool {
	return true
}

var game IGame = &Game{}

type IGame interface {
	// if game.isLocationsConnected(currLocation, targetLocatoin) {}
	isLocationsConnected(currLoc, targetLoc string) bool

	// player.commandReaction(game.getGoToLocationDescription(targetLocatoin))
	getGoToLocationDescription(location string) string

	// player.commandReaction(game.getLookupLocationDescription(currLocation))
	getLookupLocationDescription(location string) string

	// loc := game.getLocation(locName)
	getLocation(name string) *Location // TODO: make it interface
}

type ICommand interface {
	execute(game IGame, player IPlayer) error
}

type IPlayer interface {
	// HandleInput: pass command to player, player should produce output message.
	HandleInput(cmd string)

	// GetOutput: messages from player (responses to commands).
	GetOutput() chan string

	// commandReaction: produce output message
	commandReaction(msg string)

	// var currLocation string = player.getLocation()
	getLocation() string

	// player.setLocation(targetLocatoin)
	setLocation(targetLocation string)

	// if player.hasItem(a.applyItem)
	hasItem(item string) bool

	collectItem(item string)
}

type EmptyStruct struct{}