package redis

import (
	"TestCase/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClientInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
}

type CacheClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewCacheClient(client *redis.Client) *CacheClient {
	return &CacheClient{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *CacheClient) SetPerson(person models.Person) error {
	key := fmt.Sprintf("person:%d", person.ID)
	jsonData, err := json.Marshal(person)
	if err != nil {
		return err
	}

	err = c.client.Set(c.ctx, key, jsonData, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheClient) GetPerson(personID uint) (*models.Person, error) {
	key := fmt.Sprintf("person:%d", personID)
	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var person models.Person
	err = json.Unmarshal([]byte(data), &person)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (c *CacheClient) DeletePerson(personID uint) error {
	key := fmt.Sprintf("person:%d", personID)
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
