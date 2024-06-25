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

	// broadcast: сказать Пора топать в универ
	case strings.HasPrefix(cmd, "сказать "):
		phrase := strings.TrimPrefix(cmd, "сказать ")
		return &ShoutCmd{phrase: phrase}

	// personal: сказать_игроку Izolda Может ещё по чаю лучше?
	case strings.HasPrefix(cmd, "сказать_игроку "):
		name_phrase := strings.TrimPrefix(cmd, "сказать_игроку ")
		name, phrase, _ := strings.Cut(name_phrase, " ")
		return &SayToPlayerCmd{to: name, body: phrase}
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
	descr := game.getLookupLocationDescription(currLocation, player)

	// I'm alone?
	players := game.getPlayersInLocation(currLocation)
	if len(players) > 1 {
		// other players
		for _, p := range players {
			if p.getName() != player.getName() {
				descr = descr + fmt.Sprintf(". Кроме вас тут ещё %s", p.getName())
			} else {
				// it's me
			}
		}
	} else {
		// only one
	}

	player.commandReaction(descr)

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

переход на новую локацию может быть под условием.
что порождает два варианта: условие выполнено, не выполнено.

пример, условие:
{22, "применить ключи дверь", "дверь открыта"}

условие проверяет стейт игрока-в-локации,
пример: ранее игрок открыл дверь (поменял состояние объекта дверь).

игрок может хранить под ключом (локация, объект) состояние обьекта в этой локации,
после "применить ключи дверь" дверь получает стейт "дверь открыта".
команда "идти улица" проверяет стейт двери (в данных игрока для данной локации)

идти улица => дверь закрыта
применить ключи дверь => дверь открыта
идти улица => на улице весна. можно пройти - домой

	`

	var targetLocation string = g.locationName
	var currLocation string = player.getLocation()
	var currLocObj = game.getLocation(currLocation)

	if currLocation == targetLocation {
		return fmt.Errorf("GotoCmd failed, target == current: %s", currLocation)
	}

	if game.isLocationsConnected(currLocation, targetLocation) {
		canGo := true
		reaction := game.getGoToLocationDescription(targetLocation)

		// if location have condition and it is not met
		if currLocObj.isGotoConditionExists(targetLocation) {
			cond, negativeReaction := currLocObj.getGotoCondition(targetLocation) // дверь открыта, дверь закрыта
			if player.isStateInLocationState(currLocation, cond) {
				// дверь открыта
				canGo = true
			} else {
				// дверь закрыта
				canGo = false
				reaction = negativeReaction
			}
		} else {
			// no condition
			canGo = true
		}

		// if no condition or condition met
		if canGo {
			player.setLocation(targetLocation)
		}
		player.commandReaction(reaction)

	} else {
		// not connected
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

	locName := player.getLocation()
	location := game.getLocation(locName)
	location.removeItem(p.item)

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

	location := game.getLocation(player.getLocation())
	if location.isItemExists(t.item) {
		if player.hasBag() {
			player.collectItem(t.item)
			location.removeItem(t.item)
			player.commandReaction(fmt.Sprintf("предмет добавлен в инвентарь: %s", t.item))
		} else {
			// no bag
			player.commandReaction("некуда класть")
		}
	} else {
		// no item in location
		player.commandReaction("нет такого")
	}

	return nil
}

var _ ICommand = &ApplyCmd{} // type check

//   - применить $item $object => что было сделано
//     items: ключи, конспекты, телефон
//     objects: дверь, шкаф
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
		locObj := game.getLocation(locName)
		obj, err := locObj.getObject(a.toObject) // not a good idea if game state is independent for each player
		if err == nil {
			if obj.isCompatibleWith(a.applyItem) {
				obj.activateWith(a.applyItem)
				player.commandReaction(obj.getState())
				player.setStateInLocationState(locName, obj.getState())
			} else {
				return fmt.Errorf("item '%s' can't be applied to object '%s'", a.applyItem, a.toObject)
			}
		} else {
			// show("ApplyCmd.execute: ", fmt.Errorf("no such object (%s) in location: '%s'; %w", a.toObject, locObj.name, err))
			player.commandReaction("не к чему применить")
		}
	} else {
		player.commandReaction(fmt.Sprintf("нет предмета в инвентаре - %s", a.applyItem))
	}

	return nil
}

// - сказать $phrase => широковещание всем игрокам в локации, цитирование сказанного
type ShoutCmd struct {
	phrase string
}

// execute implements ICommand.
func (s *ShoutCmd) execute(game IGame, player IPlayer) error {
	// "Izolda говорит: Пора топать в универ",
	locName := player.getLocation()
	players := game.getPlayersInLocation(locName)
	for _, p := range players {
		p.commandReaction(fmt.Sprintf("%s говорит: %s", player.getName(), s.phrase))
	}
	return nil
}

var _ ICommand = &ShoutCmd{} // type check

// - сказать_игроку $player $phrase_option => адресное обращение к игроку в локации, цитирование
type SayToPlayerCmd struct {
	to   string
	body string
}

// execute implements ICommand.
func (s *SayToPlayerCmd) execute(game IGame, player IPlayer) error {
	locName := player.getLocation()
	players := game.getPlayersInLocation(locName)
	for _, p := range players {
		if p.getName() == s.to {

			if s.body == "" {
				// nothing to say
				p.commandReaction(fmt.Sprintf("%s выразительно молчит, смотря на вас", player.getName()))
			} else {
				// have phrase
				p.commandReaction(fmt.Sprintf("%s говорит вам: %s", player.getName(), s.body))
			}

			return nil
		}
	}

	player.commandReaction(fmt.Sprintf("тут нет такого игрока"))
	return nil
}

var _ ICommand = &SayToPlayerCmd{} // type check
