-- 修复系统日志权限
-- 为现有数据库添加系统日志管理权限
-- 注意：Casbin中间件会移除 /api 前缀，所以Casbin中的路径应该是 /system/log/* 而不是 /api/system/log/*

-- 1. 插入系统日志权限到权限表 (AA06)
INSERT INTO AA06 (AAF001, AAF002, AAF003, AAF004, AAF005, AAF006, AAF007, AAF008, AAF009, created_at, updated_at)
VALUES
    ('perm_system_log_list', 'system:log:list', '系统日志列表', '/api/system/log/list', 'GET', '查看系统操作日志列表', '1', 'system', 'system', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    AAF003 = VALUES(AAF003),
    AAF004 = VALUES(AAF004),
    AAF005 = VALUES(AAF005),
    AAF006 = VALUES(AAF006),
    AAF007 = VALUES(AAF007),
    AAF008 = VALUES(AAF008),
    AAF009 = VALUES(AAF009),
    updated_at = NOW();

-- 2. 为超级管理员角色添加系统日志权限到Casbin
-- 注意：Casbin中的路径需要去掉 /api 前缀
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', REPLACE(AAF004, '/api', ''), AAF005
FROM AA06
WHERE AAF001 = 'perm_system_log_list'
  AND NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' 
      AND v0 = 'super_admin' 
      AND v1 = REPLACE(AA06.AAF004, '/api', '')
      AND v2 = AA06.AAF005
  );

-- 3. 显示添加结果
SELECT 
    '权限表' AS source,
    COUNT(*) AS count
FROM AA06
WHERE AAF001 LIKE 'perm_system_log%'
UNION ALL
SELECT 
    'Casbin策略表' AS source,
    COUNT(*) AS count
FROM casbin_rule
WHERE ptype = 'p' 
  AND v0 = 'super_admin'
  AND (v1 LIKE '/system/log%');

