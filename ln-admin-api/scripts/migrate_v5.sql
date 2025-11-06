-- MySQL 数据库迁移脚本 V5
-- 新增：角色权限授予关系表 AA13

-- 角色权限授予关系表 AA13
CREATE TABLE IF NOT EXISTS `AA13` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `AAM001` varchar(50) NOT NULL COMMENT '授予者角色ID',
  `AAM002` varchar(50) NOT NULL COMMENT '被授予者角色ID',
  `AAM003` varchar(50) NOT NULL COMMENT '权限ID',
  PRIMARY KEY (`id`),
  KEY `idx_aam001` (`AAM001`),
  KEY `idx_aam002` (`AAM002`),
  KEY `idx_aam003` (`AAM003`),
  KEY `idx_grantor_role_permission` (`AAM001`, `AAM002`, `AAM003`),
  KEY `idx_grantee_role_permission` (`AAM002`, `AAM003`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限授予关系表：存储角色权限的授予关系，实现层级授权机制。通过 AAM001（授予者角色ID）、AAM002（被授予者角色ID）和 AAM003（权限ID）建立权限授予关系，记录谁授予给谁的权限。超级管理员拥有所有权限，可以为其他角色授予权限；管理员只能看到被授予的权限，并可以将这些权限授予给下一层角色。支持权限的层级传递，确保权限管理的安全性和可追溯性。';

-- 创建唯一联合索引（防止重复授予）
CREATE UNIQUE INDEX IF NOT EXISTS `uk_grantor_grantee_permission` ON `AA13` (`AAM001`, `AAM002`, `AAM003`);

