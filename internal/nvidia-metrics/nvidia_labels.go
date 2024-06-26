package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	gauge "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
)

// SetDeviceMetric sets the metric value for the given device
func SetDeviceMetric(handle nvml.Device, metricConfig config.Metric, metricValue float64) {
	metric := metricConfig.GetMetric()
	metricLabels := labelManager.GetMetricLabelValues(handle, metric)
	gauge.SetGaugeMetric(metric, metricLabels, metricValue)
}

// AddFunctions adds the label function to the map
func (lf LabelFunctions) AddFunctions() {

	lf.Add(config.GPU_ID.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		index, ret := device.GetIndex()
		return index, ret
	})

	lf.Add(config.GPU_NAME.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		name, ret := device.GetName()
		return name, ret
	})

	// GPU temperature threshold protections can shut down system when it hits the temp.limit,
	lf.Add(config.GPU_TEM_THRESHOLD.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		threshold, ret := device.GetTemperatureThreshold(nvml.TEMPERATURE_THRESHOLD_SHUTDOWN)
		return threshold, ret
	})

	//determines the rate at which the GPU can access and manipulate data stored in the VRAM
	lf.Add(config.GPU_MEM_CLOCK_MAX.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		clock, ret := device.GetMaxClockInfo(nvml.CLOCK_MEM)
		return clock, ret

	})

	lf.Add(config.GPU_SM_CLOCK_MAX.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		clock, ret := device.GetMaxClockInfo(nvml.CLOCK_SM)
		return clock, ret
	})

	lf.Add(config.GPU_GRAPHICS_CLOCK_MAX.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		clock, ret := device.GetMaxClockInfo(nvml.CLOCK_GRAPHICS)
		return clock, ret
	})

	lf.Add(config.GPU_CORES.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		cores, ret := device.GetNumGpuCores()
		return cores, ret
	})

	lf.Add(config.GPU_DRIVER_VERSION.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		driverVersion, ret := nvml.SystemGetDriverVersion()
		return driverVersion, ret
	})

	lf.Add(config.GPU_CUDA_VERSION.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		cudaVersion, ret := nvml.SystemGetCudaDriverVersion()
		return cudaVersion, ret
	})

	lf.Add(config.GPU_PEAK_FLOPS.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		id, err := device.GetIndex()
		if err != nvml.SUCCESS {
			return 0, err
		}

		// Get device flops
		flops, err := config.GetGpuFlops(id)
		if err != nvml.SUCCESS || flops == 0 {
			return 0, nvml.ERROR_UNKNOWN
		}
		// Convert flops to TFLOPS
		tflops := flops / 1e12

		return tflops, nvml.SUCCESS
	})

	// @TODO add additional label function to the map
	//lf.Add(config.GPU_POWER.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
	//	operationMode, _, r := device.GetGpuOperationMode()
	//})

	logger.Debug("Collected GPU Labels")

}
