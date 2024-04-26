package client

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/chaihaobo/gocommon/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gitlab.seakoi.net/engineer/backend/be-template/constant"
	"gitlab.seakoi.net/engineer/backend/be-template/resource"
)

type (
	Client interface {
		Ping(ctx context.Context) error
		//	Begin starts a transaction
		Begin(ctx context.Context) (context.Context, error)
		// Rollback rolls back the transaction
		// and returns the context without the transaction
		Rollback(ctx context.Context) (context.Context, error)
		// Commit commits the transaction
		// and returns the context without the transaction
		Commit(ctx context.Context) (context.Context, error)
		//	DB returns the gorm.DB instance
		//  if context has transaction, it will return the transaction instance
		//  otherwise, it will return the gorm.DB instance
		DB(ctx context.Context) *gorm.DB
		//	Close the raw sql db connections pool
		Close() error
	}
	client struct {
		rawDB  *sql.DB
		gormDB *gorm.DB
	}
)

func (c *client) Close() error {
	return c.rawDB.Close()
}

func (c *client) Ping(ctx context.Context) (err error) {
	if pinger, ok := c.DB(ctx).ConnPool.(interface{ Ping() error }); ok {
		err = pinger.Ping()
	}
	return
}

func (c *client) Begin(ctx context.Context) (context.Context, error) {
	beginedDB := c.gormDB.Begin().WithContext(ctx)
	if err := beginedDB.Error; err != nil {
		return nil, err
	}
	return context.WithValue(ctx, constant.ContextKeyTrx, beginedDB), nil
}

func (c *client) Rollback(ctx context.Context) (context.Context, error) {
	if db, ok := ctx.Value(constant.ContextKeyTrx).(*gorm.DB); ok && db != nil {
		if err := db.Rollback().Error; err != nil {
			return ctx, err
		}
		return context.WithValue(ctx, constant.ContextKeyTrx, nil), nil

	}
	return ctx, nil
}

func (c *client) Commit(ctx context.Context) (context.Context, error) {
	if db, ok := ctx.Value(constant.ContextKeyTrx).(*gorm.DB); ok && db != nil {
		if err := db.Commit().Error; err != nil {
			return ctx, err
		}
		return context.WithValue(ctx, constant.ContextKeyTrx, nil), nil

	}
	return ctx, nil
}

func (c *client) DB(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(constant.ContextKeyTrx).(*gorm.DB); ok && db != nil {
		return db

	}
	return c.gormDB.WithContext(ctx)
}

func New(res resource.Resource) (Client, error) {
	dbConfig := res.Configuration().Database
	db, err := mysql.DB(mysql.Config{
		Host:        dbConfig.Host,
		Port:        dbConfig.Port,
		User:        dbConfig.User,
		Password:    dbConfig.Password,
		Name:        dbConfig.Name,
		MaxOpen:     dbConfig.MaxOpen,
		MaxIdle:     dbConfig.MaxIdle,
		MaxLifetime: int(dbConfig.MaxLifetime.Minutes()),
		MaxIdleTime: int(dbConfig.MaxIdleTime.Minutes()),
		Location:    dbConfig.Location,
		ParseTime:   true,
	})
	if err != nil {
		return nil, err
	}
	gormDB, err := gorm.Open(gormMysql.New(gormMysql.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		return nil, err
	}
	return &client{db, gormDB}, nil
}
