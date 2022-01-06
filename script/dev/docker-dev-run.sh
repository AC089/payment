#!/bin/bash

docker-compose -f script/dev/payment-dev-compose.yml down -t 3000   

docker-compose -f script/dev/payment-dev-compose.yml up -d 