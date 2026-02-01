package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	json "github.com/adavi068/go-json-library"
)

func fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "go-json-library-example/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
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
	// Fetch location data for Washington DC
	lat, lon := 38.8894, -77.0352
	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", lat, lon)

	fmt.Println("=== NWS API Dynamic JSON Example ===")
	fmt.Printf("Fetching weather data for (%.4f, %.4f)...\n\n", lat, lon)

	// Fetch and load points data
	pointsData, err := fetch(pointsURL)
	if err != nil {
		log.Fatal(err)
	}

	points, err := json.LoadBytes(pointsData)
	if err != nil {
		log.Fatal(err)
	}

	// Extract location info using dot notation
	fmt.Println("=== Location Info ===")
	fmt.Println("Grid ID:", points.Get("properties.gridId").String())
	fmt.Println("Grid X:", points.Get("properties.gridX").Int())
	fmt.Println("Grid Y:", points.Get("properties.gridY").Int())
	fmt.Println("City:", points.Get("properties.relativeLocation.properties.city").String())
	fmt.Println("State:", points.Get("properties.relativeLocation.properties.state").String())

	// Get forecast URL and fetch forecast
	forecastURL := points.Get("properties.forecast").String()
	fmt.Println("\nForecast URL:", forecastURL)

	forecastData, err := fetch(forecastURL)
	if err != nil {
		log.Fatal(err)
	}

	forecast, err := json.LoadBytes(forecastData)
	if err != nil {
		log.Fatal(err)
	}

	// Display forecast periods
	fmt.Println("\n=== Weather Forecast ===")
	periods := forecast.Get("properties.periods")
	fmt.Printf("Found %d forecast periods\n\n", periods.Len())

	// Show first 4 periods using array access
	for i := 0; i < 4 && i < periods.Len(); i++ {
		p := periods.Index(i)
		fmt.Printf("%-15s %d°%s\n",
			p.Get("name").String(),
			p.Get("temperature").Int(),
			p.Get("temperatureUnit").String())
		fmt.Printf("  Wind: %s %s\n",
			p.Get("windSpeed").String(),
			p.Get("windDirection").String())
		fmt.Printf("  %s\n\n", p.Get("shortForecast").String())
	}

	// Demonstrate iteration with ForEach
	fmt.Println("=== All Period Names (using ForEach) ===")
	periods.ForEach(func(idx string, p *json.Value) {
		fmt.Printf("  [%s] %s\n", idx, p.Get("name").String())
	})

	// Show available keys in a period
	fmt.Println("\n=== Available Fields in Period ===")
	fmt.Println(periods.Index(0).Keys())

	// Pretty print one period
	fmt.Println("\n=== First Period (Pretty Printed) ===")
	fmt.Println(periods.Index(0).Pretty())
}
