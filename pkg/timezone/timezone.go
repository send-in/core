package timezone

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"strings"
	"sync"
)

//go:embed cities.json
var citiesJSON []byte
//go:embed countries.json
var countriesJSON []byte

type City struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
}

type Country struct {
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

type TimezoneResult struct {
	Location string `json:"location"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
}

var (
	cityIndex  map[string][]City 
	countryISO map[string]string 
	countryTZ  map[string]string
	countryName map[string]string
	initOnce   sync.Once
	loadErr    error
)

func load() {
	var cities []City
	if  err := json.Unmarshal(citiesJSON, &cities); 
		err != nil {
		loadErr = err
		return
	}

	cityIndex = make(map[string][]City, len(cities))
	for _, c := range cities {
		cityIndex[strings.ToLower(c.City)] = append(
			cityIndex[strings.ToLower(c.City)],
			c,
		)
	}

	var byCode map[string]Country
	if  err := json.Unmarshal(countriesJSON, &byCode); 
		err != nil {
		loadErr = err
		return
	}

	countryISO = make(map[string]string, len(byCode))
	countryTZ = make(map[string]string, len(byCode))
	countryName = make(map[string]string, len(byCode))

	for code, entry := range byCode {
		countryISO[strings.ToLower(entry.Name)] = code
		countryTZ[code] = entry.Timezone
		countryName[code] = entry.Name
	}
}

func InferTimezone(location string) (*TimezoneResult, error){
	initOnce.Do(load)
	if loadErr != nil {
		return nil, loadErr
	}

	parts := split(location)
	if len(parts) == 0 {
		return nil, nil
	}

	var city, country string
	switch len(parts) {
		case 1:
			if strings.HasSuffix(strings.ToLower(parts[0]), "area") {
				city = normalizeCity(parts[0])
			} else {
				country = parts[0]
			}
		default: {
			city = normalizeCity(parts[0])
			country = parts[len(parts)-1]
		}
	}

	code := countryISO[strings.ToLower(country)]

	if city != "" {
		for _, c := range cityIndex[strings.ToLower(city)] {
			if code == "" || strings.EqualFold(c.Country, code) {
				return &TimezoneResult{
					Location: location,
					Country:  countryName[c.Country],
					Timezone: c.Timezone,
				}, nil
			}
		}
	}

	if tz := countryTZ[code]; tz != "" {
		return &TimezoneResult{
			Location: location,
			Country:  country,
			Timezone: tz,
		}, nil
	}

	return nil, nil
}

func normalizeCity(s string) string {
	return strings.TrimSpace(
		regexp.
		MustCompile(`(?i)^(Greater\s+)|(\s+(Metropolitan|Metro)?\s*Area$)`).
		ReplaceAllString(s, ""),
	)
}

func split(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}