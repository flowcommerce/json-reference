# reference-json

Reference data for countries, currencies, languages, timezones, and
regions in JSON format. We collect information from multiple sources,
as well as provide our own additional metadata, and combine into a
simple, canonical set of reference data that we use throughout Flow.

## JSON Data

New reference data should be added to csv files in [data/original](https://github.com/flowcommerce/json-reference/tree/main/data/original)

To generate JSON data execute:

   `go run reference.go all`

The script might create a lot of unrelated changes (removed currencies, locales, etc) that can be breaking.
Therefore, every modification that is not desirable should be manually reverted before merging changes.

The final JSON data files can be found in the
[data/final](https://github.com/flowcommerce/json-reference/tree/main/data/final)
directory.

  - [Continents](https://github.com/flowcommerce/json-reference/blob/main/data/final/continents.json)
    A list of continents and the countries they contain

  - [Countries](https://github.com/flowcommerce/json-reference/blob/main/data/final/countries.json)
    A list of countries, including metadata on their measurement
    system, default currency, languages, and timezones

  - [Currencies](https://github.com/flowcommerce/json-reference/blob/main/data/final/currencies.json)
    A list of currencies, including metadata for localization

  - [Languages](https://github.com/flowcommerce/json-reference/blob/main/data/final/languages.json)
    A list of languages and the countries in which they are spoken

  - [Locales](https://github.com/flowcommerce/json-reference/blob/main/data/final/locales.json)
    A list of locales and specific number formats

  - [Payment Methods](https://github.com/flowcommerce/json-reference/blob/main/data/final/payment-methods.json)
    A list of all the payment methods supported by Flow

  - [Regions](https://github.com/flowcommerce/json-reference/blob/main/data/final/regions.json)
    A region represents a geographic area of the world. Regions can be countries, continents or other political areas (like the Eurozone)

  - [Timezones](https://github.com/flowcommerce/json-reference/blob/main/data/final/timezones.json)
    A list of timezones

## Related libraries

  - [Scala Library](https://github.com/flowcommerce/lib-reference-scala)
  - [JavaScript Library](https://github.com/flowcommerce/lib-reference-javascript)


## Local development

We rely on a git submodule to pull in the `cldr-json` project. Before
running the underlying commands, first run:


```
git submodule init
git submodule update
```

View commands available:

  `go run reference.go`

Run end to end process:

  `go run reference.go all`

More details can be found [here](https://www.notion.so/flow/References-bd8b9b8f5c434d21aa0bf1c0b98e6d66)
