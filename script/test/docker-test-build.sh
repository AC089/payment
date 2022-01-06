#!/bin/bash
docker build --build-arg env=test -t payment-test .

docker-compose -f script/test/payment-test-compose.yml down -t 3000   

docker-compose -f script/test/payment-test-compose.yml up -d 