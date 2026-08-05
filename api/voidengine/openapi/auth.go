package openapi

import (
	"BlackHole/api/router"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/auth"
	"BlackHole/pkg/env"

	"github.com/gin-gonic/gin"
)

func RequireAuth(authenticator auth.Authenticator) router.RouteWrapper {
	if authenticator == nil {
		return nil
	}

	return func(r router.Route) router.Route {
		return router.NewRoute(r.Method(), r.Path(), func(c *gin.Context) {
			principal, err := authenticator.Authenticate(c.Request.Context(), auth.RequestFromHTTP(c.Request))
			if err != nil {
				_ = c.Error(wrapAuthError(err))
				c.Abort()
				return
			}
			if principal == nil {
				_ = c.Error(apperror.New(errorcode.Unauthorized))
				c.Abort()
				return
			}

			ctx := c.Request.Context()
			if requestEnv, ok := env.FromContext(ctx); ok {
				ctx = env.WithContext(ctx, requestEnv.WithPrincipal(principal))
			}
			c.Request = c.Request.WithContext(ctx)
			r.Handler()(c)
		})
	}
}

func wrapAuthError(err error) error {
	if _, ok := apperror.As(err); ok {
		return err
	}
	return apperror.Wrap(errorcode.Unauthorized, err)
}
