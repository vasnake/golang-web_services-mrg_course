package main

import (
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

// --------------------------------------------------------------
// предыдущая версия игры с изменениями для нескольких игроков, тесты для 1 игрока

type game0Case struct {
	step    int
	command string
	answer  string
}

var game0cases = [][]game0Case{

	[]game0Case{
		{1, "осмотреться", "ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"}, // действие осмотреться
		{2, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},                                         // действие идти
		{3, "идти комната", "ты в своей комнате. можно пройти - коридор"},
		{4, "осмотреться", "на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор"},
		{5, "одеть рюкзак", "вы одели: рюкзак"},                   // действие одеть
		{6, "взять ключи", "предмет добавлен в инвентарь: ключи"}, // действие взять
		{7, "взять конспекты", "предмет добавлен в инвентарь: конспекты"},
		{8, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{9, "применить ключи дверь", "дверь открыта"}, // действие применить
		{11, "идти улица", "на улице весна. можно пройти - домой"},
	}, // 0
	[]game0Case{
		{1, "осмотреться", "ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор"},
		{2, "завтракать", "неизвестная команда"},  // придёт топать в универ голодным :(
		{3, "идти комната", "нет пути в комната"}, // через стены ходить нельзя
		{4, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{5, "применить ключи дверь", "нет предмета в инвентаре - ключи"},
		{6, "идти комната", "ты в своей комнате. можно пройти - коридор"},
		{7, "осмотреться", "на столе: ключи, конспекты, на стуле - рюкзак. можно пройти - коридор"},
		{8, "взять ключи", "некуда класть"}, // надо взять рюкзак сначала
		{9, "одеть рюкзак", "вы одели: рюкзак"},
		{10, "осмотреться", "на столе: ключи, конспекты. можно пройти - коридор"}, // состояние изменилось
		{11, "взять ключи", "предмет добавлен в инвентарь: ключи"},
		{12, "взять телефон", "нет такого"},                                // неизвестный предмет
		{13, "взять ключи", "нет такого"},                                  // предмента уже нет в комнатеы - мы его взяли
		{14, "осмотреться", "на столе: конспекты. можно пройти - коридор"}, // состояние изменилось
		{15, "взять конспекты", "предмет добавлен в инвентарь: конспекты"},
		{16, "осмотреться", "пустая комната. можно пройти - коридор"}, // состояние изменилось
		{17, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{18, "идти кухня", "кухня, ничего интересного. можно пройти - коридор"},
		{19, "осмотреться", "ты находишься на кухне, на столе чай, надо идти в универ. можно пройти - коридор"}, // состояние изменилось
		{20, "идти коридор", "ничего интересного. можно пройти - кухня, комната, улица"},
		{21, "идти улица", "дверь закрыта"},                                  //условие не удовлетворено
		{22, "применить ключи дверь", "дверь открыта"},                       //состояние изменилось
		{23, "применить телефон шкаф", "нет предмета в инвентаре - телефон"}, // нет предмета
		{24, "применить ключи шкаф", "не к чему применить"},                  // предмет есть, но применить его к этому нельзя
		{25, "идти улица", "на улице весна. можно пройти - домой"},
	}, // 1
}

func TestGameSingleplayer(t *testing.T) {
	for caseNum, commands := range game0cases {

		players := map[string]*Player{
			"Tristan": NewPlayer("Tristan"),
		}

		playersOutput := map[string]string{
			"Tristan": "",
		}

		// async, read player chan, write to buf
		mu := &sync.Mutex{}
		go func() {
			output := players["Tristan"].GetOutput()
			for msg := range output {
				mu.Lock()
				playersOutput["Tristan"] = msg
				mu.Unlock()
			}
		}()

		initGame()
		addPlayer(players["Tristan"])

		for _, item := range commands {
			// send inp message
			players["Tristan"].HandleInput(item.command)

			// collect out message
			time.Sleep(time.Millisecond)
			runtime.Gosched() // дадим считать ответ
			mu.Lock()
			answer := playersOutput["Tristan"]
			mu.Unlock()

			// check message
			if answer != item.answer {
				// t.Error vs t.Fatal
				t.Fatal("case:", caseNum, item.step,
					"\n\tcmd:", item.command,
					"\n\texpected:", item.answer,
					"\n\tactual  :  ", answer,
				)
			}
		}
	}

}

// --------------------------------------------------------------
// новая игра - взаимодействие двух игроков
// поскольку предполагается, что комнады отправляются асинхронно, то в ответах тстовых кейсов - последний ответ, полученный данными гроком

type game1Case struct {
	step    int
	player  string
	command string
	answers map[string]string
}

var game1Cases = [][]game1Case{

	{
		{
			1,
			"Tristan",
			"осмотреться",
			map[string]string{
				"Tristan": "ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор. Кроме вас тут ещё Izolda",
			},
		}, // действие осмотреться
		{
			2,
			"Izolda",
			"осмотреться",
			map[string]string{
				"Izolda": "ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ. можно пройти - коридор. Кроме вас тут ещё Tristan",
			},
		}, // действие осмотреться
		{
			3,
			"Izolda",
			"сказать Пора топать в универ",
			map[string]string{
				"Tristan": "Izolda говорит: Пора топать в универ",
				"Izolda":  "Izolda говорит: Пора топать в универ",
			},
		}, // действие сказать
		{
			4,
			"Tristan",
			"сказать_игроку Izolda Может ещё по чаю лучше?",
			map[string]string{
				"Izolda": "Tristan говорит вам: Может ещё по чаю лучше?",
			},
		}, // действие сказать_игроку
		{
			5,
			"Izolda",
			"сказать_игроку Tristan",
			map[string]string{
				"Tristan": "Izolda выразительно молчит, смотря на вас",
			},
		}, // действие сказать_игроку
		{
			6,
			"Tristan",
			"идти коридор",
			map[string]string{
				"Tristan": "ничего интересного. можно пройти - кухня, комната, улица",
			},
		}, // действие идти
		{
			7,
			"Izolda",
			"сказать_игроку Tristan",
			map[string]string{
				"Izolda": "тут нет такого игрока",
			},
		}, // действие сказать_игроку
	},
}

func TestGameMiltiplayer(t *testing.T) {
	for caseNum, commands := range game1Cases {

		var lastOutput = map[string]string{}

		players := map[string]*Player{
			"Tristan": NewPlayer("Tristan"),
			"Izolda":  NewPlayer("Izolda"),
		}

		mu := &sync.Mutex{}

		go func() {
			output := players["Tristan"].GetOutput()
			for msg := range output {
				mu.Lock()
				lastOutput["Tristan"] = msg
				mu.Unlock()
			}
		}()

		go func() {
			output := players["Izolda"].GetOutput()
			for msg := range output {
				mu.Lock()
				lastOutput["Izolda"] = msg
				mu.Unlock()
			}
		}()

		initGame()
		addPlayer(players["Tristan"])
		addPlayer(players["Izolda"])

		for _, item := range commands {
			lastOutput = map[string]string{}               // clear
			players[item.player].HandleInput(item.command) // player write to chan, goroutine pipe from chan to lastOutput
			time.Sleep(time.Millisecond)
			runtime.Gosched() // дадим считать ответ
			if !reflect.DeepEqual(lastOutput, item.answers) {
				t.Fatal("case:", caseNum, item.step,
					"\n\tcmd:", item.command,
					"\n\texpected:", item.answers,
					"\n\tactual  :  ", lastOutput,
				)
			}
		}
	}

}
