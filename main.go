package main

import (
	"log"
	"os"
	"sync"

	envLoader "github.com/joho/godotenv"
	"github.com/pianzm/arr/config"
	"github.com/pianzm/arr/config/redis"
	usecase "github.com/pianzm/arr/src/member/v1/usecase"
)

type HTTPServer struct {
	uc usecase.MemberUsecase
}

func main() {
	envFile := ".env-local"
	if _, err := os.Stat("/.dockerenv"); err == nil {
		envFile = ".env"
	}

	if err := envLoader.Load(envFile); err != nil {
		log.Panicf("error loading env file : %s", err.Error())
	}
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panicf("error bootstrap config: %s", err.Error())
	}
	redisConnection, err := redis.ConnectRedis(cfg)
	if err != nil {
		log.Panicf("cannot connect to redis %s", err.Error())
	}

	memberUc := usecase.NewUsecase(redisConnection)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		initHTTP(cfg, memberUc)
	}()

	go func() {
		defer wg.Done()
		subscriber(redisConnection, memberUc)
	}()

	wg.Wait()

}
