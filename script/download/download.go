package download

// Downloads source files, storing locally in the sources directory

import (
        "fmt"
        "io"
        "io/ioutil"
        "net/http"
        "os"
)

func DownloadAll() {
	tmp := download("https://raw.githubusercontent.com/bdswiss/country-language/master/data.json")
	os.Rename(tmp, "sources/languages.json")

	tmp = download("https://raw.githubusercontent.com/datasets/country-codes/master/data/country-codes.csv")
	os.Rename(tmp, "sources/countries.csv")

	tmp = download("http://dev.maxmind.com/static/csv/codes/country_continent.csv")
	os.Rename(tmp, "sources/continents.csv")
}

func exitWithError(message string) {
	fmt.Printf("*** ERROR ***: %s\n", message)
}

// Download the provided url to a temp file, returning the path to the temporary file
func download(url string) string {
	target, err := ioutil.TempFile("", "reference-download")
        if err != nil {
		exitWithError("Error creating temporary file")
        }

        response, err := http.Get(url)
        if err != nil {
		exitWithError(fmt.Sprintf("Error downloading url %s", url))
        }
	defer response.Body.Close()

	_, err = io.Copy(target, response.Body)
	if err != nil {
		exitWithError(fmt.Sprintf("Error writing to file %s", target))
	}

	return target.Name()
}
