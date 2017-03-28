package javascript

import (
	"../common"
	"fmt"
	"os"
)

type CurrencyFormat struct {
	Formats      map[string]JavascriptFormat `json:"formats"`
	Locales      map[string]LocaleDicionary  `json:"locales"`
}

type JavascriptFormat struct {
	Symbol     string `json:"symbol"`
	Decimal    string `json:"decimal"`
	Group      string `json:"group"`
	Precision  string `json:"precision"`
	Format     string `json:"format"`
}

type LocaleDicionary struct {
	Locale     string `json:"locale"`
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
		all[l.Id] = JavascriptFormat{
			Symbol: l.Symbol,
			Decimal: l.Numbers.Decimal,
			Group: l.Numbers.Group,
			Precision: 2,
			Format: "%s%v",
		}
	}

	return all;
}

func generateLocaleDictionary(data CommonData) map[string]LocaleDicionary {
	all := map[string]LocaleDicionary{}
	for _, l := range data.Locales {
		if c.DefaultCurrency != "" {
			currency := findCurrencyByIso42173(data, c.DefaultCurrency)
			if currency.DefaultLocale != "" {
				all[c.Iso_3166_3] = JavascriptCountry{
					Locale: currency.DefaultLocale,
				}
			}
		}
	}
	return all
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
