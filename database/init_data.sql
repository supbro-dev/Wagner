-- 初始化工作点
INSERT INTO workplace (`code`,
                       `name`,
                       `region_code`,
                       `industry_code`,
                       `sub_industry_code`)
VALUES ('workplace1', '1号工作点', 'NORTH', 'FOOD', null);

-- 初始化员工信息
INSERT INTO `employee` (`id`,
                        `gmt_create`,
                        `gmt_modified`,
                        `name`,
                        `number`,
                        `identity`,
                        `sensitive_info`,
                        `workplace_code`,
                        `position_code`,
                        `work_group_code`,
                        `status`)
VALUES (1, now(), now(), '大兄弟', 'A1001', '1001', NULL, 'workplace1', 'pick', NULL, 'WORKING');
INSERT INTO `employee`
VALUES (2, now(), now(), '张伟', 'A1002', '1002', NULL, 'workplace1', 'pick', 'WG1', 'WORKING');
INSERT INTO `employee`
VALUES (3, now(), now(), '王芳', 'A1003', '1003', NULL, 'workplace1', 'pick', 'WG1', 'WORKING');
INSERT INTO `employee`
VALUES (4, now(), now(), '李强', 'A1004', '1004', NULL, 'workplace1', 'pick', 'WG1', 'WORKING');
INSERT INTO `employee`
VALUES (5, now(), now(), '刘洋', 'A1005', '1005', NULL, 'workplace1', 'packaging-operator', 'WG1', 'WORKING');
INSERT INTO `employee`
VALUES (6, now(), now(), '陈敏', 'A1006', '1006', NULL, 'workplace1', 'packaging-operator', 'WG2', 'WORKING');
INSERT INTO `employee`
VALUES (7, now(), now(), '黄海燕', 'A1007', '1007', NULL, 'workplace1', 'packaging-operator', 'WG2', 'WORKING');
INSERT INTO `employee`
VALUES (8, now(), now(), '吴婷婷', 'A1008', '1008', NULL, 'workplace1', 'batch-packer', 'WG2', 'WORKING');
INSERT INTO `employee`
VALUES (9, now(), now(), '赵静', 'A1009', '1009', NULL, 'workplace1', 'batch-packer', 'WG2', 'WORKING');
INSERT INTO `employee`
VALUES (10, now(), now(), '周杰', 'A1010', '1010', NULL, 'workplace1', 'batch-packer', 'WG2', 'WORKING');

-- 初始化组织架构
INSERT INTO `standard_position`
VALUES (1, '2025-06-05 11:14:31', '2025-06-14 17:12:15', 'outbound', '出库部', '-1', 'DEPT', '1', 1, 'FOOD', NULL,
        '{\"workLoadRollUp\": true}', NULL, 1);
INSERT INTO `standard_position`
VALUES (2, '2025-06-05 11:14:31', '2025-06-10 14:34:07', 'inbound', '入库部', '-1', 'DEPT', '1', 1, 'FOOD', NULL, NULL,
        NULL, 2);
INSERT INTO `standard_position`
VALUES (3, '2025-06-05 11:14:31', '2025-06-10 14:34:07', 'inventory', '库维部', '-1', 'DEPT', '1', 1, 'FOOD', NULL,
        NULL, NULL, 3)
;
INSERT INTO `standard_position`
VALUES (4, '2025-06-05 11:24:35', '2025-06-14 17:12:15', 'picking', '拣选部', 'outbound', 'DEPT', '2', 1, 'FOOD', NULL,
        '{\"workLoadRollUp\": true}', NULL, 1);
INSERT INTO `standard_position`
VALUES (5, '2025-06-05 11:24:35', '2025-06-10 14:34:07', 'packaging', '包装部', 'outbound', 'DEPT', '2', 1, 'FOOD',
        NULL, NULL, NULL, 2);
INSERT INTO `standard_position`
VALUES (6, '2025-06-05 11:24:35', '2025-06-10 14:34:07', 'shipping', '包裹出库部', 'outbound', 'DEPT', '2', 1, 'FOOD',
        NULL, NULL, NULL, 3);
INSERT INTO `standard_position`
VALUES (7, '2025-06-05 11:24:35', '2025-06-20 10:30:47', 'picker', '拣选员', 'picking', 'POSITION', '3', 1, 'FOOD',
        NULL, '{\"workLoadRollUp\": true}', NULL, 1);
INSERT INTO `standard_position`
VALUES (8, '2025-06-05 11:24:35', '2025-06-20 10:30:47', 'packaging-operator', '包装员', 'packaging', 'POSITION', '3',
        1, 'FOOD', NULL, NULL, NULL, 2);
INSERT INTO `standard_position`
VALUES (9, '2025-06-05 11:24:35', '2025-06-14 17:30:50', 'batch-picking', '批量拣选', 'picker', 'DIRECT_PROCESS', '4',
        1, 'FOOD', NULL, NULL, 'taskType == \"T001\"? true : false', 1);
INSERT INTO `standard_position`
VALUES (10, '2025-06-05 11:24:35', '2025-06-14 17:12:15', 'single-item-picking', '单品拣选', 'picker', 'DIRECT_PROCESS',
        '4', 1, 'FOOD', NULL, '{\"workLoadRollUp\": true}', 'taskType == \"T002\"? true : false', 2);
INSERT INTO `standard_position`
VALUES (11, '2025-06-05 11:24:35', '2025-06-10 14:34:07', 'picking-support', '拣选辅助', 'picker', 'INDIRECT_PROCESS',
        '4', 1, 'FOOD', NULL, NULL, 'indirectWorkType == \"B2\"? true : false', 3);
INSERT INTO `standard_position`
VALUES (12, '2025-06-20 10:30:47', '2025-06-20 10:30:47', 'batch-packer', '秒杀员', 'packaging', 'POSITION', '3', 1,
        'FOOD', NULL, '{\"workLoadRollUp\": true}', NULL, 3);
INSERT INTO `standard_position`
VALUES (13, '2025-06-20 10:30:47', '2025-06-20 10:30:47', 'batch-packaging', '批量包装', 'batch-packer',
        'INDIRECT_PROCESS', '4', 1, 'FOOD', NULL, NULL, 'indirectWorkType == \"B3\"? true : false', 1);

-- 初始化计算参数
INSERT INTO `calc_dynamic_param`
VALUES (1, '2025-06-03 14:07:55', '2025-06-19 21:59:23', 'SINK_STORAGE', 'FOOD', NULL,
        '[{\"sinkType\": \"SUMMARY\", \"fieldColumnList\": [{\"fieldName\": \"stationType\", \"columnName\": \"station_type\"}]}, {\"sinkType\": \"EMPLOYEE_STATUS\"}]'),
       (2, '2025-06-03 14:53:39', '2025-06-06 23:39:58', 'INJECT_SOURCE', 'FOOD', NULL,
        '[{\"fieldName\": \"stationType\"}, {\"fieldName\": \"indirectWorkType\"}, {\"fieldName\": \"taskType\"}]'),
       (4, '2025-06-03 14:58:38', '2025-06-12 20:57:14', 'DYNAMIC_CALC_NODE', 'FOOD', NULL,
        '{\"nodeNames\": \"SetCrossDayAttendance,ComputeAttendanceDefaultEndTime,ComputeAttendanceDefaultStartTime,CutOffAttendanceTime,AddCrossDayData,FilterOtherDaysData,FilterExpiredData,MarchProcess,PaddingUnfinishedWorkEndTime,CutOffOvertimeWork,CutOffCrossWork,AddReasonableBreakTime,CutOffWorkByRest,CalcWorkTransitionTime,MatchRestProcess,GenerateIdleData\"}'),
       (5, '2025-06-05 17:50:14', '2025-06-05 17:50:14', 'CALC_PARAM', 'FOOD', NULL,
        '{\"AttendanceParam\": {\"AttendanceAbsencePenaltyHour\": 8}}');