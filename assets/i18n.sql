/*
 Navicat Premium Data Transfer

 Source Server         : gl2gl
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 13/07/2022 21:03:02
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for i18n
-- ----------------------------
DROP TABLE IF EXISTS "i18n";
CREATE TABLE "i18n" (
  "lang_key" varchar(255) NOT NULL,
  "key" varchar(255) NOT NULL,
  "value" varchar(3000) NOT NULL,
  "memo" varchar(255),
  PRIMARY KEY ("lang_key", "key")
);

-- ----------------------------
-- Records of i18n
-- ----------------------------
BEGIN;
INSERT INTO "i18n" VALUES ('en-us', 'ERR_AUTH_FAILD', '会话验证失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_CLIENT_PARAMETER', '客户端参数错误', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_DATETIME_FUTURE', '时间点不可超过当日', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_DB_QUERY_ERR', '数据查询错误', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_DIR_CREATE_FAILED', '文件夹创建错误', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_DIR_NOT_EXISTS', '文件夹不存在', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_EDIT_NOT', '修改失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_EXISTS_RECORDS', '已存在记录', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_FILE_NOT_EXISTS', '文件不存在', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_IMLEMENT_OK', '已执行无法删除', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_JSON_MARSHAL_NOT', 'json转义失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_LOGIN_FAILED', '登录失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_LOTTERY_COUNT_COMMAND', '发送消息队列失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_NOT_UNBIND', '解绑失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_OLD_PWD_NOT', '旧密码不符', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_PASSWORD_NOT', '密码错误', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_PATH_CHECK_FAILED', '路径检查错误', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_POST_REPEAT', '表单重复提交', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_READONLY_RECORD', '只读记录不可修改', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_READ_ONLY_RECORD', '记录不能被修改', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_RECORD_NOT_EXISTS', '记录不存在', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_SAVE_DB', '保存失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_SESSION_NOT_VERFY', '会话验证失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_SQL_COMMIT_WORNG', '事务提交失败', NULL);
INSERT INTO "i18n" VALUES ('en-us', 'ERR_SQL_ROLLBACK_WORNG', '事务回滚失败', NULL);
INSERT INTO "i18n" VALUES ('zh-cn', 'ERR_DEL_NOT', '删除失败', NULL);
INSERT INTO "i18n" VALUES ('zh-cn', 'ERR_IP_NOT', 'ip地址不合法', NULL);
INSERT INTO "i18n" VALUES ('zh-cn', 'ERR_SOCKET_UPGRADE', 'WEB SOCKET 升级失败', NULL);
COMMIT;

PRAGMA foreign_keys = true;
