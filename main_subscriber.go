package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
			log.Println("err processing: ", err.Error())
		}
	}
}

func process(payload string, uc usecase.MemberUsecase) error {
	log.Println("going to sleep..")

	// mocking long running process
	// 10 second
	time.Sleep(10000 * time.Millisecond)

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
	file, err := ioutil.TempFile(os.TempDir(), "app-*.csv")
	if err != nil {
		return err
	}
	text := []byte(`some long running processing result`)
	if _, err = file.Write(text); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	model.FilePath = file.Name()
	model.Completed = true
	log.Printf("set the status to completed with path %s\n", file.Name())
	if err := uc.SetStatus(context.Background(), &model); err != nil {
		return err
	}

	return nil
}
