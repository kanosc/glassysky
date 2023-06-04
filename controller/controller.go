package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/kanosc/glassysky/common"
)

const (
	clientCookieName = "UserCookie"
)

func init() {
	luckContent = ReadLuckContent()
	log.Println("Loading luck content from json file")
	log.Println(*luckContent)
}

const (
	MAX_DIR_SIZE = 3 * 1024 * 1024 * 1024
)

type LuckContent struct {
	Luck       []string
	LuckInfo   []string
	GoodThings []string
	BadThings  []string
}

var luckContent *LuckContent

func ReadLuckContent() *LuckContent {
	luck := new(LuckContent)
	b, err := ioutil.ReadFile("./controller/LuckContent.json")
	err = json.Unmarshal(b, luck)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Println("unmarshall success")
	return luck

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
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title":  "Main website",
				"action": "downloadVerify",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func HandleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func HandleStart(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func HandleToday(c *gin.Context) {
	log.Println("recive request")
	randNum := RollInt(int64(8))
	c.HTML(http.StatusOK, "today.html", gin.H{
		"title":     luckContent.Luck[randNum],
		"content":   luckContent.LuckInfo[randNum],
		"yiThing":   luckContent.GoodThings[randNum],
		"buyiThing": luckContent.BadThings[randNum],
		"date":      GetDateStr(),
	})
}

func HandleLogout(c *gin.Context) {
	log.Println("recive logout request")
	c.SetCookie(clientCookieName, "", -1, "", "", false, true)
	HandleIndex(c)
}

func SetCookieDefault(c *gin.Context, cookieName string, cookieValue string) {
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

func MakeAuthVerifyHandler(name string, pwd string, handler func(c *gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		HandleVerifyAuth(c, name, pwd, handler)
	}
}

func HandleVerifyAuth(c *gin.Context, name string, pwd string, next func(c *gin.Context)) {
	log.Println("recive request")
	username := c.PostForm("username")
	password := c.PostForm("password")
	cookieName := clientCookieName
	if username == name && password == pwd {
		clientUUID, _ := uuid.NewUUID()
		SetCookieDefault(c, cookieName, clientUUID.String())
		next(c)
	} else {
		c.String(http.StatusForbidden, fmt.Sprintf("Name or password is wrong"))
		c.Abort()
	}
	log.Println(username, password)
}

func HandleUpload(c *gin.Context) {

	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload"]
	if len(files) == 0 {
		c.String(http.StatusForbidden, fmt.Sprintf("No file recieved, please select files."))
		c.Abort()
		return
		log.Println("number of file is 0")

	}
	log.Printf("%d files recieved", len(files))

	var totalUploadSize int64

	for _, file := range files {
		if strings.Contains(file.Filename, "/") || len(file.Filename) > 100 {
			c.String(http.StatusForbidden, fmt.Sprintf("Filename is too long, max length is 100 characters."))
			return
		}
		totalUploadSize += file.Size
	}

	curDirSize, _ := GetDirSize("./file_storage/")
	restDirSize := MAX_DIR_SIZE - curDirSize
	log.Println("current dir size ", curDirSize, "B")
	log.Println("rest dir size ", restDirSize/1024/1024, "MB")
	if totalUploadSize > restDirSize {
		c.String(http.StatusForbidden, fmt.Sprintf("File size is too large"))
		c.Abort()
		return
	}

	for _, file := range files {
		log.Println(file.Filename)
		if flg, _ := FileExistInDir("./file_storage/", file.Filename); flg {
			c.String(http.StatusForbidden, fmt.Sprintf("Filename already exist, change the filename or delete exist file."))
			c.Abort()
			return
		}

		// 上传文件至指定目录
		c.SaveUploadedFile(file, "./file_storage/"+file.Filename)
	}
	//HandleDownload(c)
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func HandleDownload(c *gin.Context) {
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

func HandleDelete(c *gin.Context) {
	filename := c.Query("filename")
	if flg, _ := FileExistInDir("./file_storage/", filename); !flg {
		c.String(http.StatusNotFound, fmt.Sprintf("filename not exist."))
		c.Abort()
		return
	}
	err := os.Remove("./file_storage/" + filename)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, fmt.Sprintf("filename exist, but remove failed."))
		c.Abort()
		return
	}

	log.Println("Delete file " + filename + " success")
	HandleDownload(c)
	//c.String(http.StatusOK, fmt.Sprintf("%s has been deleted.", filename))

}
