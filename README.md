# reference-json

Reference data for countries, currencies, languages, timezones, and
regions in JSON format. We collect information from multiple sources,
as well as provide our own additional metadata, and combine into a
simple, canonical set of reference data that we use throughout Flow.

## JSON Data

The JSON data files can be found in the [data/final](/flowcommerce/json-reference/blob/master/data/final) directory.

## Related libraries

  - [Scala Library](/flowcommerce/lib-reference-scala)
  - [JavaScript Library](/flowcommerce/lib-reference-javascript)

  `go run script/reference.go all`

## Generating the data

  `go run script/reference.go all`
