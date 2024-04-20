# big-market-go

## Overview

This is a lottery module implemented using go-zero. The front project in https://github.com/HSn0918/big-market-front

## Environment
- Mysql:8.0
- Redis
- go1.22
  
## Notice
Before start the project, import the SQL file and change/app/strategy/raffle/CMD/API/etc/strategy-api.yaml

sql file in ./db/big_market.sql

Execute strategy_armory API before executing random_raffle API
## Quick Start
```shell
git clone https://github.com/HSn0918/big-market-go.git
cd big-market-go
go mod tidy
cd app/strategy/raffle/api
go run strategy.go
```
## apis
## POST random raffle

POST /api/v1/raffle/random_raffle

> Body request

```json
{
  "strategy_id": 100001
}
```

## GET strategy armory

GET /api/v1/raffle/strategy_armory

> Body request

```json
{
  "strategy_id": 100001
}
```
## POST query raffle award list

POST /api/v1/raffle/query_raffle_award_list

> Body request

```json
{
  "strategy_id": 100001
}
```


