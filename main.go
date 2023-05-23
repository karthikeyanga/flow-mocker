package main

import (
	"context"
	"fmt"
	"io"
	"mocker/common"
	"mocker/config"
	"mocker/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var configFileName string

	if len(os.Args) < 2 {
		fmt.Println("Needs config file path")
		os.Exit(EXITING_AS_ERROR_IN_ARGS)
	}
	configFileName = os.Args[1]
	configFileReader, err := os.Open(configFileName)
	if err != nil {
		fmt.Println("Error opening the config file", configFileName, err)
		os.Exit(EXITING_AS_ERROR_PROCESSING_CONFIG_FILE)
	}
	//Initialise Config
	appConfig, err := config.New(configFileReader)
	if err != nil {
		fmt.Println("Error reading config file", err)
		os.Exit(EXITING_AS_ERROR_PROCESSING_CONFIG_FILE)
	}
	defer appConfig.Close()
	appContext := appConfig.NewContext("main")
	defer appContext.Close()
	//Close the config file
	if err := configFileReader.Close(); err != nil {
		appContext.Log.Errorln("error while closing yml config file")
	}
	//Last step is to start the server
	ginEngine, err := initServer(appContext)
	if err != nil {
		os.Exit(EXITING_AS_ERROR_WHILE_SERVER_START)
	}
	if err := InitModules(appContext); err != nil {
		appContext.Log.Errorln("error while initialising modules", err)
		os.Exit(EXITING_AS_ERROR_INIT_MODULES)
	}

	startServerWithSignalledShutdown(ginEngine, appContext)
}

func initServer(ac *config.AppContext) (*gin.Engine, error) {
	//serverConfig:=ac.Config.ServerConfig
	//gin.DefaultWriter = io.MultiWriter()
	errorLogFile, _ := ac.GetFile("server.log.error")
	accessLogFile, _ := ac.GetFile("server.log.access")
	apimode := ac.Config.APIMode
	gin.DefaultErrorWriter = io.MultiWriter(errorLogFile)
	ac.NewContext("startup.server").Log.Infoln("Running in Mode " + util.APIMODE_GIN_MODE[apimode])
	gin.SetMode(util.APIMODE_GIN_MODE[apimode])
	ginEngine := gin.Default()
	ginEngine.Use(config.AppContextGinMiddleware(ac.AppConfig))
	ginEngine.Use(util.AccessLogGinMiddleware(accessLogFile), util.GinBodyLogMiddleware())
	//staticPath:=serverConfig.StaticPath
	//ginEngine.LoadHTMLGlob(staticPath + "/webroot/**/*[.html|tmpl]")
	//ginEngine.Use(util.SentryMiddleware(raven.DefaultClient, false))
	RouterPatternsInit(ac, ginEngine)
	return ginEngine, nil
}

func startServerWithSignalledShutdown(engine *gin.Engine, ac *config.AppContext) {
	serverConfig := ac.Config.ServerConfig
	srv := &http.Server{
		Addr:    serverConfig.Host + ":" + serverConfig.Port,
		Handler: engine,
	}

	ac.Go(func(tac common.AppContexter) {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	})
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		ac.Log.Fatal("Server Shutdown: ", err)
	}

	ac.Log.Println("Server exiting")
}

func InitModules(ac *config.AppContext) error {

	return nil
}
