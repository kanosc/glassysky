package controller

import (
	"fmt"
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
				"action": "uploadVerify",
			})
			c.Abort()
		}
		c.Next()
	}
}

func CookieChecker2() gin.HandlerFunc {
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
		"title":     todayLuck[randNum],
		"content":   todayContent[randNum],
		"yiThing":   yiThing[randNum],
		"buyiThing": buyiThing[randNum],
		"date":      GetDateStr(),
	})
}

func HandleLoginUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "Main website",
	})
}

func HandleLoginDownload(c *gin.Context) {
	log.Println("recive request")
	HandleDownload(c)
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
	cookieName := "UserCookie"
	if username == name && password == pwd {
		clientUUID, _ := uuid.NewUUID()
		SetCookieDefault(c, cookieName, clientUUID.String())
		next(c)
	} else {
		c.String(http.StatusForbidden, fmt.Sprintf("Name or password is wrong"))
	}
	log.Println(username, password)
}

func HandleFile(c *gin.Context) {

	log.Println("recive request")
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "Main website",
	})
}

func HandleUpload(c *gin.Context) {

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

	curDirSize, _ := GetDirSize("./file_storage/")
	restDirSize := MAX_DIR_SIZE - curDirSize
	log.Println("current dir size ", curDirSize, "B")
	log.Println("rest dir size ", restDirSize/1024/1024, "MB")
	if totalUploadSize > restDirSize {
		c.String(http.StatusForbidden, fmt.Sprintf("File size is too large"))
		return
	}

	for _, file := range files {
		log.Println(file.Filename)
		if flg, _ := FileExistInDir("./file_storage/", file.Filename); flg {
			c.String(http.StatusForbidden, fmt.Sprintf("Filename already exist, change the filename or delete exist file."))
		}

		// 上传文件至指定目录
		c.SaveUploadedFile(file, "./file_storage/"+file.Filename)
	}
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
		c.String(http.StatusForbidden, fmt.Sprintf("Filename not exist."))
		c.Abort()
		return
	}
	err := os.Remove("./file_storage/" + filename)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Delete file " + filename + " success")
	}
	//handleDownload(c)
	c.String(http.StatusOK, fmt.Sprintf("%s has been deleted.", filename))

}
