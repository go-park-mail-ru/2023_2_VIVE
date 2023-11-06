package app

import (
	"HnH/configs"

	"github.com/gomodule/redigo/redis"
)

func getRedis() (redis.Conn, error) {
	conn, err := redis.DialURL(configs.HnHRedisConfig.GetConnectionURL())
	if err != nil {
		return nil, err
	}

	_, err = redis.String(conn.Do("PING"))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
