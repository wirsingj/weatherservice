package weatherclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	userAgent  string
}

const nwsBaseURL = "https://api.weather.gov"

func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: "weatherservice-demo (contact: wirsingj@gmail.com)",
	}
}

func (c *Client) GetTodaysForecast(ctx context.Context, lat, lon float64) (DailyForecast, error) {
	pointsURL := fmt.Sprintf("%s/points/%f,%f", nwsBaseURL, lat, lon)

	// First fetch the points data to get forecast URL
	var points pointsResponse
	if err := c.getJSON(ctx, pointsURL, &points); err != nil {
		return DailyForecast{}, fmt.Errorf("fetch points metadata: %w", err)
	}
	if points.Properties.Forecast == "" {
		return DailyForecast{}, fmt.Errorf("points response missing forecast URL")
	}

	// Fetch forecast data with given points
	var forecast forecastResponse
	if err := c.getJSON(ctx, points.Properties.Forecast, &forecast); err != nil {
		return DailyForecast{}, fmt.Errorf("fetch forecast periods: %w", err)
	}

	// Get today's forecast period entry
	period, err := pickTodayPeriod(forecast.Properties.Periods, time.Now())
	if err != nil {
		return DailyForecast{}, err
	}

	return DailyForecast{
		ShortForecast:   period.ShortForecast,
		Temperature:     period.Temperature,
		TemperatureUnit: period.TemperatureUnit,
	}, nil
}

// getJSON is a general helper function for calling external API
func (c *Client) getJSON(ctx context.Context, url string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/geo+json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("status %d from %s: %s", resp.StatusCode, url, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return fmt.Errorf("decode response JSON: %w", err)
	}
	return nil
}

// pickTodayPeriod looks through periods provided and finds today's period
func pickTodayPeriod(periods []forecastPeriod, now time.Time) (forecastPeriod, error) {
	if len(periods) == 0 {
		return forecastPeriod{}, fmt.Errorf("no forecast periods returned")
	}

	for _, p := range periods {
		start, err := time.Parse(time.RFC3339, p.StartTime)
		if err != nil {
			continue
		}
		localNow := now.In(start.Location())
		if sameDay(start, localNow) {
			return p, nil
		}
	}

	return periods[0], nil
}

// sameDay is a helper function for checking if days match
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
