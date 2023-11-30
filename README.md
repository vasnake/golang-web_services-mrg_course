# Разработка веб-сервисов на Go

[Go course, MRG, Романов Василий](https://github.com/vasnake/golang-web_services-mrg_course/blob/main/README.md)

- [part1](part1.md)
- [part2](part2.md)
- [part3](part3.md)
- part4 не существует?

## homework

### float numbers

> Numeric constants represent exact values of arbitrary precision and do not overflow.
Consequently, there are no constants denoting the IEEE-754 negative zero, infinity, and not-a-number values.

Напиши библиотечные функции:
- Предикат `IsInvalidFloat`, используемый в фильтрации коллекции чисел. Инвалид: inf, nan
- Сравнения двух чисел float, с учетом дельты, +/- zero, nan, +/- inf

### int numbers bits

> Two's complement is achieved by:
- Step 1: starting with the equivalent positive number.
- Step 2: inverting (or flipping) all bits – changing every 0 to 1, and every 1 to 0;
- Step 3: adding 1 to the entire inverted number, ignoring any overflow.

Напиши отображение 32-бит числа в int32, uint32. Битовое представление не меняется, меняется интерпретация набора бит.
Рассмотри два случая: int32 получен как результат хеширования
и 1) надо отобразить диапазон значений на uint32. 2) надо сохранить битовое представление.

Напиши конвертеры: (u)int => bitstring, bitstring => u(int).

### concurrency patterns

Напиши параллельную обработку по трем паттернам https://habr.com/ru/companies/timeweb/articles/770912/

### week 2, dynamic workers pool

[week 2, extra homework](./sandbox/week02_homework/wp_extra/wp_extra.md)

## Info, links

- `Communicating Sequential Processes` by Tony Hoare
- https://cs.stanford.edu/people/eroberts/courses/soco/projects/2008-09/tony-hoare/csp.html
- http://www.usingcsp.com/cspbook.pdf
- CSP vs Actor model (channels vs actor/mailbox)
- [If a map isn’t a reference variable, what is it?](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it)
- [Implementing a bignum calculator with Rob Pike](https://youtu.be/PXoG0WX0r_E)
- [Lexical Scanning in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE)
- [Advanced Topics in Programming Languages: Concurrency/message passing Newsqueak](https://youtu.be/hB05UFqOtFA), prime sieve on channels
- [3.5 Years, 500k Lines of Go (Part 1)][https://npf.io/2017/03/3.5yrs-500k-lines-of-go/] `return fmt.Errorf("While doing foo: %v", err.Error())`
