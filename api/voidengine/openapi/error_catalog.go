package openapi

import (
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/pkg/apperror"
	"net/http"
)

func newErrorCatalog() (*apperror.Catalog, error) {
	return apperror.NewCatalog(
		apperror.Definition{Code: apperror.Success, HTTPStatus: http.StatusOK},
		apperror.Definition{Code: errorcode.APINotFound, HTTPStatus: http.StatusNotFound},
		apperror.Definition{Code: errorcode.InvalidParams, HTTPStatus: http.StatusBadRequest},
		apperror.Definition{Code: errorcode.SystemError, HTTPStatus: http.StatusInternalServerError},
		apperror.Definition{Code: errorcode.InvalidUserName, HTTPStatus: http.StatusBadRequest},
		apperror.Definition{Code: errorcode.UserNotFound, HTTPStatus: http.StatusNotFound},
	)
}
