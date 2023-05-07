package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

const (
	MAX_DIR_SIZE = 500 * 1024 * 1024
)

func main() {

	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.StaticFS("/pages/", http.Dir("pages"))
	router.StaticFS("/resource/", http.Dir("file_storage"))
	//router.LoadHTMLGlob("pages/*")
	router.LoadHTMLFiles("pages/loginUpload.html", "pages/loginDownload.html", "pages/download.html", "pages/upload.html", "pages/start.html", "pages/index.html", "pages/today.html")
	router.GET("/", handleStart)
	router.GET("/index", handleIndex)
	router.GET("/today", handleToday)

	router.GET("/download", handleDownload)
	router.POST("/upload", CookieChecker(), handleUpload)
	router.GET("/loginUpload", CookieChecker(), handleLoginUpload)
	router.GET("/loginDownload", handleLoginDownload)
	router.POST("/uploadVerify", handleVerifyUpload)
	router.POST("/downloadVerify", handleVerifyDownload)
	router.GET("/delete", handleDelete)
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
