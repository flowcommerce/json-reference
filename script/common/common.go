package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/flowcommerce/tools/util"
	"io/ioutil"
	"os"
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
	DefaultCurrency    string `json:"default_currency,omitempty"`
	Languages          []string `json:"languages"`
	Timezones          []string	`json:"timezones"`
}

type Currency struct {
	Name           string `json:"name"`
	Iso_4217_3     string `json:"iso_4217_3"`
	NumberDecimals int    `json:"number_decimals"`
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
	Name                     string `json:"name"`
	Description              string `json:"description"`
	Offset                   int    `json:"offset"`
}

func LoadContinents() []Continent {
	continents := []Continent{}
	err := json.Unmarshal(ReadFile("data/final/continents.json"), &continents)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal continents: %s", err))
	return continents
}

func LoadCountries() []Country {
	countries := []Country{}
	err := json.Unmarshal(ReadFile("data/final/countries.json"), &countries)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal countries: %s", err))
	return countries
}

func LoadCurrencies() []Currency {
	currencies := []Currency{}
	err := json.Unmarshal(ReadFile("data/final/currencies.json"), &currencies)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal currencies: %s", err))
	return currencies
}

func LoadLanguages() []Language {
	languages := []Language{}
	err := json.Unmarshal(ReadFile("data/final/languages.json"), &languages)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal languages: %s", err))
	return languages
}

func LoadTimezones() []Timezone {
	timezones := []Timezone{}
	err := json.Unmarshal(ReadFile("data/final/timezones.json"), &timezones)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal timezones: %s", err))
	return timezones
}

func LoadRegions() []Region {
	regions := []Region{}
	err := json.Unmarshal(ReadFile("data/final/regions.json"), &regions)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal regions: %s", err))
	return regions
}


func ReadFile(path string) []byte {
	file, err := ioutil.ReadFile(path)
	util.ExitIfError(err, fmt.Sprintf("Could not read file %s", path))

	fileStr := string(file)
	return []byte(fileStr)
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

func FilterNonEmpty(vs []string) []string {
	return Filter(vs, func(v string) bool {
		return v != ""
	})
}

func WriteJson(target string, data interface{}) {
	tmp, err := ioutil.TempFile("", "reference-csv-to-json")
	util.ExitIfError(err, "Error creating temporary file")
	defer tmp.Close()

	v, err := json.MarshalIndent(&data, "", "  ")
	util.ExitIfError(err, "Error marshalling record to json")

	w := bufio.NewWriter(tmp)
	_, err = w.Write(v)
	util.ExitIfError(err, "Error writing to tmp file")
	w.Flush()

	fi, err := tmp.Stat()
	util.ExitIfError(err, "Error checking file size")
	if fi.Size() == 0 {
		fmt.Printf("Failed to serialize json to file %s - empty file\n", target)
		os.Exit(1)
	}
	
	err = os.Rename(tmp.Name(), target)
	util.ExitIfError(err, "Error renaming tmp file")
}
