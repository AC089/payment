-- MySQL dump 10.13  Distrib 8.0.23, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: payment
-- ------------------------------------------------------
-- Server version	5.7.16

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `ali_conf`
--

DROP TABLE IF EXISTS `ali_conf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ali_conf` (
  `id` bigint(20) NOT NULL,
  `pay_app_id` bigint(20) NOT NULL,
  `app_id` varchar(32) NOT NULL COMMENT '支付宝分配给开发者的应用ID',
  `mch_id` bigint(20) NOT NULL COMMENT '支付宝账户ID',
  `app_name` varchar(45) NOT NULL COMMENT '应用名称',
  `pay_channel` tinyint(1) NOT NULL COMMENT '支付渠道(1:alipay-app 2:alipay-h5)',
  `return_url` varchar(255) DEFAULT NULL COMMENT '前台回跳地址',
  `notify_url` varchar(255) NOT NULL COMMENT '服务端异步通知地址',
  `private_key` varchar(2048) NOT NULL COMMENT '应用私钥',
  `alipay_public_key` varchar(512) NOT NULL COMMENT '支付宝公钥',
  `sign_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '商户生成签名字符串所使用的签名算法类型,1:RSA2, 2:RSA',
  `gateway` varchar(255) DEFAULT NULL COMMENT '应用网关',
  `secret_type` tinyint(1) DEFAULT NULL COMMENT '密钥类型, 1:普通公私钥, 2:公私钥证书',
  `app_cert_path` varchar(255) DEFAULT NULL COMMENT '应用公钥证书绝对路径',
  `alipay_cert_path` varchar(255) DEFAULT NULL COMMENT '支付宝公钥证书文件绝对路径',
  `alipay_root_cert_path` varchar(255) DEFAULT NULL COMMENT '支付宝 CA 根证书文件绝对路径',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='支付宝配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ali_conf`
--

LOCK TABLES `ali_conf` WRITE;
/*!40000 ALTER TABLE `ali_conf` DISABLE KEYS */;
/*!40000 ALTER TABLE `ali_conf` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pay_app`
--

DROP TABLE IF EXISTS `pay_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pay_app` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pay_app_code` varchar(45) NOT NULL COMMENT 'app代号',
  `pay_app_name` varchar(100) NOT NULL COMMENT 'app名称',
  `description` varchar(255) NOT NULL COMMENT '描述',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='app';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pay_app`
--

LOCK TABLES `pay_app` WRITE;
/*!40000 ALTER TABLE `pay_app` DISABLE KEYS */;
/*!40000 ALTER TABLE `pay_app` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `receipt_detail`
--

DROP TABLE IF EXISTS `receipt_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `receipt_detail` (
  `id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='ios内购凭据表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `receipt_detail`
--

LOCK TABLES `receipt_detail` WRITE;
/*!40000 ALTER TABLE `receipt_detail` DISABLE KEYS */;
/*!40000 ALTER TABLE `receipt_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refund_detail`
--

DROP TABLE IF EXISTS `refund_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refund_detail` (
  `id` bigint(20) NOT NULL,
  `transaction_no` varchar(32) NOT NULL COMMENT '支付渠道系统生成唯一订单号',
  `out_trade_no` varchar(32) NOT NULL COMMENT '商户网站|App唯一订单号',
  `out_refund_no` varchar(64) NOT NULL COMMENT '退款订单号',
  `reason` varchar(80) NOT NULL DEFAULT '' COMMENT '退款原因',
  `trade_detail_id` bigint(20) NOT NULL COMMENT '交易详情id',
  `amount` int(11) NOT NULL COMMENT '退款金额,单位分',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '退款状态,0-退款处理中 1-退款成功 2-退款关闭 3-退款异常',
  `client_ip` varchar(32) NOT NULL DEFAULT '' COMMENT '客户端ip',
  `device` varchar(45) NOT NULL DEFAULT '' COMMENT '设备标识',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '退款创建时间,取回调通知返回值',
  `success_time` int(11) NOT NULL DEFAULT '0' COMMENT '退款成功时间,取回调通知时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='退款详情表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refund_detail`
--

LOCK TABLES `refund_detail` WRITE;
/*!40000 ALTER TABLE `refund_detail` DISABLE KEYS */;
/*!40000 ALTER TABLE `refund_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `trade_detail`
--

DROP TABLE IF EXISTS `trade_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `trade_detail` (
  `id` bigint(20) NOT NULL,
  `pay_app_id` bigint(20) NOT NULL COMMENT '发起支付appid',
  `pay_app_code` varchar(45) NOT NULL COMMENT 'app代号',
  `pay_channel` tinyint(1) NOT NULL COMMENT '支付渠道(1:alipay-app 2:alipay-h5 3:wechat-app 4:wechat-h5 5:iosPaid)',
  `chan_app_id` varchar(32) NOT NULL COMMENT '支付渠道应用id',
  `out_trade_no` varchar(32) NOT NULL COMMENT '商户App唯一订单号',
  `transaction_no` varchar(32) DEFAULT '' COMMENT '支付渠道系统生成唯一订单号',
  `account_id` bigint(20) NOT NULL COMMENT '支付账户ID',
  `server_id` varchar(45) NOT NULL COMMENT '发起支付服务id',
  `amount` int(11) NOT NULL DEFAULT '0' COMMENT '支付金额,单位分',
  `actual_amount` int(11) NOT NULL DEFAULT '0' COMMENT '实际支付金额,单位分',
  `subject` varchar(255) NOT NULL DEFAULT '' COMMENT '商品标题',
  `body` varchar(255) NOT NULL DEFAULT '' COMMENT '商品描述信息',
  `client_ip` varchar(32) NOT NULL DEFAULT '' COMMENT '客户端ip',
  `device` varchar(45) NOT NULL DEFAULT '' COMMENT '设备标识',
  `player` varchar(128) NOT NULL DEFAULT '' COMMENT '用户标识:alipay为支付宝唯一用户号,wechat为openid',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '交易创建时间,取回调通知返回值',
  `success_time` int(11) NOT NULL DEFAULT '0' COMMENT '交易成功时间,取回调通知时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '支付状态:0-未支付 1-支付成功 2-业务处理成功',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `out_trade_no_UNIQUE` (`out_trade_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='支付交易详情表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `trade_detail`
--

LOCK TABLES `trade_detail` WRITE;
/*!40000 ALTER TABLE `trade_detail` DISABLE KEYS */;
/*!40000 ALTER TABLE `trade_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transfers_detail`
--

DROP TABLE IF EXISTS `transfers_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transfers_detail` (
  `id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='转账详情';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transfers_detail`
--

LOCK TABLES `transfers_detail` WRITE;
/*!40000 ALTER TABLE `transfers_detail` DISABLE KEYS */;
/*!40000 ALTER TABLE `transfers_detail` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `wechat_conf`
--

DROP TABLE IF EXISTS `wechat_conf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `wechat_conf` (
  `id` bigint(20) NOT NULL,
  `pay_app_id` bigint(20) NOT NULL,
  `mch_id` varchar(32) NOT NULL COMMENT '微信商户号ID',
  `app_id` varchar(32) NOT NULL COMMENT '应用ID',
  `pay_channel` tinyint(1) NOT NULL COMMENT '支付渠道(3:wechat-app 4:wechat-h5)',
  `api_key` varchar(32) NOT NULL COMMENT 'API密钥',
  `notify_url` varchar(255) NOT NULL COMMENT '服务端异步通知地址',
  `private_key_path` varchar(128) NOT NULL COMMENT '私钥证书绝对路径',
  `plat_cert_path` varchar(45) NOT NULL COMMENT '支付平台公钥证书绝对路径',
  `serial` varchar(45) NOT NULL COMMENT '商户API证书序列号',
  `apiv3` varchar(32) NOT NULL COMMENT '用商户平台上设置的APIv3密钥',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微信支付配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `wechat_conf`
--

LOCK TABLES `wechat_conf` WRITE;
/*!40000 ALTER TABLE `wechat_conf` DISABLE KEYS */;
/*!40000 ALTER TABLE `wechat_conf` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-09-13 13:41:33
