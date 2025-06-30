-- 大兄弟当天有直接/间接作业，排班，考勤无下班打卡
INSERT INTO `action` (`employee_number`,
                      `workplace_code`,
                      `start_time`,
                      `end_time`,
                      `action_type`,
                      `properties`,
                      `operate_day`,
                      `action_code`,
                      `work_load`)
VALUES ('A1001', 'workplace1', CONCAT(CURDATE(), ' 10:00:00'),
        CONCAT(CURDATE(), ' 12:00:00'), 'DirectWork', '{\"taskType\": \"T001\", \"stationType\": \"A1\"}', curDate(),
        'TB0001', '{\"itemNum\": 20}'),
       ( 'A1001', 'workplace1', CONCAT(CURDATE(), ' 13:00:00'),
         CONCAT(CURDATE(), ' 15:00:00'), 'DirectWork', '{\"taskType\": \"T002\"}', curDate(), 'TB0002',
        '{\"itemNum\": 50}'),
       ( 'A1001', 'workplace1', CONCAT(CURDATE(), ' 15:30:00'),
         CONCAT(CURDATE(), ' 18:00:00'), 'IndirectWork', '{\"indirectWorkType\": \"B2\"}', curDate(), 'IL10001', NULL),
       ( 'A1001', 'workplace1', CONCAT(CURDATE(), ' 08:20:00'), NULL,
        'Attendance', NULL, curDate(), 'a1', NULL),
       ('A1001', 'workplace1', CONCAT(CURDATE(), ' 08:00:00'),
        CONCAT(CURDATE(), ' 18:00:00'), 'Scheduling',
        concat('{\"restList\": [{\"endTime\": \"',CONCAT(CURDATE(), ' 12:10:00'), '\", \"startTime\": \"',CONCAT(CURDATE(), ' 11:30:00'), '\"}]}'),
        curDate(), 'r1', NULL);


-- 张伟这个员工是正常考勤上下班，且全天有大量的直接/间接作业
-- 排班记录 (包含两个休息时段)
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`, `work_load`)
VALUES ('A1002', 'workplace1',
        CONCAT(CURDATE(), ' 08:00:00'),
        CONCAT(CURDATE(), ' 18:00:00'),
        'Scheduling',
        CONCAT('{"restList": [',
               '{"startTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), ' 11:30:00', '", "endTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), ' 12:30:00', '"}',
               ',',
               '{"startTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), ' 15:00:00', '", "endTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), ' 15:15:00', '"}',
               ']}'),
        CURDATE(),
        'SCH002',
        NULL);

-- 考勤记录
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`, `work_load`)
VALUES ('A1002', 'workplace1',
        CONCAT(CURDATE(), ' 08:05:00'),
        CONCAT(CURDATE(), ' 18:10:00'),
        'Attendance',
        NULL,
        CURDATE(),
        'ATT002',
        NULL);

-- 直接作业 (12条，T001/T002连续分组)
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`, `work_load`)
VALUES
-- 上午 T001连续组 (3条)
('A1002', 'workplace1', CONCAT(CURDATE(), ' 08:30:00'), CONCAT(CURDATE(), ' 09:20:00'), 'DirectWork', '{"taskType": "T001", "stationType": "S1"}', CURDATE(), 'DW2001', '{"itemNum": 25}'),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 09:20:00'), CONCAT(CURDATE(), ' 10:10:00'), 'DirectWork', '{"taskType": "T001", "stationType": "S2"}', CURDATE(), 'DW2002', '{"itemNum": 30}'),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 10:10:00'), CONCAT(CURDATE(), ' 11:05:00'), 'DirectWork', '{"taskType": "T001", "stationType": "S3"}', CURDATE(), 'DW2003', '{"itemNum": 35}'),

-- 上午 T002连续组 (2条)
('A1002', 'workplace1', CONCAT(CURDATE(), ' 11:05:00'), CONCAT(CURDATE(), ' 11:30:00'), 'DirectWork', '{"taskType": "T002", "stationType": "S4"}', CURDATE(), 'DW2004', '{"itemNum": 20}'),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 12:30:00'), CONCAT(CURDATE(), ' 13:25:00'), 'DirectWork', '{"taskType": "T002", "stationType": "S5"}', CURDATE(), 'DW2005', '{"itemNum": 40}'),

-- 下午 T001连续组 (3条)
('A1002', 'workplace1', CONCAT(CURDATE(), ' 13:25:00'), CONCAT(CURDATE(), ' 14:15:00'), 'DirectWork', '{"taskType": "T001", "stationType": "S1"}', CURDATE(), 'DW2006', '{"itemNum": 30}'),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 14:15:00'), CONCAT(CURDATE(), ' 15:00:00'), 'DirectWork', '{"taskType": "T001", "stationType": "S2"}', CURDATE(), 'DW2007', '{"itemNum": 25}'),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 15:15:00'), CONCAT(CURDATE(), ' 16:05:00'), 'DirectWork', '{"taskType": "T001", "stationType": "S3"}', CURDATE(), 'DW2008', '{"itemNum": 35}'),

-- 下午 T002连续组 (4条)
('A1002', 'workplace1', CONCAT(CURDATE(), ' 16:05:00'), CONCAT(CURDATE(), ' 16:45:00'), 'DirectWork', '{"taskType": "T002", "stationType": "S4"}', CURDATE(), 'DW2009', '{"itemNum": 28}'),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 16:45:00'), CONCAT(CURDATE(), ' 17:20:00'), 'DirectWork', '{"taskType": "T002", "stationType": "S5"}', CURDATE(), 'DW2010', '{"itemNum": 32}'),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 17:20:00'), CONCAT(CURDATE(), ' 17:50:00'), 'DirectWork', '{"taskType": "T002", "stationType": "S1"}', CURDATE(), 'DW2011', '{"itemNum": 30}'),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 17:50:00'), CONCAT(CURDATE(), ' 18:10:00'), 'DirectWork', '{"taskType": "T002", "stationType": "S2"}', CURDATE(), 'DW2012', '{"itemNum": 18}');

-- 间接作业 (8条，B2/B3连续分组)
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`, `work_load`)
VALUES
-- 上午 B2连续组 (2条)
('A1002', 'workplace1', CONCAT(CURDATE(), ' 08:15:00'), CONCAT(CURDATE(), ' 08:30:00'), 'IndirectWork', '{"indirectWorkType": "B2"}', CURDATE(), 'IW2001', NULL),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 09:05:00'), CONCAT(CURDATE(), ' 09:20:00'), 'IndirectWork', '{"indirectWorkType": "B2"}', CURDATE(), 'IW2002', NULL),

-- 上午 B3连续组 (2条)
('A1002', 'workplace1', CONCAT(CURDATE(), ' 10:10:00'), CONCAT(CURDATE(), ' 10:25:00'), 'IndirectWork', '{"indirectWorkType": "B3"}', CURDATE(), 'IW2003', NULL),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 11:30:00'), CONCAT(CURDATE(), ' 12:30:00'), 'IndirectWork', '{"indirectWorkType": "B3"}', CURDATE(), 'IW2004', NULL),

-- 下午 B2连续组 (2条)
('A1002', 'workplace1', CONCAT(CURDATE(), ' 13:20:00'), CONCAT(CURDATE(), ' 13:25:00'), 'IndirectWork', '{"indirectWorkType": "B2"}', CURDATE(), 'IW2005', NULL),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 14:10:00'), CONCAT(CURDATE(), ' 14:15:00'), 'IndirectWork', '{"indirectWorkType": "B2"}', CURDATE(), 'IW2006', NULL),

-- 下午 B3连续组 (2条)
('A1002', 'workplace1', CONCAT(CURDATE(), ' 15:00:00'), CONCAT(CURDATE(), ' 15:15:00'), 'IndirectWork', '{"indirectWorkType": "B3"}', CURDATE(), 'IW2007', NULL),
('A1002', 'workplace1', CONCAT(CURDATE(), ' 16:45:00'), CONCAT(CURDATE(), ' 17:00:00'), 'IndirectWork', '{"indirectWorkType": "B3"}', CURDATE(), 'IW2008', NULL);



-- 王芳正常作业数据
-- 排班数据（全天排班，含两个休息时段）
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`)
VALUES (
           'A1003',
           'workplace1',
           CONCAT(CURDATE(), ' 08:00:00.000'),
           CONCAT(CURDATE(), ' 18:00:00.000'),
           'Scheduling',
           concat('{"restList": [{"startTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), ' 11:30:00", "endTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'),' 12:00:00"}, {"startTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), ' 15:00:00", "endTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'),' 15:15:00"}]}'),
           DATE_FORMAT(CURDATE(), '%Y-%m-%d'),
           'SCH003'
       );

-- 考勤数据（当天实际出勤）
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `operate_day`, `action_code`)
VALUES (
           'A1003',
           'workplace1',
           CONCAT(CURDATE(), ' 08:05:00.000'),
           CONCAT(CURDATE(), ' 18:05:00.000'),
           'Attendance',
           DATE_FORMAT(CURDATE(), '%Y-%m-%d'),
           'ATT003'
       );

-- 直接作业数据（6个连续任务）
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`, `work_load`)
VALUES
    ('A1003', 'workplace1', CONCAT(CURDATE(), ' 08:30:00.000'), CONCAT(CURDATE(), ' 10:00:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW3001', '{"itemNum": 30, "skuNum": 15, "packageNum": 3}'),
    ('A1003', 'workplace1', CONCAT(CURDATE(), ' 10:00:00.000'), CONCAT(CURDATE(), ' 11:30:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW3002', '{"itemNum": 45, "skuNum": 22, "packageNum": 5}'),
    ('A1003', 'workplace1', CONCAT(CURDATE(), ' 12:00:00.000'), CONCAT(CURDATE(), ' 13:30:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW3003', '{"itemNum": 50, "skuNum": 25, "packageNum": 6}'),
    ('A1003', 'workplace1', CONCAT(CURDATE(), ' 13:30:00.000'), CONCAT(CURDATE(), ' 15:00:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW3004', '{"itemNum": 55, "skuNum": 28, "packageNum": 7}'),
    ('A1003', 'workplace1', CONCAT(CURDATE(), ' 15:15:00.000'), CONCAT(CURDATE(), ' 16:30:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW3005', '{"itemNum": 40, "skuNum": 20, "packageNum": 4}'),
    ('A1003', 'workplace1', CONCAT(CURDATE(), ' 16:30:00.000'), CONCAT(CURDATE(), ' 17:50:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW3006', '{"itemNum": 60, "skuNum": 30, "packageNum": 8}');

-- 间接作业数据（集中在直接作业前后）
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`)
VALUES
    ('A1003', 'workplace1', CONCAT(CURDATE(), ' 08:05:00.000'), CONCAT(CURDATE(), ' 08:30:00.000'), 'IndirectWork', '{"indirectWorkType": "B2"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'IW3001'),
    ('A1003', 'workplace1', CONCAT(CURDATE(), ' 17:50:00.000'), CONCAT(CURDATE(), ' 18:05:00.000'), 'IndirectWork', '{"indirectWorkType": "B3"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'IW3002');


-- 李强中午午休没有休息，仍然持续作业的情况
-- 排班数据（含午休时段）
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`)
VALUES (
           'A1004',
           'workplace1',
           CONCAT(CURDATE(), ' 08:00:00.000'),
           CONCAT(CURDATE(), ' 18:00:00.000'),
           'Scheduling',
           CONCAT('{"restList": [{"startTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), ' 12:00:00", "endTime": "', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), ' 13:00:00"}]}'),
           DATE_FORMAT(CURDATE(), '%Y-%m-%d'),
           'SCH004'
       );

-- 考勤数据（下班未打卡）
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `operate_day`, `action_code`)
VALUES (
           'A1004',
           'workplace1',
           CONCAT(CURDATE(), ' 08:02:00.000'),
           NULL,  -- 下班未打卡
           'Attendance',
           DATE_FORMAT(CURDATE(), '%Y-%m-%d'),
           'ATT004'
       );

-- 直接作业数据（短任务，午休时段有工作）
INSERT INTO `action` (`employee_number`, `workplace_code`, `start_time`, `end_time`, `action_type`, `properties`, `operate_day`, `action_code`, `work_load`)
VALUES
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 08:15:00.000'), CONCAT(CURDATE(), ' 08:28:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4001', '{"itemNum": 8, "skuNum": 4, "packageNum": 1}'),
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 08:30:00.000'), CONCAT(CURDATE(), ' 08:42:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4002', '{"itemNum": 10, "skuNum": 5, "packageNum": 2}'),
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 08:44:00.000'), CONCAT(CURDATE(), ' 08:57:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4003', '{"itemNum": 9, "skuNum": 4, "packageNum": 1}'),
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 09:00:00.000'), CONCAT(CURDATE(), ' 09:15:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4004', '{"itemNum": 12, "skuNum": 6, "packageNum": 2}'),
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 09:17:00.000'), CONCAT(CURDATE(), ' 09:30:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4005', '{"itemNum": 11, "skuNum": 5, "packageNum": 1}'),
    -- 午休时段工作（12:00-12:20）
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 12:00:00.000'), CONCAT(CURDATE(), ' 12:20:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4006', '{"itemNum": 15, "skuNum": 7, "packageNum": 2}'),
    -- 下午工作
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 13:20:00.000'), CONCAT(CURDATE(), ' 13:35:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4007', '{"itemNum": 14, "skuNum": 7, "packageNum": 2}'),
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 13:37:00.000'), CONCAT(CURDATE(), ' 13:50:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4008', '{"itemNum": 13, "skuNum": 6, "packageNum": 2}'),
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 13:52:00.000'), CONCAT(CURDATE(), ' 14:05:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4009', '{"itemNum": 11, "skuNum": 5, "packageNum": 1}'),
    -- 15:00-17:00 空闲（无记录）
    -- 傍晚工作
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 17:05:00.000'), CONCAT(CURDATE(), ' 17:18:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4010', '{"itemNum": 16, "skuNum": 8, "packageNum": 2}'),
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 17:20:00.000'), CONCAT(CURDATE(), ' 17:35:00.000'), 'DirectWork', '{"taskType": "T001"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4011', '{"itemNum": 14, "skuNum": 7, "packageNum": 2}'),
    ('A1004', 'workplace1', CONCAT(CURDATE(), ' 17:37:00.000'), CONCAT(CURDATE(), ' 17:50:00.000'), 'DirectWork', '{"taskType": "T002"}', DATE_FORMAT(CURDATE(), '%Y-%m-%d'), 'DW4012', '{"itemNum": 18, "skuNum": 9, "packageNum": 3}');

