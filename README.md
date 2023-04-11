# storage
a scalable key-value storage that developed base on Clean Architecture and DDD.

* *Gin* as web framework
* *Postgres* as persistent storage
* *BigCache* as in-memory cache layer

![storage drawio](https://user-images.githubusercontent.com/45311375/230890628-dd533cd1-e176-4df5-98e1-9bc17634b61f.png)

## How To Run This Project

### Docker Compose

Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into your workspace
$ git clone https://github.com/jamshidi799/storage.git

#move to project
$ cd storage

$ docker-compose up
```

### Build Project

first you should download dependency
```bash
$ go mod download
```

then create `.env` file in root of project and put this value on it
(change values base on your setup)

```
ENVIRONMENT=dev

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=storage

JWT_SECRET=secret
```

this config assumes that a postgres database listen on `localhost:5432` with that configs.

then run the project:
```bash
$ go run main.go 
```

project will listen on http://localhost:8080

## swagger

you can find OpenApi spec on http://localhost:8080/swagger/index.html

for **updating** swagger spec you should first _swag_ with this:
```bash
$ go install github.com/swaggo/swag/cmd/swag@latest
```

then run this command in root of project:

```bash
$ swag init
```





