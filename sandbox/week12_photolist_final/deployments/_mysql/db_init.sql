-- Adminer 4.7.0 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `photos`;
CREATE TABLE `photos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `path` varchar(255) NOT NULL,
  `rating` bigint(20) NOT NULL DEFAULT '0',
  `comment` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `photos` (`id`, `user_id`, `path`, `rating`, `comment`) VALUES
(6,	5,	'5b280f8f6c549272ee8626f26bfdf8f0',	0,	'asd'),
(7,	5,	'0606d2dce38b4993ec7ab4991a88648d',	0,	'rarsfd'),
(8,	6,	'bb73d02d20c139e4adc8dc817397cb80',	0,	'view 1'),
(9,	6,	'6128956443b9c8044a9f73c325e95b2d',	0,	'view 2'),
(10,	6,	'6f89427583acff4ecb380450ed00b24b',	0,	'view 3'),
(11,	6,	'9fc84cce2b7847a1669b5a3a049fecea',	1,	'view 4'),
(12,	6,	'9740f4023540baaa685e42277677ec67',	3,	'view 5'),
(14,	5,	'c805172e4c3e4e1b35d73077ff47e3dc',	0,	'studio 2'),
(15,	7,	'93aaabaf6c9afc54965d721f108474df',	0,	'building 1'),
(16,	7,	'93e4f2d335f35df202323284642480c5',	0,	'building 2'),
(17,	7,	'df311d3c12b2542f36e0f2ec5c1735fd',	0,	'building 3'),
(18,	7,	'433b9ab5f0649e6a17da024e4efaadef',	0,	'building 4'),
(19,	7,	'aa767771763cb2b7c0f2f92909fe0fe3',	0,	'building 5'),
(20,	7,	'21370b81864074aedf3263861f7dd519',	0,	'building 6'),
(22,	7,	'ee303ac2f78c12b4705ec61b15cb815f',	0,	'building 7'),
(23,	7,	'9bf7a47d61f5e67c4c323881eb1f86fa',	0,	'building 8'),
(25,	7,	'08143736bbbd748dde28104c1c4e02f0',	0,	'building 10'),
(26,	7,	'70ba66d831ce7c42891799d246829d0d',	1,	'building 11');

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
  `id` varchar(32) NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  UNIQUE KEY `id` (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `user_follows`;
CREATE TABLE `user_follows` (
  `user_id` int(11) NOT NULL,
  `follow_id` int(11) NOT NULL,
  KEY `follow_id` (`follow_id`),
  KEY `user_id_follow_id` (`user_id`,`follow_id`),
  CONSTRAINT `user_follows_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `user_follows` (`user_id`, `follow_id`) VALUES
(5,	6);

DROP TABLE IF EXISTS `user_photos_likes`;
CREATE TABLE `user_photos_likes` (
  `photo_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  KEY `photo_id` (`photo_id`),
  KEY `user_id_photo_id` (`user_id`,`photo_id`),
  CONSTRAINT `user_photos_likes_ibfk_1` FOREIGN KEY (`photo_id`) REFERENCES `photos` (`id`),
  CONSTRAINT `user_photos_likes_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `user_photos_likes` (`photo_id`, `user_id`) VALUES
(11,	5),
(12,	5),
(26,	5);

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `login` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varbinary(100) NOT NULL,
  `ver` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `login` (`login`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `users` (`id`, `login`, `password`) VALUES
(1,	'golangcourse',	UNHEX('7362415A62716D7072AC09EAE839A4A1C95E73EC5FA3FC6EACE4D2C78BF4BF1C6906789B557F8C55'));

INSERT INTO `users` (`id`, `login`, `email`, `password`, `ver`) VALUES
(5,	'rvasily.msk',	'romanov.vasily@gmail.com',	UNHEX('6F7A6B5341575172BF22E32CBCE77A1942B342B996B0EF172673FD214512D2675C7843C6A0FA9597'),	0),
(6,	'views',	'views@example.com',	UNHEX('48444259794A5452704CB8CE7C435EED4700742D64390D4DC6EB0E68A0AEB412BB8BF4A0D5D9D7FA'),	0),
(7,	'buildings',	'buildings@example.com',	UNHEX('6645586B774251785DC303327BE1427F17F608441672EADF322F44695D0180E79880B27F0274BC66'),	0);

-- 2019-08-25 12:49:46