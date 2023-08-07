# Diksha x Lion SuperIndo challange

This guide will help you set up a  and run this project.

## About this project
this is application for simulate database transaction between two database in Golang
for candidate Hafid Abi D.

## Prerequisites

Before proceeding, make sure you have the following installed on your system:

- Go: [Official Installation Guide](https://golang.org/doc/install)
- PostgreSQL: [Official Downloads](https://www.postgresql.org/download)

## Features
- REST API with gin
- GORM ORM
- Dependency injection (not yet implement wire for enhance Dependency Injecton process)
- Swaggo for documentation

## How to run this project

   copy config.json.example into config.json then
    ```
    go build -o main && ./main
    ```
   or you can use docker to run the app
    ```
    docker-compose up
    ```