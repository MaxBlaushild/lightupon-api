# lightupon-api

![alt text](logo.png?raw=true "Lightupon")

## What is this?

Lightupon will get you off of your couch and into the streets. The app allows single or many people to band together and travel through narrated journies happening in real locations. This repo contains the current build of API. This service is written in Go. It connects to a postgres database

#### Dependenices:

- Go
- Godeps
- Postgres

## How to run:

Clone this repository into the src folder of your configured go workspace. Dependencies should be included in the Godeps folder, but just incase, navigate to the root and run:

```
go get
```

To run the app, run:

```
go run main.go
```

Tables will be automigrated into your database at that point.
