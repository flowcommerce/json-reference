package cleanse

// Reads source files, cleansing and writing all as json to data/1-cleansed

import (
	"../common"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/bradfitz/slice"
	"github.com/flowcommerce/tools/util"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Continent struct {
	Code2 string `json:"code2"`
	Code3 string `json:"code3"`
	Name  string `json:"name"`
}

type CountryContinent struct {
	ContinentCode string `json:"continent"`
	CountryCode   string `json:"country"`
}

type CountryDuty struct {
	CountryCode   string `json:"country"`
	DeliveredDuty string `json:"duty"`
}

type Country struct {
	Name       string `json:"name"`
	Continent  string `json:"continent"`
	Iso_3166_2 string `json:"iso_3166_2"`
	Iso_3166_3 string `json:"iso_3166_3"`
	Currency   string `json:"currency"`
}

type Currency struct {
	Name           string `json:"name"`
	Iso_4217_3     string `json:"iso_4217_3"`
	NumberDecimals int    `json:"number_decimals"`
}

type CurrencyLocale struct {
	CurrencyCode string `json:"currency"`
	LocaleId     string `json:"locale"`
}

type Language struct {
	Name      string   `json:"name"`
	Iso_639_2 string   `json:"iso_639_2"`
	Countries []string `json:"countries"`
	Locales   []string `json:"locales"`
}

type LocaleName struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Province struct {
	Iso_3166_2   string `json:"iso_3166_2"`
	Name         string `json:"name"`
	CountryCode  string `json:"country"`
	ProvinceType string `json:"province_type"`
}

type Timezone struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Offset      int    `json:"offset"`
}

type CountryTimezone struct {
	TimezoneCode string `json:"timezone"`
	CountryCode  string `json:"country"`
}

type IncomingLanguages struct {
	LanguageFamilies []string           `json:"languageFamilies"`
	Languages        []IncomingLanguage `json:"languages"`
}

type IncomingLanguage struct {
	Iso_639_2 string                   `json:"iso639_1"`
	Names     []string                 `json:"name"`
	Countries []string                 `json:"countries"`
	Locales   []IncomingLanguageLocale `json:"langCultureMs"`
}

type IncomingLanguageLocale struct {
	Id   string `json:"langCultureName"`
	Name string `json:"displayName"`
}

type IncomingCurrency struct {
	Name           string `json:"name"`
	Iso_4217_3     string `json:"iso_4217_3"`
	NumberDecimals int    `json:"number_decimals"`
}

type IncomingNumbers struct {
	Main map[string]IncomingNumbersMain `json:"main"`
}

type IncomingNumbersMain struct {
	Identity CldrIdentity           `json:"identity"`
	Numbers  IncomingNumbersNumbers `json:"numbers"`
}

type CldrIdentity struct {
	Language  string `json:"language"`
	Territory string `json:"territory"`
}

type IncomingNumbersNumbers struct {
	Symbols IncomingNumbersSymbols `json:"symbols-numberSystem-latn"`
}

type IncomingNumbersSymbols struct {
	Decimal string `json:"decimal"`
	Group   string `json:"group"`
}

type Number struct {
	Country    string     `json:"country"`
	Language   string     `json:"language"`
	Separators Separators `json:"separators"`
}

type PaymentMethod struct {
	Id           string   `json:"id"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	SmallWidth   int      `json:"small_width"`
	SmallHeight  int      `json:"small_height"`
	MediumWidth  int      `json:"medium_width"`
	MediumHeight int      `json:"medium_height"`
	LargeWidth   int      `json:"large_width"`
	LargeHeight  int      `json:"large_height"`
	Regions      []string `json:"regions"`
}

type Separators struct {
	Decimal string `json:"decimal"`
	Group   string `json:"group"`
}

type CldrCurrencies struct {
	Main map[string]CldrCurrenciesMain `json:"main"`
}

type CldrCurrenciesMain struct {
	Identity CldrIdentity          `json:"identity"`
	Numbers  CldrCurrenciesNumbers `json:"numbers"`
}

type CldrCurrenciesNumbers struct {
	Currencies map[string]CldrIncomingCurrency `json:"currencies"`
}

type CldrIncomingCurrency struct {
	Name            string `json:"displayName"`
	NameCountOne    string `json:"displayName-count-one"`
	NameCountOther  string `json:"displayName-count-other"`
	Symbol          string `json:"symbol,omitempty"`
	SymbolAltNarrow string `json:"symbol-alt-narrow,omitempty"`
}

type CldrCurrency struct {
	Iso_4217_3 string          `json:"iso_4217_3"`
	Symbols    CurrencySymbols `json:"symbols"`
}

type CurrencySymbols struct {
	Primary string `json:"primary"`
	Narrow  string `json:"narrow,omitempty"`
}

type convertFunction func(records map[string]string) interface{}
type acceptsFunction func(records map[string]string) bool
type idFunction func(records map[string]string) string

func Cleanse() {
	languages, localeNames := readLanguages("data/source/languages.json")
	writeJson("data/cleansed/languages.json", languages)
	writeJson("data/cleansed/locale-names.json", localeNames)

	unsupportedCurrencyCodes := common.UnsupportedCurrencyCodes()
	unsupportedCountryCodes := common.UnsupportedCountryCodes()

	countriesSource := readCsv("data/source/countries.csv")
	writeJson("data/cleansed/countries.json",
		toObjects(countriesSource,
			func(record map[string]string) bool {
				return record["ISO3166-1-Alpha-2"] != "" && record["ISO3166-1-Alpha-3"] != "" && !common.ContainsIgnoreCase(unsupportedCountryCodes, record["ISO3166-1-Alpha-3"]) && !common.ContainsIgnoreCase(unsupportedCurrencyCodes, record["ISO4217-currency_alphabetic_code"])
			},
			func(record map[string]string) interface{} {
				iso3 := record["ISO3166-1-Alpha-3"]
				currency := record["ISO4217-currency_alphabetic_code"]
				if currency == "" {
					if iso3 == "CZE" {
						currency = "CZK"
					} else if iso3 == "HKG" {
						currency = "HKD"
					} else if iso3 == "TWN" {
						currency = "TWD"
					} else {
						fmt.Printf("Country %s does not have a currency\n", iso3)
						os.Exit(1)
					}
				}

				return Country{
					Name:       countryName(record),
					Iso_3166_2: record["ISO3166-1-Alpha-2"],
					Iso_3166_3: iso3,
					Currency:   currency,
					Continent:  record["Continent"],
				}
			},
			func(record map[string]string) string {
				return countryName(record)
			},
		),
	)

	numbers := loadCldrNumbers("cldr-numbers-full/main")
	writeJson("data/cleansed/numbers.json", numbers)

	currencySymbols := readCurrencySymbols("data/source/cldr-currencies.json")
	writeJson("data/cleansed/currency-symbols.json", currencySymbols)

	currencies := readCurrencies("data/original/currencies.json")
	writeJson("data/cleansed/currencies.json", currencies)

	writeJson("data/cleansed/country-duties.json",
		toObjects(readCsv("data/original/country-duties.csv"),
			func(record map[string]string) bool {
				return record["duty"] != ""
			},
			func(record map[string]string) interface{} {
				return CountryDuty{
					CountryCode:   record["country"],
					DeliveredDuty: record["duty"],
				}
			},
			func(record map[string]string) string {
				return record["country"] + record["duty"]
			},
		),
	)

	writeJson("data/cleansed/provinces.json",
		toObjects(readCsv("data/original/provinces.csv"),
			func(record map[string]string) bool {
				return record["province"] != ""
			},
			func(record map[string]string) interface{} {
				return Province{
					Iso_3166_2:   record["province"],
					Name:         record["name"],
					CountryCode:  record["country"],
					ProvinceType: record["type"],
				}
			},
			func(record map[string]string) string {
				return record["country"] + record["province"]
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
					CountryCode:   record["iso 3166 country"],
				}
			},
			func(record map[string]string) string {
				return record["continent code"] + record["iso 3166 country"]
			},
		),
	)

	writeJson("data/cleansed/payment-methods.json",
		toObjects(readCsv("data/original/payment-methods.csv"),
			func(record map[string]string) bool {
				return record["id"] != ""
			},
			func(record map[string]string) interface{} {
				return PaymentMethod{
					Id:           record["id"],
					Type:         record["type"],
					Name:         record["name"],
					SmallWidth:   toInt32(record["small_width"]),
					SmallHeight:  toInt32(record["small_height"]),
					MediumWidth:  toInt32(record["medium_width"]),
					MediumHeight: toInt32(record["medium_height"]),
					LargeWidth:   toInt32(record["large_width"]),
					LargeHeight:  toInt32(record["large_height"]),
					Regions:      strings.Split(record["regions"], " "),
				}
			},
			func(record map[string]string) string {
				return record["id"]
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
					CountryCode:  strings.ToUpper(record["country"]),
				}
			},
			func(record map[string]string) string {
				return record["timezone"] + record["country"]
			},
		),
	)

	writeJson("data/cleansed/currency-locales.json",
		toObjects(readCsv("data/original/currency-locales.csv"),
			func(record map[string]string) bool {
				return record["currency"] != "" && record["locale"] != ""
			},
			func(record map[string]string) interface{} {
				return CurrencyLocale{
					CurrencyCode: record["currency"],
					LocaleId:     record["locale"],
				}
			},
			func(record map[string]string) string {
				return record["currency"] + record["locale"]
			},
		),
	)
}

func countryName(record map[string]string) string {
	overrides := map[string]string{
		"Bolivia (Plurinational State of)":   "Bolivia",
		"Micronesia (Federated States of)":   "Micronesia",
		"Saint Martin (French part)":         "Saint Martin",
		"Sint Maarten (Dutch part)":          "Sint Maarten",
		"Venezuela (Bolivarian Republic of)": "Venezuela",
		"Falkland Islands (Malvinas)":        "Falkland Islans",
	}

	name := record["official_name_en"]
	if name == "" {
		name = record["name"]
	}
	if name == "" {
		fmt.Printf("ERROR: Missing country name for record: %s\n", record)
		os.Exit(1)
	}
	if overrides[name] == "" {
		return name
	} else {
		return overrides[name]
	}
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
		util.ExitIfError(err, fmt.Sprintf("Error processing csv file %s: %s", file, err))
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
		util.ExitIfError(err, fmt.Sprintf("Error processing csv file %s: %s", file, err))
		if first {
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
		if v != "" {
			row[headers[i]] = v
		}
	}

	return row
}

func readLanguages(file string) ([]Language, []LocaleName) {
	lang := IncomingLanguages{}
	err := json.Unmarshal(common.ReadFile(file), &lang)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshall languages: %s", err))

	languages := []Language{}

	localeNameMap := map[string]LocaleName{}

	for _, l := range lang.Languages {
		name := l.Names[0]
		if len(l.Iso_639_2) > 0 && name != "" && len(l.Countries) > 0 {
			locales := []string{}
			for _, incomingLocale := range l.Locales {
				localeId := common.FormatLocaleId(incomingLocale.Id)
				localeNameMap[incomingLocale.Id] = LocaleName{
					Id:   localeId,
					Name: incomingLocale.Name,
				}
				locales = append(locales, localeId)
			}

			languages = append(languages, Language{
				Name:      name,
				Iso_639_2: l.Iso_639_2,
				Countries: l.Countries,
				Locales:   locales,
			})
		}
	}
	sortLanguages(languages)

	names := []LocaleName{}
	for _, v := range localeNameMap {
		names = append(names, v)
	}
	sortLocaleNames(names)

	return languages, names
}

func readCurrencySymbols(file string) map[string]CurrencySymbols {
	data := CldrCurrencies{}
	err := json.Unmarshal(common.ReadFile(file), &data)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshall cldr currencies: %s", err))

	unsupportedCurrencyCodes := common.UnsupportedCurrencyCodes()
	currencySymbols := map[string]CurrencySymbols{}

	for _, main := range data.Main {
		for code, c := range main.Numbers.Currencies {
			if !common.ContainsIgnoreCase(unsupportedCurrencyCodes, code) {
				if c.Symbol != "" && c.Symbol != code {
					var narrow string
					if c.Symbol == c.SymbolAltNarrow {
						narrow = ""
					} else {
						narrow = c.SymbolAltNarrow
					}

					currencySymbols[code] = CurrencySymbols{
						Primary: c.Symbol,
						Narrow:  narrow,
					}
				}
			}
		}
	}
	return currencySymbols
}

func readCurrencies(file string) []Currency {
	data := []IncomingCurrency{}
	err := json.Unmarshal(common.ReadFile(file), &data)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshall currencies: %s", err))

	unsupportedCurrencyCodes := common.UnsupportedCurrencyCodes()
	currencies := []Currency{}

	for _, c := range data {
		if !common.ContainsIgnoreCase(unsupportedCurrencyCodes, c.Iso_4217_3) {
			currencies = append(currencies, Currency{
				Name:           c.Name,
				Iso_4217_3:     c.Iso_4217_3,
				NumberDecimals: c.NumberDecimals,
			})
		}
	}
	sortCurrencies(currencies)

	return currencies
}

func readNumbers(file string) []Number {
	data := IncomingNumbers{}
	err := json.Unmarshal(common.ReadFile(file), &data)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshall numbers: %s", err))

	numbers := []Number{}
	for _, main := range data.Main {
		country := main.Identity.Territory
		if country == "" {
			// e.g. 'fr' where the code maps to both the country and language
			country = main.Identity.Language
		}
		numbers = append(numbers, Number{
			Language: main.Identity.Language,
			Country:  country,
			Separators: Separators{
				Decimal: main.Numbers.Symbols.Decimal,
				Group:   main.Numbers.Symbols.Group,
			},
		})
	}

	return numbers
}

func toObjects(records []map[string]string, accepts acceptsFunction, f convertFunction, id idFunction) []interface{} {
	added := map[string]interface{}{}
	for _, v := range records {
		if accepts(v) {
			id := strings.ToUpper(id(v))
			if added[id] == nil {
				added[id] = f(v)
			}
		}
	}

	return sortObjects(added)
}

func LoadProvinces() []Province {
	provinces := []Province{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/provinces.json"), &provinces)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal provinces: %s", err))
	return provinces
}

func LoadCountryDuties() []CountryDuty {
	countryDuties := []CountryDuty{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/country-duties.json"), &countryDuties)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal country duties: %s", err))
	return countryDuties
}

func LoadCountryContinents() []CountryContinent {
	countryContinents := []CountryContinent{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/country-continents.json"), &countryContinents)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal country continents: %s", err))
	return countryContinents
}

func LoadPaymentMethods() []PaymentMethod {
	paymentMethods := []PaymentMethod{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/payment-methods.json"), &paymentMethods)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal payment methods: %s", err))
	return paymentMethods
}

func LoadCurrencyLocales() map[string]string {
	currencyLocales := []CurrencyLocale{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/currency-locales.json"), &currencyLocales)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal country continents: %s", err))

	table := map[string]string{}
	for _, cl := range currencyLocales {
		table[cl.CurrencyCode] = cl.LocaleId
	}
	return table
}

func LoadContinents() []Continent {
	return []Continent{
		Continent{
			Name:  "Africa",
			Code2: "AF",
			Code3: "AFR",
		},
		Continent{
			Name:  "Antarctica",
			Code2: "AN",
			Code3: "ANT",
		},
		Continent{
			Name:  "Asia",
			Code2: "AS",
			Code3: "ASI",
		},
		Continent{
			Name:  "Europe",
			Code2: "EU",
			Code3: "EUR",
		},
		Continent{
			Name:  "North America",
			Code2: "NA",
			Code3: "NOA",
		},
		Continent{
			Name:  "Oceania",
			Code2: "OC",
			Code3: "OCE",
		},
		Continent{
			Name:  "South America",
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

func LoadCurrencySymbols() map[string]CurrencySymbols {
	symbols := map[string]CurrencySymbols{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/currency-symbols.json"), &symbols)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal symbols: %s", err))
	return symbols
}

func LoadLanguages() []Language {
	languages := []Language{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/languages.json"), &languages)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal languages: %s", err))
	return languages
}

func LoadLocaleNames() []LocaleName {
	names := []LocaleName{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/locale-names.json"), &names)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal locale names: %s", err))
	return names
}

func LoadNumbers() []Number {
	numbers := []Number{}
	err := json.Unmarshal(common.ReadFile("data/cleansed/numbers.json"), &numbers)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshal numbers: %s", err))
	return numbers
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

func loadCldrNumbers(dir string) []Number {
	numbers := []Number{}
	filepath.Walk(dir, func(path string, dirInfo os.FileInfo, err error) error {
		numbersPath := fmt.Sprintf("%s/%s/numbers.json", dir, dirInfo.Name())
		if fileExists(numbersPath) {
			for _, n := range readNumbers(numbersPath) {
				if n.Country != "" {
					numbers = append(numbers, n)
				}
			}
		}
		return nil
	})
	return numbers
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func sortObjects(data map[string]interface{}) []interface{} {
	keys := []string{}
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var all []interface{}
	for _, key := range keys {
		all = append(all, data[key])
	}
	return all
}

func sortCurrencies(currencies []Currency) []Currency {
	slice.Sort(currencies[:], func(i, j int) bool {
		return strings.ToLower(currencies[i].Name) < strings.ToLower(currencies[j].Name)
	})
	return currencies
}

func sortLanguages(languages []Language) []Language {
	slice.Sort(languages[:], func(i, j int) bool {
		return strings.ToLower(languages[i].Name) < strings.ToLower(languages[j].Name)
	})
	return languages
}

func sortLocaleNames(names []LocaleName) []LocaleName {
	slice.Sort(names[:], func(i, j int) bool {
		return strings.ToLower(names[i].Name) < strings.ToLower(names[j].Name)
	})
	return names
}

func toInt32(value string) int {
	v, err := strconv.Atoi(value)
	util.ExitIfError(err, fmt.Sprintf("Failed to convert value[%s] to int32: %s", value, err))
	return v
}
