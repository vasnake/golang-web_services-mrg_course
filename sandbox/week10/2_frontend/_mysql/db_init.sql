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
  `comment` varchar(512) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;


SET NAMES utf8mb4;

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
  `id` varchar(32) NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  UNIQUE KEY `id` (`id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `login` varchar(100) NOT NULL,
  `password` varbinary(100) NOT NULL,
  `ver` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `login` (`login`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

INSERT INTO `users` (`id`, `login`, `password`) VALUES
(1,	'golangcourse',	UNHEX('7362415A62716D7072AC09EAE839A4A1C95E73EC5FA3FC6EACE4D2C78BF4BF1C6906789B557F8C55'));

DROP TABLE IF EXISTS `user_photos_likes`;
CREATE TABLE `user_photos_likes` (
  `photo_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  KEY `photo_id` (`photo_id`),
  KEY `user_id_photo_id` (`user_id`,`photo_id`),
  CONSTRAINT `user_photos_likes_ibfk_1` FOREIGN KEY (`photo_id`) REFERENCES `photos` (`id`),
  CONSTRAINT `user_photos_likes_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 2019-07-14 15:02:20