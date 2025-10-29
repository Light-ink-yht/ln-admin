package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"

	"gorm.io/gorm"
)

var _ repository.UserRepository = (*userRepositoryImpl)(nil)

type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository() repository.UserRepository {
	return &userRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	aa01 := user.ToAA01()
	if err := r.db.WithContext(ctx).Create(aa01).Error; err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	user.FromAA01(aa01)
	return nil
}

// Update 更新用户
func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	aa01 := user.ToAA01()
	if err := r.db.WithContext(ctx).Model(&entity.AA01{}).
		Where("AAA001 = ?", user.UserID).
		Updates(aa01).Error; err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}
	return nil
}

// Delete 删除用户（软删除）
func (r *userRepositoryImpl) Delete(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Where("AAA001 = ?", userID).
		Delete(&entity.AA01{}).Error; err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}
	return nil
}

// FindByID 根据用户ID查找
func (r *userRepositoryImpl) FindByID(ctx context.Context, userID string) (*entity.User, error) {
	var aa01 entity.AA01
	if err := r.db.WithContext(ctx).Where("AAA001 = ?", userID).First(&aa01).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	user := &entity.User{}
	user.FromAA01(&aa01)
	return user, nil
}

// FindByPhone 根据手机号查找
func (r *userRepositoryImpl) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var aa01 entity.AA01
	if err := r.db.WithContext(ctx).Where("AAA003 = ?", phone).First(&aa01).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	user := &entity.User{}
	user.FromAA01(&aa01)
	return user, nil
}

// FindByEmail 根据邮箱查找
func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var aa01 entity.AA01
	if err := r.db.WithContext(ctx).Where("AAA002 = ?", email).First(&aa01).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	user := &entity.User{}
	user.FromAA01(&aa01)
	return user, nil
}

// FindByIDs 根据用户ID列表批量查找
func (r *userRepositoryImpl) FindByIDs(ctx context.Context, userIDs []string) ([]*entity.User, error) {
	var aa01List []entity.AA01
	if err := r.db.WithContext(ctx).Where("AAA001 IN ?", userIDs).Find(&aa01List).Error; err != nil {
		return nil, fmt.Errorf("批量查询用户失败: %w", err)
	}
	users := make([]*entity.User, len(aa01List))
	for i, aa01 := range aa01List {
		user := &entity.User{}
		user.FromAA01(&aa01)
		users[i] = user
	}
	return users, nil
}

// List 分页查询用户列表
func (r *userRepositoryImpl) List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.User, int64, error) {
	var aa01List []entity.AA01
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.AA01{})

	// 应用查询条件
	for key, value := range conditions {
		switch key {
		case "phone":
			query = query.Where("AAA003 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "email":
			query = query.Where("AAA002 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "status":
			query = query.Where("AAA010 = ?", value)
		case "gender":
			query = query.Where("AAA008 = ?", value)
		case "nickname":
			query = query.Where("AAA005 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		case "full_name":
			query = query.Where("AAA006 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&aa01List).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询用户失败: %w", err)
	}

	users := make([]*entity.User, len(aa01List))
	for i, aa01 := range aa01List {
		user := &entity.User{}
		user.FromAA01(&aa01)
		users[i] = user
	}

	return users, total, nil
}

// UpdateLoginInfo 更新登录信息
func (r *userRepositoryImpl) UpdateLoginInfo(ctx context.Context, userID string, loginIP string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"AAA012": gorm.Expr("AAA012 + 1"),
		"AAA013": &now,
		"AAA014": loginIP,
	}

	if err := r.db.WithContext(ctx).Model(&entity.AA01{}).
		Where("AAA001 = ?", userID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("更新登录信息失败: %w", err)
	}

	return nil
}
