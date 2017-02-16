package flow

// Reads all cleansed data, transforming into the final Flow format we desire

import (
	"../cleanse"
	"../common"
	"fmt"
	"os"
	"sort"
)

func Generate() {
	continents := cleanse.LoadContinents()
	countries := cleanse.LoadCountries()
	currencies := cleanse.LoadCurrencies()
	languages := cleanse.LoadLanguages()

	
	writeJson("data/3-flow/continents.json", commonContinents(continents, countries))
	writeJson("data/3-flow/languages.json", commonLanguages(languages, countries))
	writeJson("data/3-flow/currencies.json", commonCurrencies(currencies))
}

func writeJson(target string, objects interface{}) {
	fmt.Printf("Writing %s\n", target)
	common.WriteJson(target, objects)
}
	
func commonContinents(continents []cleanse.Continent, countries []cleanse.Country) []common.Continent {
	all := make([]common.Continent, len(continents))
	for _, c := range(continents) {
		var theseCountries []string
		for _, country := range(countries) {
			if (country.Continent != "") {
				continent := findContinentByCode(continents, country.Continent)
				if (c == continent) {
					theseCountries = append(theseCountries, country.Iso_3166_3)
				}
			}
		}
		sort.Strings(theseCountries)

		all = append(all, common.Continent{
			Name: c.Name,
			Code: c.Code3,
			Countries: theseCountries,
		})
	}
	return all
}

func commonLanguages(languages []cleanse.Language, countries []cleanse.Country) []common.Language {
	all := make([]common.Language, len(languages))
	for _, l := range(languages) {
		var theseCountries []string

		for _, countryCode := range(l.Countries) {
			country := findCountryByCode(countries, countryCode)
			theseCountries = append(theseCountries, country.Iso_3166_3)
		}
		sort.Strings(theseCountries)

		all = append(all, common.Language{
			Name: l.Name,
			Iso_639_2: l.Iso_639_2,
			Countries: theseCountries,
		})
	}
	return all
}

func commonCurrencies(currencies []cleanse.Currency) []common.Currency {
	all := make([]common.Currency, len(currencies))
	for _, c := range(currencies) {
		all = append(all, common.Currency{
			Name: c.Name,
			Iso_4217_3: c.Iso_4217_3,
			NumberDecimals: c.NumberDecimals,
		})
	}
	return all
}

func findCountryByCode(countries []cleanse.Country, code string) cleanse.Country {
	for _, c := range(countries) {
		if c.Iso_3166_2 == code || c.Iso_3166_3 == code {
			return c
		}
	}
	fmt.Printf("Invalid country code: %s\n", code)
	os.Exit(1)
	return cleanse.Country{}
}

func findContinentByCode(continents []cleanse.Continent, code string) cleanse.Continent {
	for _, c := range(continents) {
		if c.Code2 == code || c.Code3 == code {
			return c
		}
	}
	fmt.Printf("Invalid continent code: %s\n", code)
	os.Exit(1)
	return cleanse.Continent{}
}

func findLanguageByCode(languages []cleanse.Language, code string) cleanse.Language {
	for _, c := range(languages) {
		if c.Iso_639_2 == code {
			return c
		}
	}
	fmt.Printf("Invalid language code: %s\n", code)
	os.Exit(1)
	return cleanse.Language{}
}
