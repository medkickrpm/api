package middleware

import (
	"MedKick-backend/pkg/database/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := Store.Get(c.Request(), "medkick-session")
		if err != nil {
			http.Error(c.Response(), "Failed to get session", http.StatusInternalServerError)
			return err
		}

		userId, ok := session.Values["user-id"].(uint)

		if !ok {
			c.Set("x-guest", true)
			return next(c)
		}

		u := models.User{ID: &userId}

		err = u.GetUser()
		if err != nil {
			c.Set("x-guest", true)
			return next(c)
		}

		if err == nil {
			c.Set("x-guest", false)
			c.Set("x-id", userId)
			c.Set("x-user", u)
			c.Set("x-auth-type", "cookie")
			return next(c)
		}

		// If we get here, they had a cookie with an invalid user
		// so delete it.
		delete(session.Values, "user-id")
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			http.Error(c.Response(), "Failed to save session", http.StatusInternalServerError)
			return err
		}
		c.Set("x-guest", true)
		return next(c)
	}
}

func NotGuest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		guest := c.Get("x-guest").(bool)
		if guest {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "You must be logged in to perform this action."})
		}
		return next(c)
	}
}

func HasRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("x-user").(models.User)
			for _, role := range roles {
				if user.Role == role {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, map[string]string{"message": "You do not have permission to perform this action."})
		}
	}
}

func IsGuest(c echo.Context) bool {
	return c.Get("x-guest").(bool)
}

func GetSelf(c echo.Context) models.User {
	return c.Get("x-user").(models.User)
}
