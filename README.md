# CC User API

[![Build Status](https://travis-ci.org/arbolista-dev/cc-user-api.svg?branch=master)](https://travis-ci.org/arbolista-dev/cc-user-api)


## Setup

Copy and configure your environment variables:

```sh
cp .env.example .env
export $(cat .env | xargs)
```

Run in this order:
```
make images

make create-database

make create-api
```

*Address is 127.0.0.1:8082*

Enviromental variables:
```
CC_DBNAME - Name of the postgres DB
CC_DBUSER - Name of the postgres user
CC_DBPASS - Password of the user
CC_DBADDRESS - Address of the postgres service

CC_JWTSIGN - A secret string to sing JWT

CC_SPARKPOSTKEY - The key of the sparkpost app
```

## Routes
```
POST    /user              // Add a new user
POST    /user/login        // User login
GET     /user/logout       // Logout from current session
GET     /user/logoutall    // Logout user from all sessions
DELETE  /user              // Delete user
PUT     /user              // Update user (name or email)
PUT     /user/answers      // Update user answers
PUT     /user/location     // Set user location (city, county, state, country)
POST    /user/reset/req    // Request a password reset -> send email to user
POST    /user/reset        // Confirm newly set password
GET     /user/leaders      // Return leaders (paginated)
GET     /user/locations    // Return available locations
GET     /page/passreset    // Show password reset page
```

## Return types
### Success:
```
{
  success: true,
  data: ...
}
```

### Error:
```
{
  success: false,
  message: “{field: error-code}”
}
```

## Endpoints
### Adding a user
Use:
```
POST     /user
```
Body:
```
{
  "first_name": "Juan",
  "last_name" : "Perez",
  "email" : "juanp@test.testy",
  "password": "juanito",
  "city": "Austin",
  "county": "Hays",
  "state": "Texas",
  "country": "us",
  "public": true
}
```
Validation on the fields:
```
"first_name" : required and min size 4
"last_name" : required and min size 4
"password" : required and min size 4
"email" : required and must math format regex
```

### Logging in
Use:
```
POST    /user/login
```
Body:
```
{
  "email": "jdoe@test.testy",
  "password": "johnyboy"
}
```
If successful returns:
```
{
  "success": true,
  "data": {
    "answers": "eyJDTzIiOiAibG90cyIsICJjaXR5IjogIkdETCIsICJtb25leSI6ICJub25lIn0=",
    "name": "John",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0Njc2Nzc4OTEsImlhdCI6MTQ2NjQ2ODI5MSwiaWQiOjY4LCJqdGkiOiJWZE44MyJ9.u-QfbyuieTRyiuqYIbxb01F0I1qdNUamQY4yMItrMhU",
  }
}
```

### Update user
Use:
```
PUT     /user

HTTP Headers:
Authorization: <token>
```
Body:
```
{
  "email": "juanpp@test.testy",
  "password": "juanito"
}
```

### Update user answers
Use:
```
PUT     /user/answers

HTTP Headers:
Authorization: <token>
```
Body:
```
{
  "answers":{"result_food_total": "5", "result_housing_total": "6", "result_services_total": "3", "result_goods_total": "4", "result_transport_total": "8", "result_grand_total": "26"}
}
```

### Set user location
Use:
```
PUT     /user/location

HTTP Headers:
Authorization: <token>
```
Body:
```
{
  "city":"Brooklyn","county":"Kings","state":"New York","country":"us"
}
```

### List leaders
Use:
```
GET     /user/leaders

```
Body:
```
{
  limit: 20, offset: 0, category: "Footprint", city: "", state: ""
}
```

### List locations
Use:
```
GET     /user/locations

```
