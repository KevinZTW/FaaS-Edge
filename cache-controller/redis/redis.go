package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

var ctx = context.TODO()

type Client struct {
	*goredis.Client
}

func NewClient() *Client {
	return &Client{
		Client: goredis.NewClient(&goredis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (c *Client) FormCluster() {
	val, err := c.Do(ctx, "get", "key").Result()

	// 10.244.0.109 func
	// 10.244.0.99 openfass-api

	if err != nil {
		if err == goredis.Nil {
			fmt.Println("key does not exists")
			return
		}
		panic(err)
	}
	fmt.Println(val.(string))
}

func ExampleClient() {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == goredis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
