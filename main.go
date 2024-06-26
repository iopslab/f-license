package main

import (
	"crypto/tls"
	"fmt"
	"github.com/crawlab-team/f-license/controllers"
	"github.com/gin-gonic/gin"

	"github.com/crawlab-team/f-license/config"
	"github.com/crawlab-team/f-license/storage"

	"log"
	"net/http"
)

func main() {
	config.Global.Load("config.json")
	storage.Connect()

	router := GenerateRouter()

	addr := fmt.Sprintf("0.0.0.0:%d", config.Global.Port)
	certFile := config.Global.ServerOptions.CertFile
	keyFile := config.Global.ServerOptions.KeyFile

	log.Println(fmt.Sprintf("server address: %s", addr))

	if config.Global.ServerOptions.EnableTLS {
		srv := &http.Server{
			Addr:         addr,
			Handler:      router,
			TLSConfig:    &config.Global.ServerOptions.TLSConfig,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
		log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
	} else {
		log.Fatal(http.ListenAndServe(addr, router))
	}
}

func GenerateRouter() *gin.Engine {
	/**
	r := mux.NewRouter()
	// Endpoints called by product owners
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(AuthenticationMiddleware)
	adminRouter.HandleFunc("/licenses", GetAllLicenses).Methods(http.MethodGet)
	adminRouter.HandleFunc("/licenses", GenerateLicense).Methods(http.MethodPost)
	adminRouter.HandleFunc("/licenses/{id}", GetLicense).Methods(http.MethodGet)
	adminRouter.HandleFunc("/licenses/{id}/activate", ChangeLicenseActiveness).Methods(http.MethodPut)
	adminRouter.HandleFunc("/licenses/{id}/inactivate", ChangeLicenseActiveness).Methods(http.MethodPut)
	adminRouter.HandleFunc("/licenses/{id}/delete", DeleteLicense).Methods(http.MethodDelete)

	// Endpoints called by product instances having license
	r.HandleFunc("/license/verify", VerifyLicense).Methods(http.MethodPost)
	r.HandleFunc("/license/ping", Ping).Methods(http.MethodPost)
	*/

	app := gin.New()
	adminGroup := app.Group("/admin")
	adminGroup.GET("/licenses", controllers.LicenseController.GetList)
	adminGroup.POST("/licenses", controllers.LicenseController.PostGenerate)

	return app
}
