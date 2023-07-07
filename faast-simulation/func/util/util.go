package util

import (
	"fmt"
	"os"
)

func MustMapEnv(key string, value *string) {
	*value = os.Getenv(key)
	if *value == "" {
		msg := fmt.Sprintf("ENV %s not exsit", key)
		panic(msg)
	}
}
