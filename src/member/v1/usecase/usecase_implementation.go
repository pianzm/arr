package usecase

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/bgadrian/fastfaker/faker"
	redisConf "github.com/pianzm/arr/config/redis"
	"github.com/pianzm/arr/helper"
	"github.com/pianzm/arr/src/member/v1/model"
	"github.com/pianzm/arr/src/member/v1/query"
	"github.com/pianzm/arr/src/member/v1/repo"
)

type UsecaseImpl struct {
	Publisher redisConf.Client
	WriteRepo repo.MemberRepository
	QueryRead query.MemberRead
}

func NewUsecase(redisClient redisConf.Client, writeDB repo.MemberRepository, queryRead query.MemberRead) MemberUsecase {
	return &UsecaseImpl{
		Publisher: redisClient,
		WriteRepo: writeDB,
		QueryRead: queryRead,
	}
}
func (a *UsecaseImpl) Publish(ctx context.Context, params *model.QueueStatus) error {
	if err := a.Publisher.Publish(ctx, helper.DownloadChannel, params); err != nil {
		return err
	}
	if _, err := a.Publisher.Set(ctx, params.RequestID, params, 0); err != nil {
		return err
	}

	return nil
}

func (a *UsecaseImpl) GetStatus(ctx context.Context, requestID string) (*model.QueueStatus, error) {
	res, err := a.Publisher.Get(ctx, requestID)
	if err != nil {
		return nil, err
	}
	model := model.QueueStatus{}
	if err := json.Unmarshal([]byte(res), &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func (a *UsecaseImpl) SetStatus(ctx context.Context, params *model.QueueStatus) error {
	params.Completed = true
	if _, err := a.Publisher.Set(ctx, params.RequestID, params, 0); err != nil {
		return err
	}
	return nil
}

func (a *UsecaseImpl) InitData(ctx context.Context) error {
	len := 2000000
	members := []*model.Member{}
	currentTime := time.Now()

	var wg sync.WaitGroup
	chMembers := make(chan *model.Member, 200)
	for i := 0; i < len; i++ {
		wg.Add(1)
		go a.faker(&wg, chMembers)
	}

	go func() {
		wg.Wait()
		close(chMembers)
	}()

	for m := range chMembers {
		members = append(members, m)
	}
	ellapseTime := time.Since(currentTime)
	log.Printf("fake data generated in %f", ellapseTime.Seconds())
	insertTime := time.Now()

	if err := a.WriteRepo.Insert(ctx, members); err != nil {
		return err
	}
	insertEnd := time.Since(insertTime)
	log.Printf("insert data into DB takes in %f", insertEnd.Seconds())
	return nil
}

func (a *UsecaseImpl) faker(wg *sync.WaitGroup, response chan<- *model.Member) {
	defer wg.Done()
	generator := faker.NewSafeFaker()

	member := model.Member{
		FirstName: generator.FirstName(),
		LastName:  generator.LastName(),
		Email:     generator.Email(),
	}

	response <- &member
}

func (a *UsecaseImpl) GetAll(ctx context.Context) ([][]string, int64, error) {
	members, err := a.QueryRead.GetAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	total, err := a.QueryRead.GetTotal(ctx)
	if err != nil {
		return nil, 0, err
	}
	return members, total, nil
}
