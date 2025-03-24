/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1-3306
 Source Server Type    : MySQL
 Source Server Version : 80013
 Source Host           : 127.0.0.1:3306
 Source Schema         : tkerland

 Target Server Type    : MySQL
 Target Server Version : 80013
 File Encoding         : 65001

 Date: 18/09/2024 00:48:50
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for app_startup
-- ----------------------------
DROP TABLE IF EXISTS `app_startup`;
CREATE TABLE `app_startup`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `app_version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `mac` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `app_id`(`app_id` ASC, `app_version` ASC) USING BTREE,
  INDEX `mac`(`mac` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of app_startup
-- ----------------------------
INSERT INTO `app_startup` VALUES (7, 'shadow-bot-copy', '1.0.0', '127.0.0.1', '123456', '2024-09-02 15:22:54');
INSERT INTO `app_startup` VALUES (8, 'shadow-bot-copy', '1.0.0', '127.0.0.1', '123456', '2024-09-02 15:23:05');
INSERT INTO `app_startup` VALUES (9, 'shadow-bot-copy', '1.0.0', '127.0.0.1', '123456', '2024-09-02 15:26:18');
INSERT INTO `app_startup` VALUES (10, 'shadow-bot-copy', '1.0.0', '127.0.0.1', '123456', '2024-09-02 15:31:48');
INSERT INTO `app_startup` VALUES (11, 'shadow-bot-copy', '1.0.0', '127.0.0.1', '123456', '2024-09-02 15:31:56');
INSERT INTO `app_startup` VALUES (12, 'shadow-bot-copy', '1.0.0', '127.0.0.1', '123456', '2024-09-02 15:32:01');
INSERT INTO `app_startup` VALUES (13, 'shadow-bot-copy', '1.0.0', '127.0.0.1', '123456', '2024-09-02 15:42:04');

-- ----------------------------
-- Table structure for calc_complete
-- ----------------------------
DROP TABLE IF EXISTS `calc_complete`;
CREATE TABLE `calc_complete`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `user_id` bigint(20) NOT NULL DEFAULT 0,
  `uuid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `is_public` tinyint(1) NULL DEFAULT 0,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uuid`(`uuid` ASC) USING BTREE,
  INDEX `user_id`(`user_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of calc_complete
-- ----------------------------
INSERT INTO `calc_complete` VALUES (37, '旅游账单', 3, '52b21856-e7d2-4e20-9e6c-06cadaefac8b', 0, '2024-09-01 09:51:30', '2024-09-02 21:46:14');

-- ----------------------------
-- Table structure for calc_expense
-- ----------------------------
DROP TABLE IF EXISTS `calc_expense`;
CREATE TABLE `calc_expense`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `party_id` bigint(20) NOT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `amount` decimal(10, 2) NOT NULL,
  `complete_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `party_id`(`party_id` ASC) USING BTREE,
  INDEX `complete_id`(`complete_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of calc_expense
-- ----------------------------
INSERT INTO `calc_expense` VALUES (24, 13, '住宿', 600.00, 37);
INSERT INTO `calc_expense` VALUES (25, 14, '吃饭', 300.00, 37);

-- ----------------------------
-- Table structure for calc_party
-- ----------------------------
DROP TABLE IF EXISTS `calc_party`;
CREATE TABLE `calc_party`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `share` tinyint(4) NOT NULL,
  `complete_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `name`(`name` ASC, `complete_id` ASC) USING BTREE,
  INDEX `complete_id`(`complete_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of calc_party
-- ----------------------------
INSERT INTO `calc_party` VALUES (13, '单位1', 1, 37);
INSERT INTO `calc_party` VALUES (14, '单位2', 2, 37);
INSERT INTO `calc_party` VALUES (15, '单位3', 3, 37);

-- ----------------------------
-- Table structure for calc_share
-- ----------------------------
DROP TABLE IF EXISTS `calc_share`;
CREATE TABLE `calc_share`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `party_id` bigint(20) NOT NULL,
  `share` tinyint(4) NOT NULL,
  `expense_id` bigint(20) NOT NULL,
  `complete_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `expense_id_party_id`(`expense_id` ASC, `party_id` ASC) USING BTREE,
  INDEX `party_id`(`party_id` ASC) USING BTREE,
  INDEX `complete_id`(`complete_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of calc_share
-- ----------------------------
INSERT INTO `calc_share` VALUES (79, 13, 1, 24, 37);
INSERT INTO `calc_share` VALUES (80, 14, 2, 24, 37);
INSERT INTO `calc_share` VALUES (81, 15, 3, 24, 37);
INSERT INTO `calc_share` VALUES (82, 13, 1, 25, 37);
INSERT INTO `calc_share` VALUES (83, 14, 1, 25, 37);
INSERT INTO `calc_share` VALUES (84, 15, 1, 25, 37);

-- ----------------------------
-- Table structure for log
-- ----------------------------
DROP TABLE IF EXISTS `log`;
CREATE TABLE `log`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `operation` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `details` json NULL,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of log
-- ----------------------------

-- ----------------------------
-- Table structure for login
-- ----------------------------
DROP TABLE IF EXISTS `login`;
CREATE TABLE `login`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `token` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `expires_at` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `token`(`token` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of login
-- ----------------------------
INSERT INTO `login` VALUES (1, 3, '77c928c7fec0de4ff05dcf259d3aea8417c64cdb9dcde160b1a44b6d952818ed', '2024-08-31 23:58:38', '2024-09-15 15:34:46');
INSERT INTO `login` VALUES (12, 3, '88705c509de6fe23c77e4de15969eb55b4c19b1b724ba63f4bf343782a6247cc', '2024-08-31 17:33:12', '2024-09-01 17:33:12');
INSERT INTO `login` VALUES (13, 3, '0a693689aabce616cd9843262625dc9aabce476f866b920b203d7b2481b753f9', '2024-09-01 23:59:37', '2024-09-02 23:59:37');
INSERT INTO `login` VALUES (14, 3, '024e8aad055fb252a81b95333ed8c5f7d3bb2108c07f72b6eda7eca0417d5bbb', '2024-09-13 23:43:36', '2024-09-14 23:43:36');
INSERT INTO `login` VALUES (15, 3, 'db1bbd46108396e3ac163a1faab69b96f3844c514df18be33285198caf14b6e9', '2024-09-14 23:03:49', '2024-09-15 23:03:49');
INSERT INTO `login` VALUES (16, 3, '1ad748de9d82172fd759f73c59792383703891fd46e7300ed0db894ad92a87e9', '2024-09-15 15:25:12', '2024-09-16 15:25:12');
INSERT INTO `login` VALUES (17, 3, 'f1625c09eaaa4936e62baa83d6ee16743781a25608dface83b3e1eb337c9c82b', '2024-09-15 15:38:51', '2024-09-15 15:39:01');
INSERT INTO `login` VALUES (18, 3, '273184eb251418e08e579732e10185bd03d6c9ba43bdbd965ab577f2dfadadb2', '2024-09-15 15:39:35', '2024-10-16 15:39:35');
INSERT INTO `login` VALUES (19, 3, 'd512ab68707708a47c2435da2315f09ac8d6d074a0f6f3554ff8a2a9c2941a90', '2024-09-16 22:14:28', '2024-09-17 22:14:28');

-- ----------------------------
-- Table structure for product
-- ----------------------------
DROP TABLE IF EXISTS `product`;
CREATE TABLE `product`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sku` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `status` tinyint(4) NOT NULL DEFAULT 0,
  `price` decimal(10, 2) NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `sku`(`sku` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of product
-- ----------------------------
INSERT INTO `product` VALUES (1, 'calc-month', '计算器月卡', '分享计算器一个月', 1, 9.90, '2024-09-10 22:38:03', '2024-09-10 22:38:07');
INSERT INTO `product` VALUES (2, 'calc-yeal', '计算器年卡', '分享计算器一年', 1, 99.00, '2024-09-10 22:38:03', '2024-09-13 14:43:49');
INSERT INTO `product` VALUES (3, 'test', '测试订阅服务', '测试订阅服务', 1, 0.01, '2024-09-10 22:38:03', '2024-09-13 14:43:49');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `is_verified` tinyint(1) NULL DEFAULT 0,
  `verify_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `email`(`email` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (3, '373944668@qq.com', '$2a$10$UfAg2sXZaiaOaGDxsLbBg.coCiuLVkjMhztRuKY7RND7lFGAqTrf.', 1, 'e6f4750bd13063189fe94f4ca95ec2aa', '2024-08-26 23:34:29', '2024-09-15 15:15:57');

-- ----------------------------
-- Table structure for user_order
-- ----------------------------
DROP TABLE IF EXISTS `user_order`;
CREATE TABLE `user_order`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `order_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `sku` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `user_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `amount` decimal(10, 2) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT 0,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `order_no`(`order_no` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 168 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_order
-- ----------------------------
INSERT INTO `user_order` VALUES (174, '20240917012538156462', 'test', '3', 0.01, 1, '2024-09-17 01:25:38', '2024-09-18 00:46:06');

-- ----------------------------
-- Table structure for user_payment
-- ----------------------------
DROP TABLE IF EXISTS `user_payment`;
CREATE TABLE `user_payment`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `order_id` bigint(20) NOT NULL,
  `third_trade_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `local_trade_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `amount` decimal(10, 2) NOT NULL,
  `status` tinyint(4) NOT NULL DEFAULT 0,
  `paid_at` datetime NULL DEFAULT NULL,
  `details` json NULL,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `local_trade_no`(`local_trade_no` ASC) USING BTREE,
  INDEX `order_id`(`order_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_payment
-- ----------------------------
INSERT INTO `user_payment` VALUES (2, 174, '2024091722001441231428539688', '20240918004605528149', 0.01, 1, '2024-09-18 00:46:06', '{\"body\": \"\", \"sign\": \"ABvUrWnBpzbhmCGLDdJSK2117PuKqMmeunWt2AR0HnEeZ1rAGq+ldXVA0uSKZ82sQUE1yXX25nNL3Uq6ACGvK8ulWLk/oaGFa9EP+YqPjKex/r1+M23wudNBJWqtiEkPiOgZrzxWD6GWAVQbPwxnSA0ox/vJXgXPFK2apVNRlE+NIrbQDy3N2mw0FDqh2RRKMM6JkenGueIkoL0Z8Pbl+Xhtc1SV0K60eFZEWttpJ/i+ew+GYUT/ymCOt9PN7ZtBGBM8SUZ9GGTRzNb2onvl0z3CiD9BkkbAnqOzxqh5mKeBk3xmxVBZCl15XZlcSU3Whw6kAHCtffdb1YE69Vd2CA==\", \"app_id\": \"2021004172647391\", \"charset\": \"utf-8\", \"subject\": \"测试订阅服务\", \"version\": \"1.0\", \"buyer_id\": \"\", \"trade_no\": \"2024091722001441231428539688\", \"gmt_close\": \"\", \"notify_id\": \"2024091701222092607041231425498712\", \"seller_id\": \"2088002507561426\", \"sign_type\": \"RSA2\", \"gmt_create\": \"2024-09-17 09:25:42\", \"gmt_refund\": \"\", \"out_biz_no\": \"\", \"refund_fee\": \"\", \"auth_app_id\": \"2021004172647391\", \"gmt_payment\": \"2024-09-17 09:26:06\", \"notify_time\": \"2024-09-17 09:26:07\", \"notify_type\": \"trade_status_sync\", \"agreement_no\": \"\", \"dback_amount\": \"\", \"dback_status\": \"\", \"out_trade_no\": \"20240917012538156462\", \"point_amount\": \"0.00\", \"seller_email\": \"alanglvtao@163.com\", \"total_amount\": \"0.01\", \"trade_status\": \"TRADE_SUCCESS\", \"bank_ack_time\": \"\", \"refund_amount\": \"\", \"refund_reason\": \"\", \"refund_status\": \"\", \"buyer_logon_id\": \"222***@qq.com\", \"fund_bill_list\": \"[{\\\"amount\\\":\\\"0.01\\\",\\\"fundChannel\\\":\\\"ALIPAYACCOUNT\\\"}]\", \"invoice_amount\": \"0.01\", \"out_request_no\": \"\", \"receipt_amount\": \"0.01\", \"passback_params\": \"\", \"buyer_pay_amount\": \"0.01\", \"voucher_detail_list\": \"\", \"external_agreement_no\": \"\"}', '2024-09-18 00:46:06');

-- ----------------------------
-- Table structure for user_product
-- ----------------------------
DROP TABLE IF EXISTS `user_product`;
CREATE TABLE `user_product`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sku` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `begin_time` datetime NULL DEFAULT NULL,
  `end_time` datetime NULL DEFAULT NULL,
  `user_id` bigint(20) NOT NULL DEFAULT 0,
  `payment_id` bigint(20) NULL DEFAULT NULL,
  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `app_id`(`sku` ASC, `begin_time` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Records of user_product
-- ----------------------------
INSERT INTO `user_product` VALUES (14, 'calc-month', '2024-09-10 22:38:03', '2024-10-10 22:38:03', 3, NULL, '2024-09-16 12:47:43');
INSERT INTO `user_product` VALUES (15, 'calc-yeal', '2024-09-10 22:38:03', '2025-10-10 22:38:03', 3, NULL, '2024-09-16 12:47:43');
INSERT INTO `user_product` VALUES (17, 'test', '2024-09-18 00:46:06', '2024-10-18 00:46:06', 3, 2, '2024-09-18 00:46:06');

SET FOREIGN_KEY_CHECKS = 1;
