package main

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	MAX_DIR_SIZE = 500 * 1024 * 1024
)

func getDateStr() string {
	t := time.Now()
	weekday := []string{"日", "一", "二", "三", "四", "五", "六"}
	date := fmt.Sprintf("%v年%v月%v日 星期%v", t.Year(), int(t.Month()), t.Day(), weekday[int(t.Weekday())])
	return date
}

func handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func handleStart(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func rollInt(end int64) int {
	ret, _ := rand.Int(rand.Reader, big.NewInt(end+1))
	return int(ret.Int64())
}

func handleToday(c *gin.Context) {
	log.Println("recive request")
	randNum := rollInt(int64(8))
	c.HTML(http.StatusOK, "today.html", gin.H{
		"title":     todayLuck[randNum],
		"content":   todayContent[randNum],
		"yiThing":   yiThing[randNum],
		"buyiThing": buyiThing[randNum],
		"date":      getDateStr(),
	})
}

func handleLoginUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "Main website",
	})
}

func CookieChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("UserCookie")
		log.Println("client cookie is", cookie)
		session := sessions.Default(c)
		cookieExist, _ := session.Get(cookie).(bool)
		log.Println(cookieExist)
		log.Println(session.Get(cookie))

		if err != nil || cookieExist != true {
			c.HTML(http.StatusOK, "loginUpload.html", gin.H{
				"title": "Main website",
			})
			c.Abort()
		}
		c.Next()
	}
}

func handleLoginDownload(c *gin.Context) {
	log.Println("recive request")
	c.HTML(http.StatusOK, "loginDownload.html", gin.H{
		"title": "Main website",
	})
}

func setCookieDefault(c *gin.Context, cookieName string, cookieValue string) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		c.SetCookie(cookieName, cookieValue, 1800, "", "", false, true)

		session := sessions.Default(c)
		session.Options(sessions.Options{MaxAge: 1800})
		session.Set(cookieValue, true)
		session.Save()

		log.Printf("Set new cookie value: %s \n", cookieValue)
	} else {
		log.Printf("Cookie value: %s \n", cookie)

	}

}

func handleVerifyUpload(c *gin.Context) {
	log.Println("recive request")
	username := c.PostForm("username")
	password := c.PostForm("password")
	cookieName := "UserCookie"
	if username == "test" && password == "upload" {
		clientUUID, _ := uuid.NewUUID()
		setCookieDefault(c, cookieName, clientUUID.String())
		handleFile(c)
	} else {
		c.String(http.StatusForbidden, fmt.Sprintf("Name or password is wrong"))
	}
	log.Println(username, password)
}

func handleVerifyDownload(c *gin.Context) {
	log.Println("recive request")
	username := c.PostForm("username")
	password := c.PostForm("password")
	cookieName := "UserDownload"
	if username == "test" && password == "download" {
		clientUUID, _ := uuid.NewUUID()
		setCookieDefault(c, cookieName, clientUUID.String())
		handleDownload(c)
	} else {
		c.String(http.StatusForbidden, fmt.Sprintf("Name or password is wrong"))
	}
	log.Println(username, password)
}

func handleFile(c *gin.Context) {

	log.Println("recive request")
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "Main website",
	})
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != "" {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func fileExistInDir(path, filename string) (bool, error) {
	//	var fileExistFlag = false
	/* err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info.Name() == filename {
			fileExistFlag = true
		}
		return err
	}) */
	fullFileName := path + filename
	if _, err := os.Stat(fullFileName); os.IsNotExist(err) {
		return false, err
	} else {
		return true, err
	}
}

func getFileNameList(path string) []string {
	filenames := make([]string, 0)
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		log.Println("###########", len(info.Name()))
		if !info.IsDir() && len(info.Name()) != 0 {
			filenames = append(filenames, info.Name())
		}
		return err
	})
	log.Println(len(filenames))
	for _, j := range filenames {
		log.Println("[" + j + "]")
	}
	return filenames
}

func handleUpload(c *gin.Context) {

	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload"]

	var totalUploadSize int64

	for _, file := range files {
		if strings.Contains(file.Filename, "/") {
			c.String(http.StatusForbidden, fmt.Sprintf("Invalid filename."))
			return
		}
		totalUploadSize += file.Size
	}

	curDirSize, _ := getDirSize("./file_storage/")
	restDirSize := MAX_DIR_SIZE - curDirSize
	log.Println("current dir size ", curDirSize, "B")
	log.Println("rest dir size ", restDirSize/1024/1024, "MB")
	if totalUploadSize > restDirSize {
		c.String(http.StatusForbidden, fmt.Sprintf("File size is too large"))
		return
	}

	for _, file := range files {
		log.Println(file.Filename)
		if flg, _ := fileExistInDir("./file_storage/", file.Filename); flg {
			c.String(http.StatusForbidden, fmt.Sprintf("Filename already exist, change the filename or delete exist file."))
		}

		// 上传文件至指定目录
		c.SaveUploadedFile(file, "./file_storage/"+file.Filename)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func handleDownload(c *gin.Context) {
	log.Println("recive request")
	var files = []string{}
	err := filepath.Walk("./file_storage/", func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return err
	})
	if err != nil {
		log.Println(err.Error())
	}
	c.HTML(http.StatusOK, "download.html", gin.H{
		"files": files,
	})

}

func handleDelete(c *gin.Context) {
	cookie, err := c.Cookie("UserDownload")
	log.Println("client cookie is", cookie)
	session := sessions.Default(c)
	cookieExist, _ := session.Get(cookie).(bool)

	if err != nil || cookieExist != true {
		c.HTML(http.StatusOK, "loginDownload.html", gin.H{
			"title": "Main website",
		})
		return
	}

	filename := c.Query("filename")
	if flg, _ := fileExistInDir("./file_storage/", filename); !flg {
		c.String(http.StatusForbidden, fmt.Sprintf("Filename not exist."))
		return
	}
	err = os.Remove("./file_storage/" + filename)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Delete file " + filename + " success")
	}
	//handleDownload(c)
	c.String(http.StatusOK, fmt.Sprintf("%s has been deleted.", filename))

}

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
