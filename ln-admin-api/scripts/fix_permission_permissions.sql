-- 为权限管理添加权限配置，并分配给超级管理员角色
-- 执行方式：直接在数据库中运行此SQL脚本

-- 1. 确保权限表 (AA06) 中存在权限管理相关权限
INSERT IGNORE INTO AA06 (AAF001, AAF002, AAF003, AAF004, AAF005, AAF006, AAF007, AAF008, created_at, updated_at) VALUES
('perm_permission_list', 'permission:list', '权限列表', '/api/permission/list', 'GET', '查看权限列表', '1', 'system', NOW(), NOW()),
('perm_permission_detail', 'permission:detail', '权限详情', '/api/permission/*', 'GET', '查看权限详情', '1', 'system', NOW(), NOW()),
('perm_permission_create', 'permission:create', '创建权限', '/api/permission', 'POST', '创建新权限', '1', 'system', NOW(), NOW()),
('perm_permission_update', 'permission:update', '更新权限', '/api/permission/*', 'PUT', '更新权限信息', '1', 'system', NOW(), NOW()),
('perm_permission_delete', 'permission:delete', '删除权限', '/api/permission/*', 'DELETE', '删除权限', '1', 'system', NOW(), NOW());

-- 2. 获取超级管理员的角色标识
SET @super_admin_role_key = (SELECT AAE002 FROM AA05 WHERE AAE003 = '超级管理员' LIMIT 1);

-- 3. 为超级管理员角色分配所有权限管理权限到 Casbin
-- 注意：Casbin 策略中的资源路径需要与中间件处理后的路径一致，即移除 /api 前缀
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', @super_admin_role_key, '/permission/list', 'GET'),
('p', @super_admin_role_key, '/permission/*', 'GET'),
('p', @super_admin_role_key, '/permission', 'POST'),
('p', @super_admin_role_key, '/permission/*', 'PUT'),
('p', @super_admin_role_key, '/permission/*', 'DELETE');

-- 4. 确保默认用户 (18797131041) 拥有超级管理员角色
SET @default_user_id = (SELECT AAD001 FROM AA04 WHERE AAD002 = '18797131041' LIMIT 1);
SET @super_admin_role_id = (SELECT AAE001 FROM AA05 WHERE AAE002 = 'super_admin' LIMIT 1);

INSERT IGNORE INTO AA07 (AAG001, AAG002, created_at, updated_at) VALUES
(@default_user_id, @super_admin_role_id, NOW(), NOW());

-- 5. 确保 Casbin 中用户与角色的关联
INSERT IGNORE INTO casbin_rule (ptype, v0, v1) VALUES
('g', @default_user_id, @super_admin_role_key);

-- 6. 验证权限（查询结果）
-- 查看权限管理相关权限
SELECT AAF001 as permission_id, AAF002 as permission_key, AAF003 as permission_name, AAF004 as resource_path, AAF005 as method
FROM AA06
WHERE AAF002 LIKE 'permission:%'
ORDER BY AAF002;

-- 查看超级管理员角色的权限管理权限
SELECT v0 as role_key, v1 as resource_path, v2 as method
FROM casbin_rule
WHERE ptype = 'p' 
  AND v0 = 'super_admin'
  AND v1 LIKE '/permission%';

-- 查看用户拥有的角色
SELECT v0 as user_id, v1 as role_key
FROM casbin_rule
WHERE ptype = 'g' 
  AND v0 IN (SELECT AAD001 FROM AA04 WHERE AAD002 = '18797131041');

