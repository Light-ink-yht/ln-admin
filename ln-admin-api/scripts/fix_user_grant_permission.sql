-- 修复用户授权权限配置
-- 注意：此脚本会更新现有数据库中的权限路径，确保与路由匹配
-- 注意：Casbin中的路径会移除 /api 前缀，所以这里使用 /api/user/*，在添加到Casbin时会变成 /user/*

-- 添加用户授权权限（如果不存在）
INSERT INTO AA06 (AAF001, AAF002, AAF003, AAF004, AAF005, AAF006, AAF007, AAA007, AAA008)
SELECT 'perm_user_grant', 'user:grant', '用户授权', '/api/user/*', 'POST', '给用户直接分配权限', '1', 'system', 'system'
WHERE NOT EXISTS (SELECT 1 FROM AA06 WHERE AAF001 = 'perm_user_grant');

-- 如果已存在但路径不对，更新路径
UPDATE AA06 
SET AAF004 = '/api/user/*' 
WHERE AAF001 = 'perm_user_grant' AND AAF004 != '/api/user/*';

-- 添加用户授权权限到 super_admin 角色（如果不存在）
-- 注意：在Casbin中存储时会移除 /api 前缀，所以实际路径是 /user/*
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', '/user/*', 'POST'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' AND v0 = 'super_admin' AND v1 = '/user/*' AND v2 = 'POST'
    AND (v3 IS NULL OR v3 = '')
);

