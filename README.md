# Разработка веб-сервисов на Go

[Go course, MRG, Романов Василий](https://github.com/vasnake/golang-web_services-mrg_course/blob/main/README.md)

- [part1](part1.md)
- [part2](part2.md)
- [part3](part3.md)
- part4 не существует?

## homework

> Numeric constants represent exact values of arbitrary precision and do not overflow.
Consequently, there are no constants denoting the IEEE-754 negative zero, infinity, and not-a-number values.

Напиши библиотечные функции:
- Сравнения двух чисел float, с учетом дельты, +/- zero.
- Предикат `IsInvalidFloat`, используемый в фильтрации коллекции чисел. Инвалид: inf, nan.

## Info, links

- `Communicating Sequential Processes` by Tony Hoare
- https://cs.stanford.edu/people/eroberts/courses/soco/projects/2008-09/tony-hoare/csp.html
- http://www.usingcsp.com/cspbook.pdf
- CSP vs Actor model (channels vs actor/mailbox)
- [If a map isn’t a reference variable, what is it?](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it)
- [Implementing a bignum calculator with Rob Pike](https://youtu.be/PXoG0WX0r_E)
- [Lexical Scanning in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE)
- [Advanced Topics in Programming Languages: Concurrency/message passing Newsqueak](https://youtu.be/hB05UFqOtFA), prime sieve on channels
