# TODO

## Introduction

This project is a web application that allows users to create an account and login to access a task management platform.
Users can add, edit and delete tasks and can also set the completed status of tasks.

The server is written using Golang and Go Fiber along with GORM to connect and create models to a Postgres database.
The frontend utilizes Vanilla JS and Bootstrap along with GO HTML templating for server-side rendering.
Session middleware is also utilized using the Go Fiber session library and Redis storage caching.

## Getting Started

### 1. Install Docker using their website and install a Postgres image and a Redis image

docker pull postgres
docker pull redis

Run container for the first time:

`docker run -name %NAME% -e POSTGRES_PASSWORD=%PASSWORD% -p %PORT%:%PORT% -d postgres`

- Add your custom name, password, port (fill personal %variable%) and change app.env file accordingly (DB variables)

`docker run -it --name %NAME% -p %PORT%:%PORT% -d redis`

To start existing container:

`docker start -a %NAME%`

If you would like the default configurations, simply use the app.env variables

### 2. Install Golang (specific version 1.19) using their site

### 3. Start application

`go run server.go`

