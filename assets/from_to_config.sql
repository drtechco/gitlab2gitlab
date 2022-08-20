/*
 Navicat Premium Data Transfer

 Source Server         : gl2gl
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 13/07/2022 19:55:15
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for from_to_config
-- ----------------------------
DROP TABLE IF EXISTS "from_to_config";
CREATE TABLE "from_to_config" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "from_address" TEXT NOT NULL,
  "from_access_token" TEXT NOT NULL,
  "to_address" TEXT NOT NULL,
  "to_access_token" TEXT NOT NULL,
  "status" INTEGER NOT NULL,
  "delete_ branch" INTEGER NOT NULL,
  "last_sync_time" datetime NOT NULL,
  "last_sync_status" INTEGER NOT NULL
);

-- ----------------------------
-- Auto increment value for from_to_config
-- ----------------------------
UPDATE "main"."sqlite_sequence" SET seq = 1 WHERE name = 'from_to_config';

PRAGMA foreign_keys = true;
