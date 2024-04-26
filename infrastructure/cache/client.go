package cache

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"

	"gitlab.seakoi.net/engineer/backend/be-template/resource"
)

type (
	Locker interface {
		Lock(ctx context.Context) error
		Unlock(ctx context.Context) error
	}
	Client interface {
		Ping(ctx context.Context) error
		Set(ctx context.Context, key string, value interface{}, opts ...Option) error
		Get(ctx context.Context, key string, result interface{}) error
		Mutex(ctx context.Context, key string, opts ...Option) (Locker, error)
		Close() error
	}

	client struct {
		resource resource.Resource
		db       *redis.Client
		redSync  *redsync.Redsync
	}

	locker struct {
		mutex *redsync.Mutex
	}
)

func (c *client) Close() error {
	return c.db.Close()
}

func NewClient(res resource.Resource) (Client, error) {
	redisConfig := res.Configuration().Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.Index,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)

	return &client{
		resource: res,
		db:       redisClient,
		redSync:  rs,
	}, nil

}

func (l *locker) Lock(ctx context.Context) error {
	return l.mutex.LockContext(ctx)
}

func (l *locker) Unlock(ctx context.Context) error {
	_, err := l.mutex.UnlockContext(ctx)
	return err
}

func (c *client) Mutex(ctx context.Context, key string, opts ...Option) (Locker, error) {

	cacheOptions := &options{}
	for _, opt := range opts {
		opt.Apply(cacheOptions)
	}

	var lockOptions []redsync.Option
	if cacheOptions.expiry > 0 {
		lockOptions = append(lockOptions, redsync.WithExpiry(cacheOptions.expiry))
	}

	mutex := c.redSync.NewMutex(key, lockOptions...)

	return &locker{mutex}, nil
}

func (c *client) Ping(ctx context.Context) error {

	return c.db.Ping(ctx).Err()
}

func (c *client) Set(ctx context.Context, key string, value interface{}, opts ...Option) error {
	cacheOptions := &options{
		expiry: 0,
	}
	for _, opt := range opts {
		opt.Apply(cacheOptions)
	}

	return c.db.Set(ctx, key, value, cacheOptions.expiry).Err()
}

func (c *client) Get(ctx context.Context, key string, result interface{}) error {
	return c.db.Get(ctx, key).Scan(result)
}
