-- 修复系统配置权限
-- 为现有数据库添加系统配置管理权限
-- 注意：Casbin中间件会移除 /api 前缀，所以Casbin中的路径应该是 /system/config/* 而不是 /api/system/config/*

-- 1. 插入系统配置权限到权限表 (AA06)
INSERT INTO AA06 (AAF001, AAF002, AAF003, AAF004, AAF005, AAF006, AAF007, AAF008, AAF009, created_at, updated_at)
VALUES
    ('perm_system_config_list', 'system:config:list', '系统配置列表', '/api/system/config/list', 'GET', '查看系统配置列表', '1', 'system', 'system', NOW(), NOW()),
    ('perm_system_config_group', 'system:config:group', '系统配置分组', '/api/system/config/group/*', 'GET', '根据分组查看系统配置', '1', 'system', 'system', NOW(), NOW()),
    ('perm_system_config_detail', 'system:config:detail', '系统配置详情', '/api/system/config/*', 'GET', '查看系统配置详情', '1', 'system', 'system', NOW(), NOW()),
    ('perm_system_config_create', 'system:config:create', '创建系统配置', '/api/system/config', 'POST', '创建系统配置', '1', 'system', 'system', NOW(), NOW()),
    ('perm_system_config_update', 'system:config:update', '更新系统配置', '/api/system/config/*', 'PUT', '更新系统配置', '1', 'system', 'system', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    AAF003 = VALUES(AAF003),
    AAF004 = VALUES(AAF004),
    AAF005 = VALUES(AAF005),
    AAF006 = VALUES(AAF006),
    AAF007 = VALUES(AAF007),
    AAF008 = VALUES(AAF008),
    AAF009 = VALUES(AAF009),
    updated_at = NOW();

-- 2. 为超级管理员角色添加系统配置权限到Casbin
-- 注意：Casbin中的路径需要去掉 /api 前缀
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', 
    CASE 
        WHEN AAF004 = '/api/system/config/list' THEN '/system/config/list'
        WHEN AAF004 = '/api/system/config/group/*' THEN '/system/config/group/*'
        WHEN AAF004 = '/api/system/config/*' THEN '/system/config/*'
        WHEN AAF004 = '/api/system/config' THEN '/system/config'
        ELSE REPLACE(AAF004, '/api', '')
    END,
    AAF005
FROM AA06
WHERE AAF001 IN ('perm_system_config_list', 'perm_system_config_group', 'perm_system_config_detail', 'perm_system_config_create', 'perm_system_config_update')
  AND NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' 
      AND v0 = 'super_admin' 
      AND v1 = CASE 
          WHEN AA06.AAF004 = '/api/system/config/list' THEN '/system/config/list'
          WHEN AA06.AAF004 = '/api/system/config/group/*' THEN '/system/config/group/*'
          WHEN AA06.AAF004 = '/api/system/config/*' THEN '/system/config/*'
          WHEN AA06.AAF004 = '/api/system/config' THEN '/system/config'
          ELSE REPLACE(AA06.AAF004, '/api', '')
      END
      AND v2 = AA06.AAF005
  );

-- 3. 显示添加结果
SELECT 
    '权限表' AS source,
    COUNT(*) AS count
FROM AA06
WHERE AAF001 LIKE 'perm_system_config%'
UNION ALL
SELECT 
    'Casbin策略表' AS source,
    COUNT(*) AS count
FROM casbin_rule
WHERE ptype = 'p' 
  AND v0 = 'super_admin'
  AND (v1 LIKE '/system/config%');

