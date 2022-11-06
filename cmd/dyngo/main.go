package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"

	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/airtableclient"

	_ "embed"
)

//go:embed data/logtail.js
var logtailjs string

var debugFlag bool

func main() {
	flag.BoolVar(&debugFlag, "debug", false, "Output debug info")

	atc := airtableclient.New(
		os.Getenv("AIRTABLE_APIKEY"),
		os.Getenv("AIRTABLE_BASE_ID"),
	)
	// Start the initial run. This way we know the website is up to date.
	maybeStartJob(atc)

	var kickerPage = makeKickerHTML("MyName")

	e := echo.New()

	if debugFlag {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root: "public",
		},
	))

	//	e.GET("/", func(c echo.Context) error {
	//		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	//	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/logtail.js", func(c echo.Context) error {
		// Currently there is no CSS for this.
		return c.HTML(http.StatusOK, logtailjs)
	})

	e.GET("/logtail.css", func(c echo.Context) error {
		// Currently there is no CSS for this.
		return c.HTML(http.StatusOK, "")
	})

	e.GET("/.build", func(c echo.Context) error {
		maybeStartJob(atc)
		return c.HTML(http.StatusOK, kickerPage)
	})

	// route `/log.txt` comes from the filesystem.
	// This won't work because the JS wants to download ranges.
	//	e.GET("/.build/log.txt", func(c echo.Context) error {
	//		f, err := os.Open("/log.txt")
	//		if err != nil {
	//			return c.String(http.StatusOK, "Nothing yet!\n")
	//		}
	//		return c.Stream(http.StatusOK, echo.MIMETextPlainCharsetUTF8, f)
	//	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

// ?_=1667674645947

var kickerTmpl = `<head>
<title>Log viewer for {{.name}}</title>
<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0/jquery.min.js"></script>
<script type="text/javascript" src="/logtail.js"></script>
<link href="/logtail.css" rel="stylesheet" type="text/css">
</head>

<body>
<div id="header">
    <b>Click reload to regenerate {{.name}}</b> --- 
	<a href="./.build">Reversed</a> or
	<a href="./.build?noreverse">chronological</a> view.
	<a id="pause" href='#'>Pause</a>.
</div>
<pre id="data">Loading...</pre>
</body>`

func makeKickerHTML(name string) string {
	var b bytes.Buffer

	t := template.Must(template.New("kicker").Parse(kickerTmpl))
	err := t.Execute(&b, map[string]string{
		"name": name,
	})

	if err != nil {
		log.Fatal(err)
	}
	return b.String()
}
