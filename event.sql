CREATE DATABASE IF NOT EXISTS event;
use event;
CREATE TABLE if not exists `event.events` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE if not exists `event.pro_event` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `rid` varchar(256) NOT NULL DEFAULT '',
  `event_id` varchar(256) NOT NULL DEFAULT '',
  `referer` varchar(256) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE if not exists `event.stg_event` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `rid` varchar(256) NOT NULL DEFAULT '',
  `event_id` varchar(256) NOT NULL DEFAULT '',
  `referer` varchar(256) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
