NBP api
=======

API for NBP currencies - project based on Golang.

## Usage
Send GET request to https://nbp-api.herokuapp.com/ giving date, type of data, and currency code.
- `Date` - `RRRR-MM-DD`
- `Type` - `avg` or `both`
- `Code` - `*` for all or specific codes like `USD,EUR,GBP` (multiple currencies should be separated by comma)

Example calls:
- `https://nbp-api.herokuapp.com/2015-11-25/avg/*` - get's average currency rates for all currencies
- `https://nbp-api.herokuapp.com/2015-11-25/both/USD,EUR` - get's buy, sell values for USD and EUR

## Live demo
Check the [http://karolgorecki.pl/nbp-api/](http://karolgorecki.pl/nbp-api/)  
*Could be a little bit slow in the beginig (using free heroku account for API - it needs to sleep)*

## todo
there is a lot to imporeve... Still working on it ;)