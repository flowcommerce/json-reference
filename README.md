# reference-json

Reference data for countries, currencies, languages, timezones, and
regions in JSON format. We collect information from multiple sources,
as well as provide our own additional metadata, and combine into a
simple, canonical set of reference data that we use throughout Flow.

## JSON Data

JSON data files can be found in the
[data/final](/flowcommerce/json-reference/tree/master/data/final)
directory.

  - [Continents](/flowcommerce/json-reference/blob/initial_scripts/data/final/continents.json)
    A list of continents and the countries they contain

  - [Countries](/flowcommerce/json-reference/blob/initial_scripts/data/final/countries.json)
    A list of countries, including metadata on their measurement
    system, default currency, languages, and timezones

  - [Currencies](/flowcommerce/json-reference/blob/initial_scripts/data/final/currencies.json)
    A list of currencies, including metadata for localization

  - [Languages](/flowcommerce/json-reference/blob/initial_scripts/data/final/languages.json)
    A list of languages and the countries in which they are spoken

  - [Regions](/flowcommerce/json-reference/blob/initial_scripts/data/final/regions.json)
    A region represents a geographic area of the world. Regions can be countries, continents or other political areas (like the Eurozone)

  - [Timezones](/flowcommerce/json-reference/blob/initial_scripts/data/final/timezones.json)
    A list of timezones

## Related libraries

  - [Scala Library](/flowcommerce/lib-reference-scala)
  - [JavaScript Library](/flowcommerce/lib-reference-javascript)

  `go run script/reference.go all`

## Generating the data

  `go run script/reference.go all`
