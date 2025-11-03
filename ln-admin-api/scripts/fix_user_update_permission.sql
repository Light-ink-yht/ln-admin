-- 修复用户更新权限问题
-- 确保 super_admin 角色有 PUT /user/* 权限

-- 删除可能存在的错误配置
DELETE FROM casbin_rule 
WHERE ptype = 'p' 
  AND v0 = 'super_admin' 
  AND v1 IN ('/api/user', '/api/user/*', '/api/user/update')
  AND v2 = 'PUT';

-- 确保更新用户权限在数据库中正确配置
UPDATE AA06 
SET AAF004 = '/api/user/*', AAF005 = 'PUT'
WHERE AAF001 = 'perm_user_update';

-- 添加更新用户权限到 super_admin 角色（如果不存在）
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', '/user/*', 'PUT'
WHERE NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' AND v0 = 'super_admin' AND v1 = '/user/*' AND v2 = 'PUT'
);

-- 验证：查询 super_admin 的所有用户相关权限
SELECT ptype, v0, v1, v2 
FROM casbin_rule 
WHERE ptype = 'p' 
  AND v0 = 'super_admin' 
  AND v1 LIKE '/user%'
ORDER BY v1, v2;

