package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//set a cookie to both local and client if cookie doesn't exist in client
func SetCookieDefault(c *gin.Context, cookieName string, cookieValue string, expireTime int) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		c.SetCookie(cookieName, cookieValue, expireTime, "", "", false, true)

		session := sessions.Default(c)
		session.Options(sessions.Options{MaxAge: expireTime})
		session.Set(cookieValue, true)
		session.Save()

		log.Printf("Set new cookie value: %s \n", cookieValue)
	} else {
		log.Printf("Cookie value: %s \n", cookie)

	}

}

//set a cookie to client if it doesn't exist
func SetCookieClient(c *gin.Context, cookieName string, cookieValue string, expireTime int) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		c.SetCookie(cookieName, cookieValue, expireTime, "", "", false, true)

		log.Printf("Cookie doesn't exist, set new cookie value: %s \n", cookieValue)
	} else {
		log.Printf("Cookie already exists, value: %s \n", cookie)

	}

}

//set a cookie to client which expires at the end of the day
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

func DeleteCookieClient(c *gin.Context, cookieName string) {
	expiration := time.Now().AddDate(0, 0, -1) // 将过期时间设置为过去的时间，即立即失效

	cookie := http.Cookie{
		Name:    cookieName,
		Value:   "",
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(c.Writer, &cookie)

}

func DeleteCookieDefault(c *gin.Context, cookieName string) {
	cookieValue, err := c.Cookie(cookieName)
	if err != nil {

		session := sessions.Default(c)
		session.Options(sessions.Options{MaxAge: -1})
		session.Set(cookieValue, false)
		session.Save()

		expiration := time.Now().AddDate(0, 0, -1) // 将过期时间设置为过去的时间，即立即失效

		cookie := http.Cookie{
			Name:    cookieName,
			Value:   "",
			Expires: expiration,
			Path:    "/",
		}
		http.SetCookie(c.Writer, &cookie)

		log.Printf("Delete cookie: %s \n", cookieName)
	} else {
		log.Printf("Cookie not exists: %s \n", cookieName)

	}

}
