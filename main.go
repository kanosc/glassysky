package main

import (
	"log"
	"net/http"

	. "github.com/kanosc/glassysky/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func main() {

	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.StaticFS("/pages/", http.Dir("pages"))
	router.StaticFS("/resource/", http.Dir("file_storage"))
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
	//router.Run(":80")
	//router.RunTLS(":443", "./cert/cert.crt", "./cert/rsa_private.key")
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("glassysky.cn"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}
	//log.Fatal(autotls.Run(router, "example1.com"))
	log.Fatal(autotls.RunWithManager(router, &m))

}
