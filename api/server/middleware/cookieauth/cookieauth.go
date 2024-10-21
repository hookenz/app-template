package cookieauth

import (
	"github.com/hookenz/moneygo/api/db"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func Middleware(db db.Database) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, err := readSessionCookie(c)
			if err != nil {
				log.Debug().Msgf("readSessionCookie failed: err: %v", err)
				return c.Redirect(302, "/login")
			}

			// how do I read the DB from here?
			valid, err := db.GetSession(id)
			if err != nil {
				log.Debug().Msgf("db.GetSession failed: err: %v", err)
				return c.Redirect(302, "/login")
			}

			if !valid {
				log.Debug().Msg("session id is not valid")
				return c.Redirect(302, "/login")
			}

			err = next(c)
			if err != nil {
				return err
			}

			// OK!
			log.Debug().Msgf("Looks like you're authenticated")
			return nil
		}
	}
}

func readSessionCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("id")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
