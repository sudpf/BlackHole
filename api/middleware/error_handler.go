package middleware

import (
	"BlackHole/api/common/response"
	"BlackHole/pkg/apperror"
	"BlackHole/pkg/env"
	"BlackHole/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(provider *env.Provider, catalog *apperror.Catalog, systemCode apperror.Code) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		writeError(c, provider, catalog, systemCode, c.Errors.Last().Err)
	}
}

func writeError(
	c *gin.Context,
	provider *env.Provider,
	catalog *apperror.Catalog,
	systemCode apperror.Code,
	err error,
) {
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
				WithField("code", systemCode).
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

	requestEnv := provider.NewEnvFromContext(c.Request.Context())
	message, localizeErr := requestEnv.Localize(apperror.MessageID(appErr.Code()), appErr.Params())
	if localizeErr != nil {
		logger.FromContext(c.Request.Context()).
			WithError(localizeErr).
			WithField("code", appErr.Code()).
			Error("localize API error")

		appErr = apperror.Wrap(systemCode, err)
		definition, _ = catalog.Lookup(systemCode)
		message, localizeErr = requestEnv.Localize(apperror.MessageID(systemCode), nil)
		if localizeErr != nil {
			logger.FromContext(c.Request.Context()).WithError(localizeErr).Error("localize system error")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	details := appErr.Details()
	if details == nil {
		details = requestEnv.TranslateErrors(appErr.Unwrap())
	}

	response.WriteError(c, definition.HTTPStatus, appErr.Code(), message, details)
}
