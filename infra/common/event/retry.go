package event

import (
	"context"
	"fmt"
)

type RetryError interface {
	Error() string
	GetCode() int
}

const (
	DbError        = 29001001
	RedisError     = 29001003
	ConnServiceErr = 29001007
)

// RecoverToRetry Recover panic to retry 只将需要重试的错误类型返回 error,其余都内部消化不返回error
/*
need retry type
1. RedisError     = 29001003
2. DbError        = 29001001
3. ConnServiceErr = 29001007
*/
func RecoverToRetry(ctx context.Context, msg *Message, handler func(ctx context.Context, msg *Message) error) (retry bool, retErr error) {
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			if err, ok := panicErr.(RetryError); ok {
				errCode := err.GetCode()
				retErr = err
				if errCode == DbError || errCode == RedisError || errCode == ConnServiceErr {
					retry = true
				}
			} else {
				retErr = fmt.Errorf("%v", panicErr)
			}
		}
	}()
	retErr = handler(ctx, msg)
	return false, retErr
}
