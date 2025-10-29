package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error

	// Update 更新用户
	Update(ctx context.Context, user *entity.User) error

	// Delete 删除用户（软删除）
	Delete(ctx context.Context, userID string) error

	// FindByID 根据用户ID查找
	FindByID(ctx context.Context, userID string) (*entity.User, error)

	// FindByPhone 根据手机号查找
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)

	// FindByEmail 根据邮箱查找
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindByIDs 根据用户ID列表批量查找
	FindByIDs(ctx context.Context, userIDs []string) ([]*entity.User, error)

	// List 分页查询用户列表
	List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.User, int64, error)

	// UpdateLoginInfo 更新登录信息
	UpdateLoginInfo(ctx context.Context, userID string, loginIP string) error
}
