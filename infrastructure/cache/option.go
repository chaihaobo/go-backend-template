package cache

import "time"

type (
	options struct {
		expiry time.Duration
	}

	Option interface {
		Apply(option *options)
	}

	OptionFunc func(option *options)
)

func (l OptionFunc) Apply(option *options) {
	l(option)
}

func WithExpiry(expiry time.Duration) Option {
	return OptionFunc(func(option *options) {
		option.expiry = expiry
	})
}
