package flow

// Reads all cleansed data, transforming into the final Flow format we desire

import (
	"../cleanse"
	"../common"
	"fmt"
	"os"
	"sort"
)

type CleansedDataSet struct {
	Continents []cleanse.Continent
	Countries []cleanse.Country
	Currencies []cleanse.Currency
	Languages []cleanse.Language
}
	
func Generate() {
	data := CleansedDataSet{
		Continents: cleanse.LoadContinents(),
		Countries: cleanse.LoadCountries(),
		Currencies: cleanse.LoadCurrencies(),
		Languages: cleanse.LoadLanguages(),
	}

	writeJson("data/3-flow/continents.json", commonContinents(data))
	writeJson("data/3-flow/languages.json", commonLanguages(data))
	writeJson("data/3-flow/currencies.json", commonCurrencies(data))
	writeJson("data/3-flow/countries.json", commonCountries(data))
}

func writeJson(target string, objects interface{}) {
	fmt.Printf("Writing %s\n", target)
	common.WriteJson(target, objects)
}
	
func commonContinents(data CleansedDataSet) []common.Continent {
	all := make([]common.Continent, len(data.Continents))
	for _, c := range(data.Continents) {
		var theseCountries []string
		for _, country := range(data.Countries) {
			if (country.Continent != "") {
				continent := findContinentByCode(data.Continents, country.Continent)
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

func commonLanguages(data CleansedDataSet) []common.Language {
	all := make([]common.Language, len(data.Languages))
	for _, l := range(data.Languages) {
		var theseCountries []string

		for _, countryCode := range(l.Countries) {
			country := findCountryByCode(data.Countries, countryCode)
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

func commonCurrencies(data CleansedDataSet) []common.Currency {
	unsupported := []string {
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
	
	all := make([]common.Currency, len(data.Currencies))
	for _, c := range(data.Currencies) {
		if !common.Contains(unsupported, c.Iso_4217_3) {
			all = append(all, common.Currency{
				Name: c.Name,
				Iso_4217_3: c.Iso_4217_3,
				NumberDecimals: c.NumberDecimals,
			})
		}
	}
	return all
}

func commonCountries(data CleansedDataSet) []common.Country {
	unsupported := []string {
		[
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
			"ZWE"
		]
	
	all := make([]common.Country, len(data.Countries))
	for _, c := range(data.Countries) {
		if !common.Contains(unsupported, c.Iso_3166_3) {
			languages := make([]string)
			for _, l := range(data.Languages) {
				if common.Contains(l.Countries, c.Iso_3166_3) {
					languages = append(languages, c.Iso_3166_3)
				}
			}
		
			all = append(all, common.Country{
				Name: c.Name,
				Iso_3166_2: c.Iso_3166_2,
				Iso_3166_3: c.Iso_3166_3,
				MeasurementSystem: xxx,
				DefaultCurrency: xxx,
				Languages: languages,
				Timezones: timezones,

			})
		}
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
