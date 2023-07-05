package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func ChatCookieChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("chatname")
		log.Println("client cookie is", cookie)

		if err != nil {
			c.HTML(http.StatusOK, "chat_login.html", gin.H{
				"title": "Main website",
			})
			c.Abort()
			return
		}
		c.Set("chatname", cookie)
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

type RandomImage struct {
	Error  int    `json:"error"`
	Img    string `json:"img"`
	Result int    `json:"result"`
}

func GetRandomImage() (string, error) {
	// 发送 GET 请求获取 JSON 数据
	resp, err := http.Get("https://img.xjh.me/random_img.php?return=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 解析 JSON 数据
	var randomImage RandomImage
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&randomImage)
	if err != nil {
		return "", err
	}

	// 拼接图片 URL
	imgUrl := "https:" + randomImage.Img

	return imgUrl, nil
}

func HandleToday(c *gin.Context) {
	crand, err1 := c.Cookie("luck")
	cimg, err2 := c.Cookie("luckimg")
	if err1 != nil || err2 != nil {
		randNum := RollInt(int64(len(luckContent.Luck) - 1))
		randImg, err := GetRandomImage()
		if err != nil {
			log.Println(err.Error())
			randImg = "#"
		}
		rns := strconv.Itoa(randNum)
		log.Println("generate rand value:", randNum)
		SetCookieToday(c, "luck", rns)
		SetCookieToday(c, "luckimg", randImg)
		c.HTML(http.StatusOK, "luck.html", gin.H{
			"title":     luckContent.Luck[randNum],
			"content":   luckContent.LuckInfo[randNum],
			"yiThing":   luckContent.GoodThings[randNum],
			"buyiThing": luckContent.BadThings[randNum],
			"date":      GetDateStr(),
			"imageURL":  randImg,
		})
	} else {
		cnum, e := strconv.Atoi(crand)
		if e != nil || cnum < 0 || cnum > len(luckContent.Luck)-1 {
			log.Println("cookie error, invalid rand, reset cookie")
			c.Abort()
		} else {
			c.HTML(http.StatusOK, "luck.html", gin.H{
				"title":     luckContent.Luck[cnum],
				"content":   luckContent.LuckInfo[cnum],
				"yiThing":   luckContent.GoodThings[cnum],
				"buyiThing": luckContent.BadThings[cnum],
				"date":      GetDateStr(),
				"imageURL":  cimg,
			})

		}

	}

	log.Println("handle luck success")
}

func HandleResetToday(c *gin.Context) {
	DeleteCookieToday(c, "luck")
	DeleteCookieToday(c, "luckimg")
	c.Redirect(http.StatusFound, "/today")
}

func HandleSuccess(c *gin.Context) {
	c.String(http.StatusOK, "success")
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

func SetCookieToday(c *gin.Context, cookieName string, cookieValue string) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		// 计算当天的24点时间
		now := time.Now()
		t := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

		// 设置cookie
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    cookieValue,
			Path:     "/",
			Expires:  t,
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
		log.Printf("Set new cookie value: %s \n", cookieValue)
	} else {
		log.Printf("Cookie value: %s \n", cookie)

	}

}

func DeleteCookieToday(c *gin.Context, cookieName string) {
	expiration := time.Now().AddDate(0, 0, -1) // 将过期时间设置为过去的时间，即立即失效

	cookie := http.Cookie{
		Name:    cookieName,
		Value:   "",
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(c.Writer, &cookie)

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
	//HandleDownload_list(c)
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func HandleDownload_card(c *gin.Context) {
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
	c.HTML(http.StatusOK, "download_card.html", gin.H{
		"files": files,
	})

}

func HandleDownload_list(c *gin.Context) {
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
	c.HTML(http.StatusOK, "download_list.html", gin.H{
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
	HandleDownload_list(c)
	//c.String(http.StatusOK, fmt.Sprintf("%s has been deleted.", filename))

}
