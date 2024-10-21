package server

import (
	"embed"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"github.com/hookenz/app-template/api/db"
	"github.com/hookenz/app-template/api/server/handler"
	"github.com/hookenz/app-template/api/server/middleware/cookieauth"
	"github.com/hookenz/app-template/api/server/middleware/logging"
	"github.com/hookenz/app-template/web/pages"
)

type Server struct {
	e        *echo.Echo
	address  string
	staticfs embed.FS
	db       db.Database
}

func New(address string, db db.Database, staticfs embed.FS) *Server {
	s := &Server{
		e:        echo.New(),
		address:  address,
		staticfs: staticfs,
		db:       db,
	}

	s.e.HideBanner = true

	s.setupMiddleware()
	s.setupHandlers()
	s.setupStaticHandler()
	return s
}

func (s *Server) setupMiddleware() {
	logging.NewLogger()
	s.e.Use(logging.Middleware)
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
}

func (s *Server) setupHandlers() {
	s.e.GET("/", IndexHandler)
	s.e.GET("/login", LoginHandler)

	api := handler.NewHandler(s.db)
	s.e.POST("/api/auth", api.Authenticate)
	s.e.POST("/api/logout", api.Logout)

	// authenticated routes follow
	authenticated := s.e.Group("", cookieauth.Middleware(s.db))
	authenticated.GET("/home", HomeHandler)
}

func IndexHandler(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Index())
}

func LoginHandler(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Login())
}

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Home())
}

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func (s *Server) setupStaticHandler() {
	// Serve the frontend at "/"
	fs := echo.MustSubFS(s.staticfs, "")
	s.e.StaticFS("/", fs)
}

func (s *Server) Start() {
	s.e.Start(s.address)
}
