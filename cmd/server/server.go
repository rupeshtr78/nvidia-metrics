package server

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rupeshtr78/nvidia-metrics/api"
	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

func RunServer() {
	configFile := getEnv("CONFIG_FILE", "config/metrics.yaml")
	logLevel := getEnv("LOG_LEVEL", "info")
	port := getEnv("PORT", "9500")
	host := getEnv("HOST", "0.0.0.0")
	interval := getEnv("INTERVAL", "5")
	logFilePath := getEnv("LOG_FILE_PATH", "logs/gpu-metrics.log")
	logToFile := getEnv("LOG_TO_FILE", "false")

	flag.StringVar(&configFile, "config", configFile, "Path to the configuration file")
	flag.StringVar(&logLevel, "loglevel", logLevel, "Log level (debug, info, warn, error,fatal)")
	flag.StringVar(&port, "port", port, "Port to run the metrics server")
	flag.StringVar(&host, "host", host, "Host to run the metrics server")
	flag.StringVar(&interval, "interval", interval, "Time interval in seconds to scrape metrics")
	flag.StringVar(&logFilePath, "logfile", logFilePath, "Log file path")
	flag.StringVar(&logToFile, "filelog", logToFile, "Enable file logging")

	flag.Parse()

	if configFile == "" {
		log.Fatal("Config file is required")
	}

	fileLogBool, err := strconv.ParseBool(logToFile)
	if err != nil {
		log.Fatal("Failed to convert file log to boolean", err)
	}

	// Initialize the logger
	err = logger.GetLogger(logLevel, fileLogBool, logFilePath)
	if err != nil {
		log.Fatal("Failed to initialize logger", err)
	}

	metricsConfig := filepath.Join(configFile)

	ctxCreateMetrics, cancelCreateMetrics := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCreateMetrics()

	// Register the metrics with Prometheus
	err = prometheusmetrics.CreatePrometheusMetrics(ctxCreateMetrics, metricsConfig)
	if err != nil {
		logger.Fatal("Failed to create Prometheus metrics", zap.Error(err))
		os.Exit(1)
	}

	// get the address from the host and port
	address := host + ":" + port

	//  get the metrics scrape interval
	t, err := strconv.Atoi(interval)
	if err != nil {
		logger.Fatal("Failed to convert time to integer", zap.Error(err))
		os.Exit(1)
	}

	var scrapreInterval time.Duration
	if t > 0 {
		scrapreInterval = time.Duration(t) * time.Second
	}

	// Start the metrics server with a long-running context
	// ctxRunServer uses context.WithCancel to ensure the server can run indefinitely until explicitly canceled.
	ctxRunServer, cancelRunServer := context.WithCancel(context.Background())
	defer cancelRunServer()

	// start the metrics server
	api.RunPrometheusMetricsServer(ctxRunServer, address, scrapreInterval)
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// @TODO - Remove this function for testing only
func RunMetricsLocal() {

	// Initialize NVML before starting the metric collection loop
	nvidiametrics.InitNVML()
	defer nvidiametrics.ShutdownNVML()

	ctx := context.TODO()

	for {
		nvidiametrics.CollectGpuMetrics(ctx)
		time.Sleep(30 * time.Second)
		for key, label := range prometheusmetrics.RegisteredLabels {
			logger.Debug("Registered label", zap.String("key", key), zap.Any("label", label))
		}
		// "key":"gpu_power_usage","label":{"label1":"gpu_id","label2":"gpu_name"}}

	}

	// get all labels for debugging

}
