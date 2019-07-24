package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/common"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/routes"
	migrateControllers "github.com/speshiy/V-K-Alcohol-Excise-Parser/service/controllers"
	serviceRoutes "github.com/speshiy/V-K-Alcohol-Excise-Parser/service/routes"
	"github.com/speshiy/V-K-Alcohol-Excise-Parser/settings"
	"github.com/xlab/closer"
	"golang.org/x/crypto/acme/autocert"
)

var srvHTTP *http.Server
var srvHTTPS *http.Server
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

	//Запуск автомиграции
	migrateControllers.AutoMigrate()

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
	var hostPolicy autocert.HostPolicy

	hostPolicy = autocert.HostWhitelist("vkaep.tuvis.world")

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache("cert-cache"),
	}

	srvHTTP = &http.Server{
		Addr:         ":" + settings.PortHTTP,
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	srvHTTPS = &http.Server{
		Addr:         ":" + settings.PortHTTPS,
		Handler:      router,
		TLSConfig:    &tls.Config{GetCertificate: certManager.GetCertificate},
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	srvService = &http.Server{
		Addr:         ":" + settings.PortService,
		Handler:      routerService,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//Запускаем HTTP порт
	go func() {
		var err error
		//Если IsSSL объявлен то, связываем его с http сервером, чтобы получить сертификаты
		if settings.IsSSL {
			srvHTTP.Handler = certManager.HTTPHandler(srvHTTP.Handler)
		}
		err = srvHTTP.ListenAndServe()
		if err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()

	//Запускаем HTTPS порт если есть флаг
	if settings.IsSSL {
		go func() {
			var err error
			//Запускаем сервер вместе с TLS шифрованием
			err = srvHTTPS.ListenAndServeTLS("", "")
			if err != http.ErrServerClosed {
				// Error starting or closing listener:
				log.Printf("HTTPS server ListenAndServe: %v", err)
			}
		}()
	}

	//Запускаем Service порт
	go func() {
		var err error
		err = srvService.ListenAndServe()
		if err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTPService server ListenAndServe: %v", err)
		}
	}()

}

func gracefullStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srvHTTP.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server Shutdown ", srvHTTP.Addr, "successfull")

	if err := srvHTTPS.Shutdown(ctx); err != nil {
		log.Printf("HTTPS server Shutdown: %v", err)
	}
	log.Println("HTTPS server Shutdown ", srvHTTPS.Addr, "successfull")

	if err := srvService.Shutdown(ctx); err != nil {
		log.Printf("HTTPService server Shutdown: %v", err)
	}
	log.Println("HTTPService server Shutdown ", srvService.Addr, "successfull")
}
