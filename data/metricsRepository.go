package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type MetricsRepository struct {
	mutex       sync.RWMutex
	metricsDir  string
	currentDate string
	
	// In-memory counters (reset daily)
	EndpointStats   map[string]*EndpointMetric
	SeriesUsage     map[string]*SeriesMetric
	APIKeyUsage     map[string]*APIKeyMetric
	HourlyStats     [24]int64
	TotalRequests   int64
	TotalErrors     int64
}

type EndpointMetric struct {
	TotalRequests    int64     `json:"total_requests"`
	UniqueAPIKeys    int       `json:"unique_api_keys"`
	AvgResponseTimeMs float64  `json:"avg_response_time_ms"`
	ErrorCount       int64     `json:"error_count"`
	ResponseTimes    []float64 `json:"-"` // Not serialized, used for avg calculation
	APIKeysSet       map[string]bool `json:"-"` // Not serialized
}

type SeriesMetric struct {
	Requests     int64     `json:"requests"`
	LastAccessed time.Time `json:"last_accessed"`
}

type APIKeyMetric struct {
	TotalRequests int64     `json:"total_requests"`
	EndpointsUsed map[string]bool `json:"-"` // Not serialized
	LastActive    time.Time `json:"last_active"`
	ErrorCount    int64     `json:"error_count"`
}

type DailyMetrics struct {
	Date          string                    `json:"date"`
	TotalRequests int64                     `json:"total_requests"`
	TotalErrors   int64                     `json:"total_errors"`
	UniqueAPIKeys int                       `json:"unique_api_keys"`
	Endpoints     map[string]*EndpointMetric `json:"endpoints"`
	SeriesUsage   map[string]*SeriesMetric   `json:"series_usage"`
	APIKeyUsage   map[string]*APIKeyUsage    `json:"api_key_usage"`
	HourlyStats   [24]int64                 `json:"hourly_distribution"`
}

type APIKeyUsage struct {
	TotalRequests int64     `json:"total_requests"`
	EndpointsUsed []string  `json:"endpoints_used"`
	LastActive    time.Time `json:"last_active"`
	ErrorRate     float64   `json:"error_rate"`
}

func NewMetricsRepository(metricsDir string) *MetricsRepository {
	if metricsDir == "" {
		metricsDir = "./metrics"
	}
	
	// Ensure metrics directory exists
	if err := os.MkdirAll(metricsDir, 0755); err != nil {
		log.Printf("Failed to create metrics directory: %v", err)
	}
	
	repo := &MetricsRepository{
		metricsDir:    metricsDir,
		currentDate:   time.Now().Format("2006-01-02"),
		EndpointStats: make(map[string]*EndpointMetric),
		SeriesUsage:   make(map[string]*SeriesMetric),
		APIKeyUsage:   make(map[string]*APIKeyMetric),
	}
	
	// Load current day's metrics if they exist
	repo.loadCurrentMetrics()
	
	return repo
}

func (m *MetricsRepository) RecordRequest(endpoint, apiKey, seriesIdentifier string, responseTimeMs float64, isError bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// Check if date has changed (new day)
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != m.currentDate {
		m.rotateDailyMetrics()
		m.currentDate = currentDate
	}
	
	now := time.Now()
	hour := now.Hour()
	
	// Update totals
	m.TotalRequests++
	m.HourlyStats[hour]++
	if isError {
		m.TotalErrors++
	}
	
	// Update endpoint stats
	if m.EndpointStats[endpoint] == nil {
		m.EndpointStats[endpoint] = &EndpointMetric{
			APIKeysSet: make(map[string]bool),
		}
	}
	endpointStat := m.EndpointStats[endpoint]
	endpointStat.TotalRequests++
	endpointStat.APIKeysSet[apiKey] = true
	endpointStat.UniqueAPIKeys = len(endpointStat.APIKeysSet)
	endpointStat.ResponseTimes = append(endpointStat.ResponseTimes, responseTimeMs)
	if isError {
		endpointStat.ErrorCount++
	}
	
	// Calculate average response time
	total := 0.0
	for _, rt := range endpointStat.ResponseTimes {
		total += rt
	}
	endpointStat.AvgResponseTimeMs = total / float64(len(endpointStat.ResponseTimes))
	
	// Update series usage if identifier provided
	if seriesIdentifier != "" {
		if m.SeriesUsage[seriesIdentifier] == nil {
			m.SeriesUsage[seriesIdentifier] = &SeriesMetric{}
		}
		m.SeriesUsage[seriesIdentifier].Requests++
		m.SeriesUsage[seriesIdentifier].LastAccessed = now
	}
	
	// Update API key usage
	if m.APIKeyUsage[apiKey] == nil {
		m.APIKeyUsage[apiKey] = &APIKeyMetric{
			EndpointsUsed: make(map[string]bool),
		}
	}
	apiKeyStat := m.APIKeyUsage[apiKey]
	apiKeyStat.TotalRequests++
	apiKeyStat.EndpointsUsed[endpoint] = true
	apiKeyStat.LastActive = now
	if isError {
		apiKeyStat.ErrorCount++
	}
}

func (m *MetricsRepository) GetCurrentMetrics() *DailyMetrics {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	// Calculate unique API keys
	uniqueKeys := make(map[string]bool)
	for key := range m.APIKeyUsage {
		uniqueKeys[key] = true
	}
	
	// Convert API key usage to serializable format
	apiKeyUsage := make(map[string]*APIKeyUsage)
	for key, metric := range m.APIKeyUsage {
		endpoints := make([]string, 0, len(metric.EndpointsUsed))
		for endpoint := range metric.EndpointsUsed {
			endpoints = append(endpoints, endpoint)
		}
		
		errorRate := 0.0
		if metric.TotalRequests > 0 {
			errorRate = float64(metric.ErrorCount) / float64(metric.TotalRequests)
		}
		
		apiKeyUsage[key] = &APIKeyUsage{
			TotalRequests: metric.TotalRequests,
			EndpointsUsed: endpoints,
			LastActive:    metric.LastActive,
			ErrorRate:     errorRate,
		}
	}
	
	return &DailyMetrics{
		Date:          m.currentDate,
		TotalRequests: m.TotalRequests,
		TotalErrors:   m.TotalErrors,
		UniqueAPIKeys: len(uniqueKeys),
		Endpoints:     m.EndpointStats,
		SeriesUsage:   m.SeriesUsage,
		APIKeyUsage:   apiKeyUsage,
		HourlyStats:   m.HourlyStats,
	}
}

func (m *MetricsRepository) rotateDailyMetrics() {
	// Save current metrics to file
	metrics := m.GetCurrentMetrics()
	filename := filepath.Join(m.metricsDir, fmt.Sprintf("daily_metrics_%s.json", m.currentDate))
	
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal metrics: %v", err)
		return
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Failed to write metrics file: %v", err)
		return
	}
	
	log.Printf("Rotated daily metrics to %s", filename)
	
	// Reset in-memory counters
	m.EndpointStats = make(map[string]*EndpointMetric)
	m.SeriesUsage = make(map[string]*SeriesMetric)
	m.APIKeyUsage = make(map[string]*APIKeyMetric)
	m.HourlyStats = [24]int64{}
	m.TotalRequests = 0
	m.TotalErrors = 0
}

func (m *MetricsRepository) loadCurrentMetrics() {
	filename := filepath.Join(m.metricsDir, fmt.Sprintf("daily_metrics_%s.json", m.currentDate))
	
	data, err := os.ReadFile(filename)
	if err != nil {
		// File doesn't exist, start fresh
		return
	}
	
	var metrics DailyMetrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		log.Printf("Failed to unmarshal existing metrics: %v", err)
		return
	}
	
	// Restore in-memory state
	m.TotalRequests = metrics.TotalRequests
	m.TotalErrors = metrics.TotalErrors
	m.HourlyStats = metrics.HourlyStats
	
	// Restore endpoint stats with internal maps
	for endpoint, stat := range metrics.Endpoints {
		stat.APIKeysSet = make(map[string]bool)
		stat.ResponseTimes = []float64{} // Don't restore response times, they're just for avg calculation
		m.EndpointStats[endpoint] = stat
	}
	
	// Restore series usage
	m.SeriesUsage = metrics.SeriesUsage
	
	// Restore API key usage with internal maps
	for key, usage := range metrics.APIKeyUsage {
		metric := &APIKeyMetric{
			TotalRequests: usage.TotalRequests,
			EndpointsUsed: make(map[string]bool),
			LastActive:    usage.LastActive,
			ErrorCount:    int64(usage.ErrorRate * float64(usage.TotalRequests)),
		}
		for _, endpoint := range usage.EndpointsUsed {
			metric.EndpointsUsed[endpoint] = true
		}
		m.APIKeyUsage[key] = metric
	}
	
	log.Printf("Loaded existing metrics for %s", m.currentDate)
}

func (m *MetricsRepository) GetHistoricalMetrics(days int) ([]*DailyMetrics, error) {
	var results []*DailyMetrics
	
	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		filename := filepath.Join(m.metricsDir, fmt.Sprintf("daily_metrics_%s.json", date))
		
		data, err := os.ReadFile(filename)
		if err != nil {
			continue // Skip missing files
		}
		
		var metrics DailyMetrics
		if err := json.Unmarshal(data, &metrics); err != nil {
			log.Printf("Failed to unmarshal metrics for %s: %v", date, err)
			continue
		}
		
		results = append(results, &metrics)
	}
	
	return results, nil
}