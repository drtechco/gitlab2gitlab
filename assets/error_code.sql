/*
 Navicat Premium Data Transfer

 Source Server         : gl2gl
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 13/07/2022 21:02:43
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for error_code
-- ----------------------------
DROP TABLE IF EXISTS "error_code";
CREATE TABLE "error_code" (
  "code" integer NOT NULL,
  "i18n_key" text(255) NOT NULL,
  "memo" text(255),
  PRIMARY KEY ("code")
);

-- ----------------------------
-- Records of error_code
-- ----------------------------
BEGIN;
INSERT INTO "error_code" VALUES (100001, 'ERR_DB_QUERY_ERR', '数据查询错误');
INSERT INTO "error_code" VALUES (100002, 'ERR_AUTH_FAILD', '会话验证失败');
INSERT INTO "error_code" VALUES (100003, 'ERR_NOT_PERMISSION', '没有权限');
INSERT INTO "error_code" VALUES (100004, 'ERR_EXISTS_RECORDS', '已存在记录');
INSERT INTO "error_code" VALUES (100005, 'ERR_SAVE_DB', '保存失败');
INSERT INTO "error_code" VALUES (100006, 'ERR_RECORD_NOT_EXISTS', '记录不存在');
INSERT INTO "error_code" VALUES (100007, 'ERR_CLIENT_PARAMETER', '客户端参数错误');
INSERT INTO "error_code" VALUES (100011, 'ERR_LOGIN_FAILED', '登录失败');
INSERT INTO "error_code" VALUES (100012, 'ERR_SESSION_NOT_VERFY', '会话验证失败');
INSERT INTO "error_code" VALUES (100014, 'ERR_DIR_NOT_EXISTS', '文件夹不存在');
INSERT INTO "error_code" VALUES (100015, 'ERR_PATH_CHECK_FAILED', '路径检查错误');
INSERT INTO "error_code" VALUES (100016, 'ERR_DIR_CREATE_FAILED', '文件夹创建错误');
INSERT INTO "error_code" VALUES (100017, 'ERR_OLD_PWD_NOT', '旧密码不符');
INSERT INTO "error_code" VALUES (100018, 'ERR_FILE_NOT_EXISTS', '文件不存在');
INSERT INTO "error_code" VALUES (100022, 'ERR_READ_ONLY_RECORD', '记录不能被修改');
INSERT INTO "error_code" VALUES (100023, 'ERR_SOCKET_UPGRADE', 'WEB SOCKET 升级失败');
INSERT INTO "error_code" VALUES (100046, 'ERR_JSON_MARSHAL_NOT', 'json转义失败');
INSERT INTO "error_code" VALUES (100057, 'ERR_IP_NOT', 'ip地址不合法');
INSERT INTO "error_code" VALUES (100061, 'ERR_READONLY_RECORD', '只读记录不可编辑');
INSERT INTO "error_code" VALUES (100063, 'ERR_POST_REPEAT', '表单重复提交');
INSERT INTO "error_code" VALUES (100067, 'ERR_PASSWORD_NOT', '密码错误');
INSERT INTO "error_code" VALUES (100069, 'ERR_IMLEMENT_OK', '已执行无法删除');
COMMIT;

PRAGMA foreign_keys = true;
