-- MySQL 数据库迁移脚本 V2
-- 新增：系统配置表、日志表、短信验证码表

-- 系统配置表 AA02
CREATE TABLE IF NOT EXISTS `AA02` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `AAA021` varchar(50) NOT NULL COMMENT '配置键',
  `AAA022` varchar(500) DEFAULT NULL COMMENT '配置值',
  `AAA023` varchar(100) DEFAULT NULL COMMENT '配置名称',
  `AAA024` varchar(50) DEFAULT NULL COMMENT '配置分组',
  `AAA025` varchar(500) DEFAULT NULL COMMENT '配置描述',
  `AAA026` varchar(10) DEFAULT '1' COMMENT '状态 1 启用，2 禁用',
  `AAA027` datetime DEFAULT NULL COMMENT '最后修改时间',
  `AAA028` varchar(50) DEFAULT NULL COMMENT '创建人',
  `AAA029` varchar(50) DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_aaa021` (`AAA021`),
  KEY `idx_aaa024` (`AAA024`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 系统日志表 AA03
CREATE TABLE IF NOT EXISTS `AA03` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `AAA031` varchar(50) DEFAULT NULL COMMENT '日志ID',
  `AAA032` varchar(20) DEFAULT NULL COMMENT '日志级别',
  `AAA033` varchar(200) DEFAULT NULL COMMENT '模块名称',
  `AAA034` varchar(100) DEFAULT NULL COMMENT '操作类型',
  `AAA035` text COMMENT '日志内容',
  `AAA036` varchar(50) DEFAULT NULL COMMENT '用户ID',
  `AAA037` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `AAA038` varchar(500) DEFAULT NULL COMMENT '请求路径',
  `AAA039` varchar(20) DEFAULT NULL COMMENT '请求方法',
  `AAA040` int DEFAULT NULL COMMENT '响应状态码',
  `AAA041` varchar(500) DEFAULT NULL COMMENT '用户代理',
  `AAA042` varchar(200) DEFAULT NULL COMMENT '错误信息',
  `AAA043` datetime DEFAULT NULL COMMENT '日志时间',
  `AAA044` varchar(50) DEFAULT NULL COMMENT '创建人',
  PRIMARY KEY (`id`),
  KEY `idx_aaa032` (`AAA032`),
  KEY `idx_aaa043` (`AAA043`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统日志表';

-- 短信验证码表 AA04
CREATE TABLE IF NOT EXISTS `AA04` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `AAA041` varchar(50) NOT NULL COMMENT '验证码ID',
  `AAA042` varchar(20) DEFAULT NULL COMMENT '手机号',
  `AAA043` varchar(10) DEFAULT NULL COMMENT '验证码',
  `AAA044` varchar(20) DEFAULT NULL COMMENT '验证码类型',
  `AAA045` varchar(10) DEFAULT '1' COMMENT '状态 1 未使用，2 已使用，3 已过期',
  `AAA046` datetime DEFAULT NULL COMMENT '过期时间',
  `AAA047` datetime DEFAULT NULL COMMENT '使用时间',
  `AAA048` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `AAA049` int DEFAULT '0' COMMENT '发送次数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_aaa041` (`AAA041`),
  KEY `idx_aaa042` (`AAA042`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='短信验证码表';

-- 初始化短信配置数据（示例）
INSERT INTO `AA02` (`AAA021`, `AAA022`, `AAA023`, `AAA024`, `AAA025`, `AAA026`, `AAA028`) VALUES
('sms_secret_id', 'your_secret_id', '腾讯云SecretId', 'sms', '腾讯云短信服务SecretId', '1', 'system'),
('sms_secret_key', 'your_secret_key', '腾讯云SecretKey', 'sms', '腾讯云短信服务SecretKey', '1', 'system'),
('sms_app_id', 'your_app_id', '腾讯云AppId', 'sms', '腾讯云短信应用ID', '1', 'system'),
('sms_sign_name', 'your_sign_name', '默认短信签名', 'sms', '默认短信签名（如果模板中没有指定）', '1', 'system'),
('sms_code_expire_minutes', '10', '验证码过期时间', 'sms', '短信验证码过期时间（分钟），默认10分钟', '1', 'system'),
-- 注册短信模板（JSON格式，包含模板ID、标题、内容等）
('sms_template_register', '{"template_id":"your_register_template_id","title":"您的应用名称","content":"您的注册验证码是：{code}，10分钟内有效，请勿泄露给他人。","type":"register","description":"用户注册场景使用的短信模板","status":"1"}', '注册短信模板', 'sms', '用户注册场景使用的短信模板（JSON格式）', '1', 'system'),
-- 忘记密码短信模板（JSON格式，包含模板ID、标题、内容等）
('sms_template_forgot', '{"template_id":"your_forgot_template_id","title":"您的应用名称","content":"您的密码重置验证码是：{code}，10分钟内有效，请勿泄露给他人。","type":"forgot","description":"忘记密码场景使用的短信模板","status":"1"}', '忘记密码短信模板', 'sms', '忘记密码场景使用的短信模板（JSON格式）', '1', 'system')
ON DUPLICATE KEY UPDATE `AAA022`=VALUES(`AAA022`);

