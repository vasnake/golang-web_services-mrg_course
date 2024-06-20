package main

import (
	"fmt"
	"strings"
)

func parseCommand(cmd string) ICommand {
	switch {

	case strings.HasPrefix(cmd, "осмотреться"):
		return &LookAroundCmd{}

	case strings.HasPrefix(cmd, "идти "):
		loc := strings.TrimPrefix(cmd, "идти ")
		return &GotoCmd{locationName: loc}

	case strings.HasPrefix(cmd, "одеть "):
		item := strings.TrimPrefix(cmd, "одеть ")
		return &PutOnCmd{item: item}

	case strings.HasPrefix(cmd, "взять "):
		item := strings.TrimPrefix(cmd, "взять ")
		return &TakeCmd{item: item}

	case strings.HasPrefix(cmd, "применить "):
		item_object := strings.TrimPrefix(cmd, "применить ")
		pair := strings.Split(item_object, " ")
		if len(pair) == 2 {
			return &ApplyCmd{
				applyItem: pair[0],
				toObject:  pair[1],
			}
		} else {
			show("parseCommand, invalid apply command: ", cmd)
			return nil
		}

	default:
		return &UnknownCmd{}

	}
}

var _ ICommand = &UnknownCmd{} // type check
// {2, "завтракать", "неизвестная команда"}
type UnknownCmd struct{}

// execute implements ICommand.
func (u *UnknownCmd) execute(game IGame, player IPlayer) error {
	player.commandReaction("неизвестная команда")
	return nil
}

var _ ICommand = &LookAroundCmd{} // type check

// - осмотреться => описание локации (что есть, куда можно, что нужно, етц)
type LookAroundCmd struct{}

// execute implements ICommand.
func (l *LookAroundCmd) execute(game IGame, player IPlayer) error {
	var currLocation string = player.getLocation()
	player.commandReaction(game.getLookupLocationDescription(currLocation))

	// player.commandReaction("ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор")
	// case: 0 4
	// cmd: осмотреться
	// expected: на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор

	return nil
}

var _ ICommand = &GotoCmd{} // type check

//   - идти $location => описание локации (состояние, куда можно пройти дальше)
//     locations: коридор, комната, улица, кухня
type GotoCmd struct {
	locationName string
}

// execute implements ICommand.
func (g *GotoCmd) execute(game IGame, player IPlayer) error {
	var _ = `
игрок находится в текущей $location (при переходе она меняется на целевую);
текущая локация содержит список связанных (куда можно пройти) локаций;
целевая локация содержит описание, которое надо выдать (реакция);
	`
	var targetLocation string = g.locationName
	var currLocation string = player.getLocation()
	if currLocation == targetLocation {
		return fmt.Errorf("GotoCmd failed, target == current: %s", currLocation)
	}

	if game.isLocationsConnected(currLocation, targetLocation) {
		player.setLocation(targetLocation)
		player.commandReaction(game.getGoToLocationDescription(targetLocation))
	} else {
		// return fmt.Errorf("GotoCmd failed, locations are not connected, curr %s, target %s", currLocation, targetLocatoin)
		player.commandReaction(fmt.Sprintf("нет пути в %s", targetLocation))
	}

	return nil
}

var _ ICommand = &PutOnCmd{} // type check

//   - одеть $item => что было сделано
//     items: рюкзак
type PutOnCmd struct {
	// return &PutOnCmd{item: item}
	item string
}

// execute implements ICommand.
func (p *PutOnCmd) execute(game IGame, player IPlayer) error {
	// {5, "одеть рюкзак", "вы одели: рюкзак"}
	var _ = `
айтем (если есть в локации)
добавить в инвентарь плеера;
если взял успешно, убрать айтем из локации;
	`
	player.collectItem(p.item)
	// location.removeItem(p.item)
	player.commandReaction(fmt.Sprintf("вы одели: %s", p.item))
	return nil
}

var _ ICommand = &TakeCmd{} // type check

//   - взять $item => что было сделано
//     items: ключи, конспекты, телефон
type TakeCmd struct {
	item string
}

// execute implements ICommand.
func (t *TakeCmd) execute(game IGame, player IPlayer) error {
	// {6, "взять ключи", "предмет добавлен в инвентарь: ключи"}
	// {8, "взять ключи", "некуда класть"}
	var _ = `
айтем (если есть в локации)
добавить в инвентарь плеера (если есть бэг);
если взял успешно, убрать айтем из локации;
	`
	if player.hasBag() {
		player.collectItem(t.item)
		// location.removeItem(t.item)
		player.commandReaction(fmt.Sprintf("предмет добавлен в инвентарь: %s", t.item))
	} else {
		player.commandReaction("некуда класть")
	}

	return nil
}

var _ ICommand = &ApplyCmd{} // type check

//   - применить $item $object => что было сделано
//     items: ключи, конспекты, телефон
//     objects: дверь, шкаф
//     examples:
//     ключи дверь
//     телефон шкаф
//     ключи шкаф
type ApplyCmd struct {
	applyItem string
	toObject  string
}

// execute implements ICommand.
func (a *ApplyCmd) execute(game IGame, player IPlayer) error {
	// {9, "применить ключи дверь", "дверь открыта"}
	var _ = `
у плеера в инвентори (мешке) есть айтем;
в локации (где плеер находится) есть обьект;
обьект может быть "активирован" айтемом;
реакция: состояние обьекта после "активации"

	`
	if player.hasItem(a.applyItem) {
		locName := player.getLocation()
		loc := game.getLocation(locName)
		obj, err := loc.getObject(a.toObject)
		if err == nil {
			if obj.isCompatibleWith(a.applyItem) {
				obj.activateWith(a.applyItem)
				player.commandReaction(obj.getState())
			} else {
				return fmt.Errorf("item '%s' can't be applied to object '%s'", a.applyItem, a.toObject)
			}
		} else {
			return fmt.Errorf("no such object (%s) in location: '%s'; %w", a.toObject, loc.name, err)
		}
	} else {
		player.commandReaction(fmt.Sprintf("нет предмета в инвентаре - %s", a.applyItem))
	}

	return nil
}

//   - сказать $phrase => широковещание всем игрокам в локации, цитирование сказанного
//     examples: Пора топать в универ
type ShoutCmd struct{}

// execute implements ICommand.
func (s *ShoutCmd) execute(game IGame, player IPlayer) error {
	panic("unimplemented")
}

var _ ICommand = &ShoutCmd{} // type check

//   - сказать_игроку $player $phrase_option => адресное обращение к игроку в локации, цитирование
//     examples:
//     Izolda Может ещё по чаю лучше?
//     Tristan
type SayToPlayerCmd struct{}

// execute implements ICommand.
func (s *SayToPlayerCmd) execute(game IGame, player IPlayer) error {
	panic("unimplemented")
}

var _ ICommand = &SayToPlayerCmd{} // type check
