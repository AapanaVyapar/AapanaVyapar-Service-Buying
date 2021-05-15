package redis

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"os"
	"strconv"
)

func InitRedis() *redis.Client {
	dbName, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	cert, err := tls.LoadX509KeyPair("./redis-sharding/redis-tls-container/certs/client.crt", "./redis-sharding/redis-tls-container/certs/client.key")
	if err != nil {
		panic(err)
	}

	caCert, err := ioutil.ReadFile("./redis-sharding/redis-tls-container/certs/ca.crt")
	if err != nil {
		panic(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)

	redisDb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbName,
		TLSConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            pool,
			InsecureSkipVerify: true,
		},
	})
	return redisDb
}
