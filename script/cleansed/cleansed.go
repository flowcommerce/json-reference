package cleansed

// Reads source files, cleansing and writing all as json to data/1-cleansed

import (
	"github.com/flowcommerce/tools/util"
	"../common"
	"bufio"
	"encoding/csv"
	"encoding/json"
        "fmt"
        "io"
        "io/ioutil"
        "os"
	"strconv"
        "strings"
)

type CleansedContinent struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
}

type CleansedCountry struct {
	Name                     string `json:"name"`
	Continent                string `json:"continent"`
	Iso_3166_2               string `json:"iso_3166_2"`
	Iso_3166_3               string `json:"iso_3166_3"`
	Currency                 string `json:"currency"`
	Languages                []string `json:"languages"`
}

type CleansedCurrency struct {
	Name                     string `json:"name"`
	Iso_4217_3               string `json:"iso_4217_3"`
	NumberDecimals           int64  `json:"number_decimals"`
}

type CleansedLanguage struct {
	Name                     string `json:"name"`
	Iso_639_2                string `json:"iso_639_2"`
	Countries                []string `json:"countries"`
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


type convertFunction func(records map[string]string) interface{}
type acceptsFunction func(records map[string]string) bool

func Cleanse() {
	writeJson("data/2-cleansed/languages.json", readLanguages("data/1-sources/languages.json"))

	countriesSource := readCsv("data/1-sources/countries.csv")
	writeJson("data/2-cleansed/countries.json",
		toObjects(countriesSource,
			func(record map[string]string) bool {
				return record["ISO3166-1-Alpha-2"] != "" && record["ISO3166-1-Alpha-3"] != ""
			},
			func(record map[string]string) interface{} {
				return CleansedCountry {
					Name: record["official_name_en"],
					Iso_3166_2: record["ISO3166-1-Alpha-2"],
					Iso_3166_3: record["ISO3166-1-Alpha-3"],
					Currency: record["ISO4217-currency_alphabetic_code"],
					Languages: common.FilterNonEmpty(strings.Split(record["Languages"], ",")),
					Continent: record["Continent"],
					
				}
			},
		),
	)

	writeJson("data/2-cleansed/currencies.json",
		toObjects(countriesSource,
			func(record map[string]string) bool {
				return record["ISO4217-currency_name"] != "" && record["ISO4217-currency_alphabetic_code"] != ""
			},
			func(record map[string]string) interface{} {
				n := record["ISO4217-currency_number_decimals"]
				var numberDecimals int64
				var err error
				if (n != "") {
					numberDecimals, err = strconv.ParseInt(n, 10, 0)
					util.ExitIfError(err, fmt.Sprintf("Error parsing int[%s]", n))
				}

				return CleansedCurrency{
					Name: record["ISO4217-currency_name"],
					Iso_4217_3: record["ISO4217-currency_alphabetic_code"],
					NumberDecimals: numberDecimals,
				}
			},
		),
	)

	writeJson("data/2-cleansed/continents.json",
		toObjects(readCsv("data/1-sources/continents.csv"),
			func(record map[string]string) bool {
				return record["continent code"] != "" && record["continent code"] != "--"
			},
			func(record map[string]string) interface{} {
				return CleansedContinent{
					Name: record["iso 3166 country"],
					Code: record["continent code"],
				}
			},
		),
	)	

	writeJson("data/2-cleansed/languages.json", filterLanguages(readLanguages("data/1-sources/languages.json")))
}

// readCsv Reads a CSV file, returning a list of map[string]string objects
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
			row := make(map[string]string)

			for i, v := range record {
				if (v != "") {
					row[headers[i]] = v
				}
			}
			all = append(all, row)
		}
	}

	return all
}

func readLanguages(file string) []CleansedLanguage {
	lang := IncomingLanguages{}
	err := json.Unmarshal(common.ReadFile(file), &lang)
	util.ExitIfError(err, fmt.Sprintf("Failed to unmarshall languages: %s", err))

	languages := []CleansedLanguage{}
	
	for _, l := range(lang.Languages) {
		name := l.Names[0]
		if len(l.Iso_639_2) > 0 && name != "" {
			languages = append(languages, CleansedLanguage{
				Name: name,
				Iso_639_2: l.Iso_639_2,
				Countries: l.Countries,
			})
		}
	}

	return languages
}

/**
 * filterLanguages takes only the languages that have at least one country assigned to them
 */
func filterLanguages(languages []CleansedLanguage) []CleansedLanguage {
	final := []CleansedLanguage{}
	
	for _, l := range(languages) {
		if (l.Countries != nil && len(l.Countries) > 0) {
			final = append(final, l)
		}
	}

	return final
}

func writeJson(target string, data interface{}) {
	tmp, err := ioutil.TempFile("", "reference-csv-to-json")
	util.ExitIfError(err, "Error creating temporary file")
	defer tmp.Close()

	v, err := json.MarshalIndent(data, "", "  ")
	util.ExitIfError(err, "Error marshalling record to json")
	
	w := bufio.NewWriter(tmp)
	_, err = w.Write(v)
	util.ExitIfError(err, "Error writing to tmp file")
	
	err = os.Rename(tmp.Name(), target)
	util.ExitIfError(err, "Error renaming tmp file")
}

func toObjects(records []map[string]string, accepts acceptsFunction, f convertFunction) []interface{} {
	var all []interface{}
	for _, v := range records {
		if (accepts(v)) {
			all = append(all, f(v))
		}
	}
	return all
}

