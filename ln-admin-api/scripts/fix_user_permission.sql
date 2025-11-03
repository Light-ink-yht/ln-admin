-- 为手机号 18797131041 的用户添加超级管理员权限
-- 执行方式：直接在数据库中运行此SQL脚本

-- 1. 查询用户ID（查看用户信息）
SELECT AAA001 as user_id, AAA003 as phone, AAA005 as nickname
FROM aa01 
WHERE AAA003 = '18797131041';

-- 2. 查询超级管理员角色
SELECT AAD001 as role_id, AAD002 as role_key, AAD003 as role_name
FROM aa05 
WHERE AAD002 = 'super_admin';

-- 3. 为超级管理员角色添加所有权限（如果还没有）
-- 注意：中间件会移除 /api 前缀，所以策略中也要移除 /api
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', 'super_admin', 
       CASE 
         WHEN AAF005 LIKE '/api/%' THEN SUBSTRING(AAF005, 5)  -- 移除 /api 前缀
         ELSE AAF005 
       END,
       AAF006
FROM aa06
WHERE AAF008 = '1'  -- 启用状态的权限
  AND NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'p' 
      AND v0 = 'super_admin' 
      AND v1 = CASE 
                 WHEN AAF005 LIKE '/api/%' THEN SUBSTRING(AAF005, 5)
                 ELSE AAF005 
               END
      AND v2 = AAF006
  );

-- 4. 为用户分配超级管理员角色（g策略）
-- 注意：需要先查询上面得到的 user_id，替换下面的 YOUR_USER_ID
-- 方法1：如果已知user_id
INSERT INTO casbin_rule (ptype, v0, v1)
VALUES ('g', 'YOUR_USER_ID', 'super_admin')
ON DUPLICATE KEY UPDATE v1 = 'super_admin';

-- 方法2：动态查询用户ID并分配角色（MySQL）
INSERT INTO casbin_rule (ptype, v0, v1)
SELECT 'g', AAA001, 'super_admin'
FROM aa01
WHERE AAA003 = '18797131041'
  AND NOT EXISTS (
    SELECT 1 FROM casbin_rule 
    WHERE ptype = 'g' 
      AND v0 = aa01.AAA001 
      AND v1 = 'super_admin'
  );

-- 5. 验证权限（查询结果）
-- 查看用户拥有的角色
SELECT v0 as user_id, v1 as role_key
FROM casbin_rule
WHERE ptype = 'g' 
  AND v0 IN (SELECT AAA001 FROM aa01 WHERE AAA003 = '18797131041');

-- 查看超级管理员角色的所有权限
SELECT v0 as role_key, v1 as resource_path, v2 as method
FROM casbin_rule
WHERE ptype = 'p' 
  AND v0 = 'super_admin';

