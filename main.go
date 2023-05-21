package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"

	. "github.com/kanosc/glassysky/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

var router = gin.Default()
var portFlag = flag.String("port", "9090", "usage: -port 443|9090, 9090 as default")
var modeFlag = flag.String("mode", "production", "usage: -mode debug|production, production as default")

func startServerLocal(e *gin.Engine, port string) {
	e.Run(":" + port)
}

func startServerRemoteTLS(e *gin.Engine) {
	//router.RunTLS(":443", "./cert/cert.crt", "./cert/rsa_private.key")
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("glassysky.cn"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}
	//log.Fatal(autotls.Run(router, "example1.com"))
	log.Fatal(autotls.RunWithManager(e, &m))

}

func DelPoint(s string) string {
	sa := strings.Split(s, ".")
	ret := ""
	for _, a := range sa {
		ret = ret + a
	}
	return ret
}

func main() {

	flag.Parse()
	log.Println("run mode:", *modeFlag)
	log.Println("port set:", *portFlag)
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.StaticFS("/pages/", http.Dir("pages"))
	router.StaticFS("/resource/", http.Dir("file_storage"))
	router.SetFuncMap(template.FuncMap{
		"DelPoint": DelPoint,
	})
	//router.LoadHTMLGlob("pages/*")
	router.LoadHTMLFiles("pages/login.html", "pages/download.html", "pages/upload.html", "pages/start.html", "pages/index.html", "pages/today.html")
	router.GET("/", HandleStart)
	router.GET("/index", HandleIndex)
	router.GET("/today", HandleToday)

	router.GET("/download", HandleDownload)
	router.POST("/upload", CookieChecker(), HandleUpload)
	router.GET("/loginUpload", CookieChecker(), HandleLoginUpload)
	router.GET("/loginDownload", CookieChecker2(), HandleLoginDownload)
	router.POST("/uploadVerify", MakeAuthVerifyHandler("test", "file123", HandleFile))
	router.POST("/downloadVerify", MakeAuthVerifyHandler("test", "file123", HandleDownload))
	router.GET("/delete", CookieChecker2(), HandleDelete)

	if *modeFlag == "debug" {
		startServerLocal(router, *portFlag)
	} else {
		startServerRemoteTLS(router)
	}

}
