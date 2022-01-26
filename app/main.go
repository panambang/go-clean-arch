package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	_logmovieRepo "github.com/bxcodec/go-clean-arch/logmovie/repository/mysql"
	_movieHttpDelivery "github.com/bxcodec/go-clean-arch/movie/delivery/http"
	_movieHttpDeliveryMiddleware "github.com/bxcodec/go-clean-arch/movie/delivery/http/middleware"
	_movieRepo "github.com/bxcodec/go-clean-arch/movie/repository/movie"
	_movieUcase "github.com/bxcodec/go-clean-arch/movie/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	apiKey := viper.GetString(`api_key`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _movieHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	logmovieRepo := _logmovieRepo.NewMysqlLogmovieRepository(dbConn)

	ar := _movieRepo.NewMysqlMovieRepository(apiKey)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	mu := _movieUcase.NewMovieUsecase(ar, timeoutContext)

	_movieHttpDelivery.NewMovieHandler(e, mu, logmovieRepo)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
