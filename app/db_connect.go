package app

import (
	"HnH/configs"

	"database/sql"

	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/stdlib"
)

func getRedis() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:   5,
		MaxActive: 5,

		Wait: true,

		IdleTimeout:     0,
		MaxConnLifetime: 0,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(configs.HnHRedisConfig.GetConnectionURL())
			if err != nil {
				return nil, err
			}

			_, err = redis.String(conn.Do("PING"))
			if err != nil {
				conn.Close()
				return nil, err
			}

			return conn, nil
		},
	}

	return pool
}

func getPostgres() (*sql.DB, error) {
	dsn := configs.HnHPostgresConfig.GetConnectionString()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
