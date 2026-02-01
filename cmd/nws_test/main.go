// This program tests the custom JSON library by fetching and parsing
// real JSON data from the National Weather Service API.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	json "github.com/adavi068/go-json-library"
)

// NWS API response structures
type PointsResponse struct {
	Properties PointsProperties `json:"properties"`
}

type PointsProperties struct {
	GridID         string `json:"gridId"`
	GridX          int    `json:"gridX"`
	GridY          int    `json:"gridY"`
	Forecast       string `json:"forecast"`
	ForecastHourly string `json:"forecastHourly"`
	City           string `json:"relativeLocation"`
}

type ForecastResponse struct {
	Properties ForecastProperties `json:"properties"`
}

type ForecastProperties struct {
	Updated   string   `json:"updated"`
	Periods   []Period `json:"periods"`
}

type Period struct {
	Number          int     `json:"number"`
	Name            string  `json:"name"`
	Temperature     int     `json:"temperature"`
	TemperatureUnit string  `json:"temperatureUnit"`
	WindSpeed       string  `json:"windSpeed"`
	WindDirection   string  `json:"windDirection"`
	ShortForecast   string  `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}

func fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-json-library-test/1.0")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func main() {
	fmt.Println("=== National Weather Service API JSON Test ===")
	fmt.Println()

	// Test with Washington DC coordinates (38.8894, -77.0352)
	lat, lon := 38.8894, -77.0352
	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", lat, lon)

	fmt.Printf("Fetching location data for (%.4f, %.4f)...\n", lat, lon)

	pointsData, err := fetch(pointsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching points: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Received %d bytes of JSON\n", len(pointsData))
	fmt.Println()

	// Validate the JSON first
	if !json.Valid(pointsData) {
		fmt.Println("ERROR: Invalid JSON received!")
		os.Exit(1)
	}
	fmt.Println("✓ JSON is valid")

	// Parse into a generic map first to demonstrate flexibility
	var rawData map[string]interface{}
	if err := json.Unmarshal(pointsData, &rawData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Successfully parsed into map[string]interface{}")

	// Extract some data from the raw map
	if props, ok := rawData["properties"].(map[string]interface{}); ok {
		fmt.Printf("  Grid ID: %v\n", props["gridId"])
		fmt.Printf("  Grid X: %v\n", props["gridX"])
		fmt.Printf("  Grid Y: %v\n", props["gridY"])

		if relLoc, ok := props["relativeLocation"].(map[string]interface{}); ok {
			if locProps, ok := relLoc["properties"].(map[string]interface{}); ok {
				fmt.Printf("  Location: %v, %v\n", locProps["city"], locProps["state"])
			}
		}
	}
	fmt.Println()

	// Get forecast URL from the response
	forecastURL := ""
	if props, ok := rawData["properties"].(map[string]interface{}); ok {
		if url, ok := props["forecast"].(string); ok {
			forecastURL = url
		}
	}

	if forecastURL == "" {
		fmt.Println("Could not find forecast URL in response")
		os.Exit(1)
	}

	fmt.Printf("Fetching forecast from: %s\n", forecastURL)
	forecastData, err := fetch(forecastURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching forecast: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Received %d bytes of forecast JSON\n", len(forecastData))
	fmt.Println()

	// Parse forecast into struct
	var forecast map[string]interface{}
	if err := json.Unmarshal(forecastData, &forecast); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing forecast JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Successfully parsed forecast data")
	fmt.Println()

	// Display forecast periods
	if props, ok := forecast["properties"].(map[string]interface{}); ok {
		if periods, ok := props["periods"].([]interface{}); ok {
			fmt.Println("=== Weather Forecast ===")
			fmt.Println()

			// Show first 4 periods
			count := 4
			if len(periods) < count {
				count = len(periods)
			}

			for i := 0; i < count; i++ {
				if period, ok := periods[i].(map[string]interface{}); ok {
					name := period["name"]
					temp := period["temperature"]
					unit := period["temperatureUnit"]
					shortForecast := period["shortForecast"]
					wind := period["windSpeed"]
					windDir := period["windDirection"]

					fmt.Printf("%-20s %v°%v\n", name, temp, unit)
					fmt.Printf("  Wind: %v %v\n", wind, windDir)
					fmt.Printf("  %v\n", shortForecast)
					fmt.Println()
				}
			}
		}
	}

	// Test marshaling back to JSON
	fmt.Println("=== Marshal Test ===")
	fmt.Println()

	testData := map[string]interface{}{
		"location": "Washington DC",
		"lat":      lat,
		"lon":      lon,
		"success":  true,
	}

	marshaled, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Marshaled test data:")
	fmt.Println(string(marshaled))
	fmt.Println()
	fmt.Println("✓ All tests passed!")
}
