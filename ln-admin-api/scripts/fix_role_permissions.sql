-- 修复角色管理权限脚本
-- 用途：为超级管理员角色添加所有角色管理相关的权限
-- 执行方式：mysql -u root -p ln_admin < scripts/fix_role_permissions.sql

-- 1. 添加角色管理权限到权限表（如果不存在）
INSERT IGNORE INTO AA04 (AAF001, AAF002, AAF003, AAF004, AAF005, AAF006, AAF007, status, created_at, updated_at) VALUES
('perm_role_list', 'role:list', '角色列表', '/api/role/list', 'GET', 'system', 'system', '1', NOW(), NOW()),
('perm_role_detail', 'role:detail', '角色详情', '/api/role/*', 'GET', 'system', 'system', '1', NOW(), NOW()),
('perm_role_create', 'role:create', '创建角色', '/api/role', 'POST', 'system', 'system', '1', NOW(), NOW()),
('perm_role_update', 'role:update', '更新角色', '/api/role/*', 'PUT', 'system', 'system', '1', NOW(), NOW()),
('perm_role_delete', 'role:delete', '删除角色', '/api/role/*', 'DELETE', 'system', 'system', '1', NOW(), NOW());

-- 2. 获取超级管理员角色Key
SET @super_admin_role_key = (SELECT AAE002 FROM AA05 WHERE AAE002 = 'super_admin' LIMIT 1);

-- 3. 为超级管理员角色添加所有角色管理权限到Casbin（移除/api前缀，因为中间件会移除）
-- 注意：Casbin策略格式为 (ptype, v0, v1, v2) = ('p', role_key, resource_path, method)
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', @super_admin_role_key, '/role/list', 'GET'),
('p', @super_admin_role_key, '/role/*', 'GET'),
('p', @super_admin_role_key, '/role', 'POST'),
('p', @super_admin_role_key, '/role/*', 'PUT'),
('p', @super_admin_role_key, '/role/*', 'DELETE');

-- 4. 确保用户 18797131041 拥有 super_admin 角色
SET @user_id = (SELECT AAE001 FROM AA01 WHERE AAE003 = '18797131041' LIMIT 1);
SET @super_admin_role_id = (SELECT AAE001 FROM AA05 WHERE AAE002 = 'super_admin' LIMIT 1);

-- 4.1 在用户角色关联表中添加角色（如果不存在）
INSERT IGNORE INTO AA06 (AAE001, AAF001) VALUES (@user_id, @super_admin_role_id);

-- 4.2 在Casbin中添加用户角色关系（如果不存在）
INSERT IGNORE INTO casbin_rule (ptype, v0, v1) VALUES
('g', @user_id, @super_admin_role_key);

-- 5. 显示结果
SELECT '修复完成！' AS message;
SELECT CONCAT('用户ID: ', @user_id) AS user_info;
SELECT CONCAT('超级管理员角色Key: ', @super_admin_role_key) AS role_info;
SELECT '请重启应用服务器以使权限生效。' AS note;

