package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/routes"
	serviceRoutes "github.com/speshiy/V-K-Alcohol-Excise-Parser/service/routes"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/settings"
	"github.com/xlab/closer"
)

var srv *http.Server
var srvService *http.Server

func main() {
	closer.Bind(gracefullStop)

	numcpu := runtime.NumCPU()
	fmt.Println("CPU count:", numcpu)
	fmt.Println("Tuvis Server use GOMAXPROCS(1)")
	runtime.GOMAXPROCS(1)

	//Initialize global variables
	err := common.InitGlobalVars()
	if err != nil {
		log.Fatalln(err.Error())
	}

	//Stop program if flag noRun
	if settings.IsRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	//Initilize default routes
	router := gin.Default()
	routerService := gin.Default()

	//Setting CORS params for request Headers

	//Add cors to middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "X-Server-Date",
			"X-Token"},
		ExposeHeaders:    []string{"X-Server-Date", "X-Token"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	}))

	//Initializing app routes
	router = routes.InitRoutes(router)
	routerService = serviceRoutes.InitRoutes(routerService)

	//Starting API and Service Servers
	StartServers(router, routerService)

	//Wait to finishing all workers during 3 seconds
	closer.Hold()
}

//StartServers just start
func StartServers(router *gin.Engine, routerService *gin.Engine) {
	srv = &http.Server{
		Addr:    ":" + settings.Port,
		Handler: router,
	}

	srvService = &http.Server{
		Addr:    ":" + settings.PortService,
		Handler: routerService,
	}

	go func() {
		var err error

		if settings.IsSSL {
			err = srv.ListenAndServeTLS(settings.CertPath, settings.KeyPath)
		} else {
			err = srv.ListenAndServe()
		}

		if err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP/S server ListenAndServe: %v", err)
		}
	}()

	go func() {
		var err error

		if settings.IsSSL {
			err = srvService.ListenAndServeTLS(settings.CertPath, settings.KeyPath)
		} else {
			err = srvService.ListenAndServe()
		}

		if err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP/S Service server ListenAndServe: %v", err)
		}
	}()

}

func gracefullStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server Shutdown ", srv.Addr, "successfull")

	if err := srvService.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Println("Service server Shutdown ", srvService.Addr, "successfull")

}
