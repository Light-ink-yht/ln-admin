-- 初始化默认系统配置
-- 为现有数据库添加默认的系统配置数据

INSERT INTO AA02 (AAB001, AAB002, AAB003, AAB004, AAB005, AAB006, AAB008, AAB009, created_at, updated_at)
VALUES
    ('site_name', 'LN Admin', '网站名称', 'system', '网站/系统的名称，显示在页面标题、LOGO等位置', '1', 'system', 'system', NOW(), NOW()),
    ('site_logo', '', '网站LOGO', 'system', '网站LOGO图片URL地址', '1', 'system', 'system', NOW(), NOW()),
    ('site_favicon', '', '网站图标', 'system', '网站Favicon图标URL地址', '1', 'system', 'system', NOW(), NOW()),
    ('site_copyright', '© 2024 LN Admin 基于 Vue 3 + Ant Design Vue', '版权信息', 'system', '网站底部版权信息', '1', 'system', 'system', NOW(), NOW()),
    ('site_description', 'LN Admin 管理系统', '网站描述', 'system', '网站/系统的描述信息', '1', 'system', 'system', NOW(), NOW()),
    ('site_keywords', 'LN Admin,管理系统,后台管理', '网站关键词', 'system', '网站SEO关键词，多个关键词用逗号分隔', '1', 'system', 'system', NOW(), NOW()),
    ('site_beian', '', '备案号', 'system', '网站ICP备案号', '1', 'system', 'system', NOW(), NOW()),
    ('site_contact_email', '', '联系邮箱', 'system', '系统联系邮箱地址', '1', 'system', 'system', NOW(), NOW()),
    ('site_contact_phone', '', '联系电话', 'system', '系统联系电话', '1', 'system', 'system', NOW(), NOW()),
    ('site_address', '', '公司地址', 'system', '公司/组织地址', '1', 'system', 'system', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    AAB003 = VALUES(AAB003),
    AAB004 = VALUES(AAB004),
    AAB005 = VALUES(AAB005),
    AAB006 = VALUES(AAB006),
    AAB009 = VALUES(AAB009),
    updated_at = NOW();

-- 显示添加结果
SELECT 
    COUNT(*) AS total_configs,
    COUNT(CASE WHEN AAB004 = 'system' THEN 1 END) AS system_configs
FROM AA02
WHERE AAB001 LIKE 'site_%';

