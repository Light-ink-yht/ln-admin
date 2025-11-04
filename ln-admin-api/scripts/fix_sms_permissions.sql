-- 修复短信管理权限
-- 为现有数据库添加短信管理权限
-- 注意：Casbin中间件会移除 /api 前缀，所以Casbin中的路径应该是 /sms/* 而不是 /api/sms/*

-- 1. 插入短信管理权限到权限表 (AA06)
INSERT INTO AA06 (AAF001, AAF002, AAF003, AAF004, AAF005, AAF006, AAF007, AAF008, AAF009, created_at, updated_at)
VALUES
    ('perm_sms_template_list', 'sms:template:list', '短信模板列表', '/api/sms/template/list', 'GET', '查看短信模板列表', '1', 'system', 'system', NOW(), NOW()),
    ('perm_sms_template_detail', 'sms:template:detail', '短信模板详情', '/api/sms/template/*', 'GET', '查看短信模板详情', '1', 'system', 'system', NOW(), NOW()),
    ('perm_sms_template_create', 'sms:template:create', '创建短信模板', '/api/sms/template', 'POST', '创建短信模板', '1', 'system', 'system', NOW(), NOW()),
    ('perm_sms_template_update', 'sms:template:update', '更新短信模板', '/api/sms/template/*', 'PUT', '更新短信模板', '1', 'system', 'system', NOW(), NOW()),
    ('perm_sms_template_delete', 'sms:template:delete', '删除短信模板', '/api/sms/template/*', 'DELETE', '删除短信模板', '1', 'system', 'system', NOW(), NOW()),
    ('perm_sms_code_list', 'sms:code:list', '短信验证码列表', '/api/sms/code/list', 'GET', '查看短信验证码列表', '1', 'system', 'system', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    AAF003 = VALUES(AAF003),
    AAF004 = VALUES(AAF004),
    AAF005 = VALUES(AAF005),
    AAF006 = VALUES(AAF006),
    AAF007 = VALUES(AAF007),
    AAF008 = VALUES(AAF008),
    AAF009 = VALUES(AAF009),
    updated_at = NOW();

-- 2. 为超级管理员角色添加短信管理权限到Casbin
-- 注意：Casbin中的路径需要去掉 /api 前缀
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', REPLACE(AAF004, '/api', ''), AAF005
FROM AA06
WHERE AAF001 LIKE 'perm_sms%'
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
WHERE AAF001 LIKE 'perm_sms%'
UNION ALL
SELECT 
    'Casbin策略表' AS source,
    COUNT(*) AS count
FROM casbin_rule
WHERE ptype = 'p' 
  AND v0 = 'super_admin'
  AND (v1 LIKE '/sms%');

