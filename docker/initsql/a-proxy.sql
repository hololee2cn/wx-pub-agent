-- ----------------------------
-- 初始化数据库pub_platform_mgr
-- ----------------------------

CREATE DATABASE IF NOT EXISTS pub_platform_mgr
    DEFAULT CHARSET utf8mb4
    COLLATE utf8mb4_general_ci;

USE pub_platform_mgr;
SET NAMES utf8mb4;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
                         `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
                        `open_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'open id',
                        `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
                        `phone` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户号码',
                        `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
                        `delete_time` int(11) NOT NULL DEFAULT 0 COMMENT '删除时间',
                        PRIMARY KEY (`id`) USING BTREE,
                        UNIQUE INDEX `open_id`(`open_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户user表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for msg_log
-- ----------------------------
DROP TABLE IF EXISTS `msg_log`;
CREATE TABLE `msg_log`  (
                         `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
                         `request_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '请求id',
                         `msg_id` bigint(64) NOT NULL DEFAULT 0 COMMENT '微信消息id',
                         `to_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接收者openid',
                         `phone` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户号码',
                         `template_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '模板id',
                         `content` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '发送模板内容',
                         `cause` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '失败原因',
                         `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '发送状态，0为pending判定中，1为sending发送中，2为success成功，3为failure失败',
                         `count` tinyint(1) NOT NULL DEFAULT 0 COMMENT '发送次数',
                         `create_time` int(11) NOT NULL DEFAULT 0 COMMENT '创建时间',
                         `update_time` int(11) NOT NULL DEFAULT 0 COMMENT '更新时间',
                         PRIMARY KEY (`id`) USING BTREE,
                         INDEX `request_id`(`request_id`) USING BTREE,
                         INDEX `msg_id`(`msg_id`) USING BTREE,
                         INDEX `status`(`status`) USING BTREE,
                         INDEX `count`(`count`) USING BTREE,
                         INDEX `create_time`(`create_time`) USING BTREE,
                         INDEX `update_time`(`update_time`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '消息日志表' ROW_FORMAT = DYNAMIC;
