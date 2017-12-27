package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZachEddy/gokit-todolist/pkg/todolist"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	db := initDatabase("localhost", "zacharyeddy", "todolist")
	db.AutoMigrate(todolist.Task{})
	defer db.Close()
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s = todolist.TodoListService{DB: db}

	var h http.Handler
	{
		h = todolist.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}

func initDatabase(host, user, dbname string) *gorm.DB {
	params := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable", host, user, dbname)
	db, err := gorm.Open("postgres", params)
	if err != nil {
		msg := fmt.Sprintf("Unable to initialize database: %v", err.Error())
		panic(msg)
	}
	return db
}
