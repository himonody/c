DROP TABLE IF EXISTS `admin_sys_config`;
CREATE TABLE `admin_sys_config` (
                                    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键编码',
                                    `config_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'ConfigName',
                                    `config_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'ConfigKey',
                                    `config_value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'ConfigValue',
                                    `config_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'ConfigType',
                                    `is_frontend` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '是否前台',
                                    `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'Remark',
                                    `create_by` bigint DEFAULT NULL COMMENT '创建者',
                                    `update_by` bigint DEFAULT NULL COMMENT '更新者',
                                    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后更新时间',
                                    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='配置管理';


DROP TABLE IF EXISTS `admin_sys_user`;
CREATE TABLE `admin_sys_user` (
                                  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '编码',
                                  `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '用户名',
                                  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '密码',
                                  `nick_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '昵称',
                                  `role` int DEFAULT 1 COMMENT '1:superadmin 2:user',
                                  `salt` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '加盐',
                                  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '备注',
                                  `status`  int DEFAULT 1 COMMENT '1:启用 2:禁用',
                                  `create_by` bigint DEFAULT NULL COMMENT '创建者',
                                  `update_by` bigint DEFAULT NULL COMMENT '更新者',
                                  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后更新时间',
                                  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='系统用户管理';



DROP TABLE IF EXISTS `app_user`;
CREATE TABLE `app_user` (
                            `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户id',
                            `level_id` int NOT NULL DEFAULT '1' COMMENT '用户等级编号',
                            `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '账号名称/用户名',
                            `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '用户昵称',
                            `avatar` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '头像路径',
                            `pwd` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '登录密码',
                            `ref_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '推荐码',
                            `ref_id` int NOT NULL DEFAULT '0' COMMENT '推荐id',
                            `friend_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '邀请码',
                            `friend_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '邀请码',
                            `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '1' COMMENT '状态(1-正常 2-异常)',
                            `online_status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '1' COMMENT '在线状态(1-离线 2-在线)',
                            `register_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
                            `register_ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '注册IP',
                            `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
                            `last_login_ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT '最后登录IP',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',

                            PRIMARY KEY (`id`) USING BTREE,
                            UNIQUE KEY `uk_username` (`username`),
                            UNIQUE KEY `uk_ref_code` (`ref_code`),
                            KEY `idx_ref_id` (`ref_id`),
                            KEY `idx_status` (`status`),
                            KEY `idx_online_status` (`online_status`),
                            KEY `idx_register_at` (`register_at`),
                            KEY `idx_last_login_at` (`last_login_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户管理';

DROP TABLE IF EXISTS `app_user_wallet`;
CREATE TABLE `app_user_wallet` (
                            `user_id` bigint DEFAULT 0  COMMENT '用户id',
                            `pay_pwd` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '提现密码',
                            `pay_status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '1' COMMENT '提现状态(1-启用 2-禁用)',
                            `balance`decimal(30,2) NOT NULL DEFAULT '0.00' COMMENT '余额',
                            `frozen`decimal(30,2) NOT NULL DEFAULT '0.00' COMMENT '冻结金额',
                            `total_r`decimal(30,2) NOT NULL DEFAULT '0.00' COMMENT '总充值',
                            `total_w`decimal(30,2) NOT NULL DEFAULT '0.00' COMMENT '总提现',
                            `total_re`decimal(30,2) NOT NULL DEFAULT '0.00' COMMENT '打卡总收益',
                            `total_i`decimal(30,2) NOT NULL DEFAULT '0.00' COMMENT '邀请总收益',
                            `address` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '提现地址',
                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                            UNIQUE KEY `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户钱包';


DROP TABLE IF EXISTS `app_user_withdraw`;
CREATE TABLE `app_user_withdraw` (
                                     `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键自增ID',
                                     `biz_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '业务唯一订单号(对外展示)',
                                     `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                                     `amount` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000' COMMENT '提现金额',
                                     `fee` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000' COMMENT '提现手续费',
                                     `actual_amount` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000' COMMENT '实际到账金额(amount-fee)',
                                     `address` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '提现目标地址(钱包/银行卡号)',
                                     `apply_ip` VARCHAR(45) NOT NULL DEFAULT '' COMMENT '申请人IP',
                                     `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1待审核, 2审核通过(待打款), 3拒绝, 4打款成功, 5打款失败',
                                     `reject_reason` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '拒绝或失败原因',
                                     `tx_hash` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '第三方转账流水号/区块链哈希',

                                     `review_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '审核人管理员ID',
                                     `review_ip` VARCHAR(45) NOT NULL DEFAULT '' COMMENT '审核人IP',
                                     `reviewed_at` DATETIME DEFAULT NULL COMMENT '审核完成时间',

                                     `version` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
                                     `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                     `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',

                                     PRIMARY KEY (`id`),
                                     UNIQUE KEY `unq_biz_id` (`biz_id`),
                                     INDEX `idx_user_updated` (`user_id`, `created_at` DESC),
                                     INDEX `idx_sort_logic` (`reviewed_at`, `created_at` DESC),
                                     INDEX `idx_status_updated` (`status`, `updated_at` DESC),
                                     INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='提现申请表';

DROP TABLE IF EXISTS `app_user_balance_log`;
CREATE TABLE `app_user_balance_log` (
                                        `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增ID',
                                        `biz_id` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '外部业务关联ID(对应提现表的biz_id)',
                                        `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                                        `type` TINYINT NOT NULL DEFAULT 0 COMMENT '账变类型: 1提现申请冻结, 2提现成功扣除, 3提现拒绝退回, 4充值, 5活动奖励...',

                                        `amount` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000' COMMENT '变动金额(带符号, 如-10.00)',
                                        `before_balance` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000' COMMENT '变动前余额',
                                        `after_balance` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000' COMMENT '变动后余额',
                                        `before_frozen` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000' COMMENT '变动前冻结金额',
                                        `after_frozen` DECIMAL(32, 8) NOT NULL DEFAULT '0.00000000' COMMENT '变动后冻结金额',

                                        `remark` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '业务备注',
                                        `operator_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '操作人ID(0为系统, 其余为管理员ID)',
                                        `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

                                        PRIMARY KEY (`id`),
    -- 核心索引：用户查询账单流水
                                        INDEX `idx_user_type_created` (`user_id`, `type`, `created_at` DESC),
    -- 业务索引：用于根据提现单号回溯账变记录
                                        INDEX `idx_biz_id` (`biz_id`),
                                        INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户资金账变表';

-- ----------------------------
-- Table structure for app_challenge_config
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge`;
CREATE TABLE `app_challenge` (
                                 `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '活动配置ID',
                                 `day_count` INT NOT NULL DEFAULT 0 COMMENT '挑战天数 1/7/21',
                                 `amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '单人挑战金额',
                                 `checkin_start` SMALLINT UNSIGNED DEFAULT 0 COMMENT '每日打卡开始时间',
                                 `checkin_end` SMALLINT UNSIGNED DEFAULT 0 COMMENT '每日打卡结束时间',
                                 `platform_bonus` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '平台补贴金额',
                                 `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1启用 2停用',
                                 `sort` TINYINT NOT NULL DEFAULT 1 COMMENT '排序 1 2 3 4',
                                 `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                 PRIMARY KEY (`id`),
                                 UNIQUE KEY `uk_day_amount` (`day_count`,`amount`),
                                 KEY `idx_status` (`status`),
                                 KEY `idx_sort` (`sort` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin  COMMENT='打卡挑战活动配置';

-- ----------------------------
-- Table structure for app_challenge_user
-- ----------------------------
DROP TABLE IF EXISTS `app_user_challenge`;
CREATE TABLE `app_user_challenge` (
                                      `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户挑战ID',
                                      `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                                      `challenge_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '活动配置ID',
                                      `pool_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '奖池ID',
                                      `challenge_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '用户挑战金额',
                                      `start_date` DATE NOT NULL COMMENT '活动开始日期 YYYYMMDD',
                                      `end_date`  DATE NOT NULL COMMENT '活动结束日期 YYYYMMDD',
                                      `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1进行中 2成功 3失败',
                                      `fail_reason` TINYINT NOT NULL DEFAULT 2 COMMENT '失败原因 0无 1已打卡 2未打卡',
                                      `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '报名时间',
                                      `finished_at` datetime  DEFAULT NULL COMMENT '完成时间',
                                      PRIMARY KEY (`id`),
                                      UNIQUE KEY `uk_user_active` (`user_id`,`status`),
                                      KEY `idx_pool` (`pool_id`),
                                      KEY `idx_challenge_id` (`challenge_id`),
                                      KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户参与挑战记录';

-- ----------------------------
-- Table structure for app_challenge_checkin
-- ----------------------------
DROP TABLE IF EXISTS `app_user_challenge_checkin`;
CREATE TABLE `app_user_challenge_checkin` (
                                              `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '打卡ID',
                                              `user_challenge_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户挑战ID',
                                              `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                                              `checkin_date` DATE NOT NULL COMMENT '打卡日期 YYYYMMDD',
                                              `checkin_time` datetime  DEFAULT NULL COMMENT '打卡时间',
                                              `mood_code` TINYINT NOT NULL DEFAULT 0 COMMENT '心情枚举 1开心 2平静 3一般 4疲惫 5低落 6爆棚',
                                              `mood_text` TEXT NOT NULL DEFAULT '' COMMENT '用户心情文字描述（最多200字）',
                                              `content_type` TINYINT NOT NULL DEFAULT 1 COMMENT '打卡内容类型 1图片 2视频广告',
                                              `status` TINYINT NOT NULL DEFAULT 2 COMMENT '状态 1打卡成功 2未打卡 3打卡失败',
                                              `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',

                                              PRIMARY KEY (`id`),
                                              UNIQUE KEY `uk_challenge_date` (`user_challenge_id`,`checkin_date`),
                                              KEY `idx_user_date` (`user_id`,`checkin_date`),
                                              KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='每日打卡记录（心情 + 内容）';

-- ----------------------------
-- Table structure for app_challenge_checkin_image
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge_checkin_image`;
CREATE TABLE `app_challenge_checkin_image` (
                                               `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '图片ID',
                                               `checkin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '打卡ID',
                                               `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                                               `image_url` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '图片URL',
                                               `image_hash` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '图片Hash（防重复）',
                                               `sort_no` TINYINT NOT NULL DEFAULT 1 COMMENT '图片顺序',
                                               `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1正常 2屏蔽 3审核中',
                                               `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
                                               PRIMARY KEY (`id`),
                                               UNIQUE KEY `uk_checkin_hash` (`checkin_id`,`image_hash`),
                                               KEY `idx_checkin` (`checkin_id`),
                                               KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='打卡图片表';

-- ----------------------------
-- Table structure for app_challenge_checkin_video_ad
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge_checkin_video_ad`;
CREATE TABLE `app_challenge_checkin_video_ad` (
                                                  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '视频广告打卡ID',

                                                  `checkin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联打卡ID',
                                                  `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',

                                                  `ad_platform` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '广告平台 如：csj、gdt、unity',
                                                  `ad_unit_id` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '广告位ID',
                                                  `ad_order_no` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '广告联盟返回的订单号（唯一）',

                                                  `video_duration` INT NOT NULL DEFAULT 0 COMMENT '视频时长（秒）',
                                                  `watch_duration` INT NOT NULL DEFAULT 0 COMMENT '实际观看时长（秒）',

                                                  `reward_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '该广告产生的收益',
                                                  `verify_status` TINYINT NOT NULL DEFAULT 0 COMMENT '校验状态 0待校验 1成功 2失败',

                                                  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '观看完成时间',
                                                  `verified_at` datetime  DEFAULT NULL COMMENT '校验完成时间',

                                                  PRIMARY KEY (`id`),
                                                  UNIQUE KEY `uk_ad_order` (`ad_order_no`),
                                                  UNIQUE KEY `uk_checkin` (`checkin_id`),
                                                  KEY `idx_user` (`user_id`),
                                                  KEY `idx_verify_status` (`verify_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='视频广告打卡记录';



-- ----------------------------
-- Table structure for app_challenge_pool
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge_pool`;
CREATE TABLE `app_challenge_pool` (
                                      `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '奖池ID',
                                      `challenge_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '活动配置ID',
                                      `start_date` datetime  DEFAULT NULL COMMENT '活动开始日期',
                                      `end_date` datetime  DEFAULT NULL COMMENT '活动结束日期',
                                      `total_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '奖池当前总金额',
                                      `settled` TINYINT NOT NULL DEFAULT 0 COMMENT '是否已结算 0否 1是',
                                      `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                      PRIMARY KEY (`id`),
                                      UNIQUE KEY `uk_config_date` (`challenge_id`,`start_date`),
                                      KEY `idx_settled` (`settled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='活动奖池表';

-- ----------------------------
-- Table structure for app_challenge_pool_flow
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge_pool_flow`;
CREATE TABLE `app_challenge_pool_flow` (
                                           `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '奖池流水ID',
                                           `pool_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '奖池ID',
                                           `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                                           `amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '变动金额',
                                           `type` TINYINT NOT NULL DEFAULT 0 COMMENT '类型 1报名 2失败 3平台补贴 4结算',
                                           `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                           PRIMARY KEY (`id`),
                                           KEY `idx_pool` (`pool_id`),
                                           KEY `idx_user` (`user_id`),
                                           KEY `idx_type_time` (`type`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='奖池资金流水';

-- ----------------------------
-- Table structure for app_challenge_settlement
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge_settlement`;
CREATE TABLE `app_challenge_settlement` (
                                            `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '结算ID',
                                            `user_challenge_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户挑战ID',
                                            `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                                            `reward` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '最终获得金额',
                                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结算时间',
                                            PRIMARY KEY (`id`),
                                            UNIQUE KEY `uk_challenge_user` (`user_challenge_id`,`user_id`),
                                            KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='挑战结算结果';

-- ----------------------------
-- Table structure for app_challenge_daily_stat
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge_daily_stat`;
CREATE TABLE `app_challenge_daily_stat` (
                                            `stat_date` DATE NOT NULL COMMENT '统计日期 YYYYMMDD',
                                            `join_user_cnt` INT NOT NULL DEFAULT 0 COMMENT '参与人数',
                                            `success_user_cnt` INT NOT NULL DEFAULT 0 COMMENT '成功人数',
                                            `fail_user_cnt` INT NOT NULL DEFAULT 0 COMMENT '失败人数',
                                            `join_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '参与总金额',
                                            `success_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '成功金额',
                                            `fail_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '失败金额',
                                            `platform_bonus` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '平台补贴',
                                            `pool_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '奖池金额',
                                            `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                                            PRIMARY KEY (`stat_date`),
                                            KEY `idx_date` (`stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='每日运营统计';

-- ----------------------------
-- Table structure for app_challenge_total_stat
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge_total_stat`;
CREATE TABLE `app_challenge_total_stat` (
                                            `id` TINYINT NOT NULL DEFAULT 1 COMMENT '固定主键',
                                            `total_user_cnt` INT NOT NULL DEFAULT 0 COMMENT '累计用户数',
                                            `total_join_cnt` INT NOT NULL DEFAULT 0 COMMENT '累计参与人次',
                                            `total_success_cnt` INT NOT NULL DEFAULT 0 COMMENT '累计成功人次',
                                            `total_fail_cnt` INT NOT NULL DEFAULT 0 COMMENT '累计失败人次',
                                            `total_join_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '累计参与金额',
                                            `total_success_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '累计成功金额',
                                            `total_fail_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '累计失败金额',
                                            `total_platform_bonus` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '累计平台补贴',
                                            `total_pool_amount` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '累计奖池金额',
                                            `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='平台累计统计';

-- ----------------------------
-- Table structure for app_challenge_rank_daily
-- ----------------------------
DROP TABLE IF EXISTS `app_challenge_rank_daily`;
CREATE TABLE `app_challenge_rank_daily` (
                                            `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '排行ID',
                                            `rank_date` DATE NOT NULL COMMENT '排行日期',
                                            `rank_type` TINYINT NOT NULL DEFAULT 0 COMMENT '1邀请 2收益 3毅力',
                                            `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
                                            `value` DECIMAL(30,2) NOT NULL DEFAULT 0.00 COMMENT '排行值',
                                            `rank_no` INT NOT NULL DEFAULT 0 COMMENT '排名',
                                            PRIMARY KEY (`id`),
                                            UNIQUE KEY `uk_rank` (`rank_date`,`rank_type`,`user_id`),
                                            KEY `idx_rank_type` (`rank_type`,`rank_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='排行榜日快照';





