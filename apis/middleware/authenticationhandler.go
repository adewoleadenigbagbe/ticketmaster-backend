package middlewares

import (
	"net/http"

	jwtauth "github.com/Wolechacho/ticketmaster-backend/shared/helpers/utilities/auth"
	"github.com/labstack/echo/v4"
)

func AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		err = jwtauth.ValidateJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "Authentication required")
		}

		err = jwtauth.ValidateUserRoleJWT(context)
		if err != nil {
			return context.JSON(http.StatusForbidden, "You are not allowed to access this resource")
		}
		return next(context)
	}
}

func AuthorizeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		err = jwtauth.ValidateJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "Authentication required")
		}

		err = jwtauth.ValidateAdminRoleJWT(context)
		if err != nil {
			return context.JSON(http.StatusForbidden, "You are not allowed to access this resource")
		}
		return next(context)
	}
}
