package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	weatherclient "weatherservice/weatherClient"
)

type handler struct {
	weatherClient *weatherclient.Client
}

func NewServer(addr string, weatherClient *weatherclient.Client) *http.Server {
	mux := http.NewServeMux()
	h := &handler{
		weatherClient: weatherClient,
	}
	mux.HandleFunc("/weather", h.handleGetWeather)
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

// handleGetWeather is the handler for getting weather data
// responsible for parsing coordinates parameters, calling the weather client, and building the response.
func (h *handler) handleGetWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	lat, lon, err := parseCoordinates(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	forecast, err := h.weatherClient.GetTodaysForecast(r.Context(), lat, lon)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, errorResponse{Error: fmt.Sprintf("unable to fetch forecast: %v", err)})
		return
	}

	response := WeatherResponse{
		ShortForecast:               forecast.ShortForecast,
		Temperature:                 forecast.Temperature,
		TemperatureUnit:             forecast.TemperatureUnit,
		TemperatureCharacterization: characterizeTemperature(forecast.Temperature, forecast.TemperatureUnit),
	}

	writeJSON(w, http.StatusOK, response)
}

// parseCoordinates parses and verifies lat and long parameters
func parseCoordinates(r *http.Request) (float64, float64, error) {
	latRaw := strings.TrimSpace(r.URL.Query().Get("lat"))
	lonRaw := strings.TrimSpace(r.URL.Query().Get("lon"))

	if latRaw == "" || lonRaw == "" {
		return 0, 0, fmt.Errorf("lat and lon query params are required")
	}

	lat, err := strconv.ParseFloat(latRaw, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid lat value")
	}

	lon, err := strconv.ParseFloat(lonRaw, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid lon value")
	}

	// Check if lat and long are valid
	if lat < -90 || lat > 90 {
		return 0, 0, fmt.Errorf("lat must be between -90 and 90")
	}
	if lon < -180 || lon > 180 {
		return 0, 0, fmt.Errorf("lon must be between -180 and 180")
	}

	return lat, lon, nil
}

// characterizeTemperature converts to F if needed, and maps the temp to hot, cold, or moderate.
func characterizeTemperature(temp int, unit string) string {
	// convert to F if needed
	tempF := float64(temp)
	if strings.EqualFold(unit, "C") {
		tempF = (tempF * 9 / 5) + 32
	}

	// characterize temperature
	switch {
	case tempF >= 85:
		return "hot"
	case tempF <= 50:
		return "cold"
	default:
		return "moderate"
	}
}

// writeJSON is a helper function to build and write a JSON response with the given status code and value.
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"encode response: %v"}`, err), http.StatusInternalServerError)
	}
}
