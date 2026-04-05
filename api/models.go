package api

// API response
type WeatherResponse struct {
	ShortForecast               string `json:"short_forecast"`
	Temperature                 int    `json:"temperature"`
	TemperatureUnit             string `json:"temperature_unit"`
	TemperatureCharacterization string `json:"temperature_characterization"`
}

type errorResponse struct {
	Error string `json:"error"`
}
