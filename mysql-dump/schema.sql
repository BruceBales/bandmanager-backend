CREATE SCHEMA IF NOT EXISTS prim;
USE prim;
CREATE TABLE IF NOT EXISTS `users` (  `id` int(11) NOT NULL AUTO_INCREMENT,
`name` varchar(45) DEFAULT NULL,
`email` varchar(45) DEFAULT NULL,
`password` varchar(2000) DEFAULT NULL,
PRIMARY KEY (`id`),
UNIQUE KEY `id_UNIQUE` (`id`)) ENGINE=InnoDB AUTO_INCREMENT=77 DEFAULT CHARSET=latin1