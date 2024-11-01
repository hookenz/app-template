package server

import (
	"embed"
	"fmt"
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
	s.e.GET("/", NewHandler(pages.Index()))
	s.e.GET("/login", NewHandler(pages.Login()))

	// TODO: get custom error handler working
	// s.e.HTTPErrorHandler = customHTTPErrorHandler

	api := handler.NewHandler(s.db)
	s.e.POST("/api/auth", api.Authenticate)
	s.e.GET("/api/logout", api.Logout)

	// authenticated routes follow
	s.e.Use(middleware.RequestID())
	authenticated := s.e.Group("", cookieauth.Middleware(s.db))
	authenticated.GET("/home", NewHandler(pages.Page()))
	authenticated.GET("/users", NewHandler(pages.Users()))
	authenticated.GET("/posts", NewHandler(pages.Posts()))
}

// NewHandler that accepts a templ.Component directly
func NewHandler(component templ.Component) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Render the provided component
		return Render(c, http.StatusOK, component)
	}
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

func (s *Server) Start() error {
	return s.e.Start(s.address)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}
