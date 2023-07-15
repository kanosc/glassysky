package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	. "github.com/kanosc/glassysky/controller"

	"encoding/hex"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/acme/autocert"
)

var router = gin.Default()
var portFlag = flag.String("port", "9090", "usage: -port 443|9090, 9090 as default")
var modeFlag = flag.String("mode", "production", "usage: -mode debug|production, production as default")
var redisClient *redis.Client
var ctx = context.Background()

func redisInit() {
	if *modeFlag == "debug" {
		redisClient = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			//Password: "929319", // no password set
			Password: "myredis6379", // no password set
			DB:       0,
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "929319", // no password set
			//Password: "myredis6379", // no password set
			DB: 0,
		})
	}
	val, err := redisClient.Get(ctx, "testLink").Result()
	if err != nil {
		panic(err)
	}
	log.Println("testLink:" + val)
	log.Println("redis init success")

}

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
	redisInit()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.StaticFS("/assets/", http.Dir("assets"))
	router.StaticFS("/resource/", http.Dir("file_storage"))
	//router.StaticFile("/a.txt", "./a.txt")
	router.SetFuncMap(template.FuncMap{
		"DelPoint": DelPoint,
	})
	//router.LoadHTMLGlob("pages/*")
	assetDir := "assets"
	htmls := []string{"login.html", "download_card.html", "download_list.html", "index.html", "luck.html", "nav.html", "footer.html",
		"download_left.html", "download_right.html", "download_right2.html", "upload_frame.html", "chat.html", "chat_http.html", "chat_login.html",
		"chat_nav.html", "chat_room.html", "chat_room_login.html", "chat_left.html"}
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
	router.GET("/chat_login", ChatCookieChecker(), func(c *gin.Context) {
		//c.HTML(http.StatusOK, "chat_login.html", gin.H{})

		username, exist := c.Get("username")
		if !exist {
			c.HTML(http.StatusOK, "chat_login.html", gin.H{})
			return
		}

		rooms, err := redisClient.LRange(ctx, "chat:rooms", 0, -1).Result()
		log.Println("chat rooms 1 ***************", rooms)
		if err != nil {
			log.Println(err.Error())
		}
		for _, r := range rooms {
			_, err := redisClient.Get(ctx, "chat:roomname:"+r+":messages").Result()
			if err == redis.Nil {
				log.Println("room expires, start to delete ", r)
				_ = redisClient.LRem(ctx, "chat:rooms", 1, r)
			}
		}

		/*
			for i, _ := range historyMsg {
				historyMsg[i] += "\n"
			}
		*/
		c.HTML(http.StatusOK, "chat_room.html", gin.H{
			"username": username,
			"rooms":    rooms,
		})
	})

	router.POST("/chat_room", func(c *gin.Context) {
		username := c.PostForm("username")
		SetCookieClient(c, "username", username, 7200)
		rooms, err := redisClient.LRange(ctx, "chat:rooms", 0, -1).Result()
		log.Println("chat rooms 2 ***************", rooms)
		if err != nil {
			log.Println(err.Error())
		}

		c.HTML(http.StatusOK, "chat_room.html", gin.H{
			"username": username,
			"rooms":    rooms,
		})
	})

	router.GET("/chat_room", ChatCookieChecker(), func(c *gin.Context) {
		username,_ := c.Cookie("username")
		rooms, err := redisClient.LRange(ctx, "chat:rooms", 0, -1).Result()
		if err != nil {
			log.Println(err.Error())
		}

		c.HTML(http.StatusOK, "chat_room.html", gin.H{
			"username": username,
			"rooms":    rooms,
		})
	})

	router.POST("/create_room", func(c *gin.Context) {
		roomname := c.PostForm("roomname")
		maxuser := c.PostForm("maxuser")
		password := c.PostForm("password")

		log.Println("################", roomname, maxuser, password)
		roomk := "chat:roomname:" + roomname + ":messages"
		roomp := "chat:roomname:" + roomname + ":password"
		rooms := "chat:rooms"
		_, err := redisClient.Get(ctx, roomk).Result()
		if err != redis.Nil {
			log.Println("room already exist, return")
			c.String(http.StatusForbidden, "room has existed")
			c.Abort()
			return
		}

		_, err = redisClient.LPush(ctx, roomk, "room created success").Result()
		if err != nil {
			log.Println("create room", roomname, "success", maxuser)
		}

		_, err = redisClient.Expire(ctx, roomk, 24*time.Hour).Result()
		if err != nil {
			log.Println("Set expire time for room success")
		}

		if password != "" {
			_, err = redisClient.Set(ctx, roomp, password, 24*time.Hour).Result()
			if err != nil {
				log.Println("set password", roomp, "failed")
			}

			log.Println("set password", roomp, "success")
		} else {
			log.Println("no password for room", roomname)
		}

		_, err = redisClient.LPush(ctx, rooms, roomname).Result()

		c.String(http.StatusOK, "create room seccessful")
	})

	router.POST("/chat/:roomname", ChatCookieChecker(), func(c *gin.Context) {
		roomname := c.Param("roomname")
		password := c.PostForm("password")
		username, uexist:= c.Cookie("username")
		log.Printf("%v is logging in room [%v] with password [%v]\n", username, roomname, password)
		if roomname == "" || uexist != nil {
			c.HTML(http.StatusOK, "chat_login.html", gin.H{})
			return
		}

		pwd, err := redisClient.Get(ctx, "chat:roomname:"+roomname+":password").Result()
		log.Printf("password in redis %v\n", pwd)
		if err == redis.Nil || password != pwd {
			log.Println("password is wrong or not exist")
			c.String(http.StatusForbidden, "password is wrong or not exist")
			c.Abort()
			return
		}

		roomnamestr := hex.EncodeToString([]byte(roomname))
		log.Println("^^^^^^^^^^^^^^^^^^^^^^", roomnamestr)
		SetCookieClient(c, roomnamestr, password, 7200)

		historyMsg, err := redisClient.LRange(ctx, "chat:roomname:"+roomname+":messages", 0, -1).Result()
		if err != nil {
			log.Println(err.Error())
		}
		/*
			for i, _ := range historyMsg {
				historyMsg[i] += "\n"
			}
		*/

		if *modeFlag == "debug" {
			c.HTML(http.StatusOK, "chat_http.html", gin.H{
				"roomname": roomname,
				"username": username,
				"chatmsg":  historyMsg,
			})
		} else {
			c.HTML(http.StatusOK, "chat.html", gin.H{
				"roomname": roomname,
				"username": username,
				"chatmsg":  historyMsg,
			})
		}

	})

	router.GET("/chat/:roomname", func(c *gin.Context) {
		roomname := c.Param("roomname")
		username, uexist := c.Cookie("username")
		roomnamestr := hex.EncodeToString([]byte(roomname))
		password, _ := c.Cookie(roomnamestr)
		//log.Println("********** roomname, username", roomname, username)
		if roomname == "" || uexist != nil {
			c.HTML(http.StatusOK, "chat_login.html", gin.H{})
			return
		}

		_, err := redisClient.Get(ctx, "chat:roomname:"+roomname+":messages").Result()
		if err == redis.Nil {
			log.Println("room does not exist, return")
			c.String(http.StatusForbidden, "room does not exist")
			c.Abort()
			return
		}

		pwd, err := redisClient.Get(ctx, "chat:roomname:"+roomname+":password").Result()
		if err != redis.Nil && pwd != password {
			log.Println("password of room: " + pwd)
			log.Println("password of in cookie: " + password)
			c.HTML(http.StatusOK, "chat_room_login.html", gin.H{"roomname": roomname, "username": username})
			return
		}

		historyMsg, err := redisClient.LRange(ctx, "chat:roomname:"+roomname+":messages", 0, -1).Result()
		if err != nil {
			log.Println(err.Error())
		}
		/*
			for i, _ := range historyMsg {
				historyMsg[i] += "\n"
			}
		*/

		if *modeFlag == "debug" {
			c.HTML(http.StatusOK, "chat_http.html", gin.H{
				"roomname": roomname,
				"username": username,
				"chatmsg":  historyMsg,
			})
		} else {
			c.HTML(http.StatusOK, "chat.html", gin.H{
				"roomname": roomname,
				"username": username,
				"chatmsg":  historyMsg,
			})
		}

	})

	router.GET("/wschat/:roomname", func(c *gin.Context) {
		var roomname interface{} = c.Param("roomname")
		//roomname := "testroom"
		log.Println("handling msg request 1")
		if roomname != "" {
			log.Println("handling msg request 2")
			k := map[string]interface{}{"roomname": roomname}
			m.HandleRequestWithKeys(c.Writer, c.Request, k)
		}
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		log.Println("handling msg request 3")
		roomname, _ := s.Keys["roomname"].(string)
		msgkey := "chat:roomname:" + roomname + ":messages"
		roompwd := "chat:roomname:" + roomname + ":password"
		//	msgkey := "chat:roomname:" + "testroom" + ":messages"
		msgCount, err := redisClient.LLen(ctx, msgkey).Result()
		if err != nil {
			log.Println(err.Error())
		}
		if msgCount >= 500 {
			log.Println("length of messages is out of limit 500, delete one first")
			_, err = redisClient.RPop(ctx, msgkey).Result()
			if err != nil {
				log.Println(err.Error())
			}
		}
		_, err = redisClient.LPush(ctx, msgkey, string(msg)).Result()
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("inserting msg", string(msg))
		_, err = redisClient.Expire(ctx, msgkey, 24*time.Hour).Result()
		if err != nil {
			log.Println("Set expire time for room msg fail")
		}
		_, err = redisClient.Expire(ctx, roompwd, 24*time.Hour).Result()
		if err != nil {
			log.Println("Set expire time for room password fail")
		}

		//m.Broadcast(msg)
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	if *modeFlag == "debug" {
		startServerLocal(router, *portFlag)
	} else {
		startServerRemoteTLS(router)
	}

}
