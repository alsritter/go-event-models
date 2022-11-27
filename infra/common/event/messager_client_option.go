package event

import (
	"context"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type mqOptions struct {
	addr             string
	adminAddr        string
	groupName        string
	retry            int
	consumerInstance string
	ctx              context.Context
}

func defaultOptions() mqOptions {
	opts := mqOptions{
		addr:             "0.0.0.0:9876",
		adminAddr:        "0.0.0.0:10911",
		groupName:        "default",
		retry:            1,
		consumerInstance: UUIDShort(),
		ctx:              context.Background(),
	}
	return opts
}

// Option set option by close func
type Option func(*mqOptions)

func WithAdminAddr(addr string) Option {
	return func(options *mqOptions) {
		if addr == "" {
			return
		}
		options.adminAddr = addr
	}
}

func WithAddr(addr string) Option {
	return func(options *mqOptions) {
		if addr == "" {
			return
		}
		options.addr = addr
	}
}

func WithGroupName(g string) Option {
	return func(options *mqOptions) {
		if g == "" {
			return
		}
		options.groupName = g
	}
}

func WithRetry(t int) Option {
	return func(options *mqOptions) {
		if t < 0 {
			return
		}
		options.retry = t
	}
}

func WithContext(ctx context.Context) Option {
	return func(opts *mqOptions) {
		if ctx == nil {
			return
		}
		opts.ctx = ctx
	}
}

func WithInstanceName(s string) Option {
	return func(opts *mqOptions) {
		if s == "" {
			return
		}
		opts.consumerInstance = s
	}
}

func UUIDShort() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
