--
-- Database: tdycwym
--

-- DROP DATABASE IF EXISTS `nju_iot`;
-- CREATE DATABASE `nju_iot` ;

--
-- Table structure for table `edgex_related_user`
--

DROP TABLE IF EXISTS `edgex_related_user`;

CREATE TABLE `edgex_related_user` (
	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
	`user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
	`username` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
	`edgex_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'edgex服务id',
	`edgex_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'edgex服务名',
	`status` tinyint NOT NULL DEFAULT '0' COMMENT '状态: 0-无关系 1-follow 2-unfollow',
	`created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`modified_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	PRIMARY KEY (`id`),
	UNIQUE KEY `idx_relation` (`user_id`,`edgex_id`),
	KEY `idx_created_time` (`created_time`),
	KEY `idx_modify_time` (`modified_time`),
	KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户收藏edgex记录表';

--
-- Dumping data for table `edgex_related_user`
--

LOCK TABLES `edgex_related_user` WRITE;
/*!40000 ALTER TABLE `edgex_related_user` DISABLE KEYS */;
INSERT INTO `edgex_related_user` (`id`, `user_id`, `username`, `edgex_id`, `edgex_name`, `status`, `created_time`, `modified_time`) VALUES (5,654321,'徐志乐',100000,'edgex inactive',1,'2021-04-22 10:59:51','2021-04-22 10:59:51'),(6,123456,'徐志乐',100004,'edgex test',1,'2021-04-22 11:02:50','2021-04-22 11:02:50'),(8,123456,'徐志乐',100000,'edgex-inactive',1,'2021-04-22 11:06:23','2021-04-22 11:06:23');
/*!40000 ALTER TABLE `edgex_related_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `edgex_service_item`
--

DROP TABLE IF EXISTS `edgex_service_item`;

CREATE TABLE `edgex_service_item` (
	`id` bigint unsigned NOT NULL AUTO_INCREMENT,
	`user_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建者',
	`edgex_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'edgex服务名',
	`prefix` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'edgex网关前缀',
	`status` tinyint NOT NULL DEFAULT '0' COMMENT '状态: 0-inactive 1-active',
	`deleted` tinyint NOT NULL DEFAULT '0' COMMENT '0-未删除 1-已删除',
	`address` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'edgex网关域名,e.g. 106.15.79.230:8080',
	`created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
	`modified_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
	`description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '描述',
	`location` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'edgex服务位置信息',
	`extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '额外信息',
	PRIMARY KEY (`id`),
	UNIQUE KEY `idx_prefix` (`prefix`,`deleted`),
	KEY `idx_created_time` (`created_time`),
	KEY `idx_modify_time` (`modified_time`),
	KEY `idx_user_id` (`user_id`),
	KEY `idx_edgex_name` (`edgex_name`)
) ENGINE=InnoDB AUTO_INCREMENT=100005 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='edgex服务表';

--
-- Dumping data for table `edgex_service_item`
--

LOCK TABLES `edgex_service_item` WRITE;
/*!40000 ALTER TABLE `edgex_service_item` DISABLE KEYS */;
INSERT INTO `edgex_service_item` (`id`, `user_id`, `edgex_name`, `prefix`, `status`, `deleted`, `address`, `created_time`, `modified_time`, `description`, `location`, `extra`) VALUES (100000,654321,'edgex inactive','edgex-inactive',0,0,'106.15.79.230:8080','2021-04-22 10:59:51','2021-04-22 10:59:51','edgex服务创建测试-inactive','{\"province\":\"江苏\",\"city\":\"南京市\"}',''),(100004,123456,'edgex test','edgex-test',1,0,'106.15.79.230:8080','2021-04-22 11:02:50','2021-04-22 11:02:52','edgex服务创建测试','{\"province\":\"江苏\",\"city\":\"南京市\"}','');
/*!40000 ALTER TABLE `edgex_service_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `edgex_user`
--

CREATE TABLE `edgex_user` (
	`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
	`username` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
	`password` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
	`phone_number` varchar (200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '电话号码',
	`email` varchar (200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',
	`deleted` tinyint NOT NULL DEFAULT '0' COMMENT '0-未删除 1-已删除',
	`verified` tinyint NOT NULL DEFAULT '0' COMMENT '0-未验证 1-已验证',
	`entrypted` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码保护问题',
	`created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modified_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=251 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;