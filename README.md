# TradingACE in Golang

This project uses Gin, the performance and productivity Golang Framework.

If you want to learn more about Gin, please visit its website: https://gin-gonic.com/.

## Preparing requirements

1. Set ethereum env:
    1. Find file for config.yaml
    2. Set `your_project_key` of infura in ethereum url: wss://mainnet.infura.io/ws/v3/{your_project_key}

2. Set campaign starttime:
    1. Find file for config.yaml
    2. Set `starttime` of campaign in future or past that will switch different mode between `CurrentActiveMode` and
       `PastBacktestMode`

3. Setup database:
    * When you use docker to run the application, docker can automate create database and data table.
    * you only need to set config
        1. Find file for config.yaml
        2. Set `host` of database: `postgres`

    * if you doesn't use docker to run the application
        1. Setup PostgresSQL
        2. Create Database: trading_ace
        3. Running init creating data table in path: ./schema
        4. Set `host` of database: `localhost`

## Packaging and running the application

The application can be build using:

```shell script
go build 
```

* and then will generate a executable file

Running the application using exe

```shell script
./tradingACE.exe
```

## Running the application

You can run your application directly:

```shell script
go run main.go
```

## Running the application with docker

Build docker image:

```shell script
docker-compose build --no-cache
```

Start your docker container:

```shell script
docker-compose up -d
```

Rebuild docker image without downtime:

```shell script
docker-compose up -d --no-deps --build
```

Display system log in docker container:

```shell script
docker-compose logs -f --tail 100
```

Stop your docker container

```shell script
docker-compose down
```