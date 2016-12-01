# CC User API

[![Build Status](https://travis-ci.org/arbolista-dev/cc-user-api.svg?branch=master)](https://travis-ci.org/arbolista-dev/cc-user-api)


## Setup

1. Copy and configure your environment variables (see definition below):

```sh
cp .env.example .env
export $(cat .env | xargs)
```
2. Build and start stack (Postgres & User API)

```
make stack-up
```

or manually run in this order:
```
make images
make create-db
make migrate-db
make create-api
```

When using Rancher the container names are named differently. The container id's specified in the Makefile can easily be replaced doing a ('r-Default_userapi_1' specifying the existing name):
```sh
sed -i 's/user_api/r-Default_userapi_1/g' Makefile
```

*Address is 127.0.0.1:8082*

## Enviromental variables:
```
CC_ENV - Defines environment to decide which commands are run in entrypoint [dev,test,prod]
CC_DBHOST - Name of database host (should be PG docker link name)
CC_DBNAME - Name of the postgres DB
CC_DBUSER - Name of the postgres user
CC_DBPASS - Password of the user
CC_DBADDRESS - Address of the postgres service

CC_JWTSIGN - A secret string to sign JWT

CC_SENDGRID_APIKEY - The key of the SendGrid account
CC_SENDGRID_TEMPLATE_CONFIRM - The template id for confirmation emails, must have a subtitution param named link
CC_SENDGRID_TEMPLATE_RESET - The template id for reset passwords, must have a subtitution param named link
CC_CONFIRMATION_MAIL - The email used as sender by SendGrid

CC_SERVER_HOST - Host used to create the URL
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
GET     /user/passreset    // Show password reset page
GET     /user/confirm      // Confirm the email account of the users

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
