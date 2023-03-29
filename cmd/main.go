package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/postikus/go-starter/internal/balance"
	"github.com/postikus/go-starter/internal/user"
	userHandler "github.com/postikus/go-starter/internal/user/handler"
	"github.com/postikus/go-starter/pkg/database"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, cancelFn := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelFn()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT)

		select {
		case s := <-sigint:
			logger.Warn("os signal received", zap.String("signal", s.String()))
			cancelFn()
		case <-ctx.Done():
		}
	}()

	viper.AddConfigPath("configs")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err = viper.ReadInConfig(); err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@%s/%s",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.server"),
		viper.GetString("database.name"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	dbSqlx := sqlx.NewDb(db, "mysql")

	commiter := database.NewCommitter()

	balanceRepository := balance.NewRepository(logger, dbSqlx)

	userRepository := user.NewRepository(logger, dbSqlx, dbSqlx, commiter)
	userService := user.NewService(logger, userRepository, balanceRepository)

	userPost := userHandler.NewPost(logger, userService)

	r := gin.New()

	r.Handle(http.MethodPost, "/users", gin.WrapH(userPost))

	r.Handle(http.MethodGet, "/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	server := http.Server{
		Addr:    ":" + viper.GetString("http.port"),
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	ctx, cancelFn = context.WithTimeout(context.Background(), viper.GetDuration("http.shutdown_timeout"))
	defer cancelFn()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
