package task

import (
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/letterScape/backend/services"
	"log"
	"time"
)

type ChainTxTask struct{}

func (task *ChainTxTask) PollTx(c *gin.Context) {
	scheduler := gocron.NewScheduler(time.UTC)
	service := &services.WNFTInfoService{}
	job1, err := scheduler.Every(30).Seconds().Do(service.PollTx, c)
	if err != nil {
		log.Printf("PollTx job:%v err:%v", job1, err)
	}

	job2, err := scheduler.Every(30).Seconds().Do(service.SyncWnftStatus, c)
	if err != nil {
		log.Printf("SyncReadonlyStatus job:%v err:%v", job2, err)
	}
	scheduler.StartAsync()

	time.Sleep(10 * time.Second)
}
