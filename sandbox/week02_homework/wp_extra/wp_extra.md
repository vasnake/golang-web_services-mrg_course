Это экстра задание, его можно делать после того как вы сделали signer

Напишите воркер-пул в котором можно динамически ( во время работы программы ) добавлять и удалять воркеров

Масштабирование воркеров в зависимости от "нагрузки" - те от какого-то сигнала из-вне

Допустим у нас есть какой-то "датчик" ( нагрузка на цпу, потребление памяти, диск, количество обработаннызх задач - что угодно ) и в зависимости от этого датчика мы можем или сказать пулу  чтобы он завершил работу 1 воркера или добавить в пул еще 1 воркер

Только без космолетов 🙂

Там совсем небольшое решение выходит, 150-200 строчек.

Задание на каналы, select и горутины