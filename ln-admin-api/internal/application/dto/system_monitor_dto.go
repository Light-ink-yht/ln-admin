package dto

// SystemMonitorResponse 系统监控响应
type SystemMonitorResponse struct {
	CPU    CPUInfo    `json:"cpu"`
	Memory MemoryInfo `json:"memory"`
	Disk   DiskInfo   `json:"disk"`
	Uptime string     `json:"uptime"` // 系统运行时间
}

// CPUInfo CPU信息
type CPUInfo struct {
	Usage       float64 `json:"usage"`       // CPU使用率（百分比）
	Count       int     `json:"count"`       // CPU核心数
	ModelName   string  `json:"model_name"`  // CPU型号
	Temperature float64 `json:"temperature"` // CPU温度（如果可获取）
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total       uint64  `json:"total"`        // 总内存（字节）
	Used        uint64  `json:"used"`         // 已使用内存（字节）
	Available   uint64  `json:"available"`    // 可用内存（字节）
	Usage       float64 `json:"usage"`        // 内存使用率（百分比）
	TotalGB     float64 `json:"total_gb"`     // 总内存（GB）
	UsedGB      float64 `json:"used_gb"`      // 已使用内存（GB）
	AvailableGB float64 `json:"available_gb"` // 可用内存（GB）
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Total       uint64          `json:"total"`        // 总磁盘空间（字节）
	Used        uint64          `json:"used"`         // 已使用磁盘空间（字节）
	Available   uint64          `json:"available"`    // 可用磁盘空间（字节）
	Usage       float64         `json:"usage"`        // 磁盘使用率（百分比）
	TotalGB     float64         `json:"total_gb"`     // 总磁盘空间（GB）
	UsedGB      float64         `json:"used_gb"`      // 已使用磁盘空间（GB）
	AvailableGB float64         `json:"available_gb"` // 可用磁盘空间（GB）
	Disks       []DiskPartition `json:"disks"`        // 磁盘分区列表
}

// DiskPartition 磁盘分区信息
type DiskPartition struct {
	Device      string  `json:"device"`       // 设备名称
	Mountpoint  string  `json:"mountpoint"`   // 挂载点
	Fstype      string  `json:"fstype"`       // 文件系统类型
	Total       uint64  `json:"total"`        // 总空间（字节）
	Used        uint64  `json:"used"`         // 已使用空间（字节）
	Available   uint64  `json:"available"`    // 可用空间（字节）
	Usage       float64 `json:"usage"`        // 使用率（百分比）
	TotalGB     float64 `json:"total_gb"`     // 总空间（GB）
	UsedGB      float64 `json:"used_gb"`      // 已使用空间（GB）
	AvailableGB float64 `json:"available_gb"` // 可用空间（GB）
}
