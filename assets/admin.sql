/*
 Navicat Premium Data Transfer

 Source Server         : gl2gl
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 13/07/2022 21:03:20
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for admin
-- ----------------------------
DROP TABLE IF EXISTS "admin";
CREATE TABLE "admin" (
  "login_name" varchar NOT NULL,
  "password" varchar NOT NULL,
  PRIMARY KEY ("login_name")
);

-- ----------------------------
-- Records of admin
-- ----------------------------
BEGIN;
COMMIT;

PRAGMA foreign_keys = true;
