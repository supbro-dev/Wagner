-- 员工表

CREATE TABLE `employee`
(
    `id`              bigint      NOT NULL AUTO_INCREMENT COMMENT 'Id',
    `gmt_create`      datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `name`            varchar(45) NOT NULL COMMENT '员工姓名',
    `number`          varchar(45) NOT NULL COMMENT '员工工号',
    `identity`        varchar(45) NOT NULL COMMENT '身份证号',
    `sensitive_info`  json                 DEFAULT NULL,
    `workplace_code`  varchar(45)          DEFAULT NULL COMMENT '工作点编码',
    `position_code`   varchar(45)          DEFAULT NULL COMMENT '岗位编码',
    `work_group_code` varchar(45)          DEFAULT NULL COMMENT '工作组编码',
    `status`          varchar(45) NOT NULL COMMENT '员工状态',
    PRIMARY KEY (`id`),
    UNIQUE KEY `number_UNIQUE` (`number`),
    UNIQUE KEY `identity_UNIQUE` (`identity`),
    KEY               `idx_workplace_status` (`workplace_code`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='员工信息';

-- 员工动作表

CREATE TABLE `action`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `gmt_create`      datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 (毫秒精度)',
    `gmt_modified`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间 (毫秒精度)',
    `employee_number` varchar(64) NOT NULL COMMENT '员工工号',
    `workplace_code`  varchar(64) NOT NULL COMMENT '工作点编码',
    `start_time`      datetime(3) NOT NULL COMMENT '动作开始时间',
    `end_time`        datetime(3) DEFAULT NULL COMMENT '动作结束时间',
    `action_type`     varchar(64) NOT NULL COMMENT '动作类型',
    `properties`      json                 DEFAULT NULL COMMENT '动作属性 (JSON格式)',
    `operate_day`     date        NOT NULL COMMENT '动作日期',
    `action_code`     varchar(64) NOT NULL COMMENT '动作编码',
    `work_load`       json                 DEFAULT NULL COMMENT '工作量',
    PRIMARY KEY (`id`),
    KEY               `idx_employee_start_time_workplace` (`employee_number`,`start_time`,`workplace_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户动作表';

-- 工作点表

CREATE TABLE `workplace`
(
    `id`                int         NOT NULL AUTO_INCREMENT,
    `gmt_create`        datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified`      datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `code`              varchar(45) NOT NULL COMMENT '工作点编码',
    `name`              varchar(45) NOT NULL COMMENT '工作点名称',
    `region_code`       varchar(45)          DEFAULT NULL COMMENT '区域编码',
    `industry_code`     varchar(45) NOT NULL COMMENT '行业编码',
    `sub_industry_code` varchar(45)          DEFAULT NULL COMMENT '子行业编码',
    PRIMARY KEY (`id`),
    UNIQUE KEY `code_UNIQUE` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='工作点';

-- 动态计算配置表

CREATE TABLE `calc_dynamic_param`
(
    `id`                bigint      NOT NULL AUTO_INCREMENT,
    `gmt_create`        datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified`      datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `type`              varchar(45) NOT NULL COMMENT '参数类型',
    `industry_code`     varchar(45) NOT NULL COMMENT '行业编码',
    `sub_industry_code` varchar(45)          DEFAULT NULL COMMENT '子行业编码',
    `content`           json        NOT NULL COMMENT '配置内容',
    PRIMARY KEY (`id`),
    KEY                 `idx_ic_sic` (`industry_code`,`sub_industry_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='动态计算配置';

-- 人效环节
CREATE TABLE `standard_position`
(
    `id`                bigint      NOT NULL AUTO_INCREMENT,
    `gmt_create`        datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified`      datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `code`              varchar(45) NOT NULL COMMENT '编码',
    `name`              varchar(45) NOT NULL COMMENT '名称',
    `parent_code`       varchar(45) NOT NULL COMMENT '上一级编码',
    `type`              varchar(45) NOT NULL COMMENT '类型',
    `level`             varchar(45) NOT NULL COMMENT '层级',
    `version`           int         NOT NULL COMMENT '版本号',
    `industry_code`     varchar(45) NOT NULL COMMENT '行业编码',
    `sub_industry_code` varchar(45)          DEFAULT NULL COMMENT '子行业编码',
    `properties`        json                 DEFAULT NULL COMMENT '其他属性',
    `script`            text COMMENT '环节匹配脚本(EL表达式)',
    `order`             int                  DEFAULT NULL COMMENT '顺序',
    PRIMARY KEY (`id`),
    KEY                 `idx_version` (`version`),
    KEY                 `idx_ic_sic` (`industry_code`,`sub_industry_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='人效环节';

-- 小时聚合结果表

CREATE TABLE `hour_summary_result`
(
    `id`                     bigint       NOT NULL AUTO_INCREMENT,
    `gmt_create`             datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified`           datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `operate_time`           datetime     NOT NULL COMMENT '作业时间（小时）',
    `operate_day`            date         NOT NULL COMMENT '作业日期',
    `process_code`           varchar(45)  NOT NULL COMMENT '环节编码',
    `position_code`          varchar(45)           DEFAULT NULL COMMENT '作业岗位编码',
    `workplace_code`         varchar(45)  NOT NULL COMMENT '工作点编码',
    `workplace_name`         varchar(45)  NOT NULL DEFAULT '""',
    `employee_number`        varchar(45)  NOT NULL COMMENT '员工工号',
    `employee_name`          varchar(45)  NOT NULL COMMENT '员工姓名',
    `employee_position_code` varchar(45)           DEFAULT NULL COMMENT '员工归属岗位',
    `work_group_code`        varchar(45)           DEFAULT NULL COMMENT '员工工作组编码',
    `region_code`            varchar(45)           DEFAULT NULL COMMENT '工作点所属区域',
    `industry_code`          varchar(45)  NOT NULL COMMENT '工作点所属行业',
    `sub_industry_code`      varchar(45)           DEFAULT NULL COMMENT '工作点所属子行业',
    `work_load`              json                  DEFAULT NULL COMMENT '员工工作量',
    `direct_work_time`       int          NOT NULL DEFAULT '0' COMMENT '直接作业时长（秒）',
    `indirect_work_time`     int          NOT NULL DEFAULT '0' COMMENT '间接作业时长',
    `idle_time`              int          NOT NULL DEFAULT '0' COMMENT '闲置时长',
    `rest_time`              int          NOT NULL DEFAULT '0',
    `attendance_time`        int          NOT NULL DEFAULT '0',
    `process_property`       json         NOT NULL COMMENT '环节其他属性',
    `properties`             json                  DEFAULT NULL COMMENT '其他信息',
    `unique_key`             varchar(128) NOT NULL DEFAULT '""' COMMENT '业务唯一键',
    `is_deleted`             tinyint      NOT NULL DEFAULT '0' COMMENT '是否删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_unique_key` (`unique_key`),
    KEY                      `idx_en_od_pc_ot_wp` (`employee_number`,`operate_day`,`position_code`,`operate_time`,`workplace_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='小时聚合结果';

CREATE TABLE `employee_status`
(
    `id`               bigint      NOT NULL AUTO_INCREMENT COMMENT 'Id',
    `gmt_create`       datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `gmt_modified`     datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `employee_number`  varchar(45) NOT NULL COMMENT '员工工号',
    `employee_name`    varchar(45) NOT NULL COMMENT '员工姓名',
    `operate_day`      date        NOT NULL COMMENT '工作日期',
    `workplace_code`   varchar(45) NOT NULL COMMENT '工作点编码',
    `status`           varchar(45) NOT NULL COMMENT '员工状态',
    `last_action_time` datetime    NOT NULL COMMENT '上次动作发生时间',
    `last_action_code` varchar(45) NOT NULL COMMENT '上次动作发生编码',
    `work_group_code`  varchar(45)          DEFAULT NULL COMMENT '工作组编码',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_wc_od_en` (`workplace_code`,`operate_day`,`employee_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='员工状态表';

-- 脚本定义表（暂时没用到）
CREATE TABLE `script`
(
    `id`           bigint      NOT NULL AUTO_INCREMENT,
    `gmt_create`   datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `gmt_modified` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `name`         varchar(45) NOT NULL COMMENT '脚本名称',
    `type`         varchar(45) NOT NULL COMMENT '脚本类型',
    `desc`         varchar(45)          DEFAULT NULL COMMENT '脚本描述',
    `content`      text COMMENT '脚本内容',
    `version`      int         NOT NULL COMMENT '版本',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name_version` (`name`,`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='脚本';