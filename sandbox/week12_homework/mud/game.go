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
- игрок может посылать другим игрокам (в своей локации) сообщения: бродкаст или персональное
`

// addPlayer: add player to game
func addPlayer(p IPlayer) {
	show("addPlayer: ", p)
	p.setLocation("кухня")
	game.addPlayer(p)
}

// initGame: create initial game state (locations, ...)
func initGame() {
	show("initGame ...")

	g := &Game{
		locations: make(map[string]*Location, 16),
		players:   make(map[string]IPlayer, 16),
	}

	loc := &Location{
		name:              "кухня",
		gotoDescription:   "кухня, ничего интересного. можно пройти - коридор",
		lookupDescription: "ты находишься на кухне, %s, надо собрать рюкзак и идти в универ. %s",
		conditionalLookupDescription: map[string]string{
			"рюкзак": "ты находишься на кухне, %s, надо идти в универ. %s",
		},
		connectedLocations: []string{"коридор"},
		objects:            make([]*ObjectInLocation, 0, 16),
		items: []*ItemInLocation{
			{name: "чай", prefix: "на столе "},
		},
	}
	g.addLocation(loc)

	loc = &Location{
		name:            "коридор",
		gotoDescription: "ничего интересного. можно пройти - кухня, комната, улица",
		gotoConditions: map[string]GotoCondition{
			"улица": {
				condition:        "дверь открыта",
				negativeReaction: "дверь закрыта",
				action:           "применить ключи дверь",
			},
		},
		lookupDescription: "ты находишься в коридоре. можно пройти - кухня, комната, улица",
		connectedLocations: []string{
			"кухня", "комната", "улица",
		},
		objects: []*ObjectInLocation{
			{
				name:           "дверь",
				compatibleWith: []string{"ключи"},
				currentState:   "дверь закрыта",
				activatedState: "дверь открыта", // {9, "применить ключи дверь", "дверь открыта"}
			},
		},
		items: []*ItemInLocation{},
	}
	g.addLocation(loc)

	// {3, "идти комната", "ты в своей комнате. можно пройти - коридор"},
	// {4, "осмотреться", "на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор"},
	loc = &Location{
		name:               "комната",
		gotoDescription:    "ты в своей комнате. можно пройти - коридор",
		lookupDescription:  "%s. %s",
		connectedLocations: []string{"коридор"},
		objects:            []*ObjectInLocation{},
		items: []*ItemInLocation{
			{name: "ключи", prefix: "на столе: "},
			{name: "конспекты", prefix: "на столе: "},
			{name: "рюкзак", prefix: "на стуле - "},
		},
	}
	g.addLocation(loc)

	loc = &Location{
		name:               "улица",
		gotoDescription:    "на улице весна. можно пройти - домой",
		lookupDescription:  "на улице весна. можно пройти - домой",
		connectedLocations: []string{"домой"},
		objects:            make([]*ObjectInLocation, 0, 16),
		items:              []*ItemInLocation{},
	}
	g.addLocation(loc)

	game = g

	_ = `
=== RUN   TestGameSingleplayer
2024-06-25T17:45:28.337Z: new Player "Tristan";
2024-06-25T17:45:28.337Z: initGame ...
2024-06-25T17:45:28.337Z: addPlayer: &main.Player{name:"Tristan", cmdResponses:(chan string)(0xc000022360), currentLocationName:"unknown", isBagReady:false, collectedItems:map[string]main.EmptyStruct{}, locationsStates:map[string]main.EmptyStruct{}};
2024-06-25T17:45:28.337Z: Player.HandleInput: "Tristan"; "кухня"; "осмотреться";
2024-06-25T17:45:28.337Z: Player.GetOutput ...
2024-06-25T17:45:28.338Z: Player.HandleInput: "Tristan"; "кухня"; "идти коридор";
2024-06-25T17:45:28.340Z: Player.HandleInput: "Tristan"; "коридор"; "идти комната";
2024-06-25T17:45:28.341Z: Player.HandleInput: "Tristan"; "комната"; "осмотреться";
2024-06-25T17:45:28.342Z: Player.HandleInput: "Tristan"; "комната"; "одеть рюкзак";
2024-06-25T17:45:28.343Z: Player.HandleInput: "Tristan"; "комната"; "взять ключи";
2024-06-25T17:45:28.345Z: Player.HandleInput: "Tristan"; "комната"; "взять конспекты";
2024-06-25T17:45:28.346Z: Player.HandleInput: "Tristan"; "комната"; "идти коридор"; 
2024-06-25T17:45:28.347Z: Player.HandleInput: "Tristan"; "коридор"; "применить ключи дверь";
2024-06-25T17:45:28.348Z: Player.HandleInput: "Tristan"; "коридор"; "идти улица";

2024-06-25T17:45:28.349Z: new Player "Tristan";
2024-06-25T17:45:28.349Z: initGame ...
2024-06-25T17:45:28.349Z: addPlayer: &main.Player{name:"Tristan", cmdResponses:(chan string)(0xc0000223c0), currentLocationName:"unknown", isBagReady:false, collectedItems:map[string]main.EmptyStruct{}, locationsStates:map[string]main.EmptyStruct{}};
2024-06-25T17:45:28.349Z: Player.HandleInput: "Tristan"; "кухня"; "осмотреться";
2024-06-25T17:45:28.349Z: Player.GetOutput ...
2024-06-25T17:45:28.351Z: Player.HandleInput: "Tristan"; "кухня"; "завтракать";
2024-06-25T17:45:28.352Z: Player.HandleInput: "Tristan"; "кухня"; "идти комната";
2024-06-25T17:45:28.353Z: Player.HandleInput: "Tristan"; "кухня"; "идти коридор";
2024-06-25T17:45:28.355Z: Player.HandleInput: "Tristan"; "коридор"; "применить ключи дверь"; 
2024-06-25T17:45:28.356Z: Player.HandleInput: "Tristan"; "коридор"; "идти комната";
2024-06-25T17:45:28.357Z: Player.HandleInput: "Tristan"; "комната"; "осмотреться";
2024-06-25T17:45:28.358Z: Player.HandleInput: "Tristan"; "комната"; "взять ключи";
2024-06-25T17:45:28.360Z: Player.HandleInput: "Tristan"; "комната"; "одеть рюкзак";
2024-06-25T17:45:28.361Z: Player.HandleInput: "Tristan"; "комната"; "осмотреться";
2024-06-25T17:45:28.362Z: Player.HandleInput: "Tristan"; "комната"; "взять ключи";
2024-06-25T17:45:28.364Z: Player.HandleInput: "Tristan"; "комната"; "взять телефон";
2024-06-25T17:45:28.365Z: Player.HandleInput: "Tristan"; "комната"; "взять ключи";
2024-06-25T17:45:28.366Z: Player.HandleInput: "Tristan"; "комната"; "осмотреться";
2024-06-25T17:45:28.368Z: Player.HandleInput: "Tristan"; "комната"; "взять конспекты";
2024-06-25T17:45:28.369Z: Player.HandleInput: "Tristan"; "комната"; "осмотреться";
2024-06-25T17:45:28.370Z: Player.HandleInput: "Tristan"; "комната"; "идти коридор";
2024-06-25T17:45:28.371Z: Player.HandleInput: "Tristan"; "коридор"; "идти кухня"; 
2024-06-25T17:45:28.373Z: Player.HandleInput: "Tristan"; "кухня"; "осмотреться";
2024-06-25T17:45:28.374Z: Player.HandleInput: "Tristan"; "кухня"; "идти коридор";
2024-06-25T17:45:28.376Z: Player.HandleInput: "Tristan"; "коридор"; "идти улица";
2024-06-25T17:45:28.377Z: Player.HandleInput: "Tristan"; "коридор"; "применить ключи дверь";
2024-06-25T17:45:28.378Z: Player.HandleInput: "Tristan"; "коридор"; "применить телефон шкаф";
2024-06-25T17:45:28.379Z: Player.HandleInput: "Tristan"; "коридор"; "применить ключи шкаф";
2024-06-25T17:45:28.381Z: Player.HandleInput: "Tristan"; "коридор"; "идти улица";
--- PASS: TestGameSingleplayer (0.04s)

=== RUN   TestGameMiltiplayer
2024-06-25T17:45:28.382Z: new Player "Tristan";
2024-06-25T17:45:28.382Z: new Player "Izolda";
2024-06-25T17:45:28.382Z: initGame ...
2024-06-25T17:45:28.382Z: addPlayer: &main.Player{name:"Tristan", cmdResponses:(chan string)(0xc0000ae120), currentLocationName:"unknown", isBagReady:false, collectedItems:map[string]main.EmptyStruct{}, locationsStates:map[string]main.EmptyStruct{}};
2024-06-25T17:45:28.382Z: addPlayer: &main.Player{name:"Izolda", cmdResponses:(chan string)(0xc0000ae180), currentLocationName:"unknown", isBagReady:false, collectedItems:map[string]main.EmptyStruct{}, locationsStates:map[string]main.EmptyStruct{}};
2024-06-25T17:45:28.382Z: Player.HandleInput: "Tristan"; "кухня"; "осмотреться";
2024-06-25T17:45:28.382Z: Player.GetOutput ...
2024-06-25T17:45:28.382Z: Player.GetOutput ...
2024-06-25T17:45:28.384Z: Player.HandleInput: "Izolda"; "кухня"; "осмотреться";
2024-06-25T17:45:28.385Z: Player.HandleInput: "Izolda"; "кухня"; "сказать Пора топать в универ";
2024-06-25T17:45:28.387Z: Player.HandleInput: "Tristan"; "кухня"; "сказать_игроку Izolda Может ещё по чаю лучше?"; 
2024-06-25T17:45:28.388Z: Player.HandleInput: "Izolda"; "кухня"; "сказать_игроку Tristan";
2024-06-25T17:45:28.389Z: Player.HandleInput: "Tristan"; "кухня"; "идти коридор";
2024-06-25T17:45:28.391Z: Player.HandleInput: "Izolda"; "кухня"; "сказать_игроку Tristan";
--- PASS: TestGameMiltiplayer (0.01s)

	`

}

type Game struct {
	locations map[string]*Location
	players   map[string]IPlayer
}

// getPlayersInLocation implements IGame.
func (g *Game) getPlayersInLocation(locName string) []IPlayer {
	res := make([]IPlayer, 0, len(g.players))
	for _, p := range g.players {
		if p.getLocation() == locName {
			res = append(res, p)
		}
	}
	return res
}

// addPlayer implements IGame.
func (g *Game) addPlayer(p IPlayer) {
	g.players[p.getName()] = p
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
func (g *Game) getLookupLocationDescription(location string, player IPlayer) string {
	var _ = `
кухня:
надо собрать рюкзак и идти в универ
взял рюкзак =>
надо идти в универ	

нужен стейт игрока и проверка этого стейта (рюкзак собран) на условие.
в зависимости от проверки меняется дескрипшн локации.
правила игры, стейт локации, стейт игрока.
стейт игры зависит от стейта (игрок, локация).
	`

	locObj, err := g.getLocationObj(location)
	panicOnError("game getLocationObj failed", err)

	for key, descr := range locObj.conditionalLookupDescription {
		conditionMet := player.hasItem(key)
		if conditionMet {
			return locObj.getLookupDescription(descr)
		}
	}
	return locObj.getLookupDescription(locObj.lookupDescription)
}

// getGoToLocationDescription implements IGame.
func (g *Game) getGoToLocationDescription(location string) string {
	locObj, err := g.getLocationObj(location)
	panicOnError("game getLocationObj failed", err)
	return locObj.getGoToDescription()
}

func (g *Game) getLocationObj(locName string) (*Location, error) {
	loc, isExists := g.locations[locName]
	if isExists {
		return loc, nil
	}

	return nil, fmt.Errorf("g.locations doesn't have location with name '%s'", locName)
}

// isLocationsConnected implements IGame.
func (g *Game) isLocationsConnected(currLoc string, targetLoc string) bool {
	cLoc := g.getLocation(currLoc)
	tLoc := g.getLocation(targetLoc)
	return cLoc.isLocationsConnected(tLoc)
}

var game IGame = &Game{} // global state

type IGame interface {
	// if game.isLocationsConnected(currLocation, targetLocatoin) {}
	isLocationsConnected(currLoc, targetLoc string) bool

	// player.commandReaction(game.getGoToLocationDescription(targetLocatoin))
	getGoToLocationDescription(location string) string

	getLookupLocationDescription(location string, player IPlayer) string

	// loc := game.getLocation(locName)
	getLocation(name string) *Location // TODO: make it interface

	addPlayer(p IPlayer)

	getPlayersInLocation(locName string) []IPlayer
}

type ICommand interface {
	execute(game IGame, player IPlayer) error
}

type IPlayer interface {
	getName() string

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

	hasBag() bool

	// if player.isStateInLocationState("коридор", "дверь открыта")
	isStateInLocationState(locName, state string) bool
	setStateInLocationState(locName string, state string)
}

type EmptyStruct struct{}
