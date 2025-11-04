package service

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/Light-ink-yht/ln-admin/internal/application/dto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMonitorService 系统监控服务
type SystemMonitorService struct{}

// NewSystemMonitorService 创建系统监控服务
func NewSystemMonitorService() *SystemMonitorService {
	return &SystemMonitorService{}
}

// GetSystemMonitor 获取系统监控信息
func (s *SystemMonitorService) GetSystemMonitor(ctx context.Context) (*dto.SystemMonitorResponse, error) {
	// 获取CPU信息
	cpuInfo, err := s.getCPUInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取CPU信息失败: %w", err)
	}

	// 获取内存信息
	memoryInfo, err := s.getMemoryInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取内存信息失败: %w", err)
	}

	// 获取磁盘信息
	diskInfo, err := s.getDiskInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取磁盘信息失败: %w", err)
	}

	// 获取系统运行时间
	uptime := s.getUptime()

	return &dto.SystemMonitorResponse{
		CPU:    *cpuInfo,
		Memory: *memoryInfo,
		Disk:   *diskInfo,
		Uptime: uptime,
	}, nil
}

// getCPUInfo 获取CPU信息
func (s *SystemMonitorService) getCPUInfo(ctx context.Context) (*dto.CPUInfo, error) {
	// 获取CPU使用率（1秒内的平均使用率）
	percentages, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return nil, err
	}

	usage := 0.0
	if len(percentages) > 0 {
		usage = percentages[0]
	}

	// 获取CPU核心数
	count := runtime.NumCPU()

	// 获取CPU信息（型号等）
	cpuInfoList, err := cpu.InfoWithContext(ctx)
	modelName := "Unknown"
	if err == nil && len(cpuInfoList) > 0 {
		modelName = cpuInfoList[0].ModelName
	}

	return &dto.CPUInfo{
		Usage:     usage,
		Count:     count,
		ModelName: modelName,
	}, nil
}

// getMemoryInfo 获取内存信息
func (s *SystemMonitorService) getMemoryInfo(ctx context.Context) (*dto.MemoryInfo, error) {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	usage := 0.0
	if vmStat.Total > 0 {
		usage = (float64(vmStat.Used) / float64(vmStat.Total)) * 100
	}

	return &dto.MemoryInfo{
		Total:       vmStat.Total,
		Used:        vmStat.Used,
		Available:   vmStat.Available,
		Usage:       usage,
		TotalGB:     float64(vmStat.Total) / (1024 * 1024 * 1024),
		UsedGB:      float64(vmStat.Used) / (1024 * 1024 * 1024),
		AvailableGB: float64(vmStat.Available) / (1024 * 1024 * 1024),
	}, nil
}

// getDiskInfo 获取磁盘信息
func (s *SystemMonitorService) getDiskInfo(ctx context.Context) (*dto.DiskInfo, error) {
	// 获取所有磁盘分区
	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, err
	}

	var totalDisk uint64 = 0
	var usedDisk uint64 = 0
	var availableDisk uint64 = 0
	var diskPartitions []dto.DiskPartition

	// 遍历每个分区
	for _, partition := range partitions {
		// 跳过特殊文件系统
		if partition.Fstype == "" || partition.Fstype == "tmpfs" || partition.Fstype == "devtmpfs" {
			continue
		}

		// 获取分区使用情况
		usage, err := disk.UsageWithContext(ctx, partition.Mountpoint)
		if err != nil {
			continue // 跳过无法访问的分区
		}

		diskUsage := 0.0
		if usage.Total > 0 {
			diskUsage = (float64(usage.Used) / float64(usage.Total)) * 100
		}

		partitionInfo := dto.DiskPartition{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Fstype:      partition.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Available:   usage.Free,
			Usage:       diskUsage,
			TotalGB:     float64(usage.Total) / (1024 * 1024 * 1024),
			UsedGB:      float64(usage.Used) / (1024 * 1024 * 1024),
			AvailableGB: float64(usage.Free) / (1024 * 1024 * 1024),
		}

		diskPartitions = append(diskPartitions, partitionInfo)

		// 累加总磁盘空间（只统计根分区，避免重复计算）
		if partition.Mountpoint == "/" || partition.Mountpoint == "C:\\" {
			totalDisk = usage.Total
			usedDisk = usage.Used
			availableDisk = usage.Free
		}
	}

	// 如果没有找到根分区，使用第一个分区
	if totalDisk == 0 && len(diskPartitions) > 0 {
		totalDisk = diskPartitions[0].Total
		usedDisk = diskPartitions[0].Used
		availableDisk = diskPartitions[0].Available
	}

	usage := 0.0
	if totalDisk > 0 {
		usage = (float64(usedDisk) / float64(totalDisk)) * 100
	}

	return &dto.DiskInfo{
		Total:       totalDisk,
		Used:        usedDisk,
		Available:   availableDisk,
		Usage:       usage,
		TotalGB:     float64(totalDisk) / (1024 * 1024 * 1024),
		UsedGB:      float64(usedDisk) / (1024 * 1024 * 1024),
		AvailableGB: float64(availableDisk) / (1024 * 1024 * 1024),
		Disks:       diskPartitions,
	}, nil
}

// getUptime 获取系统运行时间
func (s *SystemMonitorService) getUptime() string {
	// 获取系统启动时间
	bootTime, err := host.BootTime()
	if err != nil {
		return "未知"
	}

	uptimeSeconds := time.Now().Unix() - int64(bootTime)
	days := uptimeSeconds / 86400
	hours := (uptimeSeconds % 86400) / 3600
	minutes := (uptimeSeconds % 3600) / 60

	return fmt.Sprintf("%d 天 %d 小时 %d 分钟", days, hours, minutes)
}
