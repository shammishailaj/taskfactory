package routes

import (
	"github.com/go-chi/chi"
	"github.com/shammishailaj/taskfactory/internal/app/taskfactory/server/handlers"
	"github.com/shammishailaj/taskfactory/internal/app/taskfactory/server/middleware"
	"github.com/shammishailaj/taskfactory/pkg/utils"

	"github.com/go-chi/httprate"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
	"net/http"
	"regexp"
	"time"
)

func GetRoutesList(logger *logrus.Logger, router *chi.Mux, lm *middleware.Log, u *utils.Utils) *Routes {
	router.Use(lm.Handler)

	// Enable httprate request limiter of 100 requests per minute.
	//
	// In the code example below, rate-limiting is bound to the request IP address
	// via the LimitByIP middleware handler.
	//
	// To have a single rate-limiter for all requests, use httprate.LimitAll(..).
	//
	// Please see _example/main.go for other more, or read the library code.
	router.Use(httprate.LimitByIP(100, 1*time.Minute))

	staticFilesDir := http.Dir("./web/template/" + viper.GetString("theme") + "/static/")

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)

	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(staticFilesDir)))
	return &Routes{
		Routes: []*Route{
			{
				Pattern: "/",
				Method:  "GET",
				Handler: handlers.NewHome(logger).Handler,
			},
			{
				Pattern: "/schedule", // curl -X POST -d "0 * * * * *\necho\nHello, world!" http://localhost:11111/schedule
				Method:  "POST",
				Handler: handlers.NewSchedule(u).Handler,
			},
			{
				Pattern: "/get/crons", // curl -X POST -d "0 * * * * *\necho\nHello, world!" http://localhost:11111/schedule
				Method:  "GET",
				Handler: handlers.NewCronsList(u).Handler,
			},
			{
				Pattern: "/list", // curl -X POST -d "0 * * * * *\necho\nHello, world!" http://localhost:11111/schedule
				Method:  "GET",
				Handler: handlers.NewCronsListPage(u).Handler,
			},
		},
		Router: router,
		Log:    logger,
	}
}
