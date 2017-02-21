package final

// Reads all cleansed data, transforming into the final format we desire

import (
	"../cleanse"
	"../common"
	"fmt"
	"github.com/bradfitz/slice"
	"os"
	"regexp"
	"sort"
	"strings"
)

type CleansedDataSet struct {
	Continents []cleanse.Continent
	Countries []cleanse.Country
	CountryContinents []cleanse.CountryContinent
	Currencies []cleanse.Currency
	Languages []cleanse.Language
	Timezones []cleanse.Timezone
	CountryTimezones []cleanse.CountryTimezone
}
	
func Generate() {
	data := CleansedDataSet{
		Continents: cleanse.LoadContinents(),
		Countries: cleanse.LoadCountries(),
		CountryContinents: cleanse.LoadCountryContinents(),
		Currencies: cleanse.LoadCurrencies(),
		Languages: cleanse.LoadLanguages(),
		Timezones: cleanse.LoadTimezones(),
		CountryTimezones: cleanse.LoadCountryTimezones(),
	}

	continents := commonContinents(data)
	countries := commonCountries(data)

	writeJson("data/final/continents.json", continents)
	writeJson("data/final/languages.json", commonLanguages(data))
	writeJson("data/final/currencies.json", commonCurrencies(data))
	writeJson("data/final/timezones.json", commonTimezones(data))
	writeJson("data/final/countries.json", countries)
	writeJson("data/final/regions.json", createRegions(countries, continents))
}

func writeJson(target string, objects interface{}) {
	fmt.Printf("Writing %s\n", target)
	common.WriteJson(target, objects)
}
	
func commonContinents(data CleansedDataSet) []common.Continent {
	var all []common.Continent
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
	var all []common.Language
	for _, l := range(data.Languages) {
		theseCountries := []string{}
		for _, country := range(data.Countries) {
			if common.ContainsIgnoreCase(l.Countries, country.Iso_3166_3) {
				theseCountries = append(theseCountries, country.Iso_3166_3)
			}
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

func commonTimezones(data CleansedDataSet) []common.Timezone {
	var all []common.Timezone
	for _, t := range(data.Timezones) {
		all = append(all, common.Timezone{
			Name: t.Name,
			Description: t.Description,
			Offset: t.Offset,
		})
	}
	return all
}

func commonCurrencies(data CleansedDataSet) []common.Currency {
	var all []common.Currency
	for _, c := range(data.Currencies) {
		all = append(all, common.Currency{
			Name: c.Name,
			Iso_4217_3: c.Iso_4217_3,
			NumberDecimals: c.NumberDecimals,
		})
	}
	return all
}

func commonCountries(data CleansedDataSet) []common.Country {
	var all []common.Country
	for _, c := range(data.Countries) {
		languages := []string {}
		for _, l := range(data.Languages) {
			if common.ContainsIgnoreCase(l.Countries, c.Iso_3166_3) {
				languages = append(languages, l.Iso_639_2)
			}
		}
		
		timezones := []string {}
		for _, ct := range(data.CountryTimezones) {
			if ct.CountryCode == c.Iso_3166_3 {
				tz := findTimezone(data.Timezones, ct.TimezoneCode)
				timezones = append(timezones, tz.Name)
			}
		}

		var defaultCurrency string
		if c.Currency != "" {
			defaultCurrency = findCurrencyByCode(data.Currencies, c.Currency).Iso_4217_3
		}

		sort.Strings(languages)
		sort.Strings(timezones)
		all = append(all, common.Country{
			Name: formatCountryName(c.Iso_3166_3, c.Name),
			Iso_3166_2: c.Iso_3166_2,
			Iso_3166_3: c.Iso_3166_3,
			MeasurementSystem: getMeasurementSystem(c.Iso_3166_3),
			DefaultCurrency: defaultCurrency,
			Languages: languages,
			Timezones: timezones,
			
		})
	}
        return all
}

func createRegions(countries []common.Country, continents []common.Continent) []common.Region {
	regions := []common.Region{}

	for _, c := range(countries) {
		id := generateId(c.Iso_3166_3)

		regions = append(regions, common.Region{
			Id: id,
			Name: c.Name,
			Countries: []string { c.Iso_3166_3 },
			Currencies: currenciesForCountries([]common.Country { c }),
			Languages: languagesForCountries([]common.Country { c }),
			MeasurementSystems: measurementSystemsForCountries([]common.Country { c }),
			Timezones: timezonesForCountries([]common.Country { c }),
		})
	}
	
	for _, c := range(continents) {
		if c.Code != "ANT" {
			id := generateId(c.Name)

			theseCountries := findCountries(countries, c.Countries)
			regions = append(regions, common.Region{
				Id: id,
				Name: c.Name,
				Countries: toCountryCodes(theseCountries),
				Currencies: currenciesForCountries(theseCountries),
				Languages: languagesForCountries(theseCountries),
				MeasurementSystems: measurementSystemsForCountries(theseCountries),
				Timezones: timezonesForCountries(theseCountries),
			})
		}
	}

	regions = append(regions, eurozone(countries), world(countries))
	assertUniqueRegionIds(regions)
	sortRegions(regions)
	return regions
}

func eurozone(countries []common.Country) common.Region {
	countryCodes := findCountriesByCurrency(countries, "EUR")
	theseCountries := findCountries(countries, countryCodes)	
	return common.Region{
		Id: "eurozone",
		Name: "Eurozone",
		Countries: countryCodes,
		Currencies: currenciesForCountries(theseCountries),
		Languages: languagesForCountries(theseCountries),
		MeasurementSystems: measurementSystemsForCountries(theseCountries),
		Timezones: timezonesForCountries(theseCountries),
	}
}

func world(countries []common.Country) common.Region {
	var codes []string
	for _, c := range(countries) {
		codes = append(codes, c.Iso_3166_3)
	}
	sort.Strings(codes)

	return common.Region{
		Id: "world",
		Name: "World",
		Countries: codes,
		Currencies: currenciesForCountries(countries),
		Languages: languagesForCountries(countries),
		MeasurementSystems: []string{"metric", "imperial"},
		Timezones: timezonesForCountries(countries),
	}
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

func findCurrencyByCode(currencies []cleanse.Currency, code string) cleanse.Currency {
	for _, c := range(currencies) {
		if c.Iso_4217_3 == code {
			return c
		}
	}
	fmt.Printf("Invalid currency code: %s\n", code)
	os.Exit(1)
	return cleanse.Currency{}
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

func findTimezone(timezones []cleanse.Timezone, name string) cleanse.Timezone {
	for _, c := range(timezones) {
		if strings.ToUpper(c.Name) == strings.ToUpper(name) || strings.ToUpper(c.Description) == strings.ToUpper(name) {
			return c
		}
	}
	fmt.Printf("Invalid timezone name: %s\n", name)
	os.Exit(1)
	return cleanse.Timezone{}
}

func getMeasurementSystem(iso_3166_3 string) string {
	if iso_3166_3 == "USA" || iso_3166_3 == "LBR" || iso_3166_3 == "MMR" {
		return "imperial"
	}
	return "metric"
}

func formatCountryName(iso3 string, defaultName string) string {
	switch iso3 {
	case "USA": {
		return "United States of America"
	}
	case "GBR": {
		return "United Kingdom"
	}
	default: {
		return defaultName
	}
	}
}

func toCountryCodes(countries []common.Country) []string {
	codes := []string{}

	for _, country := range(countries) {
		codes = append(codes, country.Iso_3166_3)
	}

	sort.Strings(codes)
	return codes
}

func findCountries(countries []common.Country, codes []string) []common.Country {
	matching := []common.Country{}

	for _, country := range(countries) {
		if common.ContainsIgnoreCase(codes, country.Iso_3166_3) {
			matching = append(matching, country)
		}
	}
		
	return matching
}

/**
 * Filters the list of countries to those with a matching default currency.
 */
func findCountriesByCurrency(countries []common.Country, currency string) []string {
	codes := []string{}

	for _, country := range(countries) {
		if country.DefaultCurrency == currency && !common.ContainsIgnoreCase(codes, country.Iso_3166_3) {
			codes = append(codes, country.Iso_3166_3)
		}
	}

	sort.Strings(codes)
	return codes
}

func currenciesForCountries(countries []common.Country) []string {
	codes := []string{}

	for _, country := range(countries) {
		if country.DefaultCurrency != "" && !common.ContainsIgnoreCase(codes, country.DefaultCurrency) {
			codes = append(codes, country.DefaultCurrency)
		}
	}

	sort.Strings(codes)
	return codes
}

func languagesForCountries(countries []common.Country) []string {
	codes := []string{}

	for _, country := range(countries) {
		for _, l := range(country.Languages) {
			if !common.ContainsIgnoreCase(codes, l) {
				codes = append(codes, l)
			}
		}
	}

	sort.Strings(codes)
	return codes
}

func timezonesForCountries(countries []common.Country) []string {
	codes := []string{}

	for _, country := range(countries) {
		for _, tz := range(country.Timezones) {
			if !common.ContainsIgnoreCase(codes, tz) {
				codes = append(codes, tz)
			}
		}
	}

	sort.Strings(codes)
	return codes
}

func measurementSystemsForCountries(countries []common.Country) []string {
	codes := []string{}

	for _, country := range(countries) {
		if country.MeasurementSystem != "" && !common.ContainsIgnoreCase(codes, country.MeasurementSystem) {
			codes = append(codes, country.MeasurementSystem)
		}
	}

	sort.Strings(codes)
	return codes
}

func assertUniqueRegionIds(regions []common.Region) {
	found := make(map[string]bool)

	for _, r := range(regions) {
		if found[r.Id] {
			fmt.Printf("ERROR: Duplicate region id[%s]\n", r.Id)
			os.Exit(1)
		}
		found[r.Id] = true
	}
}

func generateId(name string) string {
	safe := regexp.MustCompile("[^A-Za-z0-9]+").ReplaceAllString(name, "-")
	return strings.ToLower(safe)
}

func sortRegions(regions []common.Region) []common.Region {
	slice.Sort(regions[:], func(i, j int) bool {
		return strings.ToLower(regions[i].Name) < strings.ToLower(regions[j].Name)
	})
	return regions
}
