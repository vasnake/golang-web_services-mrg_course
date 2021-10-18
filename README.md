# go-web_services-mrg_course
Go language course from MRG: Разработка веб-сервисов на Go

## part 1, week 1
Введение в Go.
[Код, домашки, литература](./week_01/materials.zip)
https://cloud.mail.ru/public/pc7H/6Vx4txWWr

Зачем нужен еще один языка программирования?
Эффективность: компиляции, выполнения, разработки.
Многопроцессорные системы (легкие потоки, асинхронность);
простой и понятный язык, читабельный, простая и быстрая сборка, с четким стилем.

Для realtime не подходит (сборка мусора), но всё остальное OK.

- Первая программа: play.golang.com, hello_world
- Переменные, базовые типы данных: vars_1, vars_2, strings, const, types, pointers
- Переменные, составные типы данных: array, slice_1, slice_2, map
- Управляющие конструкции: control, loop
- Основы функций: functions
- Функция как объект первого класса, анонимные функции: firstclass
- Отложенное выполнение и обработка паники: defer, recover
- Основы работы со структурами: structs
- Методы структур: methods
- Пакеты и область видимости: dir.txt, visibility/
- Основы работы с интерфейсами: basic, many, cast
- Пустой интерфейс: empty_1, empty_2
- Композиция интерфейсов: embed
- Написание программы уникализации (ПУ): uniq, data_map.txt
- Написание тестов для ПУ: unique/unique, unique/unique_test


## part 1, week 2
Асинхронная работа.
[Код, домашки, литература](./week_02/w2_materials.zip)
https://cloud.mail.ru/public/YDEX/Dau2wVWuw/

- Методы обработки запросов и плюсы неблокирующего подхода:
  асинхронное выполнение, скорость процессор-кеш-память, время на переключение контекста, современные тенденции на многоядерность и параллельность,
  тяжелые процессы, потоки легче, асинхронные сопрограммы (green threads) еще легче. Невытесняющая многозадачность (eventloop, Windows 3.0) vs preemptive.
  Ввод-вывод и ожидание возврата из syscall. Время ожидания можно потратить на другие задачи, non-blocking IO. IO-bound vs CPU-bound.
  `Communicating Sequential Processes` by Tony Hoare. Горутины перемещаются между системными потоками.
- Горутины -- легковесные процессы: goroutines
### TODO
- Каналы -- передаём данные между горутинами: chan_1, chan_2
- Мультиплексирование каналов через оператор `select`
- Таймеры и таймауты
- Пакет `context` и отмена выполнения
- Асинхронное получение данных
- Пул воркеров
- `sync.Waitgroup` -- ожидание завершения работы
- Ограничение по ресурсам
- Ситуация гонки на примере конкурентной записи в map
- `sync.Mutex` для синхронизации данных
- `sync.Atomic`

## week 3
Работа с динамическими данными и производительность

## week 4
Основы HTTP

## Info
- [If a map isn’t a reference variable, what is it?](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it)

---

Developed using JetBrains GoLand IDE
[![JetBrains GoLand](./icon-goland.svg)](https://jb.gg/OpenSource)

With thanks to JetBrains and their support for open source communities
[![JetBrains Open Source Support](./jetbrains-variant-3.svg)](https://jb.gg/OpenSource)
