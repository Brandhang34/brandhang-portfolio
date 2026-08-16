package server

import (
	"net/http"

	"github.com/Brandhang34/brandhang-portfolio/cmd/web"
	"github.com/Brandhang34/brandhang-portfolio/internal/handler"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

	// Web Pages
	e.GET("/", func(c echo.Context) error {
		return render(c, web.Home())
	})
	e.GET("/about", func(c echo.Context) error {
		return render(c, web.About())
	})
	e.GET("/portfolio", func(c echo.Context) error {
		return render(c, web.Portfolio(handler.LoadAllPortfolioItems()))
	})
	e.GET("/contact", func(c echo.Context) error {
		return render(c, web.Contact())
	})

	// Search Projects functionality
	e.POST("/search-portfolio", func(c echo.Context) error {
		projects := handler.SearchPortfolio(
			c.FormValue("search-portfolio"),
			c.FormValue("filter-tags"),
		)
		return render(c, web.PortfolioList(projects))
	})

	// The contact endpoint relays into Discord, so it is rate limited to keep it
	// from being used as a spam pipe.
	e.POST("/submit-contact-msg", submitContactMsg,
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(1))))

	return e
}

// render writes a templ component to the response.
func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func submitContactMsg(c echo.Context) error {
	// Honeypot: a real browser leaves this hidden field empty, bots fill it in.
	// Return the success message so the bot has nothing to learn from the reply.
	if c.FormValue("website") != "" {
		return render(c, web.ContactResult(true, ""))
	}

	msg, err := handler.NewContactMsg(
		c.FormValue("email"),
		c.FormValue("subject"),
		c.FormValue("message"),
	)
	if err != nil {
		c.Logger().Warnf("contact form rejected: %v", err)
		return render(c, web.ContactResult(false, "Please check your email address and message, then try again."))
	}

	if err := handler.SendDiscordMsg(msg); err != nil {
		// Never fatal here: this handler runs on an unauthenticated public
		// route, so a Discord outage must not take the site down with it.
		c.Logger().Errorf("contact form delivery failed: %v", err)
		return render(c, web.ContactResult(false, "Sorry, something went wrong sending your message. Please email me directly."))
	}

	return render(c, web.ContactResult(true, ""))
}
