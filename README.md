# small-api
A small API sample using Goland, MySql and Redis

```
git clone git@github.com:panospet/small-api.git
```
## Preparation
### Mysql and Redis setup

#### Using Docker and docker-compose

If you have `docker` and `docker-compose` installed on your machine, the easiest way to start the `mysql` and `redis`
instances is by simply running the command:
```
docker-compose up -d
```
After that:
- a mysql instance with an empty database `bestprice` runs at port 3305. You can use username `bestprice`, password 
`bestprice` to connect to it, for example: 
```
mysql -h 127.0.0.1 -ubestprice -P 3305 -pbestprice bestprice
``` 
- an empty `redis` instance runs at port 6980
```
redis-cli -p 6380
```

#### Alternatively
Supposing you have already running MySql and Redis instances:
- you only need to create a MySql database `bestprice`:
Connect to your MySql:
```
mysql -h {host} -u{username} -P {post} -p{password}
```
and then give the command:
```
mysql> CREATE DATABASE bestprice;
```

And that's it. Next thing we need to do is to populate our MySql and Redis with data.

### Migrate and populate db
#### Use go migrate
Note: [You can install go migrate from here](https://github.com/golang-migrate/migrate)
Simply, from the parent folder of the project, run:
```
{directory where migrate executable is}/migrate -path ./migrations -database "mysql://{username}:{password}@tcp({mysql host}:{post})/{database name}" up
# example: migrate -path ./migrations -database "mysql://bestprice:bestprice@tcp(localhost:3305)/bestprice" up 
```
The command above creates all necessary tables (`user`, `product`, `category`).
After that, to populate MySql and Redis with some data, simply run the populate script below. 
- `-workers` is an int parameter which represents the number of goroutines to use to speed up the procedure 
(default value 10). 
- `-amount` is an int variable which represents the total amount of products to add (default value 1000).
- Will also add 18 different categories. 
```
cd cmd/populate
go run main.go -workers {number of workers} -amount {number of products}
# example: go run main.go -workers 10 -amount 1000
```

#### Alternatively, to migrate/populate MySql, import dump file
`populate.sql` is a `mysqldump` file which creates all necessary tables, together with 18 categories and 1000 products 
inside.
```
mysql -h {host} -u{username} -P {post} -p{password} bestprice < populate.sql
# example: mysql -h 127.0.0.1 -ubestprice -P 3305 -pbestprice bestprice < populate.sql
```
Now that we populated MySql with some data, we also need to fill Redis with data as well.
To do so, simply run the populate script, but only for Redis (cause MySql already has some data inside):
```
cd cmd/populate
go run main.go -workers {number of workers} -redis-only
example: go run main.go -workers 10 -redis-only
```

After this step, our MySql has products and categories ready to work with. Additionally, Redis has two keys, `product`
and `category`, with their IDs as keys and their JSON representation as values.

#### Concerning Redis caching
*For testing purposes, in the above steps I decided to cache all products/categories to Redis individually. This 
practically means that, if there are millions of them, the memory consumption might increase a lot. You can read more
at the "Caching drawbacks" section at the end of this readme file.*

### Create User
We need to create a user who's able to perform POST, PUT, DELETE requests. To do that, simply run:
```
cd cmd/user-add
go run main.go -username {username} -password {password}
# example: go run main.go -username admin -password admin
```
Passwords are stored encrypted in the database.

### Tests
To see if everything works, we can run the project unit tests. There are currently unit tests for API, pagination, and
database functionality. The command below runs all of them at once:
```
go test -v ./...
```
If all tests are green, it means we are ready to go.

## Run API
### Environment variables for MySql and Redis
To give the API the ability to perform requests to MySql and Redis, we need to set two environment variables, 
`MYSQL_PATH`, `REDIS_PATH` for the corresponding paths.
```
MYSQL_PATH="{username}:{password}@({host}:{port})/{database}?parseTime=true"
REDIS_PATH="redis://{host}:{port}/1"
# example: MYSQL_PATH="bestprice:bestprice@(localhost:3305)/bestprice?parseTime=true"
# example: REDIS_PATH="redis://localhost:6380/1"
```


### Finally, let's start the API! 
```
cd cmd/api
MYSQL_PATH="{mysql_path}" REDIS_PATH="{redis_path}" go run main.go
```
The API starts at port 8080.
### Basic requests
#### Health check
```
curl -XGET "http://localhost:8080/"
```
```
{
    "Message": "health good!"
}
```
And this is our main check that our API is up and running!

#### Authentication
All POST, PATCH, DELETE requests need basic authentication. Please use the username/password of the user you created
in the previous steps. 

### Categories requests
#### Get Categories
```
curl -XGET "http://localhost:8080/v1/categories"
```
#### Get Category By Id
```
curl -XGET "http://localhost:8080/v1/categories/1"
```
#### Create Category
```
curl -XPOST -u admin:admin 'http://localhost:8080/v1/categories' -H 'Content-Type: application/json' \
 -d '{"title":"test category", "image_url":"http:\/\/www.bestprice.gr/test_cat.png", "position":10}'
```
If response code is 201, then category has been created successfully.
#### Update Category
```
curl -XPATCH -u admin:admin 'http://localhost:8080/v1/categories/1' -H 'Content-Type: application/json' -d '{"title":"updated"}'
```
#### Delete Category
### Products requests
#### Get Products
```
curl -XGET "http://localhost:8080/v1/products"
```
#### Get Product By Id
```
curl -XGET "http://localhost:8080/v1/products/{product_uuid}"
```
#### Create Product
```
curl -XPOST -u admin:admin 'http://localhost:8080/v1/products' -H 'Content-Type: application/json' \
-d '{"category_id":12, "title":"my test product", "image_url":"http:\/\/www.bestprice.gr/test.png", "price":10, "description":"test description"}'
```
If response code is 201, then product has been created successfully.
#### Update Product
#### Delete Product

### Pagination examples
### Caching method explained

#### Caching drawbacks
- Currently we store ALL individual products + categories. If there are millions of them, this may lead to huge memory
allocation. Possible solutions: smarter caching of specific requests needed, Redis eviction policy
- Max of 15min interval between update individual requests and list requests for the same product
- In case of too many requests, redis may overload (too many goroutines). Possible solutions: job queue, rate limit