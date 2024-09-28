-- MySQL dump 10.13  Distrib 8.0.39, for Win64 (x86_64)
--
-- Host: localhost    Database: go_final_project
-- ------------------------------------------------------
-- Server version	8.0.39

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `admins`
--

DROP TABLE IF EXISTS `admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admins` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `name` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admins`
--

LOCK TABLES `admins` WRITE;
/*!40000 ALTER TABLE `admins` DISABLE KEYS */;
INSERT INTO `admins` VALUES (1,'aebdfc46-cc91-42f3-988c-7647d2380572','Alif Maulidanar','alif@example.com','$2a$14$c9icsBB8orG.tuPg6cKQ1Oj1nbwN942fnYokZek7wJP06pGempuqC','2024-09-28 12:21:44','2024-09-28 12:21:44'),(2,'7f4cd998-ed3d-421d-9da3-cf54b329ca80','John Doe','johndoe@example.com','$2a$14$GW8TEKK0bR6kwuI9/Wlh1O.h70EULbukDUAX0aFbVLPZiRUtcaxgy','2024-09-28 12:23:51','2024-09-28 12:23:51');
/*!40000 ALTER TABLE `admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `name` varchar(50) NOT NULL,
  `image_url` varchar(255) DEFAULT NULL,
  `admin_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `admin_id` (`admin_id`),
  CONSTRAINT `products_ibfk_1` FOREIGN KEY (`admin_id`) REFERENCES `admins` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
INSERT INTO `products` VALUES (1,'b7a7bbdd-7196-4a0f-891a-c04b6606a63a','Toyota Raize 2024','https://res.cloudinary.com/dnu1vesdo/image/upload/v1727531528/products/faupq73jw311afjtm4gi.png',1,'2024-09-28 13:50:52','2024-09-28 13:58:12'),(2,'c16c72f8-8366-4ab3-9f19-aa806ba501c6','Mitsubishi Xpander 2024','https://res.cloudinary.com/dnu1vesdo/image/upload/v1727531712/products/vp7qbwiz4cg2ct1jphep.jpg',1,'2024-09-28 13:53:55','2024-09-28 13:58:48'),(3,'76b626fa-5a2c-44e3-9295-4507c2aabddd','Hyundai Palisade 2024','https://res.cloudinary.com/dnu1vesdo/image/upload/v1727531809/products/dgu2kyyqdb9nkhtmegpa.png',2,'2024-09-28 13:55:32','2024-09-28 13:59:23'),(4,'eb554526-2771-4217-bca0-b8208c2375a7','Nissan Serena 2024','https://res.cloudinary.com/dnu1vesdo/image/upload/v1727531879/products/bhu0prlyzfqf30zu2idu.jpg',2,'2024-09-28 13:56:42','2024-09-28 13:59:44');
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `variants`
--

DROP TABLE IF EXISTS `variants`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `variants` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `variant_name` varchar(50) NOT NULL,
  `quantity` int DEFAULT '0',
  `product_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `variants_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `variants`
--

LOCK TABLES `variants` WRITE;
/*!40000 ALTER TABLE `variants` DISABLE KEYS */;
INSERT INTO `variants` VALUES (1,'539aedff-651e-45d4-a636-4ef08a945621','GR Sport',650,1,'2024-09-28 14:04:16','2024-09-28 14:04:16'),(2,'84ba4355-e6fb-4dbb-bcac-3408f300c605','1.0 Non-Turbo',650,1,'2024-09-28 14:04:41','2024-09-28 14:04:41'),(3,'cb61c6ed-3378-4e68-910f-c9c6a58702d7','Cross',500,2,'2024-09-28 14:05:02','2024-09-28 14:05:02'),(4,'2642ea23-d0b5-4581-9571-c727b07821c0','Ultimate',800,2,'2024-09-28 14:05:16','2024-09-28 14:05:16'),(5,'55a6d2a2-cf13-41d0-b660-32f272626001','Exceed',950,2,'2024-09-28 14:05:44','2024-09-28 14:05:44'),(6,'469022c0-512f-45fc-99ca-a2a295593d2a','Signature',250,3,'2024-09-28 14:07:13','2024-09-28 14:07:13'),(7,'c5cf7c09-8dae-4cef-97ce-b67d42581a64','Prime',400,3,'2024-09-28 14:07:22','2024-09-28 14:07:22'),(8,'b14ffcf2-7100-4a18-8a14-c2df86bb9c9a','HWS One Tone',450,4,'2024-09-28 14:09:07','2024-09-28 14:09:07'),(9,'aaa8e0f0-06b9-40dc-9b0e-e9370b650c80','HWS Two Tone',300,4,'2024-09-28 14:09:20','2024-09-28 14:09:20'),(10,'d2e13da1-2706-4b10-b71f-98d83d9dbb93','2.0 XV',450,4,'2024-09-28 14:12:22','2024-09-28 14:12:22');
/*!40000 ALTER TABLE `variants` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-09-28 21:24:53
