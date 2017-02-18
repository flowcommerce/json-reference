package cleanse

// Reads source files, cleansing and writing all as json to data/1-cleansed

import (
	"github.com/flowcommerce/tools/util"
	"../common"
	"encoding/csv"
	"encoding/json"
        "fmt"
        "io"
        "os"
	"sort"
	"strconv"
	"strings"
)

type Continent struct {
	Code2             string `json:"code2"`
	Code3             string `json:"code3"`
	Name              string `json:"name"`
}

type CountryContinent struct {
	ContinentCode     string `json:"continent"`
	CountryCode       string `json:"country"`
}

type Country struct {
	Name                     string `json:"name"`
	Continent                string `json:"continent"`
	Iso_3166_2               string `json:"iso_3166_2"`
	Iso_3166_3               string `json:"iso_3166_3"`
	Currency                 string `json:"currency"`
}

type Currency struct {
	Name                     string `json:"name"`
	Iso_4217_3               string `json:"iso_4217_3"`
	NumberDecimals           int    `json:"number_decimals"`
}

type Language struct {
	Name                     string `json:"name"`
	Iso_639_2                string `json:"iso_639_2"`
	Countries                []string `json:"countries"`
}

type Timezone struct {
	Name                     string `json:"name"`
	Description              string `json:"description"`
	Offset                   int    `json:"offset"`
}

type CountryTimezone struct {
	TimezoneCode      string `json:"timezone"`
	CountryCode       string `json:"country"`
}

type IncomingLanguages struct {
	LanguageFamilies    []string `json:"languageFamilies"`
	Languages           []IncomingLanguage `json:"languages"`
}

type IncomingLanguage struct {
	Iso_639_2           string `json:"iso639_1"`
	Names               []string `json:"name"`
	Countries           []string `json:"countries"`
}

type IncomingCurrencies struct {
	Main        IncomingCurrenciesMain`json:"main"`
}

type IncomingCurrenciesMain struct {
	EnUsPosix   IncomingCurrenciesPosix `json:"en-US-POSIX"`
}

type IncomingCurrenciesPosix struct {
	Numbers     IncomingCurrenciesNumbers `json:"numbers"`
}

type IncomingCurrenciesNumbers struct {
	Currencies     map[string]IncomingCurrency `json:"currencies"`
}

type IncomingCurrency struct {
	Name           string `json:"displayName"`
	NameSingle     string `json:"displayName-count-one"`
	NameMultiple   string `json:"displayName-count-other"`
	Symbol         string `json:"symbol"`
}

type convertFunction func(records map[string]string) interface{}
type acceptsFunction func(records map[string]string) bool
type idFunction func(records map[string]string) string

func Cleanse() {
	writeJson("data/cleansed/languages.json", readLanguages("data/source/languages.json"))

	unsupportedCurrencyCodes := unsupportedCurrencyCodes()
	unsupportedCountryCodes := unsupportedCountryCodes()

	countriesSource := readCsv("data/source/countries.csv")
	writeJson("data/cleansed/countries.json",
		toObjects(countriesSource,
			func(record map[string]string) bool {
				return record["ISO3166-1-Alpha-2"] != "" && record["ISO3166-1-Alpha-3"] != "" && record["official_name_en"] != "" && !common.ContainsIgnoreCase(unsupportedCountryCodes, record["ISO3166-1-Alpha-3"]) && !common.ContainsIgnoreCase(unsupportedCurrencyCodes, record["ISO4217-currency_alphabetic_code"])
			},
			func(record map[string]string) interface{} {
				return Country {
					Name: record["official_name_en"],
					Iso_3166_2: record["ISO3166-1-Alpha-2"],
					Iso_3166_3: record["ISO3166-1-Alpha-3"],
					Currency: record["ISO4217-currency_alphabetic_code"],
					Continent: record["Continent"],
					
				}
			},
			func(record map[string]string) string {
				return record["official_name_en"]
			},
		),
	)

	currencies := readCurrencies("data/source/currencies.json")
	fmt.Println(len(currencies))
	writeJson("data/cleansed/currencies.json",
		toObjects(countriesSource,
			func(record map[string]string) bool {
				return record["ISO4217-currency_name"] != "" && record["ISO4217-currency_alphabetic_code"] != "" && !common.ContainsIgnoreCase(unsupportedCurrencyCodes, record["ISO4217-currency_alphabetic_code"])
			},
			func(record map[string]string) interface{} {
				n := record["ISO4217-currency_number_decimals"]
				var numberDecimals int64
				var err error
				if (n != "") {
					numberDecimals, err = strconv.ParseInt(n, 10, 0)
					util.ExitIfError(err, fmt.Sprintf("Error parsing int[%s]", n))
				}

				return Currency{
					Name: record["ISO4217-currency_name"],
					Iso_4217_3: record["ISO4217-currency_alphabetic_code"],
					NumberDecimals: int(numberDecimals),
				}
			},
			func(record map[string]string) string {
				return record["ISO4217-currency_name"]
			},
		),
	)

	writeJson("data/cleansed/country-continents.json",
		toObjects(readCsv("data/source/country-continents.csv"),
			func(record map[string]string) bool {
				return record["continent code"] != "" && record["continent code"] != "--"
			},
			func(record map[string]string) interface{} {
				return CountryContinent{
					ContinentCode: record["continent code"],
					CountryCode: record["iso 3166 country"],
				}
			},
			func(record map[string]string) string {
				return record["continent code"] + record["iso 3166 country"]
			},
		),
	)	

	writeJson("data/cleansed/timezones.json", loadTimezonesFromPath("data/original/timezones.json"))

	writeJson("data/cleansed/country-timezones.json",
		toObjects(readCsv("data/original/country-timezones.csv"),
			func(record map[string]string) bool {
				return record["country"] != "" && record["timezone"] != ""
			},
			func(record map[string]string) interface{} {
				return CountryTimezone{
					TimezoneCode: record["timezone"],
					CountryCode: strings.ToUpper(record["country"]),
				}
			},
			func(record map[string]string) string {
				return record["timezone"] + record["country"]
			},
		),
	)	
	
	writeJson("data/cleansed/languages.json", filterLanguages(readLanguages("data/source/languages.json")))
}

func writeJson(target string, objects interface{}) {
	fmt.Printf("Writing %s\n", target)
	common.WriteJson(target, objects)
}
	
// readCsv Reads a CSV file, returning a list of map[string]string objects
func readCsvWithHeaders(file string, headers []string) []map[string]string {
	input, err := os.Open(file)
	util.ExitIfError(err, fmt.Sprintf("Error opening file %s", file))

        r := csv.NewReader(input)
	var all []map[string]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		util.ExitIfError(err, fmt.Sprintf("Error processing csv file %s", file))
		all = append(all, toMap(headers, record))
	}

	return all
}

// readCsv Reads a CSV file, assuming first line is header row, returning a list of map[string]string objects
func readCsv(file string) []map[string]string {
	input, err := os.Open(file)
	util.ExitIfError(err, fmt.Sprintf("Error opening file %s", file))

        r := csv.NewReader(input)
	var headers []string
	first := true
	var all []map[string]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		util.ExitIfError(err, fmt.Sprintf("Error processing csv file %s", file))
		if (first) {
			headers = record
			first = false
		} else {
			all = append(all, toMap(headers, record))
		}
	}

	return all
}

func toMap(headers []string, record []string) map[string]string {
	row := make(map[string]string)

	for i, v := range record {
		if (v != "") {
			row[headers[i]] = v
		}
	}

	return row
}


func readLanguages(file string) []Language {
	lang := IncomingLanguages{}
	err := json.Unmarshal(common.ReadFile(file), &lang)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshall languages: %s", err))

	languages := []Language{}
	
	for _, l := range(lang.Languages) {
		name := l.Names[0]
		if len(l.Iso_639_2) > 0 && name != "" {
			languages = append(languages, Language{
				Name: name,
				Iso_639_2: l.Iso_639_2,
				Countries: l.Countries,
			})
		}
	}

	return languages
}

func readCurrencies(file string) []Currency {
	data := IncomingCurrencies{}
	err := json.Unmarshal(common.ReadFile(file), &data)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshall currencies: %s", err))

	currencies := []Currency{}

	for _, c := range(data.Main.EnUsPosix.Numbers.Currencies) {
		currencies = append(currencies, Currency{
			Name: c.Name,
			Iso_639_2: c.Symbol,
		})
	}

	return currencies
}

/**
 * filterLanguages takes only the languages that have at least one country assigned to them
 */
func filterLanguages(languages []Language) []Language {
	final := []Language{}
	
	for _, l := range(languages) {
		if (l.Countries != nil && len(l.Countries) > 0) {
			final = append(final, l)
		}
	}

	return final
}

func toObjects(records []map[string]string, accepts acceptsFunction, f convertFunction, id idFunction) []interface{} {
	added := map[string]interface{} {}
	for _, v := range records {
		if (accepts(v)) {
			id := strings.ToUpper(id(v))
			if added[id] == nil {
				added[id] = f(v)
			}
		}
	}

	return sortObjects(added)
}

func LoadCountryContinents() []CountryContinent {
	countryContinents := []CountryContinent{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/country-continents.json"), &countryContinents)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal country continents: %s", err))
	return countryContinents
}

func LoadContinents() []Continent {
	return []Continent{
		Continent{
			Name: "Africa",
			Code2: "AF",
			Code3: "AFR",
		},
		Continent{
			Name: "Antarctica",
			Code2: "AN",
			Code3: "ANT",
		},
		Continent{
			Name: "Asia",
			Code2: "AS",
			Code3: "ASI",
		},
		Continent{
			Name: "Europe",
			Code2: "EU",
			Code3: "EUR",
		},
		Continent{
			Name: "North America",
			Code2: "NA",
			Code3: "NOA",
		},
		Continent{
			Name: "Oceania",
			Code2: "OC",
			Code3: "OCE",
		},
		Continent{
			Name: "South America",
			Code2: "SA",
			Code3: "SOA",
		},
	}
}

func LoadCountries() []Country {
	countries := []Country{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/countries.json"), &countries)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal countries: %s", err))
	return countries
}

func LoadCurrencies() []Currency {
	currencies := []Currency{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/currencies.json"), &currencies)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal currencies: %s", err))
	return currencies
}

func LoadLanguages() []Language {
	languages := []Language{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/languages.json"), &languages)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal languages: %s", err))
	return languages
}

func LoadTimezones() []Timezone {
	return loadTimezonesFromPath("data/cleansed/timezones.json")
}

func loadTimezonesFromPath(path string) []Timezone {
	timezones := []Timezone{}
	err := json.Unmarshal(common.ReadFile(path), &timezones)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal timezones: %s", err))
	return timezones
}

func LoadCountryTimezones() []CountryTimezone {
	countryTimezones := []CountryTimezone{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/country-timezones.json"), &countryTimezones)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal country timezones: %s", err))
	return countryTimezones
}

func unsupportedCountryCodes() []string {
	return []string {
		"AFG",
		"AGO",
		"ATF",
		"BDI",
		"BLR",
		"BVT",
		"CCK",
		"COD",
		"CUB",
		"CXR",
		"ERI",
		"FRO",
		"HMD",
		"IOT",
		"IRN",
		"IRQ",
		"LBR",
		"MDG",
		"MKD",
		"MMR",
		"MOZ",
		"PSE",
		"SDN",
		"SGS",
		"SUR",
		"SYR",
		"TJK",
		"TKM",
		"UMI",
		"ZWE",
	}
}

func unsupportedCurrencyCodes() []string {
	return []string {
		"AFN",
		"AOA",
		"BIF",
		"BYR",
		"CUP",
		"ERN",
		"IQD",
		"IRR",
		"KPW",
		"LRD",
		"MGA",
		"MKD",
		"MMK",
		"MZN",
		"SDG",
		"SRD",
		"SSP",
		"SYP",
		"TJS",
		"TMT",
		"ZWL",
	}	
}

func sortObjects(data map[string]interface{}) []interface{} {
	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var all []interface{}
	for _, key := range(keys) {
		all = append(all, data[key])
	}
	return all
}
