package service

import (
	"context"

	log "github.com/go-kit/kit/log"
	io "github.com/jscastelblancoh/users_service/user/pkg/io"
)

// Middleware describes a service middleware.
type Middleware func(UserService) UserService

type loggingMiddleware struct {
	logger log.Logger
	next   UserService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a UserService Middleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next UserService) UserService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) GetbyId(ctx context.Context, id string) (t []io.User, err error) {
	defer func() {
		l.logger.Log("method", "GetbyId", "id", id, "t", t, "err", err)
	}()
	return l.next.GetbyId(ctx, id)
}
func (l loggingMiddleware) Post(ctx context.Context, statistic io.User) (t io.User, err error) {
	defer func() {
		l.logger.Log("method", "Post", "statistic", statistic, "t", t, "err", err)
	}()
	return l.next.Post(ctx, statistic)
}
func (l loggingMiddleware) Delete(ctx context.Context, id string) (err error) {
	defer func() {
		l.logger.Log("method", "Delete", "id", id, "err", err)
	}()
	return l.next.Delete(ctx, id)
}

func (l loggingMiddleware) Put(ctx context.Context, id string, update io.Update) (error error) {
	defer func() {
		l.logger.Log("method", "Put", "id", id, "update", update, "error", error)
	}()
	return l.next.Put(ctx, id, update)
}
