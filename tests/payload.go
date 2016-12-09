package tests

var userBody = `{
  "first_name": "Juan",
  "last_name" : "Perez",
  "email" : "juanp@test.testy",
  "password": "juanito",
  "city": "Austin",
  "county": "Hays",
  "state": "Texas",
  "country": "us",
  "public": true
}`

var userBody_update = `{
  "first_name": "Juanito",
  "last_name" : "Perezin",
  "profile_data": {
    "facebook": "UC-Berkeley-CoolClimate-Network-161909000540511",
    "twitter": "coolclimatenw",
    "instagram": "coolcalifornia",
    "linkedin": "company/1293183",
    "medium": "@nature_org",
    "intro": "Hello, I'm a climate activist"
  }
}`

var loginBody = `{
  "email": "juanp@test.testy",
  "password": "juanito"
}`

var loginBody_badEmail = `{
  "email": "juanpp@test.testy",
  "password": "juanito"
}`

var loginBody_badPassword = `{
  "email": "juanp@test.testy",
  "password": "juanito2"
}`

var answers_update = `{"answers":{"result_food_total": "5", "result_housing_total": "6", "result_services_total": "3", "result_goods_total": "4", "result_transport_total": "8", "result_grand_total": "26", "input_size": "3"}}`

var location_set = `{"city":"Brooklyn","county":"Kings","state":"New York","country":"us"}`

var profile = `{
  "first_name": "Juanito",
  "last_name" : "Perezin",
  "city": "Brooklyn",
  "county": "Kings",
  "state": "New York",
  "household_size": "3",
  "total_footprint":{"result_food_total": "5", "result_housing_total": "6", "result_services_total": "3", "result_goods_total": "4", "result_shopping_total": "7", "result_transport_total": "8", "result_grand_total": "26"}
}`
