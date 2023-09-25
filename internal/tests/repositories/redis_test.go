package repositories

import (
	"TestCase/internal/models"
	redis2 "TestCase/internal/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisClientImpl(t *testing.T) {
	redisAddr := "localhost:6379"

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
	}
	cacheClient := redis2.NewCacheClient(client)

	t.Run("SetPerson", func(t *testing.T) {
		person := models.Person{
			Name:    "John",
			Surname: "Doe",
			Age:     30,
			Gender:  "Male",
		}

		err := cacheClient.SetPerson(person)
		assert.NoError(t, err)

		key := fmt.Sprintf("person:%d", person.ID)
		data, err := client.Get(context.Background(), key).Result()
		assert.NoError(t, err)

		var retrievedPerson models.Person
		err = json.Unmarshal([]byte(data), &retrievedPerson)
		assert.NoError(t, err)

		assert.Equal(t, person, retrievedPerson)

		err = cacheClient.DeletePerson(person.ID)
		assert.NoError(t, err)
	})

	t.Run("GetPerson", func(t *testing.T) {
		person := models.Person{
			Name:    "Alice",
			Surname: "Smith",
			Age:     25,
			Gender:  "Female",
		}

		key := fmt.Sprintf("person:%d", person.ID)
		jsonData, _ := json.Marshal(person)
		err := client.Set(context.Background(), key, jsonData, 24*time.Hour).Err()
		assert.NoError(t, err)

		retrievedPerson, err := cacheClient.GetPerson(person.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedPerson)
		assert.Equal(t, person, *retrievedPerson)
	})

	t.Run("DeletePerson", func(t *testing.T) {
		person := models.Person{
			Name:    "Jane",
			Surname: "Doe",
			Age:     35,
			Gender:  "Female",
		}
		key := fmt.Sprintf("person:%d", person.ID)
		jsonData, _ := json.Marshal(person)
		err := client.Set(context.Background(), key, jsonData, 24*time.Hour).Err()
		assert.NoError(t, err)
		err = cacheClient.DeletePerson(person.ID)
		assert.NoError(t, err)
		_, err = cacheClient.GetPerson(person.ID)
		assert.Error(t, err)
	})
}
