package cache

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

type RedisCacher struct {
	connected bool
	Client    *redis.Client
}

func NewRedisCache(redisUrl string) (*RedisCacher, error) {
	cachier := &RedisCacher{
		connected: false,
		Client:    redis.NewClient(&redis.Options{}),
	}
	parsedUrl, err := redis.ParseURL(redisUrl)
	if err != nil {
		return cachier, errors.New(fmt.Sprintf("error parsing redis url: %s", err))
	}
	client := redis.NewClient(parsedUrl)
	_, err = client.Ping().Result()
	if err != nil {
		return cachier, errors.New(fmt.Sprintf("error creating new redis client: %s", err))
	}
	cachier.connected = true
	cachier.Client = client
	return cachier, nil
}
func (c *RedisCacher) SetProduct(id string, prodStr string) error {
	if !c.connected {
		return errors.New("redis client is currently not connected")
	}
	if err := c.Client.HSet("product", id, prodStr).Err(); err != nil {
		return errors.New(fmt.Sprintf("error in redis hset: %s", err.Error()))
	}
	return nil
}

func (c *RedisCacher) GetProduct(id string) (string, error) {
	if !c.connected {
		return "", errors.New("redis client is currently not connected")
	}
	hget := c.Client.HGet("product", id)
	if hget.Err() != nil {
		return "", errors.New(fmt.Sprintf("error getting product with id %s from redis: %s", id, hget.Err()))
	}
	return hget.Val(), nil
}

func (c *RedisCacher) DeleteProduct(id string) error {
	if !c.connected {
		return errors.New("redis client is currently not connected")
	}
	hget := c.Client.HDel("product", id)
	if hget.Err() != nil {
		return errors.New(fmt.Sprintf("error deleting product with id %s from redis: %s", id, hget.Err()))
	}
	return nil
}

func (c *RedisCacher) GetAllProducts() (map[string]string, error) {
	if !c.connected {
		return map[string]string{}, errors.New("redis client is currently not connected")
	}
	hget := c.Client.HGetAll("product")
	if hget.Err() != nil {
		return map[string]string{}, errors.New(fmt.Sprintf("error getting products from redis: %s", hget.Err()))
	}
	return hget.Val(), nil
}

func (c *RedisCacher) SetCategory(id string, catStr string) error {
	if !c.connected {
		return errors.New("redis client is currently not connected")
	}
	if err := c.Client.HSet("category", id, catStr).Err(); err != nil {
		return errors.New(fmt.Sprintf("error in redis hset: %s", err.Error()))
	}
	return nil
}

func (c *RedisCacher) GetCategory(id string) (string, error) {
	if !c.connected {
		return "", errors.New("redis client is currently not connected")
	}
	hget := c.Client.HGet("category", id)
	if hget.Err() != nil {
		return "", errors.New(fmt.Sprintf("error getting category with id %s from redis: %s", id, hget.Err()))
	}
	return hget.Val(), nil
}

func (c *RedisCacher) DeleteCategory(id string) error {
	if !c.connected {
		return errors.New("redis client is currently not connected")
	}
	hget := c.Client.HDel("category", id)
	if hget.Err() != nil {
		return errors.New(fmt.Sprintf("error deleting category with id %s from redis: %s", id, hget.Err()))
	}
	return nil
}

func (c *RedisCacher) GetAllCategories() (map[string]string, error) {
	if !c.connected {
		return map[string]string{}, errors.New("redis client is currently not connected")
	}
	hget := c.Client.HGetAll("category")
	if hget.Err() != nil {
		return map[string]string{}, errors.New(fmt.Sprintf("error getting categories from redis: %s", hget.Err()))
	}
	return hget.Val(), nil
}

func (c *RedisCacher) SetApiRequest(path string, serializedResponse string) error {
	return errors.New("method not implemented")
}
func (c *RedisCacher) GetApiRequest(path string) (string, error) {
	return "", errors.New("method not implemented")
}
