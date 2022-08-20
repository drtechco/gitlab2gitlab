/*
 Navicat Premium Data Transfer

 Source Server         : gl2gl
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 13/07/2022 21:02:52
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for lang
-- ----------------------------
DROP TABLE IF EXISTS "lang";
CREATE TABLE "lang" (
  "key" text(255) NOT NULL,
  "memo" text(255),
  PRIMARY KEY ("key")
);

-- ----------------------------
-- Records of lang
-- ----------------------------
BEGIN;
INSERT INTO "lang" VALUES ('en-us', 'US English');
INSERT INTO "lang" VALUES ('zh-cn', '简体中文');
INSERT INTO "lang" VALUES ('zh-tw', '繁体中文');
COMMIT;

PRAGMA foreign_keys = true;
