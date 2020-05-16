# small-api
A small API sample using Goland, MySql and Redis

```
git clone git@github.com:panospet/small-api.git
cd small-api
go mod tidy
```

if machine has docker-compose then
```
docker-compose up -d
```
else create redis + mysql

```
migrate -path ./migrations -database "mysql://bestprice:bestprice@tcp(localhost:3305)/bestprice" up
```


```
mysql -h 127.0.0.1 -ubestprice -P 3305 -pbestprice bestprice
```

populate database
```
cd cmd/populate
go run main.go
```
or import mysqldump (todo)


to run api:
```
cd cmd/api
go run main.go
```

todos:
- order by category position also
- validation and 400
- DDL
- fix pagination testing
- make readme a proper readme