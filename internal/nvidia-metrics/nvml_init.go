package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

type GpuDevice interface {
	Init() nvml.Return
	Shutdown() nvml.Return
	GetIndex() (int, nvml.Return)
	GetDeviceCount() (int, nvml.Return)
	GetDeviceHandleByIndex(int) (nvml.Device, nvml.Return)
	GetUtilizationRates() (nvml.Utilization, nvml.Return)
	GetMemoryInfo() (nvml.Memory, nvml.Return)
	GetPowerUsage() (uint32, nvml.Return)
	GetRunningProcesses() ([]nvml.ProcessInfo, nvml.Return)
	GetTemperature() (uint, nvml.Return)
	GetClockInfo() (uint32, nvml.Return)
	GetEccErrors() (nvml.EccErrorCounts, nvml.Return)
	GetFanSpeed() (uint32, nvml.Return)
	GetPeakFlops() (float64, nvml.Return)
	GetPerformanceState() (nvml.Pstates, nvml.Return)
}

// GPUDeviceMetrics represents the collected metrics for a GPU device.
type GPUDeviceMetrics struct {
	DeviceIndex         int
	GPUTemperature      float64
	GPUCPUUtilization   float64
	GPUMemUtilization   float64
	GPUPowerUsage       float64
	GPURunningProcesses int
	GPUMemoryUsed       uint64
	GPUMemoryTotal      uint64
	GPUMemoryFree       uint64
	GpuPState           int32
	GpuClock            uint32
	GpuEccErrors        uint64
	GpuFanSpeed         uint32
	GpuPeakFlops        float64
}

func NewGPUDeviceMetrics() *GPUDeviceMetrics {
	return &GPUDeviceMetrics{}
}

// InitNVML initializes the NVML library.
func InitNVML() {
	if err := nvml.Init(); err != nvml.SUCCESS {
		logger.Fatal("Failed to initialize NVML", zap.Error(err))
	}
	logger.Info("Initialized NVML")
}

// ShutdownNVML shuts down the NVML library.
func ShutdownNVML() {
	if err := nvml.Shutdown(); err != nvml.SUCCESS {
		logger.Fatal("Failed to shutdown NVML", zap.Error(err))
	}
	logger.Info("Shutdown NVML")
}
