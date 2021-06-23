package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
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
	now := time.Now()
	ctx := context.Background()
	model := model.QueueStatus{}
	if err := json.Unmarshal([]byte(payload), &model); err != nil {
		return err
	}
	log.Printf("start processing requestID %s, ", model.RequestID)

	currentStatus, err := uc.GetStatus(ctx, model.RequestID)
	if err != nil {
		return err
	}
	if currentStatus.Completed {
		return nil
	}
	members, total, err := uc.GetAll(ctx)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("app-*-%d.csv", total)

	file, err := ioutil.TempFile(os.TempDir(), fileName)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(members); err != nil {
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
	ellapse := time.Since(now)
	log.Printf("completed in %f ms", ellapse.Seconds())

	return nil
}
