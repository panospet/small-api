# small-api
A small API sample using Goland, MySql and Redis

```
git clone https://github.com/panospet/small-api.git
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
If response code is 201, then category has been updated successfully.
#### Delete Category
```
curl -XDELETE -u admin:admin "http://localhost:8080/v1/categories/20"
```
Please notice, that if there are still products in the database that use this category ID, then the request results
to a conflict error. In other cases, a 200 response is returned.
 
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
```
curl -XPATCH -u admin:admin 'http://localhost:8080/v1/products/059f9348-86a3-40c9-a2b0-a586f776619c' \
 -H 'Content-Type: application/json' -d '{"price":100, "description":"updated"}'
```
If response code is 201, then product has been created successfully.
#### Delete Product
```
curl -XDELETE -u admin:admin "http://localhost:8080/v1/products/1c8c7393-5ccd-4270-9e1e-aa6ba5c43dae"
```
If response code is 200, then product has been updated successfully.

### Pagination, orderBy, limit, offset examples
There are 5 different query parameters that we can use, while performing `GET` requests for products or categories.
- `perPage`: How many elements per page will be showed. Default value is 10. Example: `/v1/products?perPage=20`
- `page`: The page number to show. For example, if we have 20 elements to show, and perPage value is 5, our data will 
be spread in 4 pages. Example: `/v1/products?page=4`. Default value is 1. 
- `orderBy`: The field based on which our elements are sorted. Its format is `field:order`. For example: 
`/v1/products/orderBy=price:asc`. Default value is `id`. Basically products/categories can be sorted based on any field.
Note: products can also be sorted by their category position (example: `/v1/products?orderBy=position`)
- `limit`: Limit the results to a specific amount. Example: `/v1/products?limit=100`
- `offset`: The amount of results at the beginning of the list that will be "ignored". Example: `/v1/categories?offset=10`,
the first 10 categories will not be showed.

Of source, all above query parameters can be combined. Example request: *Give me the 100 cheapest products, divided to
20 products per page:* 
```
http://localhost:8080/v1/products?perPage=20&orderBy=price:asc&limit=100
```

### Caching method explained
First of all, let's start by saying that caching is always a long and difficult discussion. To find the optimal way of
caching your data, it needs analysis of the usage of the application, where and when the majority of the requests happen,
 use cases etc.
 
The caching methods that I'm using in this project, are two:
- individual category / product caching by ID as key
- serialized response caching by request path as key

#### Individual category / product caching by ID
- `cmd/populate` script, apart from the database, also populates Redis with all products / categories. Key used is their 
ID, and the value is their serialized json.
- After each create/update/delete successful database operation in a product/category, a goroutine is fired to do the 
same in our cache as well.

`redis-cli` command and result:
```
127.0.0.1:6380[1]> hget category 1
"{\"id\":1,\"title\":\"updated\",\"position\":17,\"image_url\":\"http://www.bestprice.gr/sports.png\",
\"created_at\":\"2020-05-17T10:57:00Z\",\"updated_at\":\"2020-05-18T08:07:40Z\"}"
```
```
127.0.0.1:6380[1]> hget product 855246ed-cd39-4392-9a3c-decf12c49cab
"{\"id\":\"855246ed-cd39-4392-9a3c-decf12c49cab\",\"category_id\":6,\"title\":\"product 437\",\"image_url\":
\"http://www.bestprice.gr/product437.png\",\"price\":0.0217202,\"description\":\"Description for product 437\",\"created_at\":\"2020-05-17T10:57:01Z\",\"updated_at\":\"2020-05-17T10:57:01Z\"}"

```
- Benefits: each individual product/category is fetched from cache every time, which makes our application faster.
- Drawbacks: millions of data can result to memory problems. It's ok for the current mini-bestprice-version API, but
in production, other caching methods should be followed, in case our machine is not that powerful.

#### Serialized response caching by request
- For "list" requests, we cache API request responses, based on request path as key. For example, if a user performs a GET request to
`v1/products`, this request path is stored as key in Redis, together with the string serialized response as value. This
key-value pair has a TTL of 15min.
- Example: Let's say we do a GET request to `/v1/products?limit=2`. The first time, we'll have a "miss" in cache for 
this key, so, the result will come from MySql. Right after that, a goroutine will be invoked storing the serialized
response inside our cache. If a second request to the same path happens within 15 minutes, then the answer will be 
retrieved from cache instead of database. 
`redis-cli` command and result:
```
127.0.0.1:6380[1]> GET /v1/products?limit=2
"[{\"id\":\"01b41f0d-bd5b-4c0a-8432-cd97bc57cc8b\",\"category_id\":1,\"title\":\"product 408\",\"image_url\":
\"http://www.bestprice.gr/product408.png\",\"price\":32.5052,\"description\":\"Description for product 408\",
\"created_at\":\"2020-05-17T10:57:01Z\",\"updated_at\":\"2020-05-17T10:57:01Z\"},{\"id\":
\"0912af7c-b139-42a4-8b52-bc33e4a9d124\",\"category_id\":1,\"title\":\"product 493\",\"image_url\":
\"http://www.bestprice.gr/product493.png\",\"price\":104.142,\"description\":\"Description for product 493\",
\"created_at\":\"2020-05-17T10:57:01Z\",\"updated_at\":\"2020-05-17T10:57:01Z\"}]"

```

#### Caching drawbacks
- Currently we store ALL individual products + categories. If there are millions of them, this may lead to huge memory
allocation. Possible solutions: smarter caching of specific requests needed, Redis eviction policy
- Max of 15min interval between update individual requests and list requests for the same product
- In case of too many requests, redis may overload (too many goroutines). Possible solutions: job queue, rate limit