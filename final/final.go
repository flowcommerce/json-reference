package final

// Reads all cleansed data, transforming into the final format we desire

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"../cleanse"
	"../common"
	"github.com/bradfitz/slice"
)

type CleansedDataSet struct {
	Carriers             []cleanse.Carrier
	CarrierServices      []cleanse.CarrierService
	Continents           []cleanse.Continent
	Countries            []cleanse.Country
	CountryContinents    []cleanse.CountryContinent
	CountryDuties        []cleanse.CountryDuty
	Currencies           []cleanse.Currency
	CurrencySymbols      map[string]cleanse.CurrencySymbols
	Numbers              []cleanse.Number
	Languages            []cleanse.Language
	LocaleNames          []cleanse.LocaleName
	PaymentMethods       []cleanse.PaymentMethod
	Provinces            []cleanse.Province
	ProvinceTranslations []cleanse.ProvinceTranslation
	Timezones            []cleanse.Timezone
	CountryTimezones     []cleanse.CountryTimezone
	CountryDefaultLanguages     []cleanse.CountryDefaultLanguage
}

func Generate() {
	data := CleansedDataSet{
		Carriers:             cleanse.LoadCarriers(),
		CarrierServices:      cleanse.LoadCarrierServices(),
		Continents:           cleanse.LoadContinents(),
		Countries:            cleanse.LoadCountries(),
		CountryContinents:    cleanse.LoadCountryContinents(),
		CountryDuties:        cleanse.LoadCountryDuties(),
		Currencies:           cleanse.LoadCurrencies(),
		CurrencySymbols:      cleanse.LoadCurrencySymbols(),
		Languages:            cleanse.LoadLanguages(),
		LocaleNames:          cleanse.LoadLocaleNames(),
		PaymentMethods:       cleanse.LoadPaymentMethods(),
		Provinces:            cleanse.LoadProvinces(),
		ProvinceTranslations: cleanse.LoadProvinceTranslations(),
		Numbers:              cleanse.LoadNumbers(),
		Timezones:            cleanse.LoadTimezones(),
		CountryTimezones:     cleanse.LoadCountryTimezones(),
		CountryDefaultLanguages: cleanse.LoadCountryDefaultLanguages(),
	}

	continents := commonContinents(data)
	countries := commonCountries(data)
	locales := commonLocales(data)
	regions := createRegions(countries, continents)
	provinces := createProvinces(data, locales)

	writeJson("data/final/carriers.json", commonCarriers(data))
	writeJson("data/final/carrier-services.json", commonCarrierServices(data))
	writeJson("data/final/continents.json", continents)
	writeJson("data/final/payment-methods.json", commonPaymentMethods(data, regions))
	writeJson("data/final/languages.json", commonLanguages(data))
	writeJson("data/final/locales.json", locales)
	writeJson("data/final/currencies.json", commonCurrencies(data, locales))
	writeJson("data/final/timezones.json", commonTimezones(data))
	writeJson("data/final/countries.json", countries)
	writeJson("data/final/regions.json", regions)
	writeJson("data/final/provinces.json", provinces)
}

func writeJson(target string, objects interface{}) {
	fmt.Printf("Writing %s\n", target)
	common.WriteJson(target, objects)
}

func commonContinents(data CleansedDataSet) []common.Continent {
	var all []common.Continent
	for _, c := range data.Continents {
		var theseCountries []string
		for _, country := range data.Countries {
			if country.Continent != "" {
				continent := findContinentByCode(data.Continents, country.Continent)
				if c == continent {
					theseCountries = append(theseCountries, country.Iso_3166_3)
				}
			}
		}
		sort.Strings(theseCountries)

		all = append(all, common.Continent{
			Name:      c.Name,
			Code:      c.Code3,
			Countries: theseCountries,
		})
	}
	return all
}

func commonCarriers(data CleansedDataSet) []common.Carrier {
	var all []common.Carrier

	for _, carrier := range data.Carriers {
		all = append(all, common.Carrier{
			Id:          carrier.Id,
			Name:        carrier.Name,
			TrackingUrl: carrier.TrackingUrl,
		})
	}

	return all
}

func commonCarrierServices(data CleansedDataSet) []common.CarrierService {
	var all []common.CarrierService

	for _, carrierService := range data.CarrierServices {
		// find the correct carrier for this service
		var carrier cleanse.Carrier

		for _, c := range data.Carriers {
			if strings.ToUpper(carrierService.CarrierId) == strings.ToUpper(c.Id) {
				carrier = c
			}
		}

		// append the service with the carrier
		all = append(all, common.CarrierService{
			Id:   carrierService.Id,
			Name: carrierService.Name,
			Carrier: common.Carrier{
				Id:          carrier.Id,
				Name:        carrier.Name,
				TrackingUrl: carrier.TrackingUrl,
			},
		})
	}

	return all
}

func commonLocales(data CleansedDataSet) []common.Locale {
	var all []common.Locale

	englishCountries := []string{
		"ASM",
		"BOL",
		"IND",
		"KEN",
		"TON",
		"UGA",
		"GRD",
		"GRL",
		"KIR",
		"KNA",
	}

	countryMap := map[string]string{
		"zh": "cn",
		"el": "gr",
	}

	unsupportedLanguages := []string{
		"mas",
	}

	languageMap := map[string]string{
		"br":  "pt",
		"fo":  "da",
		"kw":  "ar",
		"ml":  "fr",
		"mr":  "ar",
		"nds": "nl",
		"om":  "ar",
		"os":  "ru",
		"se":  "sw",
	}

	unsupportedCountryCodes := common.UnsupportedCountryCodes()

	for _, n := range data.Numbers {
		if common.ContainsIgnoreCase(unsupportedCountryCodes, n.Country) {
			continue
		}

		var originalCountry string
		if countryMap[n.Country] == "" {
			originalCountry = n.Country
		} else {
			originalCountry = countryMap[n.Country]
		}

		var originalLanguage string
		if languageMap[n.Language] == "" {
			originalLanguage = n.Language
		} else {
			originalLanguage = languageMap[n.Language]
		}

		countryCode := normalizeCountryCode(data.Countries, originalCountry)
		languageCode := normalizeLanguageCode(data.Languages, originalLanguage)
		if countryCode == "" && languageCode == "" {
			// fmt.Printf(" - skipping locale as neither country code[%s] nor language code[%s] is known\n", originalCountry, originalLanguage)
		} else {
			if countryCode == "" {
				language := findLanguageByCode(data.Languages, languageCode)
				if len(language.Countries) == 1 {
					// Language mapped to exactly 1 country
					countryCode = normalizeCountryCode(data.Countries, language.Countries[0])
				}
				if countryCode == "" {
					// fmt.Printf(" - unknown country[%s] for language[%s] - skipping\n", originalCountry, languageCode)
					continue
				}
			}

			if languageCode == "" {
				if common.Contains(englishCountries, countryCode) {
					languageCode = normalizeLanguageCode(data.Languages, "en")
				}

				if languageCode == "" {
					if originalLanguage == "gsw" {
						// Missing locale for swiss german
						continue
					} else if common.ContainsIgnoreCase(unsupportedLanguages, originalLanguage) {
						// don't report error
						continue
					} else {
						fmt.Printf("ERROR: Unknown language[%s] w/ country[%s]\n", originalLanguage, countryCode)
						continue
					}
				}
			}

			separator := ""
			if n.Separators.Group == "," {
				separator = ","
			} else if n.Separators.Group == " " {
				// Weird encoding from cldr-json
				separator = " "
			} else if n.Separators.Group == "." {
				separator = "."
			} else if n.Separators.Group == "'" {
				separator = "'"
			} else if n.Separators.Group == "’" {
				separator = "’"
			} else {
				fmt.Printf("Invalid group separator[%s]\n", n.Separators.Group)
				os.Exit(1)
			}

			language := findLanguageByCode(data.Languages, languageCode)
			country := findCountryByCode(data.Countries, countryCode)
			id := common.FormatLocaleId(fmt.Sprintf("%s-%s", language.Iso_639_2, country.Iso_3166_2))
			name := findLocaleNameById(data.LocaleNames, id)
			if name == "" {
				name = fmt.Sprintf("%s - %s", language.Name, country.Name)
			}

			all = append(all, common.Locale{
				Id:       id,
				Name:     name,
				Country:  country.Iso_3166_3,
				Language: language.Iso_639_2,
				Numbers: common.LocaleNumbers{
					Decimal: n.Separators.Decimal,
					Group:   separator,
				},
			})
		}
	}

	uniqueLocales := uniqueLocaleIds(all)
	sortLocales(uniqueLocales)

	return uniqueLocales
}

func commonLanguages(data CleansedDataSet) []common.Language {
	var all []common.Language
	for _, l := range data.Languages {
		theseCountries := []string{}
		for _, country := range data.Countries {
			if common.ContainsIgnoreCase(l.Countries, country.Iso_3166_3) {
				theseCountries = append(theseCountries, country.Iso_3166_3)
			}
		}
		sort.Strings(theseCountries)

		theseLocales := []string{}
		for _, locale := range l.Locales {
			// TODO: Validate locale is known
			theseLocales = append(theseLocales, locale)
		}
		sort.Strings(theseLocales)

		all = append(all, common.Language{
			Name:      l.Name,
			Iso_639_2: l.Iso_639_2,
			Countries: theseCountries,
			Locales:   theseLocales,
		})
	}
	return all
}

func commonPaymentMethods(data CleansedDataSet, regions []common.Region) []common.PaymentMethod {
	var all []common.PaymentMethod
	for _, pm := range data.PaymentMethods {
		theseRegions := []string{}

		hasWorld := false
		for _, regionId := range pm.Regions {
			r := findRegion(regions, regionId)
			if r.Id == "world" {
				hasWorld = true
			} else {
				theseRegions = append(theseRegions, r.Id)
			}
		}
		sort.Strings(theseRegions)
		if hasWorld {
			// Make sure world is last
			theseRegions = append(theseRegions, "world")
		}

		all = append(all, common.PaymentMethod{
			Id:   pm.Id,
			Type: pm.Type,
			Name: pm.Name,
			Images: common.PaymentMethodImages{
				Small:  toPaymentMethodImage(pm.Id, pm.SmallWidth, pm.SmallHeight, "30"),
				Medium: toPaymentMethodImage(pm.Id, pm.MediumWidth, pm.MediumHeight, "60"),
				Large:  toPaymentMethodImage(pm.Id, pm.LargeWidth, pm.LargeHeight, "120"),
			},
			Regions: theseRegions,
		})
	}
	return all
}

func commonTimezones(data CleansedDataSet) []common.Timezone {
	var all []common.Timezone
	for _, t := range data.Timezones {
		all = append(all, common.Timezone{
			Name:        t.Name,
			Description: t.Description,
			Offset:      t.Offset,
		})
	}
	sortTimezones(all)
	return all
}

func commonCurrencies(data CleansedDataSet, locales []common.Locale) []common.Currency {
	currencyLocales := cleanse.LoadCurrencyLocales()

	var all []common.Currency
	for _, c := range data.Currencies {
		symbols := data.CurrencySymbols[c.Iso_4217_3]

		commonSymbols := &common.CurrencySymbols{}

		if symbols.Primary == "" {
			commonSymbols = nil
		} else {
			commonSymbols = &common.CurrencySymbols{
				Primary: symbols.Primary,
				Narrow:  symbols.Narrow,
			}
		}

		defaultLocale := currencyLocales[c.Iso_4217_3]
		if defaultLocale == "" {
			defaultLocale = defaultLocaleIdForCurrency(data, locales, c)
		}

		all = append(all, common.Currency{
			Name:           c.Name,
			Iso_4217_3:     c.Iso_4217_3,
			NumberDecimals: c.NumberDecimals,
			Symbols:        commonSymbols,
			DefaultLocale:  defaultLocale,
		})
	}
	return all
}

func commonCountries(data CleansedDataSet) []common.Country {
	var all []common.Country
	for _, c := range data.Countries {
		languages := []string{}
		for _, l := range data.Languages {
			if common.ContainsIgnoreCase(l.Countries, c.Iso_3166_3) {
				languages = append(languages, l.Iso_639_2)
			}
		}

		timezones := []string{}
		for _, ct := range data.CountryTimezones {
			if ct.CountryCode == c.Iso_3166_3 {
				tz := findTimezone(data.Timezones, ct.TimezoneCode)
				timezones = append(timezones, tz.Name)
			}
		}

		var defaultLanguage string
		for _, cl := range data.CountryDefaultLanguages {
			if cl.CountryCode == c.Iso_3166_3 {
                lang := findLanguageByCode(data.Languages, cl.LanguageCode)
                if (defaultLanguage != "") {
                    fmt.Printf("ERROR: invalid multiple default language codes for country[%s]\n", cl.CountryCode)
                    os.Exit(1)
                }
                if (!common.Contains(languages, lang.Iso_639_2)) {
                    fmt.Printf("ERROR: default language[%s] is not listed in languages for country[%s]\n", lang.Iso_639_2, cl.CountryCode)
                    os.Exit(1)
                }
				defaultLanguage = lang.Iso_639_2
			}
		}
		if (defaultLanguage == "" && len(languages) > 0) {
    		defaultLanguage = languages[0]
		}

		var defaultCurrency string
		if c.Currency != "" {
			defaultCurrency = findCurrencyByCode(data.Currencies, c.Currency).Iso_4217_3
		}

		var defaultDeliveredDuty string
		for _, d := range data.CountryDuties {
			if strings.ToUpper(d.CountryCode) == strings.ToUpper(c.Iso_3166_3) {
				defaultDeliveredDuty = d.DeliveredDuty
			}
		}

		sort.Strings(languages)
		sort.Strings(timezones)
		all = append(all, common.Country{
			Name:                 formatCountryName(c.Iso_3166_3, c.Name),
			Iso_3166_2:           c.Iso_3166_2,
			Iso_3166_3:           c.Iso_3166_3,
			MeasurementSystem:    getMeasurementSystem(c.Iso_3166_3),
			DefaultCurrency:      defaultCurrency,
			DefaultLanguage:      defaultLanguage,
			Languages:            languages,
			Timezones:            timezones,
			DefaultDeliveredDuty: defaultDeliveredDuty,
		})
	}
	return all
}

func createRegions(countries []common.Country, continents []common.Continent) []common.Region {
	regions := []common.Region{}

	for _, c := range countries {
		id := generateId(c.Iso_3166_3)

		regions = append(regions, common.Region{
			Id:                 id,
			Name:               c.Name,
			Countries:          []string{c.Iso_3166_3},
			Currencies:         currenciesForCountries([]common.Country{c}),
			Languages:          languagesForCountries([]common.Country{c}),
			MeasurementSystems: measurementSystemsForCountries([]common.Country{c}),
			Timezones:          timezonesForCountries([]common.Country{c}),
		})
	}

	for _, c := range continents {
		if c.Code != "ANT" {
			id := generateId(c.Name)

			theseCountries := findCountries(countries, c.Countries)
			regions = append(regions, common.Region{
				Id:                 id,
				Name:               c.Name,
				Countries:          toCountryCodes(theseCountries),
				Currencies:         currenciesForCountries(theseCountries),
				Languages:          languagesForCountries(theseCountries),
				MeasurementSystems: measurementSystemsForCountries(theseCountries),
				Timezones:          timezonesForCountries(theseCountries),
			})
		}
	}

	regions = append(regions, eurozone(countries), world(countries), europeanUnion(countries), europeanEconomicArea(countries))
	assertUniqueRegionIds(regions)
	sortRegions(regions)

	return regions
}

func createProvinces(data CleansedDataSet, locales []common.Locale) []common.Province {
	provinces := []common.Province{}

	// There are too many provinces - only do these countries for now
	validCountries := []string{
		"AUS",
		"CAN",
		"JPN",
		"USA",
		"CHN",
		"ESP",
	}

	for _, p := range data.Provinces {
		// find the country of this province
		country := findCountryByCode(data.Countries, p.CountryCode)

		// construct the province id which we need to lookup the locale translations
		provinceId := country.Iso_3166_3 + "-" + p.Iso_3166_2

		// look for all items in data.ProvinceTranslations where provinceId match and append
		translations := []common.LocalizedTranslation{}
		for _, pt := range data.ProvinceTranslations {
			if common.EqualsIgnoreCase(pt.ProvinceId, provinceId) {
				// find the locale of this translation
				locale := findLocaleById(locales, pt.LocaleId)

				// we have all data now, append the localized translation
				translations = append(translations, common.LocalizedTranslation{
					Locale: locale,
					Name:   pt.Translation,
				})
			}
		}

		// now create
		if common.ContainsIgnoreCase(validCountries, country.Iso_3166_3) {
			provinces = append(provinces, common.Province{
				Id:           provinceId,
				Iso_3166_2:   p.Iso_3166_2,
				Name:         p.Name,
				Country:      country.Iso_3166_3,
				ProvinceType: p.ProvinceType,
				Translations: translations,
			})
		}
	}

	return provinces
}

func eurozone(countries []common.Country) common.Region {
	// EUROZONE is an explicit set of Countries
	// other countries may choose to use the EUR,
	// but that doesn't mean they are part of the EUROZONE!
	// see: https://app.clubhouse.io/flow/story/12716/country-picker-incorrect-countries-are-showing-for-eurzone
	countryCodes := []string{
		"AUT",
		"BEL",
		"CYP",
		"EST",
		"FIN",
		"FRA",
		"DEU",
		"GRC",
		"IRL",
		"ITA",
		"LVA",
		"LTU",
		"LUX",
		"MLT",
		"NLD",
		"PRT",
		"SVK",
		"SVN",
		"ESP",
	}

	theseCountries := findCountries(countries, countryCodes)

	return common.Region{
		Id:                 "eurozone",
		Name:               "Eurozone",
		Countries:          countryCodes,
		Currencies:         currenciesForCountries(theseCountries),
		Languages:          languagesForCountries(theseCountries),
		MeasurementSystems: measurementSystemsForCountries(theseCountries),
		Timezones:          timezonesForCountries(theseCountries),
	}
}

func europeanUnionCountryCodes() []string {
	return []string{"AUT", "BEL", "BGR", "HRV", "CYP", "CZE", "DNK", "EST", "FIN", "FRA", "DEU", "GRC", "HUN", "IRL", "ITA", "LVA", "LTU", "LUX", "MLT", "NLD", "POL", "PRT", "ROU", "SVK", "SVN", "ESP", "SWE", "GBR"}
}

func europeanEconomicAreaCountryCodes() []string {
	return append(europeanUnionCountryCodes(), []string{"ISL", "LIE", "NOR"}...)
}

func europeanUnion(countries []common.Country) common.Region {
	countryCodes := europeanUnionCountryCodes()
	theseCountries := findCountries(countries, countryCodes)
	return common.Region{
		Id:                 "europeanunion",
		Name:               "European Union",
		Countries:          countryCodes,
		Currencies:         currenciesForCountries(theseCountries),
		Languages:          languagesForCountries(theseCountries),
		MeasurementSystems: measurementSystemsForCountries(theseCountries),
		Timezones:          timezonesForCountries(theseCountries),
	}
}

func europeanEconomicArea(countries []common.Country) common.Region {
	countryCodes := europeanEconomicAreaCountryCodes()
	theseCountries := findCountries(countries, countryCodes)
	return common.Region{
		Id:                 "europeaneconomicarea",
		Name:               "European Economic Area",
		Countries:          countryCodes,
		Currencies:         currenciesForCountries(theseCountries),
		Languages:          languagesForCountries(theseCountries),
		MeasurementSystems: measurementSystemsForCountries(theseCountries),
		Timezones:          timezonesForCountries(theseCountries),
	}
}

func world(countries []common.Country) common.Region {
	var codes []string
	for _, c := range countries {
		codes = append(codes, c.Iso_3166_3)
	}
	sort.Strings(codes)

	return common.Region{
		Id:                 "world",
		Name:               "World",
		Countries:          codes,
		Currencies:         currenciesForCountries(countries),
		Languages:          languagesForCountries(countries),
		MeasurementSystems: []string{"metric", "imperial"},
		Timezones:          timezonesForCountries(countries),
	}
}

func findLocaleById(locales []common.Locale, localeId string) common.Locale {
	for _, l := range locales {
		if l.Id == localeId {
			return l
		}
	}
	return common.Locale{}
}

func findCountryByCode(countries []cleanse.Country, code string) cleanse.Country {
	formatted := strings.ToUpper(code)
	for _, c := range countries {
		if c.Iso_3166_2 == formatted || c.Iso_3166_3 == formatted {
			return c
		}
	}
	fmt.Printf("Invalid country code: %s\n", code)
	os.Exit(1)
	return cleanse.Country{}
}

func normalizeCountryCode(countries []cleanse.Country, code string) string {
	formatted := strings.ToUpper(code)
	for _, c := range countries {
		if c.Iso_3166_2 == formatted || c.Iso_3166_3 == formatted {
			return c.Iso_3166_3
		}
	}
	return ""
}

func findCurrencyByCode(currencies []cleanse.Currency, code string) cleanse.Currency {
	for _, c := range currencies {
		if c.Iso_4217_3 == code {
			return c
		}
	}
	fmt.Printf("Invalid currency code: %s\n", code)
	os.Exit(1)
	return cleanse.Currency{}
}

func findContinentByCode(continents []cleanse.Continent, code string) cleanse.Continent {
	for _, c := range continents {
		if c.Code2 == code || c.Code3 == code {
			return c
		}
	}
	fmt.Printf("Invalid continent code: %s\n", code)
	os.Exit(1)
	return cleanse.Continent{}
}

func findLanguageByCode(languages []cleanse.Language, code string) cleanse.Language {
	for _, c := range languages {
		if c.Iso_639_2 == code {
			return c
		}
	}
	fmt.Printf("Invalid language code: %s\n", code)
	os.Exit(1)
	return cleanse.Language{}
}

func normalizeLanguageCode(languages []cleanse.Language, code string) string {
	formatted := strings.ToLower(code)
	for _, c := range languages {
		if c.Iso_639_2 == formatted {
			return c.Iso_639_2
		}
	}
	return ""
}

func findTimezone(timezones []cleanse.Timezone, name string) cleanse.Timezone {
	for _, c := range timezones {
		if strings.ToUpper(c.Name) == strings.ToUpper(name) || strings.ToUpper(c.Description) == strings.ToUpper(name) {
			return c
		}
	}
	fmt.Printf("Invalid timezone name: %s\n", name)
	os.Exit(1)
	return cleanse.Timezone{}
}

func findRegion(regions []common.Region, id string) common.Region {
	for _, c := range regions {
		if strings.ToUpper(c.Id) == strings.ToUpper(id) {
			return c
		}
	}
	fmt.Printf("Invalid region id: %s\n", id)
	os.Exit(1)
	return common.Region{}
}

func findLocaleNameById(names []cleanse.LocaleName, id string) string {
	localeId := common.FormatLocaleId(id)
	for _, n := range names {
		if n.Id == localeId {
			return n.Name
		}
	}
	return ""
}

func getMeasurementSystem(iso_3166_3 string) string {
	if iso_3166_3 == "USA" || iso_3166_3 == "LBR" || iso_3166_3 == "MMR" {
		return "imperial"
	}
	return "metric"
}

func formatCountryName(iso3 string, defaultName string) string {
	switch iso3 {
	case "USA":
		{
			return "United States of America"
		}
	case "GBR":
		{
			return "United Kingdom"
		}
	default:
		{
			return defaultName
		}
	}
}

func toCountryCodes(countries []common.Country) []string {
	codes := []string{}

	for _, country := range countries {
		codes = append(codes, country.Iso_3166_3)
	}

	sort.Strings(codes)
	return codes
}

func findCountries(countries []common.Country, codes []string) []common.Country {
	matching := []common.Country{}

	for _, country := range countries {
		if common.ContainsIgnoreCase(codes, country.Iso_3166_3) {
			matching = append(matching, country)
		}
	}

	return matching
}

func filterCountries(countries []common.Country, filter []string) []common.Country {
	filtered := []common.Country{}

	for _, country := range countries {
		if !common.ContainsIgnoreCase(filter, country.Iso_3166_3) {
			filtered = append(filtered, country)
		}
	}

	return filtered
}

func filterCodes(codes []string, filter []string) []string {
	filtered := []string{}

	for _, code := range codes {
		if !common.ContainsIgnoreCase(filter, code) {
			filtered = append(filtered, code)
		}
	}

	return filtered
}

/**
 * Filters the list of countries to those with a matching default currency.
 */
func findCountriesByCurrency(countries []common.Country, currency string) []string {
	codes := []string{}

	for _, country := range countries {
		if country.DefaultCurrency == currency && !common.ContainsIgnoreCase(codes, country.Iso_3166_3) {
			codes = append(codes, country.Iso_3166_3)
		}
	}

	sort.Strings(codes)
	return codes
}

func currenciesForCountries(countries []common.Country) []string {
	codes := []string{}

	for _, country := range countries {
		if country.DefaultCurrency != "" && !common.ContainsIgnoreCase(codes, country.DefaultCurrency) {
			codes = append(codes, country.DefaultCurrency)
		}
	}

	sort.Strings(codes)
	return codes
}

func languagesForCountries(countries []common.Country) []string {
	codes := []string{}

	for _, country := range countries {
		for _, l := range country.Languages {
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

	for _, country := range countries {
		for _, tz := range country.Timezones {
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

	for _, country := range countries {
		if country.MeasurementSystem != "" && !common.ContainsIgnoreCase(codes, country.MeasurementSystem) {
			codes = append(codes, country.MeasurementSystem)
		}
	}

	sort.Strings(codes)
	return codes
}

func defaultLocaleIdForCurrency(data CleansedDataSet, locales []common.Locale, currency cleanse.Currency) string {
	countries := []string{}
	languages := []string{}

	for _, c := range data.Countries {
		if c.Currency == currency.Iso_4217_3 {
			countries = append(countries, c.Iso_3166_3)

			for _, l := range data.Languages {
				if common.Contains(l.Countries, c.Iso_3166_3) {
					languages = append(languages, l.Iso_639_2)
				}
			}
		}
	}

	// Return first matching locale
	for _, l := range locales {
		if common.Contains(countries, l.Country) && common.Contains(languages, l.Language) {
			return l.Id
		}
	}

	return ""
}

func assertUniqueRegionIds(regions []common.Region) {
	found := make(map[string]bool)

	for _, r := range regions {
		if found[r.Id] {
			fmt.Printf("ERROR: Duplicate region id[%s]\n", r.Id)
			os.Exit(1)
		}
		found[r.Id] = true
	}
}

func uniqueLocaleIds(locales []common.Locale) []common.Locale {
	unique := []common.Locale{}

	found := make(map[string]bool)

	for _, l := range locales {
		if !found[l.Id] {
			unique = append(unique, l)
		}
		found[l.Id] = true
	}

	return unique
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

func sortLocales(locales []common.Locale) []common.Locale {
	slice.Sort(locales[:], func(i, j int) bool {
		return strings.ToLower(locales[i].Name) < strings.ToLower(locales[j].Name)
	})
	return locales
}

// https://stackoverflow.com/a/36122804/2297665
type byDescriptionAndName []common.Timezone

func (tz byDescriptionAndName) Len() int      { return len(tz) }
func (tz byDescriptionAndName) Swap(i, j int) { tz[i], tz[j] = tz[j], tz[i] }

// Sort by description first, then name (description is not unique)
func (tz byDescriptionAndName) Less(i, j int) bool {
	if strings.ToLower(tz[i].Description) < strings.ToLower(tz[j].Description) {
		return true
	}
	if strings.ToLower(tz[i].Description) > strings.ToLower(tz[j].Description) {
		return false
	}

	return strings.ToLower(tz[i].Name) < strings.ToLower(tz[j].Name)
}

func sortTimezones(timezones []common.Timezone) []common.Timezone {
	sort.Sort(byDescriptionAndName(timezones))
	return timezones
}

func toPaymentMethodImage(id string, width int, height int, size string) common.PaymentMethodImage {
	url := fmt.Sprintf("https://flowcdn.io/util/logos/payment-methods/%s/%s/original.png", id, size)

	return common.PaymentMethodImage{
		Url:    url,
		Width:  width,
		Height: height,
	}
}
