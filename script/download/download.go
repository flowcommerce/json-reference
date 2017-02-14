package download

// Downloads source files, storing locally in the sources directory

import (
	"./common"
	"encoding/csv"
	"encoding/json"
        "fmt"
        "io"
        "io/ioutil"
        "net/http"
        "os"

)

type Currency struct {
	Iso4217        string
	Name           string
	NumberDecimal  string
}


func DownloadAll() {
	// tmp := download("https://raw.githubusercontent.com/bdswiss/country-language/master/data.json")
	// os.Rename(tmp, "sources/languages.json")

	tmp := csvToJson(
		download("https://raw.githubusercontent.com/datasets/country-codes/master/data/country-codes.csv"),
		map[string]string{
			"official_name_en": "name",
			"ISO3166-1-Alpha-2": "iso31662",
			"ISO3166-1-Alpha-3": "iso31663",
			"Languages": "languages",
			"ISO4217-currency_alphabetic_code": "currency_code",
			"ISO4217-currency_minor_unit": "currency_number_decimals",
			"ISO4217-currency_name": "currency_name",
			"Continent": "continent",
		},
	)
	os.Rename(tmp, "sources/countries.json")

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

// csvToJson reads the CSV data from the specified file, writing to a file in
// json format. Returns a path to the json file
func csvToJson(file string, nameMap map[string]string) string {
	jsonFile, err := ioutil.TempFile("", "reference-csv-to-json")
	exitIfError(err, "Error creating temporary file")
	defer jsonFile.Close()
	enc := json.NewEncoder(jsonFile)

	input, err := os.Open(file)
	exitIfError(err, fmt.Sprintf("Error opening file %s", file))

        r := csv.NewReader(input)
	var headers []string
	first := true
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
			fmt.Println(record)
			json := make(map[string]interface{})

			for i, v := range record {
				if (v != "") {
					header := nameMap[headers[i]]
					if (header != "") {
						json[header] = v
					}
				}
			}
			fmt.Println(json)
			err := enc.Encode(&json)
			exitIfError(err, "Error marshalling record to json")
		}
	}

	return jsonFile.Name()
}
