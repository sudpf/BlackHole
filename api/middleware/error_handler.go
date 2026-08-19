package middleware

import (
	"BlackHole/api/common/response"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/env"
	"BlackHole/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationTranslator interface {
	TranslateErrors(requestEnv *env.Env, err error) map[string]string
}

func ErrorHandler(registry apperror.ErrorRegistry, translator ValidationTranslator) gin.HandlerFunc {
	catalog := registry.Catalog()
	systemCode := registry.SystemErrorCode()
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		appErr, known := apperror.As(err)
		if !known {
			appErr = apperror.Wrap(systemCode, err)
		}

		definition, exists := catalog.Lookup(appErr.Code())
		if !exists {
			appErr = apperror.Wrap(systemCode, err)
			definition, exists = catalog.Lookup(systemCode)
			if !exists {
				logger.FromContext(c.Request.Context()).
					WithField("code", int(systemCode)).
					Error("system error code is not registered")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		if definition.HTTPStatus >= 500 {
			logger.FromContext(c.Request.Context()).
				WithError(err).
				WithField("code", appErr.Code()).
				Error("handle API request")
		}

		if c.Writer.Written() {
			return
		}

		requestEnv, ok := env.FromContext(c.Request.Context())
		if !ok {
			requestEnv = env.NewFromContext(c.Request.Context())
		}
		message, localizeErr := catalog.Localize(requestEnv, apperror.MessageID(appErr.Code()), appErr.Params())
		if localizeErr != nil {
			logger.FromContext(c.Request.Context()).
				WithError(localizeErr).
				WithField("code", appErr.Code()).
				Error("localize API error")

			appErr = apperror.Wrap(systemCode, err)
			definition, _ = catalog.Lookup(systemCode)
			message, localizeErr = catalog.Localize(requestEnv, apperror.MessageID(systemCode), nil)
			if localizeErr != nil {
				logger.FromContext(c.Request.Context()).WithError(localizeErr).Error("localize system error")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		details := appErr.Details()
		if details == nil && translator != nil {
			details = translator.TranslateErrors(requestEnv, appErr.Unwrap())
		}

		response.WriteError(c, definition.HTTPStatus, appErr.Code(), message, details)
	}
}
