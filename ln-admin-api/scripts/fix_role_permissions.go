package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 修复角色权限的脚本
// 用途：为超级管理员角色添加所有角色管理相关的权限

func main() {
	// 数据库连接配置（请根据实际情况修改）
	dsn := "root:123456@tcp(localhost:3306)/ln_admin?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	ctx := context.Background()

	fmt.Println("开始修复角色管理权限...")

	// 1. 添加角色管理权限到权限表（如果不存在）
	rolePermissions := []struct {
		PermissionID   string
		PermissionKey  string
		PermissionName string
		ResourcePath   string
		Method         string
		Description    string
	}{
		{
			PermissionID:   "perm_role_list",
			PermissionKey:  "role:list",
			PermissionName: "角色列表",
			ResourcePath:   "/api/role/list",
			Method:         "GET",
			Description:    "查看角色列表",
		},
		{
			PermissionID:   "perm_role_detail",
			PermissionKey:  "role:detail",
			PermissionName: "角色详情",
			ResourcePath:   "/api/role/*",
			Method:         "GET",
			Description:    "查看角色详情",
		},
		{
			PermissionID:   "perm_role_create",
			PermissionKey:  "role:create",
			PermissionName: "创建角色",
			ResourcePath:   "/api/role",
			Method:         "POST",
			Description:    "创建新角色",
		},
		{
			PermissionID:   "perm_role_update",
			PermissionKey:  "role:update",
			PermissionName: "更新角色",
			ResourcePath:   "/api/role/*",
			Method:         "PUT",
			Description:    "更新角色信息",
		},
		{
			PermissionID:   "perm_role_delete",
			PermissionKey:  "role:delete",
			PermissionName: "删除角色",
			ResourcePath:   "/api/role/*",
			Method:         "DELETE",
			Description:    "删除角色",
		},
	}

	for _, perm := range rolePermissions {
		// 检查权限是否已存在
		var count int
		err := db.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM AA04 WHERE AAF002 = ?",
			perm.PermissionKey,
		).Scan(&count)
		if err != nil {
			log.Printf("检查权限失败 %s: %v", perm.PermissionKey, err)
			continue
		}

		if count == 0 {
			// 插入新权限
			_, err := db.ExecContext(ctx,
				`INSERT INTO AA04 (AAF001, AAF002, AAF003, AAF004, AAF005, AAF006, AAF006, AAF007) 
				 VALUES (?, ?, ?, ?, ?, '1', 'system', 'system')`,
				perm.PermissionID,
				perm.PermissionKey,
				perm.PermissionName,
				perm.ResourcePath,
				perm.Method,
			)
			if err != nil {
				log.Printf("插入权限失败 %s: %v", perm.PermissionKey, err)
			} else {
				fmt.Printf("✓ 添加权限: %s\n", perm.PermissionName)
			}
		} else {
			fmt.Printf("- 权限已存在: %s\n", perm.PermissionName)
		}
	}

	// 2. 获取超级管理员角色
	var superAdminRoleKey string
	err = db.QueryRowContext(ctx,
		"SELECT AAE002 FROM AA05 WHERE AAE002 = 'super_admin'",
	).Scan(&superAdminRoleKey)
	if err != nil {
		log.Fatalf("查询超级管理员角色失败: %v", err)
	}

	fmt.Printf("\n找到超级管理员角色: %s\n", superAdminRoleKey)

	// 3. 为超级管理员角色添加所有角色管理权限到Casbin
	for _, perm := range rolePermissions {
		// 移除路径前缀 /api（因为中间件会移除这个前缀）
		resourcePath := strings.TrimPrefix(perm.ResourcePath, "/api")
		// 将通配符路径中的 * 转换为 Casbin 的匹配模式
		// Casbin 的 keyMatch 使用 * 作为通配符

		// 检查策略是否已存在
		var count int
		err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM casbin_rule 
			 WHERE ptype = 'p' AND v0 = ? AND v1 = ? AND v2 = ?`,
			superAdminRoleKey, resourcePath, perm.Method,
		).Scan(&count)
		if err != nil {
			log.Printf("检查Casbin策略失败 %s %s: %v", resourcePath, perm.Method, err)
			continue
		}

		if count == 0 {
			// 插入新策略
			_, err := db.ExecContext(ctx,
				`INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', ?, ?, ?)`,
				superAdminRoleKey, resourcePath, perm.Method,
			)
			if err != nil {
				log.Printf("插入Casbin策略失败 %s %s: %v", resourcePath, perm.Method, err)
			} else {
				fmt.Printf("✓ 添加Casbin策略: %s %s %s\n", superAdminRoleKey, resourcePath, perm.Method)
			}
		} else {
			fmt.Printf("- Casbin策略已存在: %s %s %s\n", superAdminRoleKey, resourcePath, perm.Method)
		}
	}

	// 4. 确保用户 18797131041 拥有 super_admin 角色
	fmt.Println("\n检查用户 18797131041 的角色...")
	var userID string
	err = db.QueryRowContext(ctx,
		"SELECT AAE001 FROM AA01 WHERE AAE003 = '18797131041'",
	).Scan(&userID)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
	} else {
		fmt.Printf("找到用户ID: %s\n", userID)

		// 检查用户是否已有 super_admin 角色（在用户角色关联表中）
		var userRoleCount int
		err = db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM AA06 ur 
			 INNER JOIN AA05 r ON ur.AAF001 = r.AAE001 
			 WHERE ur.AAE001 = ? AND r.AAE002 = 'super_admin'`,
			userID,
		).Scan(&userRoleCount)
		if err != nil {
			log.Printf("检查用户角色失败: %v", err)
		} else if userRoleCount == 0 {
			// 获取 super_admin 的 role_id
			var superAdminRoleID string
			err = db.QueryRowContext(ctx,
				"SELECT AAE001 FROM AA05 WHERE AAE002 = 'super_admin'",
			).Scan(&superAdminRoleID)
			if err != nil {
				log.Printf("获取超级管理员角色ID失败: %v", err)
			} else {
				// 分配角色
				_, err = db.ExecContext(ctx,
					`INSERT INTO AA06 (AAE001, AAF001) VALUES (?, ?)`,
					userID, superAdminRoleID,
				)
				if err != nil {
					log.Printf("分配角色失败: %v", err)
				} else {
					fmt.Printf("✓ 为用户分配超级管理员角色\n")
				}
			}
		} else {
			fmt.Printf("- 用户已拥有超级管理员角色\n")
		}

		// 确保Casbin中有用户角色关系
		var casbinRoleCount int
		err = db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM casbin_rule 
			 WHERE ptype = 'g' AND v0 = ? AND v1 = ?`,
			userID, superAdminRoleKey,
		).Scan(&casbinRoleCount)
		if err != nil {
			log.Printf("检查Casbin用户角色关系失败: %v", err)
		} else if casbinRoleCount == 0 {
			_, err = db.ExecContext(ctx,
				`INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', ?, ?)`,
				userID, superAdminRoleKey,
			)
			if err != nil {
				log.Printf("添加Casbin用户角色关系失败: %v", err)
			} else {
				fmt.Printf("✓ 添加Casbin用户角色关系\n")
			}
		} else {
			fmt.Printf("- Casbin用户角色关系已存在\n")
		}
	}

	fmt.Println("\n修复完成！请重启应用服务器以使权限生效。")
}
