package client

import (
	"context"
	"database/sql"

	"github.com/chaihaobo/gocommon/mysql"
	"gorm.io/gorm"

	"github.com/chaihaobo/be-template/constant"
	"github.com/chaihaobo/be-template/resource"
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
	gormDB, err := mysql.GormDB(mysql.Config{
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
	}, nil)
	if err != nil {
		return nil, err
	}
	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	return &client{db, gormDB}, nil
}
