package main

import (
    "database/sql"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/swaggo/http-swagger"


    "data-processor-project/internal/api"
    "data-processor-project/internal/repository"
    "data-processor-project/internal/domain/services"
    "data-processor-project/internal/redis"
    "data-processor-project/internal/logs"
    "data-processor-project/config"
   
    _ "data-processor-project/docs"
    _ "github.com/mattn/go-sqlite3"
)

// @title Data Processor API
// @version 1.0
// @description API For data processing
// @termsOfService http://swagger.io/terms/
// @contact.name Support team
// @contact.email support@data-processor.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {

    r := mux.NewRouter()

    cfg := config.LoadConfig() 

    logs.InitLogger()
    defer logs.Sync() 

    db, err := sql.Open("sqlite3", "data_processor.db")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    repository.Migrate(db)


	kafkaService, err := services.NewKafkaService(cfg.KafkaHost)
	if err != nil {
		log.Fatalf("Failed to create Kafka service: %v", err)
	}

    rdb, err := redis.InitRedis(cfg.RedisHost,cfg.RedisPort)
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer rdb.Close() 

    userService := services.NewUserService(db)
	requestService := services.NewRequestService(db, kafkaService, rdb)

    userAPI := api.NewUserAPI(userService)
    requestAPI := api.NewRequestAPI(requestService)


    r.HandleFunc("/request", requestAPI.AddRequest).Methods("POST")
    r.HandleFunc("/user", userAPI.CreateUser).Methods("POST")


    r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    log.Println("Server started at :8080")
    http.ListenAndServe(":8080", r)
}