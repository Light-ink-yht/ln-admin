-- MySQL 数据库迁移脚本 V4
-- 新增：短信模板表

-- 短信模板表 AA08
CREATE TABLE IF NOT EXISTS `AA08` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `AAA081` varchar(50) NOT NULL COMMENT '模板ID',
  `AAA082` varchar(50) NOT NULL COMMENT '模板类型（register, forgot等）',
  `AAA083` varchar(100) DEFAULT NULL COMMENT '模板名称',
  `AAA084` varchar(200) DEFAULT NULL COMMENT '腾讯云模板ID',
  `AAA085` varchar(100) DEFAULT NULL COMMENT '短信签名/标题',
  `AAA086` text COMMENT '模板内容（支持占位符，如 {code}）',
  `AAA087` varchar(500) DEFAULT NULL COMMENT '模板描述',
  `AAA088` varchar(10) DEFAULT '1' COMMENT '状态 1 启用，2 禁用',
  `AAA089` varchar(50) DEFAULT NULL COMMENT '创建人',
  `AAA090` varchar(50) DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_aaa081` (`AAA081`),
  UNIQUE KEY `uk_aaa082` (`AAA082`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信模板表';

-- 初始化默认短信模板（示例）
INSERT INTO `AA08` (`AAA081`, `AAA082`, `AAA083`, `AAA084`, `AAA085`, `AAA086`, `AAA087`, `AAA088`, `AAA089`) VALUES
('template_register_001', 'register', '用户注册模板', 'your_register_template_id', '您的应用名称', '您的注册验证码是：{code}，10分钟内有效，请勿泄露给他人。', '用户注册场景使用的短信模板', '1', 'system'),
('template_forgot_001', 'forgot', '忘记密码模板', 'your_forgot_template_id', '您的应用名称', '您的密码重置验证码是：{code}，10分钟内有效，请勿泄露给他人。', '忘记密码场景使用的短信模板', '1', 'system')
ON DUPLICATE KEY UPDATE `AAA082`=`AAA082`;

