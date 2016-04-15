CREATE DATABASE IF NOT EXISTS event;
use event;
CREATE TABLE if not exists `events` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE if not exists `pro_event` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `rid` varchar(256) NOT NULL DEFAULT '',
  `event_id` int(11) NOT NULL DEFAULT '',
  `referer` varchar(256) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `pro_event` MODIFY COLUMN `event_id` int(11);
ALTER TABLE `pro_event` MODIFY COLUMN `referer` varchar(1024);
ALTER TABLE `pro_event` ADD `client_id` int(11);
ALTER TABLE `pro_event` ADD `genre_id` int(11);
ALTER TABLE `pro_event` ADD `other` text;
CREATE TABLE if not exists `stg_event` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `rid` varchar(256) NOT NULL DEFAULT '',
  `event_id` varchar(256) NOT NULL DEFAULT '',
  `referer` varchar(256) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
ALTER TABLE `stg_event` MODIFY COLUMN `event_id` int(11);
ALTER TABLE `stg_event` MODIFY COLUMN `referer` varchar(1024);
ALTER TABLE `stg_event` ADD `client_id` int(11);
ALTER TABLE `stg_event` ADD `genre_id` int(11);
ALTER TABLE `stg_event` ADD `other` text;
CREATE TABLE IF NOT EXISTS `client_master` (
  `id` int(11) unsigned not null AUTO_INCREMENT,
  `name` varchar(256) not null DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE IF NOT EXISTS `genre_master` (
  `id` int(11) unsigned not null AUTO_INCREMENT,
  `name` varchar(1024) not null DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
