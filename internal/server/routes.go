package server

import (
	"net/http"

	"brandhang-portfolio/cmd/web"
	"brandhang-portfolio/cmd/web/components/portfolio"
	"brandhang-portfolio/cmd/web/pages"
	"brandhang-portfolio/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	// Web Pages
	e.GET("/", func(c echo.Context) (err error) {
		return pages.Home().Render(c.Request().Context(), c.Response())
	})
	e.GET("/about", func(c echo.Context) (err error) {
		return pages.About().Render(c.Request().Context(), c.Response())
	})
	e.GET("/portfolio", func(c echo.Context) (err error) {
		return pages.Portfolio(handler.LoadAllPortfolioItems()).Render(c.Request().Context(), c.Response())
	})
	e.GET("/contact", func(c echo.Context) (err error) {
		return pages.Contact().Render(c.Request().Context(), c.Response())
	})

	// Search Projects functionality
	e.POST("/search-portfolio", func(c echo.Context) (err error) {
		searchQuery := c.FormValue("search-portfolio")
		tagsQuery := c.FormValue("filter-tags")
		projects := handler.SearchProjectsHandler(searchQuery, tagsQuery)
		return portfolio.PortfolioList(projects).Render(c.Request().Context(), c.Response())
	})

	e.POST("/submit-contact-msg", func(c echo.Context) (err error) {
		email := c.FormValue("email")
		subject := c.FormValue("subject")
		message := c.FormValue("message")
		handler.SendDiscordMsg("Email: " + email + "\nSubject: " + subject + "\nMessage:\n" + message)
		return
	})

	return e
}
