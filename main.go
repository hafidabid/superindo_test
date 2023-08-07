package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"superindo_diksha_test/config"
	"superindo_diksha_test/database"
	"superindo_diksha_test/routers"
	"time"
)
import _ "superindo_diksha_test/docs"

func main() {
	conf := config.InitConfig()
	srv := initServer(*conf)
	srv.Run()
}

type Server struct {
	config     config.Config
	httpServer *http.Server
	db         database.AppDatabase
}

func (s *Server) shutdown() {
	sign := make(chan os.Signal)
	signal.Notify(sign, os.Interrupt)
	<-sign

	log.Infof("Shutting down server!!!")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// Flushing database if needed
	if s.config.FlushDbAfterUse {

		err := s.db.Flush()
		if err != nil {
			log.Fatalf("Error flushing db due to %v", err)
		}
	}

	select {
	case <-ctx.Done():
		log.Infof("Timeout of 3 seconds.")
	}
	log.Infof("Goodbye server, Au Revoir!!!")
}

func initServer(config config.Config) *Server {
	db, err := database.NewAppDatabase(config.Database1, config.Database2)
	if err != nil {
		panic(fmt.Sprintf("Error initializing database due to %v", err))
	}

	err = db.Migrate()
	db.Flush()
	if err != nil {
		panic(fmt.Sprintf("Error migrating due to %v", err))
	}

	return &Server{
		config: config,
		httpServer: &http.Server{
			Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		},
		db: *db,
	}
}

// @title Test Lion SuperIndo
// @version 1.0
// @description This is a API for test.
// @termsOfService http://swagger.io/terms/

// @contact.name Hafid Abi

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func (s *Server) Run() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Lion SuperIndo")
	})
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routers.AppRouter(router, s.config, s.db)

	//bind router to handler
	s.httpServer.Handler = router

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	//handle shutdown
	s.shutdown()
}
