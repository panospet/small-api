# small-api
A small API sample using Goland, MySql and Redis

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