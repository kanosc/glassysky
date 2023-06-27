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
	"github.com/olahol/melody"
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

func makePath(ps []string, d string) {
	for i, p := range ps {
		ps[i] = d + "/" + p
	}
}

func main() {

	flag.Parse()
	log.Println("run mode:", *modeFlag)
	log.Println("port set:", *portFlag)
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.StaticFS("/assets/", http.Dir("assets"))
	router.StaticFS("/resource/", http.Dir("file_storage"))
	router.SetFuncMap(template.FuncMap{
		"DelPoint": DelPoint,
	})
	//router.LoadHTMLGlob("pages/*")
	assetDir := "assets"
	htmls := []string{"login.html", "download_card.html", "download_list.html", "index.html", "luck.html", "nav.html", "footer.html",
		"download_left.html", "download_right.html", "download_right2.html", "upload_frame.html", "chat.html"}
	makePath(htmls, assetDir)
	router.LoadHTMLFiles(htmls...)

	router.GET("/", HandleStart)
	router.GET("/index", HandleIndex)
	router.GET("/today", HandleToday)
	router.GET("/resetToday", HandleResetToday)

	router.GET("/download_card", CookieChecker(), HandleDownload_card)
	router.GET("/download_list", CookieChecker(), HandleDownload_list)
	router.GET("/delete", CookieChecker(), HandleDelete)
	router.POST("/upload", CookieChecker(), HandleUpload)

	//router.GET("/loginDownload", CookieChecker(), HandleDownload)
	router.POST("/downloadVerify", MakeAuthVerifyHandler("test", "file123", HandleSuccess))
	router.GET("/logout", HandleLogout)

	m := melody.New()
	router.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{
        "title": "Chatroom v1.0",
		})
    })

    router.GET("/wschat", func(c *gin.Context) {
        m.HandleRequest(c.Writer, c.Request)
    })
    m.HandleMessage(func(s *melody.Session, msg []byte) {
        m.Broadcast(msg)
    })

	if *modeFlag == "debug" {
		startServerLocal(router, *portFlag)
	} else {
		startServerRemoteTLS(router)
	}

}
