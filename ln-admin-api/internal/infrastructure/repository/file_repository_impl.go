package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
	"github.com/Light-ink-yht/ln-admin/internal/domain/repository"
	"github.com/Light-ink-yht/ln-admin/internal/infrastructure/database"

	"gorm.io/gorm"
)

var _ repository.FileRepository = (*fileRepositoryImpl)(nil)

type fileRepositoryImpl struct {
	db *gorm.DB
}

// NewFileRepository 创建文件仓库实例
func NewFileRepository() repository.FileRepository {
	return &fileRepositoryImpl{
		db: database.GetDB(),
	}
}

// Create 创建文件记录
func (r *fileRepositoryImpl) Create(ctx context.Context, file *entity.File) error {
	aa10 := file.ToAA10()
	if err := r.db.WithContext(ctx).Create(aa10).Error; err != nil {
		return fmt.Errorf("创建文件记录失败: %w", err)
	}
	file.FromAA10(aa10)
	return nil
}

// Update 更新文件记录
func (r *fileRepositoryImpl) Update(ctx context.Context, file *entity.File) error {
	aa10 := file.ToAA10()
	if err := r.db.WithContext(ctx).Model(&entity.AA10File{}).
		Where("AAJ001 = ?", file.FileID).
		Updates(aa10).Error; err != nil {
		return fmt.Errorf("更新文件记录失败: %w", err)
	}
	return nil
}

// Delete 删除文件记录（软删除）
func (r *fileRepositoryImpl) Delete(ctx context.Context, fileID string) error {
	if err := r.db.WithContext(ctx).Where("AAJ001 = ?", fileID).
		Delete(&entity.AA10{}).Error; err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}
	return nil
}

// FindByID 根据文件ID查找
func (r *fileRepositoryImpl) FindByID(ctx context.Context, fileID string) (*entity.File, error) {
	var aa10 entity.AA10File
	if err := r.db.WithContext(ctx).Where("AAJ001 = ?", fileID).First(&aa10).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}
	file := &entity.File{}
	file.FromAA10(&aa10)
	return file, nil
}

// FindByFilePath 根据文件路径查找
func (r *fileRepositoryImpl) FindByFilePath(ctx context.Context, filePath string) (*entity.File, error) {
	var aa10 entity.AA10File
	if err := r.db.WithContext(ctx).Where("AAJ004 = ?", filePath).First(&aa10).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}
	file := &entity.File{}
	file.FromAA10(&aa10)
	return file, nil
}

// List 分页查询文件列表
func (r *fileRepositoryImpl) List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.File, int64, error) {
	var aa10List []entity.AA10File
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.AA10File{})

	// 应用查询条件
	for key, value := range conditions {
		switch key {
		case "file_name":
			if values, ok := value.([]interface{}); ok && len(values) > 0 {
				var orConditions []string
				var args []interface{}
				for _, v := range values {
					orConditions = append(orConditions, "AAJ002 LIKE ?")
					args = append(args, "%"+fmt.Sprintf("%v", v)+"%")
				}
				query = query.Where("("+strings.Join(orConditions, " OR ")+")", args...)
			} else {
				query = query.Where("AAJ002 LIKE ?", "%"+fmt.Sprintf("%v", value)+"%")
			}
		case "storage_type":
			query = query.Where("AAJ005 = ?", value)
		case "category":
			query = query.Where("AAJ009 = ?", value)
		case "extension":
			if values, ok := value.([]interface{}); ok && len(values) > 0 {
				query = query.Where("AAJ008 IN ?", values)
			} else {
				query = query.Where("AAJ008 = ?", value)
			}
		case "creator_id":
			query = query.Where("AAJ010 = ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询文件总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&aa10List).Error; err != nil {
		return nil, 0, fmt.Errorf("分页查询文件失败: %w", err)
	}

	files := make([]*entity.File, len(aa10List))
	for i, aa10 := range aa10List {
		file := &entity.File{}
		file.FromAA10(&aa10)
		files[i] = file
	}

	return files, total, nil
}
