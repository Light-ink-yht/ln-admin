package repository

import (
	"context"

	"github.com/Light-ink-yht/ln-admin/internal/domain/entity"
)

// FileRepository 文件仓库接口
type FileRepository interface {
	// Create 创建文件记录
	Create(ctx context.Context, file *entity.File) error

	// Update 更新文件记录
	Update(ctx context.Context, file *entity.File) error

	// Delete 删除文件记录（软删除）
	Delete(ctx context.Context, fileID string) error

	// FindByID 根据文件ID查找
	FindByID(ctx context.Context, fileID string) (*entity.File, error)

	// List 分页查询文件列表
	List(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]*entity.File, int64, error)

	// FindByFilePath 根据文件路径查找
	FindByFilePath(ctx context.Context, filePath string) (*entity.File, error)
}
