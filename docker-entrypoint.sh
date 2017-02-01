#!/bin/bash

if [ "$1" == "dev" ]; then
  cd /go/src/github.com/arbolista-dev/cc-user-api || exit; goose -env development up
  revel run github.com/arbolista-dev/cc-user-api
elif [ "$1" == "test" ]; then
  cd /go/src/github.com/arbolista-dev/cc-user-api || exit; goose -env test up
  revel test github.com/arbolista-dev/cc-user-api
elif [ "$1" == "prod" ]; then
  cd /go/src/github.com/arbolista-dev/cc-user-api || exit; goose -env prod up
  revel run github.com/arbolista-dev/cc-user-api prod
fi
