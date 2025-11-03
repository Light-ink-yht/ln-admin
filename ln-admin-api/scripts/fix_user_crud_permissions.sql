-- 修复用户CRUD权限路径，使其与实际路由匹配
-- 注意：Casbin中间件会移除 /api 前缀，所以Casbin中的路径应该是 /user/* 而不是 /api/user/*

-- 更新用户详情权限路径（GET /api/user/:userId）
UPDATE AA06 
SET AAF004 = '/api/user/*' 
WHERE AAF001 = 'perm_user_detail' AND AAF004 != '/api/user/*';

-- 如果用户详情权限不存在，则创建它
INSERT INTO AA06 (AAF001, AAF002, AAF003, AAF004, AAF005, AAF006, AAF007, AAA007, AAA008)
SELECT 'perm_user_detail', 'user:detail', '用户详情', '/api/user/*', 'GET', '查看用户详情', '1', 'system', 'system'
WHERE NOT EXISTS (SELECT 1 FROM AA06 WHERE AAF001 = 'perm_user_detail');

-- 更新创建用户权限路径（POST /api/user）
UPDATE AA06 
SET AAF004 = '/api/user' 
WHERE AAF001 = 'perm_user_create' AND AAF004 != '/api/user';

-- 更新更新用户权限路径（PUT /api/user/:userId）
UPDATE AA06 
SET AAF004 = '/api/user/*' 
WHERE AAF001 = 'perm_user_update' AND AAF004 != '/api/user/*';

-- 更新删除用户权限路径（DELETE /api/user/:userId）
UPDATE AA06 
SET AAF004 = '/api/user/*' 
WHERE AAF001 = 'perm_user_delete' AND AAF004 != '/api/user/*';

-- 删除旧的Casbin规则（如果有错误的路径）
DELETE FROM casbin_rule 
WHERE ptype = 'p' 
  AND v0 = 'super_admin' 
  AND v1 IN ('/api/user/create', '/api/user/update', '/api/user/delete', '/api/user/*')
  AND v2 IN ('GET', 'POST', 'PUT', 'DELETE');

-- 添加用户详情权限到 super_admin 角色（如果不存在）- 注意：路径是 /user/*，不是 /api/user/*
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', '/user/*', 'GET'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' AND v0 = 'super_admin' AND v1 = '/user/*' AND v2 = 'GET'
);

-- 添加创建用户权限到 super_admin 角色（如果不存在）
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', '/user', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' AND v0 = 'super_admin' AND v1 = '/user' AND v2 = 'POST'
);

-- 添加更新用户权限到 super_admin 角色（如果不存在）
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', '/user/*', 'PUT'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' AND v0 = 'super_admin' AND v1 = '/user/*' AND v2 = 'PUT'
);

-- 添加删除用户权限到 super_admin 角色（如果不存在）
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', '/user/*', 'DELETE'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' AND v0 = 'super_admin' AND v1 = '/user/*' AND v2 = 'DELETE'
);

-- 注意：创建用户和授权用户都使用 POST /user，但由于路由优先级
-- POST /user 是创建用户，POST /user/:userId/permissions 是授权
-- 由于中间件会移除 /api 前缀，/api/user/:userId/permissions 变成 /user/:userId/permissions
-- 通配符 /user/* 可以匹配 /user/:userId/permissions，所以授权权限使用 /user/* POST 应该能工作
-- 但是创建用户需要精确匹配 /user POST，所以我们需要确保两个规则都存在

-- 检查是否已有 /user POST（创建用户）
-- 如果授权权限需要更精确的路径，可能需要单独处理
-- 但由于 keyMatch 支持通配符，/user/* 应该能匹配 /user/:userId/permissions

