package common

import (
	"encoding/json"
	"fmt"
	"github.com/flowcommerce/tools/util"
	"io/ioutil"
	"sort"
	"strings"
)

type Continent struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	Countries  []string `json:"countries"`
}

type Country struct {
	Name               string `json:"name"`
	Iso_3166_2         string `json:"iso_3166_2"`
	Iso_3166_3         string `json:"iso_3166_3"`
	MeasurementSystem  string `json:"measurement_system"`
	DefaultCurrency    string `json:"default_currency"`
	Languages          []string `json:"languages"`
	Timezones          []string	`json:"timezones"`
}

type Currency struct {
	Name           string `json:"name"`
	Iso_4217_3     string `json:"iso_4217_3"`
	NumberDecimals int `json:"number_decimals"`
}

type Language struct {
	Name               string `json:"name"`
	Iso_639_2          string `json:"iso_639_2"`
	Countries          []string `json:"countries"`
}

type Region struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Countries           []string `json:"countries"`
	Currencies          []string `json:"currencies"`
	Languages           []string `json:"languages"`
	MeasurementSystems  []string `json:"measurement_systems"`
	Timezones           []string `json:"timezones"`
}

type Timezone struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Offset      int `json:"offset"`
}

type CountryTimezone struct {
	Code        string `json:"code"`
	Timezones   []string `json:"timezones"`
}

func LoadContinents() []Continent {
	continents := []Continent{}
	err := json.Unmarshal(ReadFile("data/continents.json"), &continents)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal continents: %s", err))
	return continents
}

func LoadCountries() []Country {
	countries := []Country{}
	err := json.Unmarshal(ReadFile("data/countries.json"), &countries)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal countries: %s", err))
	return countries
}

func LoadUnsupportedCountryAndRegionCodes() []string {
	codes := []string{}
	err := json.Unmarshal(ReadFile("data/countries_regions_unsupported.json"), &codes)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal unsupported country and region codes: %s", err))
	return codes
}

func LoadUnsupportedCurrencyCodes() []string {
	codes := []string{}
	err := json.Unmarshal(ReadFile("data/currencies_unsupported.json"), &codes)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal unsupported currency codes: %s", err))
	return codes
}

func LoadCurrencies() []Currency {
	currencies := []Currency{}
	err := json.Unmarshal(ReadFile("data/currencies.json"), &currencies)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal currencies: %s", err))
	return currencies
}

func LoadLanguages() []Language {
	languages := []Language{}
	err := json.Unmarshal(ReadFile("data/languages.json"), &languages)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal languages: %s", err))
	return languages
}

func LoadTimezones() []Timezone {
	timezones := []Timezone{}
	err := json.Unmarshal(ReadFile("data/timezones.json"), &timezones)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal timezones: %s", err))
	return timezones
}

func LoadRegions() []Region {
	regions := []Region{}
	err := json.Unmarshal(ReadFile("data/regions.json"), &regions)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal regions: %s", err))
	return regions
}

func LoadCountryTimezones() []CountryTimezone {
	countryTimezones := []CountryTimezone{}
	err := json.Unmarshal(ReadFile("data/country_timezones.json"), &countryTimezones)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal country timezones: %s", err))
	return countryTimezones
}

func ReadFile(path string) []byte {
	file, err := ioutil.ReadFile(path)
	util.ExitIfError(err, fmt.Sprintf("Could not read file %s", path))

	fileStr := string(file)
	return []byte(fileStr)
}

/**
 * LanguagesForCountry filters the list of languages to those spoken in
 * the given country, returning a unique sorted list of language ISO-639-2
 * codes
 */
func LanguagesForCountry(languages []Language, iso31663 string) []string {
	codes := []string{}

	for _, l := range(languages) {
		if Contains(l.Countries, iso31663) && !Contains(codes, l.Iso_639_2) {
			codes = append(codes, l.Iso_639_2)
		}
	}
	sort.Strings(codes)

	return codes
}

func TimezonesForCountry(countryTimezones []CountryTimezone, iso31663 string) []string {
	codes := []string{}

	for _, ctz := range(countryTimezones) {
		if ctz.Code == iso31663 {
			codes = ctz.Timezones
		}
	}
	sort.Strings(codes)

	return codes
}

func Contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsIgnoreCase(list []string, value string) bool {
	for _, v := range list {
		if (strings.ToUpper(v) == strings.ToUpper(value)) {
			return true
		}
	}
	return false
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
