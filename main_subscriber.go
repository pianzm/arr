package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	redisConf "github.com/pianzm/arr/config/redis"
	"github.com/pianzm/arr/helper"
	"github.com/pianzm/arr/src/member/v1/model"
	usecase "github.com/pianzm/arr/src/member/v1/usecase"
)

func subscriber(redisClient redisConf.Client, uc usecase.MemberUsecase) {
	ctx := context.Background()
	_, err := redisClient.Ping(ctx)
	if err != nil {
		time.Sleep(3 * time.Second)
		_, err := redisClient.Ping(ctx)
		if err != nil {
			log.Fatalf("cannot subscribe to topic with error %s", err.Error())
		}
	}
	log.Println("Redis ready ..")

	topic := redisClient.Subscribe(ctx, helper.DownloadChannel)

	for msg := range topic.Channel() {
		if err := process(msg.Payload, uc); err != nil {
			log.Println("err: ", err.Error())
		}
	}
}

func process(payload string, uc usecase.MemberUsecase) error {
	log.Println("going to sleep..")

	// mocking long running process
	// 30 second
	time.Sleep(30000 * time.Millisecond)

	log.Println("wake up from sleep..")
	model := model.QueueStatus{}
	if err := json.Unmarshal([]byte(payload), &model); err != nil {
		return err
	}

	currentStatus, err := uc.GetStatus(context.Background(), model.RequestID)
	if err != nil {
		return err
	}
	if currentStatus.Completed {
		return nil
	}

	// do some processing
	model.FilePath = "/tmp/somefile.txt"
	model.Completed = true
	log.Println("set the status to completed..")
	if err := uc.SetStatus(context.Background(), &model); err != nil {
		return nil
	}

	return nil
}
