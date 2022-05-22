package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func NewRate(req string) func() limiter.Rate {
	return func() limiter.Rate {
		rate, err := limiter.NewRateFromFormatted(req)
		if err != nil {
			panic(err)
		}
		return rate
	}
}

func NewStore() limiter.Store {
	return memory.NewStore()
}

func NewLimiter(config *configs.Config) *limiter.Limiter {
	rate := NewRate(config.RateLimit)()
	store := NewStore()
	return limiter.New(store, rate)
}

func NewMiddleware(limiter *limiter.Limiter) gin.HandlerFunc {
	return mgin.NewMiddleware(limiter)
}
