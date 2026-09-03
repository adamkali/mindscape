package responses

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type CoolifyCPUMetric struct {
	Percent float64 `json:"percent"`
	Time    string  `json:"time"`
} // @name CoolifyCPUMetric

type CoolifyMemoryMetric struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
	Time        string  `json:"time"`
} // @name CoolifyMemoryMetric

type CoolifyStorageMetric struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"used_percent"`
	// Source names the machine the reading came from. Neither Coolify's API nor
	// Sentinel reports disk, so this is always the host running Mindscape.
	Source string `json:"source"`
} // @name CoolifyStorageMetric

// CoolifyWidgetMetricsData carries whatever could be collected. CPU and Memory
// are nullable and Storage may be empty on purpose: Sentinel and each
// configured filesystem are independent sources, and one being unreachable
// should not blank the remaining tiles. Anything that failed is explained in
// Warnings.
type CoolifyWidgetMetricsData struct {
	ServerName  string                 `json:"server_name"`
	ServerUUID  string                 `json:"server_uuid"`
	SentinelURL string                 `json:"sentinel_url"`
	CPU         *CoolifyCPUMetric      `json:"cpu"`
	Memory      *CoolifyMemoryMetric   `json:"memory"`
	Storage     []CoolifyStorageMetric `json:"storage"`
	Warnings    []string               `json:"warnings"`
} // @name CoolifyWidgetMetricsData

type CoolifyWidgetMetricsResponse struct {
	Data    *CoolifyWidgetMetricsData `json:"data"`
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
} // @name CoolifyWidgetMetricsResponse

func NewCoolifyWidgetMetricsResponse() *CoolifyWidgetMetricsResponse {
	return &CoolifyWidgetMetricsResponse{
		Data: &CoolifyWidgetMetricsData{
			Storage:  []CoolifyStorageMetric{},
			Warnings: []string{},
		},
		Success: false,
		Message: "",
	}
}

func (w *CoolifyWidgetMetricsResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *CoolifyWidgetMetricsResponse) Successful(
	ctx echo.Context,
	serverName string,
	serverUUID string,
	sentinelURL string,
	cpu *clients.SentinelCPU,
	memory *clients.SentinelMemory,
	storage []clients.Disk,
	warnings []string,
) error {
	data := &CoolifyWidgetMetricsData{
		ServerName:  serverName,
		ServerUUID:  serverUUID,
		SentinelURL: sentinelURL,
		Storage:     []CoolifyStorageMetric{},
		Warnings:    warnings,
	}
	if data.Warnings == nil {
		data.Warnings = []string{}
	}

	if cpu != nil {
		data.CPU = &CoolifyCPUMetric{
			Percent: cpu.Percent,
			Time:    cpu.Time,
		}
	}
	if memory != nil {
		data.Memory = &CoolifyMemoryMetric{
			Total:       memory.Total,
			Used:        memory.Used,
			Available:   memory.Available,
			Free:        memory.Free,
			UsedPercent: memory.UsedPercent,
			Time:        memory.Time,
		}
	}
	for _, disk := range storage {
		data.Storage = append(data.Storage, CoolifyStorageMetric{
			Path:        disk.Path,
			Total:       disk.Total,
			Used:        disk.Used,
			Available:   disk.Available,
			UsedPercent: disk.UsedPercent,
			Source:      "mindscape-host",
		})
	}

	w.Success = true
	w.Data = data
	return ctx.JSON(200, w)
}
