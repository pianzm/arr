package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	defaultPort = 8081
)

// os.Getenv("REDIS_HOST"), os.Getenv("REDIS_TLS"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_PORT"), redisDB
//Config base struct
type Config struct {
	DefaultPort            int
	RedisHost              string
	RedisTLS               string
	RedisPassword          string
	RedisPort              string
	RedisDB                string
	PostgreWriteDBHost     string
	PostgreWriteDBUser     string
	PostgreWriteDBPassword string
	PostresWriteDBName     string
}

func NewConfig() (*Config, error) {
	var (
		err     error
		appPort int
	)
	defer func(err error) {
		if err != nil {
			log.Panicf("%s", err.Error())
		}
	}(err)

	redisHost, err := LookUpConf("REDIS_HOST")
	if err != nil {
		return nil, err
	}
	redisTls, err := LookUpConf("REDIS_TLS")
	if err != nil {
		return nil, err
	}
	redisPassword, err := LookUpConf("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}
	redisPort, err := LookUpConf("REDIS_PORT")
	if err != nil {
		return nil, err
	}
	redisDB, err := LookUpConf("REDIS_DB")
	if err != nil {
		return nil, err
	}
	pgdbHost, err := LookUpConf("WRITE_DB_HOST")
	if err != nil {
		return nil, err
	}
	pgdbUser, err := LookUpConf("WRITE_DB_USER")
	if err != nil {
		return nil, err
	}
	pgdbPassword, err := LookUpConf("WRITE_DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	pgdbName, err := LookUpConf("WRITE_DB_NAME")
	if err != nil {
		return nil, err
	}

	// ommit error for brevity, will mapping through docker
	portEnv, _ := LookUpConf("PORT")
	portInt, err := strconv.Atoi(portEnv)
	if err != nil {
		appPort = defaultPort
	} else {
		appPort = int(portInt)
	}

	return &Config{
		DefaultPort:            appPort,
		RedisHost:              redisHost,
		RedisPassword:          redisPassword,
		RedisTLS:               redisTls,
		RedisPort:              redisPort,
		RedisDB:                redisDB,
		PostgreWriteDBHost:     pgdbHost,
		PostgreWriteDBUser:     pgdbUser,
		PostgreWriteDBPassword: pgdbPassword,
		PostresWriteDBName:     pgdbName,
	}, nil
}

func LookUpConf(input string) (string, error) {
	key, ok := os.LookupEnv(input)
	if !ok {
		return "", fmt.Errorf("need to specify %s in your env file", input)
	}
	return key, nil
}
