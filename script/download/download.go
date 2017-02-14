package download

// Downloads source files, storing locally in the sources directory

import (
	"../common"
	"bufio"
	"encoding/csv"
	"encoding/json"
        "fmt"
        "io"
        "io/ioutil"
        "net/http"
        "os"
        "strings"
)

type Continent struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	Countries  []string `json:"countries"`
}

type Country struct {
	Name                     string `json:"name"`
	Continent                string `json:"continent"`
	Iso_3166_2               string `json:"iso_3166_2"`
	Iso_3166_3               string `json:"iso_3166_3"`
	CurrencyCode             string `json:"currency_code"`
	CurrencyName             string `json:"currency_name"`
	CurrencyNumberDecimals   string `json:"currency_number_decimals"`
	Languages                []string `json:"languages"`
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

type IncomingLanguages struct {
	LanguageFamilies    []string `json:"languageFamilies"`
	Languages           []IncomingLanguage `json:"languages"`
}

type IncomingLanguage struct {
	Iso_639_2           string `json:"iso639_1"`
	Names               []string `json:"name"`
	Countries           []string `json:"countries"`
}


type convert func(records map[string]string) interface{}

func DownloadAll() {
	tmp := download("https://raw.githubusercontent.com/bdswiss/country-language/master/data.json")
	writeJson("sources/languages.json", readLanguages(tmp))

	tmp = download("https://raw.githubusercontent.com/datasets/country-codes/master/data/country-codes.csv")
	writeJson("sources/countries.json", toObjects(readCsv(tmp), func(record map[string]string) interface{} {
		return Country {
			Name: record["official_name_en"],
			Iso_3166_2: record["ISO3166-1-Alpha-2"],
			Iso_3166_3: record["ISO3166-1-Alpha-3"],
			CurrencyCode: record["ISO4217-currency_alphabetic_code"],
			CurrencyNumberDecimals: record["ISO4217-currency_number_decimals"],
			CurrencyName: record["ISO4217-currency_name"],
			Languages: strings.Split(record["Languages"], ","),
			Continent: record["Continent"],
		}
	}))

	//os.Rename(tmp, "sources/countries.json")

	// tmp = download("http://dev.maxmind.com/static/csv/codes/country_continent.csv")
	// os.Rename(tmp, "sources/continents.csv")
}

func exitIfError(err error, message string) {
        if err != nil {
		fmt.Printf("*** ERROR ***: %s\n", message)
		fmt.Println(err)
		os.Exit(1)
	}
}

// Download the provided url to a temp file, returning the file
func download(url string) string {
	target, err := ioutil.TempFile("", "reference-download")
	exitIfError(err, "Error creating temporary file")
	defer target.Close()
	
        response, err := http.Get(url)
	exitIfError(err, fmt.Sprintf("Error downloading url %s", url))
	defer response.Body.Close()

	_, err = io.Copy(target, response.Body)
	exitIfError(err, fmt.Sprintf("Error writing to file %s", target))

	return target.Name()
}

// readCsv Reads a CSV file, returning a list of map[string]string objects
func readCsv(file string) []map[string]string {
	input, err := os.Open(file)
	exitIfError(err, fmt.Sprintf("Error opening file %s", file))

        r := csv.NewReader(input)
	var headers []string
	first := true
	var all []map[string]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		exitIfError(err, fmt.Sprintf("Error processing csv file %s", file))
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

func readLanguages(file string) []common.Language {
	lang := IncomingLanguages{}
	err := json.Unmarshal(common.ReadFile(file), &lang)
	exitIfError(err, fmt.Sprintf("Failed to unmarshall languages: %s", err))

	languages := []common.Language{}
	
	for _, l := range(lang.Languages) {
		name := l.Names[0]
		if len(l.Iso_639_2) > 0 && name != "" {
			languages = append(languages, common.Language{
				Name: name,
				Iso_639_2: l.Iso_639_2,
				Countries: l.Countries,
			})
		}
	}

	return languages
}

func writeJson(target string, data interface{}) {
	tmp, err := ioutil.TempFile("", "reference-csv-to-json")
	exitIfError(err, "Error creating temporary file")
	defer tmp.Close()

	v, err := json.MarshalIndent(data, "", "  ")
	exitIfError(err, "Error marshalling record to json")
	//enc := json.NewEncoder(tmp)
	//err = enc.Encode(&data)
	
	w := bufio.NewWriter(tmp)
	_, err = w.Write(v)
	exitIfError(err, "Error writing to tmp file")
	
	err = os.Rename(tmp.Name(), target)
	exitIfError(err, "Error renaming tmp file")
}

func toObjects(records []map[string]string, f convert) []interface{} {
	var all []interface{}
	for _, v := range records {
		all = append(all, f(v))
	}
	return all
}

