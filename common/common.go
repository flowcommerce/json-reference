package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/flowcommerce/tools/util"
)

type Carrier struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	TrackingUrl string `json:"tracking_url"`
}

type CarrierService struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Carrier Carrier `json:"carrier"`
}

type Continent struct {
	Name      string   `json:"name"`
	Code      string   `json:"code"`
	Countries []string `json:"countries"`
}

type Country struct {
	Name                 string   `json:"name"`
	Iso_3166_2           string   `json:"iso_3166_2"`
	Iso_3166_3           string   `json:"iso_3166_3"`
	MeasurementSystem    string   `json:"measurement_system"`
	DefaultCurrency      string   `json:"default_currency,omitempty"`
	DefaultLanguage      string   `json:"default_language,omitempty"`
	Languages            []string `json:"languages"`
	Timezones            []string `json:"timezones"`
	DefaultDeliveredDuty string   `json:"default_delivered_duty,omitempty"`
}

type Currency struct {
	Name           string           `json:"name"`
	Iso_4217_3     string           `json:"iso_4217_3"`
	NumberDecimals int              `json:"number_decimals"`
	Symbols        *CurrencySymbols `json:"symbols,omitempty"`
	DefaultLocale  string           `json:"default_locale,omitempty"`
}

type CurrencySymbols struct {
	Primary string `json:"primary"`
	Narrow  string `json:"narrow,omitempty"`
}

type Language struct {
	Name      string   `json:"name"`
	Iso_639_2 string   `json:"iso_639_2"`
	Countries []string `json:"countries"`
	Locales   []string `json:"locales"`
}

type PaymentMethod struct {
	Id      				string              		`json:"id"`
	Type    				string              		`json:"type"`
	Name    				string              		`json:"name"`
	Images  				PaymentMethodImages 		`json:"images"`
	Regions 				[]string            		`json:"regions"`
	Capabilities		[]string		            `json:"capabilities"`
}

type PaymentMethodImages struct {
	Small  PaymentMethodImage `json:"small"`
	Medium PaymentMethodImage `json:"medium"`
	Large  PaymentMethodImage `json:"large"`
}

type PaymentMethodImage struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Province struct {
	Id           string                 `json:"id"`
	Iso_3166_2   string                 `json:"iso_3166_2"`
	Name         string                 `json:"name"`
	Country      string                 `json:"country"`
	ProvinceType string                 `json:"province_type"`
	Translations []LocalizedTranslation `json:"translations,omitempty"`
}

type LocalizedTranslation struct {
	Locale Locale `json:"locale"`
	Name   string `json:"name"`
}

type Region struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	Countries          []string `json:"countries"`
	Currencies         []string `json:"currencies"`
	Languages          []string `json:"languages"`
	MeasurementSystems []string `json:"measurement_systems"`
	Timezones          []string `json:"timezones"`
}

type Timezone struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Offset      int    `json:"offset"`
}

type Locale struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	Country  string        `json:"country"`
	Language string        `json:"language,omitempty"`
	Numbers  LocaleNumbers `json:"numbers"`
}

type LocaleNumbers struct {
	Decimal string `json:"decimal"`
	Group   string `json:"group"`
}

func Carriers() []Carrier {
	carriers := []Carrier{}
	err := json.Unmarshal(readDataFileFromUrl("carriers.json"), &carriers)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal carriers: %s", err))
	return carriers
}

func CarrierServices() []CarrierService {
	carrierServices := []CarrierService{}
	err := json.Unmarshal(readDataFileFromUrl("carrier-services.json"), &carrierServices)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal carrier services: %s", err))
	return carrierServices
}

func Continents() []Continent {
	continents := []Continent{}
	err := json.Unmarshal(readDataFileFromUrl("continents.json"), &continents)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal continents: %s", err))
	return continents
}

func Countries() []Country {
	countries := []Country{}
	err := json.Unmarshal(readDataFileFromUrl("countries.json"), &countries)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal countries: %s", err))
	return countries
}

func Currencies() []Currency {
	currencies := []Currency{}
	err := json.Unmarshal(readDataFileFromUrl("currencies.json"), &currencies)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal currencies: %s", err))
	return currencies
}

func Languages() []Language {
	languages := []Language{}
	err := json.Unmarshal(readDataFileFromUrl("languages.json"), &languages)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal languages: %s", err))
	return languages
}

func Locales() []Locale {
	locales := []Locale{}
	err := json.Unmarshal(readDataFileFromUrl("locales.json"), &locales)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal locales: %s", err))
	return locales
}

func Timezones() []Timezone {
	timezones := []Timezone{}
	err := json.Unmarshal(readDataFileFromUrl("timezones.json"), &timezones)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal timezones: %s", err))
	return timezones
}

func PaymentMethods() []PaymentMethod {
	paymentMethods := []PaymentMethod{}
	err := json.Unmarshal(readDataFileFromUrl("payment-methods.json"), &paymentMethods)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal payment methods: %s", err))
	return paymentMethods
}

func Provinces() []Province {
	provinces := []Province{}
	err := json.Unmarshal(readDataFileFromUrl("provinces.json"), &provinces)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal provinces: %s", err))
	return provinces
}

func Regions() []Region {
	regions := []Region{}
	err := json.Unmarshal(readDataFileFromUrl("regions.json"), &regions)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal regions: %s", err))
	return regions
}

func readDataFileFromUrl(name string) []byte {
	// Add a random query parameter to flush cache in github
	url := fmt.Sprintf("https://raw.githubusercontent.com/flowcommerce/json-reference/master/data/final/%s?r=%i", name, rand.Float64())
	return ReadUrl(url)
}

func ReadUrl(url string) []byte {
	res, err := http.Get(url)
	util.ExitIfError(err, fmt.Sprintf("Could not read url %s", url))

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data
}

func ReadFile(path string) []byte {
	file, err := ioutil.ReadFile(path)
	util.ExitIfError(err, fmt.Sprintf("Could not read file %s", path))

	fileStr := string(file)
	return []byte(fileStr)
}

func EqualsIgnoreCase(text1 string, text2 string) bool {
	if strings.ToUpper(text1) == strings.ToUpper(text2) {
		return true
	}
	return false
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
		if strings.ToUpper(v) == strings.ToUpper(value) {
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

func FormatLocaleId(value string) string {
	formatted := regexp.MustCompile("_").ReplaceAllString(value, "-")
	distinct := []string{}
	for _, v := range strings.Split(formatted, "-") {
		if !ContainsIgnoreCase(distinct, v) {
			distinct = append(distinct, v)
		}
	}

	return strings.Join(distinct, "-")
}

func FormatUnderscore(value string) string {
	distinct := []string{}
	for _, v := range strings.Split(value, " ") {
		if !ContainsIgnoreCase(distinct, v) {
			distinct = append(distinct, v)
		}
	}

	return strings.Join(distinct, "_")
}

func UnsupportedCountryCodes() []string {
	return []string{
		"AF",
		"AGQ",
		"ER",
		"IQ",
		"PS",
		"SD",
		"SS",
		"SR",
		"SY",
		"TJ",
	}
}

// Need to map some currency codes into the ones supported  by
// most payment processors
var remappedCurrencyCodes = map[string]string{
	"AFN": "EUR",
	"ALK": "EUR",
	"BIF": "EUR",
	"BYR": "EUR",
	"CNH": "EUR",
	"CNX": "EUR",
	"CUP": "EUR",
	"CUP,CUC": "EUR",
	"CDF": "EUR",
	"ERN": "EUR",
	"ILR": "EUR",
	"IQD": "EUR",
	"IRR": "EUR",
	"ISJ": "EUR",
	"KPW": "EUR",
	"LRD": "EUR",
	"MGA": "EUR",
	"MKD": "EUR",
	"MRU": "EUR",
	"MVP": "EUR",
	"MZN": "EUR",
	"SDG": "EUR",
	"SRD": "EUR",
	"SSP": "EUR",
	"STN": "EUR",
	"SYP": "EUR",
	"TJS": "EUR",
	"TMT": "EUR",
	"ZWL": "EUR",
}

func RemapCurrencyCodeToSupported(code string) string {
	newCurrency, _ := remappedCurrencyCodes[code]
	if newCurrency == "" {
		return code
	} else {
		return newCurrency
	}
}
