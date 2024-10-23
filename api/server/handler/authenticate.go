package handler

import (
	"net/http"

	u "github.com/hookenz/app-template/api/services/user"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

//
// TODO: https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
//

type User struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	RememberMe bool   `form:"rememberMe"`
}

func (h *Handler) Authenticate(c echo.Context) error {
	var user User

	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	// TODO: don't show the password in the logs
	log.Debug().Msgf("user.name=%s, user.password=%s rememberMe=%t", user.Username, user.Password, user.RememberMe)
	u, err := u.Authenticate(h.db, user.Username, user.Password)
	if err != nil {
		// c.Redirect(302, "/login")
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	}

	log.Debug().Msgf("User authenticated %v", u.Name)

	// Generate an id
	id, err := h.db.CreateSession(u.UserID, c.RealIP())
	if err != nil {
		return err
	}

	// Create a session cookie
	writeSessionCookie(c, id)
	return c.Redirect(302, "/home")
}

func (h *Handler) Logout(c echo.Context) error {
	writeSessionCookie(c, "")
	// var user User
	// err := c.Bind(&user)
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "bad request")
	// }

	return c.Redirect(302, "/login")
}

func writeSessionCookie(c echo.Context, sessionid string) {
	log.Debug().Msg("Set Session Cookie")
	cookie := new(http.Cookie)
	cookie.Name = "id"
	cookie.Value = sessionid
	cookie.Path = "/"
	cookie.MaxAge = 24 * 60 * 60
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)
}
