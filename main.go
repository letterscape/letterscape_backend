package main

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/letterScape/backend/conf"
	"github.com/letterScape/backend/global"
	"github.com/letterScape/backend/router"
	"github.com/letterScape/backend/task"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
	defer lib.Destroy()
	log.Println("BlockChain RpcUrl: ", global.BlockChainConfig.RpcUrl)

	(&task.ChainTxTask{}).PollTx(&gin.Context{})
	router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}

// init()是go内置方法，默认会在调用类加载之前自动执行
func init() {
	conf.SetupConfig()
}
