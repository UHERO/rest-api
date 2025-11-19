package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/UHERO/rest-api/data"
	"github.com/gorilla/mux"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(metricsRepo *data.MetricsRepository) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		
		// Extract API key from Authorization header
		apiKey := "unknown"
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			// Truncate for privacy - only keep first 8 characters
			if len(apiKey) > 8 {
				apiKey = apiKey[:8] + "..."
			}
		}
		
		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
		
		// Extract series identifier from request
		seriesID := extractSeriesIdentifier(r)
		
		// Call next handler
		next(wrapped, r)
		
		// Calculate response time
		duration := time.Since(start)
		responseTimeMs := float64(duration.Nanoseconds()) / 1000000.0
		
		// Determine if this was an error
		isError := wrapped.statusCode >= 400
		
		// Get normalized endpoint path
		endpoint := normalizeEndpoint(r.URL.Path)
		
		// Record metrics
		metricsRepo.RecordRequest(endpoint, apiKey, seriesID, responseTimeMs, isError)
	}
}

func extractSeriesIdentifier(r *http.Request) string {
	vars := mux.Vars(r)
	
	// Try to get series ID
	if id, exists := vars["id"]; exists {
		return "id:" + id
	}
	
	// Try to get series name
	if name, exists := vars["name"]; exists {
		// Replace URL-encoded characters back
		name = strings.Replace(name, "-", "&", -1)
		return "name:" + name
	}
	
	// Check query parameters as fallback
	if id := r.URL.Query().Get("id"); id != "" {
		return "id:" + id
	}
	
	if name := r.URL.Query().Get("name"); name != "" {
		name = strings.Replace(name, "-", "&", -1)
		return "name:" + name
	}
	
	return ""
}

func normalizeEndpoint(path string) string {
	// Normalize paths to group similar endpoints together
	// Remove trailing slashes
	path = strings.TrimSuffix(path, "/")
	
	// Group common endpoint patterns
	if strings.HasPrefix(path, "/v1/series/siblings") {
		return "/v1/series/siblings"
	}
	if strings.HasPrefix(path, "/v1/series/observations") {
		return "/v1/series/observations"
	}
	if strings.HasPrefix(path, "/v1/series") {
		return "/v1/series"
	}
	if strings.HasPrefix(path, "/v1/categories") {
		return "/v1/categories"
	}
	if strings.HasPrefix(path, "/v1/search") {
		return "/v1/search"
	}
	if strings.HasPrefix(path, "/v1/geo") {
		return "/v1/geo"
	}
	if strings.HasPrefix(path, "/v1/package") {
		return "/v1/package"
	}
	if strings.HasPrefix(path, "/v1/census") {
		return "/v1/census"
	}
	
	return path
}

// MetricsSummary represents the summary data returned by the metrics endpoint
type MetricsSummary struct {
	Period            string                    `json:"period"`
	Summary           *MetricsSummaryData       `json:"summary"`
	TopEndpoints      []*EndpointSummary        `json:"top_endpoints"`
	TopSeries         []*SeriesSummary          `json:"top_series"`
	ActiveAPIKeys     int                       `json:"active_api_keys"`
	HourlyDistribution [24]int64               `json:"hourly_distribution"`
}

type MetricsSummaryData struct {
	TotalRequests     int64   `json:"total_requests"`
	UniqueAPIKeys     int     `json:"unique_api_keys"`
	AvgResponseTimeMs float64 `json:"avg_response_time_ms"`
	ErrorRate         float64 `json:"error_rate"`
}

type EndpointSummary struct {
	Path     string `json:"path"`
	Requests int64  `json:"requests"`
}

type SeriesSummary struct {
	Identifier string `json:"identifier"`
	Requests   int64  `json:"requests"`
}

func GetMetrics(metricsRepo *data.MetricsRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := metricsRepo.GetCurrentMetrics()
		
		// Calculate overall average response time
		totalResponseTime := 0.0
		totalRequests := int64(0)
		for _, endpoint := range metrics.Endpoints {
			totalResponseTime += endpoint.AvgResponseTimeMs * float64(endpoint.TotalRequests)
			totalRequests += endpoint.TotalRequests
		}
		avgResponseTime := 0.0
		if totalRequests > 0 {
			avgResponseTime = totalResponseTime / float64(totalRequests)
		}
		
		// Calculate error rate
		errorRate := 0.0
		if metrics.TotalRequests > 0 {
			errorRate = float64(metrics.TotalErrors) / float64(metrics.TotalRequests)
		}
		
		// Get top endpoints
		type endpointWithRequests struct {
			path     string
			requests int64
		}
		var endpoints []endpointWithRequests
		for path, stat := range metrics.Endpoints {
			endpoints = append(endpoints, endpointWithRequests{path, stat.TotalRequests})
		}
		
		// Sort endpoints by request count (simple bubble sort since we don't expect many endpoints)
		for i := 0; i < len(endpoints)-1; i++ {
			for j := 0; j < len(endpoints)-i-1; j++ {
				if endpoints[j].requests < endpoints[j+1].requests {
					endpoints[j], endpoints[j+1] = endpoints[j+1], endpoints[j]
				}
			}
		}
		
		topEndpoints := make([]*EndpointSummary, 0, len(endpoints))
		for i, ep := range endpoints {
			if i >= 10 { // Limit to top 10
				break
			}
			topEndpoints = append(topEndpoints, &EndpointSummary{
				Path:     ep.path,
				Requests: ep.requests,
			})
		}
		
		// Get top series
		type seriesWithRequests struct {
			identifier string
			requests   int64
		}
		var series []seriesWithRequests
		for identifier, stat := range metrics.SeriesUsage {
			series = append(series, seriesWithRequests{identifier, stat.Requests})
		}
		
		// Sort series by request count
		for i := 0; i < len(series)-1; i++ {
			for j := 0; j < len(series)-i-1; j++ {
				if series[j].requests < series[j+1].requests {
					series[j], series[j+1] = series[j+1], series[j]
				}
			}
		}
		
		topSeries := make([]*SeriesSummary, 0, len(series))
		for i, s := range series {
			if i >= 10 { // Limit to top 10
				break
			}
			topSeries = append(topSeries, &SeriesSummary{
				Identifier: s.identifier,
				Requests:   s.requests,
			})
		}
		
		summary := &MetricsSummary{
			Period: metrics.Date,
			Summary: &MetricsSummaryData{
				TotalRequests:     metrics.TotalRequests,
				UniqueAPIKeys:     metrics.UniqueAPIKeys,
				AvgResponseTimeMs: avgResponseTime,
				ErrorRate:         errorRate,
			},
			TopEndpoints:       topEndpoints,
			TopSeries:          topSeries,
			ActiveAPIKeys:      metrics.UniqueAPIKeys,
			HourlyDistribution: metrics.HourlyStats,
		}
		
		WriteJSONResponse(w, summary)
	}
}

func GetHistoricalMetrics(metricsRepo *data.MetricsRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		daysStr := r.URL.Query().Get("days")
		days := 7 // default to 7 days
		
		if daysStr != "" {
			if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 30 {
				days = d
			}
		}
		
		historicalMetrics, err := metricsRepo.GetHistoricalMetrics(days)
		if err != nil {
			http.Error(w, "Failed to retrieve historical metrics", http.StatusInternalServerError)
			return
		}
		
		WriteJSONResponse(w, historicalMetrics)
	}
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, `{"error": "failed to marshal JSON"}`, http.StatusInternalServerError)
		return
	}
	
	w.Write(jsonData)
}