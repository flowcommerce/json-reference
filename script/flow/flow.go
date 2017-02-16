package flow

// Reads all cleansed data, transforming into the final Flow format we desire

import (
	"../cleansed"
	"../common"
	"fmt"
	"os"
)

func Generate() {
	continents := cleansed.LoadContinents()
	countries := cleansed.LoadCountries()

	fmt.Println(commonContinents(continents, countries))
}

func commonContinents(continents []cleansed.Continent, countries []cleansed.Country) []common.Continent {
	all := make([]common.Continent, len(continents))
	for _, c := range(continents) {
		var theseCountries []string
		for _, country := range(countries) {
			if (country.Continent == c.Code) {
				theseCountries = append(theseCountries, country.Iso_3166_3)
			}
		}
		fmt.Printf("%s theseCountries: %s\n", c, theseCountries)

		all = append(all, common.Continent{
			Name: c.Name,
			Code: c.Code,
			Countries: theseCountries,
		})
	}
	return all
}

func findCountryByIso31662(countries []cleansed.Country, code string) cleansed.Country {
	for _, c := range(countries) {
		if c.Iso_3166_2 == code {
			return c
		}
	}
	fmt.Printf("Invalid country code: %s\n", code)
	os.Exit(1)
	return cleansed.Country{}
}

func findContinentByCode(continents []cleansed.Continent, code string) cleansed.Continent {
	for _, c := range(continents) {
		if c.Code == code {
			return c
		}
	}
	fmt.Printf("Invalid continent code: %s\n", code)
	os.Exit(1)
	return cleansed.Continent{}
}
