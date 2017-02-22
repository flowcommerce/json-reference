package download

// Downloads source files, storing locally in the sources directory

import (
	"github.com/flowcommerce/tools/util"
	"fmt"
	"io"
        "io/ioutil"
        "net/http"
        "os"
)

func DownloadAll() {
	download("data/source/languages.json", "https://raw.githubusercontent.com/bdswiss/country-language/master/data.json")
	download("data/source/countries.csv", "https://raw.githubusercontent.com/datasets/country-codes/master/data/country-codes.csv")
	download("data/source/country-continents.csv", "http://dev.maxmind.com/static/csv/codes/country_continent.csv")
	download("data/source/cldr-currencies.json", "https://raw.githubusercontent.com/unicode-cldr/cldr-numbers-full/master/main/en-US-POSIX/currencies.json")
}

// Download the provided url to a temp file, returning the file
func download(target string, url string) {
	fmt.Printf("Downloading %s...\n", url)
	tmp, err := ioutil.TempFile("", "reference-download")
	util.ExitIfError(err, "Error creating temporary file")
	defer tmp.Close()
	
        response, err := http.Get(url)
	util.ExitIfError(err, fmt.Sprintf("Error downloading url %s", url))
	defer response.Body.Close()

	_, err = io.Copy(tmp, response.Body)
	util.ExitIfError(err, fmt.Sprintf("Error writing to file %s", tmp))

	os.Rename(tmp.Name(), target)
	fmt.Printf("  -> Stored in %s\n", target)
}
