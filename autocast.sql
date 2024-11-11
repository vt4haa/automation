-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Хост: 127.0.0.1:3306
-- Время создания: Ноя 11 2024 г., 11:29
-- Версия сервера: 5.7.36
-- Версия PHP: 7.4.26

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `autocast`
--

-- --------------------------------------------------------

--
-- Структура таблицы `brands`
--

DROP TABLE IF EXISTS `brands`;
CREATE TABLE IF NOT EXISTS `brands` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `brands`
--

INSERT INTO `brands` (`id`, `name`) VALUES
(1, 'Samsung'),
(2, 'Apple'),
(3, 'Sony'),
(4, 'LG');

-- --------------------------------------------------------

--
-- Структура таблицы `categories`
--

DROP TABLE IF EXISTS `categories`;
CREATE TABLE IF NOT EXISTS `categories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `categories`
--

INSERT INTO `categories` (`id`, `name`) VALUES
(1, 'Smartphone'),
(2, 'Laptop'),
(3, 'Television'),
(4, 'Home Appliance');

-- --------------------------------------------------------

--
-- Структура таблицы `clients`
--

DROP TABLE IF EXISTS `clients`;
CREATE TABLE IF NOT EXISTS `clients` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `contact` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `clients`
--

INSERT INTO `clients` (`id`, `name`, `contact`) VALUES
(1, 'John Doe', 'john@example.com'),
(2, 'Alice Brown', 'alice@example.com'),
(3, 'Bob Smith', 'bob@example.com');

-- --------------------------------------------------------

--
-- Структура таблицы `items`
--

DROP TABLE IF EXISTS `items`;
CREATE TABLE IF NOT EXISTS `items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `idItems` int(11) DEFAULT NULL,
  `quality` int(11) DEFAULT NULL,
  `idOrders` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idOrders` (`idOrders`),
  KEY `idItems` (`idItems`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `items`
--

INSERT INTO `items` (`id`, `idItems`, `quality`, `idOrders`) VALUES
(1, 1, 2, 1),
(2, 2, 1, 2),
(3, 3, 1, 3);

-- --------------------------------------------------------

--
-- Структура таблицы `orders`
--

DROP TABLE IF EXISTS `orders`;
CREATE TABLE IF NOT EXISTS `orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `idClients` int(11) DEFAULT NULL,
  `comment` varchar(255) DEFAULT NULL,
  `datatime` datetime DEFAULT NULL,
  `idManager` int(11) DEFAULT NULL,
  `idCollector` int(11) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idClients` (`idClients`),
  KEY `idManager` (`idManager`),
  KEY `idCollector` (`idCollector`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `orders`
--

INSERT INTO `orders` (`id`, `idClients`, `comment`, `datatime`, `idManager`, `idCollector`, `status`) VALUES
(1, 1, 'Order for iPhone 14', '2024-11-05 10:30:00', 1, 2, 'Pending'),
(2, 2, 'Order for Samsung Galaxy S21', '2024-11-05 11:00:00', 1, 2, 'Completed'),
(3, 3, 'Order for Sony Bravia 55\"', '2024-11-05 11:30:00', 1, 2, 'Shipped');

-- --------------------------------------------------------

--
-- Структура таблицы `positions`
--

DROP TABLE IF EXISTS `positions`;
CREATE TABLE IF NOT EXISTS `positions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `positions`
--

INSERT INTO `positions` (`id`, `name`) VALUES
(1, 'Admin'),
(2, 'Manager'),
(3, 'User');

-- --------------------------------------------------------

--
-- Структура таблицы `product`
--

DROP TABLE IF EXISTS `product`;
CREATE TABLE IF NOT EXISTS `product` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `photo` varchar(255) DEFAULT NULL,
  `idCategories` int(11) DEFAULT NULL,
  `idBrands` int(11) DEFAULT NULL,
  `quality` int(11) DEFAULT NULL,
  `price` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idCategories` (`idCategories`),
  KEY `idBrands` (`idBrands`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `product`
--

INSERT INTO `product` (`id`, `name`, `photo`, `idCategories`, `idBrands`, `quality`, `price`) VALUES
(1, 'iPhone 14', 'iphone14.jpg', 1, 2, 100, 999),
(2, 'Samsung Galaxy S21', 'galaxy_s21.jpg', 1, 1, 50, 799),
(3, 'Sony Bravia 55\"', 'bravia55.jpg', 3, 3, 20, 1200),
(4, 'LG Refrigerator', 'lg_fridge.jpg', 4, 4, 10, 1500);

-- --------------------------------------------------------

--
-- Структура таблицы `workers`
--

DROP TABLE IF EXISTS `workers`;
CREATE TABLE IF NOT EXISTS `workers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fio` varchar(255) DEFAULT NULL,
  `post` int(11) DEFAULT NULL,
  `login` varchar(255) NOT NULL,
  `pass` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `post` (`post`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

--
-- Дамп данных таблицы `workers`
--

INSERT INTO `workers` (`id`, `fio`, `post`, `login`, `pass`) VALUES
(1, 'Michael Johnson', 1, 'michael_j', '$2a$10$JDTdQwti0QhmtSLqDiNYHui/.8GDBaJ7hWFoY2ORNJXdj2GSca5Bm'),
(2, 'Susan Green', 2, 'susan_g', '$2a$10$xCg2I8igrx7vzD21zonjMeYx/y1a9BZ9KNgYqhuYJjFbb05LQ3T5K');

--
-- Ограничения внешнего ключа сохраненных таблиц
--

--
-- Ограничения внешнего ключа таблицы `items`
--
ALTER TABLE `items`
  ADD CONSTRAINT `items_ibfk_1` FOREIGN KEY (`idOrders`) REFERENCES `orders` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `items_ibfk_2` FOREIGN KEY (`idItems`) REFERENCES `product` (`id`) ON DELETE SET NULL;

--
-- Ограничения внешнего ключа таблицы `orders`
--
ALTER TABLE `orders`
  ADD CONSTRAINT `orders_ibfk_1` FOREIGN KEY (`idClients`) REFERENCES `clients` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `orders_ibfk_2` FOREIGN KEY (`idManager`) REFERENCES `workers` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `orders_ibfk_3` FOREIGN KEY (`idCollector`) REFERENCES `workers` (`id`) ON DELETE SET NULL;

--
-- Ограничения внешнего ключа таблицы `product`
--
ALTER TABLE `product`
  ADD CONSTRAINT `product_ibfk_1` FOREIGN KEY (`idCategories`) REFERENCES `categories` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `product_ibfk_2` FOREIGN KEY (`idBrands`) REFERENCES `brands` (`id`) ON DELETE SET NULL;

--
-- Ограничения внешнего ключа таблицы `workers`
--
ALTER TABLE `workers`
  ADD CONSTRAINT `workers_ibfk_1` FOREIGN KEY (`post`) REFERENCES `positions` (`id`) ON DELETE SET NULL;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
