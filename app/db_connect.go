package app

import (
	"HnH/configs"

	"github.com/gomodule/redigo/redis"
	"database/sql"

	_ "github.com/jackc/pgx/stdlib"
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
