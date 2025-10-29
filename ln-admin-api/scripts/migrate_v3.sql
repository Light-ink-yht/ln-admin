-- MySQL 数据库迁移脚本 V3
-- 新增：角色表、权限表、用户角色关联表

-- 角色表 AA05
CREATE TABLE IF NOT EXISTS `AA05` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `AAA051` varchar(50) NOT NULL COMMENT '角色ID',
  `AAA052` varchar(100) NOT NULL COMMENT '角色标识（唯一）',
  `AAA053` varchar(100) DEFAULT NULL COMMENT '角色名称',
  `AAA054` varchar(500) DEFAULT NULL COMMENT '角色描述',
  `AAA055` varchar(10) DEFAULT '1' COMMENT '状态 1 启用，2 禁用',
  `AAA056` varchar(50) DEFAULT NULL COMMENT '创建人',
  `AAA057` varchar(50) DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_aaa051` (`AAA051`),
  UNIQUE KEY `uk_aaa052` (`AAA052`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 权限表 AA06
CREATE TABLE IF NOT EXISTS `AA06` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `AAA061` varchar(50) NOT NULL COMMENT '权限ID',
  `AAA062` varchar(100) NOT NULL COMMENT '权限标识（唯一）',
  `AAA063` varchar(100) DEFAULT NULL COMMENT '权限名称',
  `AAA064` varchar(200) DEFAULT NULL COMMENT '资源路径（API路径）',
  `AAA065` varchar(10) DEFAULT NULL COMMENT '请求方法（GET, POST, PUT, DELETE等）',
  `AAA066` varchar(500) DEFAULT NULL COMMENT '权限描述',
  `AAA067` varchar(10) DEFAULT '1' COMMENT '状态 1 启用，2 禁用',
  `AAA068` varchar(50) DEFAULT NULL COMMENT '创建人',
  `AAA069` varchar(50) DEFAULT NULL COMMENT '修改人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_aaa061` (`AAA061`),
  UNIQUE KEY `uk_aaa062` (`AAA062`),
  KEY `idx_aaa064` (`AAA064`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 用户角色关联表 AA07
CREATE TABLE IF NOT EXISTS `AA07` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `AAA071` varchar(50) DEFAULT NULL COMMENT '用户ID',
  `AAA072` varchar(50) DEFAULT NULL COMMENT '角色ID',
  PRIMARY KEY (`id`),
  KEY `idx_aaa071` (`AAA071`),
  KEY `idx_aaa072` (`AAA072`),
  KEY `idx_deleted_at` (`deleted_at`),
  UNIQUE KEY `uk_user_role` (`AAA071`, `AAA072`, `deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- Casbin规则表（由Casbin适配器自动创建）
-- casbin_rule 表由 gorm-adapter 自动管理，无需手动创建

-- 初始化默认角色（示例）
INSERT INTO `AA05` (`AAA051`, `AAA052`, `AAA053`, `AAA054`, `AAA055`, `AAA056`) VALUES
('role_admin', 'admin', '超级管理员', '系统超级管理员，拥有所有权限', '1', 'system'),
('role_user', 'user', '普通用户', '普通用户角色', '1', 'system')
ON DUPLICATE KEY UPDATE `AAA052`=`AAA052`;

-- 初始化默认权限（示例）
INSERT INTO `AA06` (`AAA061`, `AAA062`, `AAA063`, `AAA064`, `AAA065`, `AAA066`, `AAA067`, `AAA068`) VALUES
('perm_user_list', 'user:list', '用户列表', '/api/user/list', 'GET', '查看用户列表', '1', 'system'),
('perm_user_info', 'user:info', '用户信息', '/api/user/userinfo', 'GET', '查看用户信息', '1', 'system'),
('perm_user_create', 'user:create', '创建用户', '/api/user/create', 'POST', '创建新用户', '1', 'system'),
('perm_user_update', 'user:update', '更新用户', '/api/user/update', 'PUT', '更新用户信息', '1', 'system'),
('perm_user_delete', 'user:delete', '删除用户', '/api/user/delete', 'DELETE', '删除用户', '1', 'system')
ON DUPLICATE KEY UPDATE `AAA062`=`AAA062`;

