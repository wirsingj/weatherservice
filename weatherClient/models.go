package weatherclient

type DailyForecast struct {
	ShortForecast   string `json:"short_forecast"`
	Temperature     int    `json:"temperature"`
	TemperatureUnit string `json:"temperature_unit"`
}

type pointsResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

type forecastResponse struct {
	Properties struct {
		Periods []forecastPeriod `json:"periods"`
	} `json:"properties"`
}

type forecastPeriod struct {
	Name            string `json:"name"`
	StartTime       string `json:"startTime"`
	Temperature     int    `json:"temperature"`
	TemperatureUnit string `json:"temperatureUnit"`
	ShortForecast   string `json:"shortForecast"`
}
