package javascript

import (
	"../common"
	"errors"
	"fmt"
	"os"
)

type CurrencyFormat struct {
	Formats      map[string]JavascriptFormat `json:"formats"`
	Locales      map[string]LocaleDictionary  `json:"locales"`
}

type JavascriptFormat struct {
	Symbol     string `json:"symbol"`
	Decimal    string `json:"decimal"`
	Group      string `json:"group"`
	Precision  int    `json:"precision"`
	Format     string `json:"format"`
}

type LocaleDictionary struct {
	Format     string `json:"format"`
}

type CommonData struct {
	Countries         []common.Country
	Currencies        []common.Currency
	Locales           []common.Locale
}

func Generate() {
	data := CommonData{
		Countries: common.Countries(),
		Currencies: common.Currencies(),
		Locales: common.Locales(),
	}

	format := CurrencyFormat{
		Formats: generateFormats(data),
		Locales: generateLocaleDictionary(data),
	}
	common.WriteJson("data/javascript/currency-format.json", format)
}

func generateFormats(data CommonData) map[string]JavascriptFormat {
	all := map[string]JavascriptFormat{}

	for _, l := range data.Locales {
		currency, err := findCurrencyByLocale(data, l)
		if err == nil && currency.Symbols != nil {
			all[l.Id] = JavascriptFormat{
				Symbol: currency.Symbols.Primary,
				Decimal: l.Numbers.Decimal,
				Group: l.Numbers.Group,
				Precision: currency.NumberDecimals,
				Format: "%s%v",
			}
		}
	}

	return all;
}

func generateLocaleDictionary(data CommonData) map[string]LocaleDictionary {
	all := map[string]LocaleDictionary{}
	for _, l := range data.Locales {
		all[l.Id] = LocaleDictionary{
			Format: "1",
		}
	}
	return all
}

func findCurrencyByLocale(data CommonData, locale common.Locale) (common.Currency, error) {
	country := findCountryByIso31663(data, locale.Country)
	if country.DefaultCurrency == "" {
		return common.Currency{}, errors.New("Country has no default currency")
	} else {
		c := findCurrencyByIso42173(data, country.DefaultCurrency)
		return c, nil
	}
}

func findCountryByIso31663(data CommonData, country string) common.Country {
	for _, c := range data.Countries {
		if c.Iso_3166_3 == country {
			return c
		}
	}

	fmt.Printf("ERROR: Country[%s] not found\n", country)
	os.Exit(1)
	return common.Country{}
}

func findCurrencyByIso42173(data CommonData, currency string) common.Currency {
	for _, c := range data.Currencies {
		if c.Iso_4217_3 == currency {
			return c
		}
	}

	fmt.Printf("ERROR: Currency[%s] not found\n", currency)
	os.Exit(1)
	return common.Currency{}
}
