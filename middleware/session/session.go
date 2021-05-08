package session

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/tdycwym/edgex_admin/config"
)

type userInfo struct {
	UserID   int64
	UserName string
}

// KEY gin session key
const KEY = "bmp1LWVkZ2V4"

// CookieName ...
const CookieName = "session_id"

// // EnableCookieSession ...
// func EnableCookieSession() gin.HandlerFunc {
// 	store := cookie.NewStore([]byte(KEY))
// 	return sessions.Sessions("session", store)
// }

// EnableRedisSession ...
func EnableRedisSession() gin.HandlerFunc {
	store, _ := redis.NewStore(10, "tcp", config.RedisConf.Address, "", []byte(KEY))
	return sessions.Sessions("session_token", store)
}

// SessionMiddleware ...
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, _ := c.Cookie(CookieName)
		if sessionID == "" {
			sessionID = uuid.NewV4().String()
		}
		session := sessions.Default(c)
		sessionValue := session.Get(sessionID)
		if sessionValue != nil {
			session.Set(sessionID, sessionValue)
			_ = session.Save()
		}
		c.SetCookie(CookieName, sessionID, 20*60, "/", "", false, true)
		c.Set(CookieName, sessionID)
		c.Next()
	}
}

// AuthSessionMiddle ...
func AuthSessionMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, _ := c.Cookie(CookieName)
		session := sessions.Default(c)
		sessionValue := session.Get(sessionID)
		if sessionID == "" || sessionValue == nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			rawData, _ := json.Marshal(gin.H{"error": "Unauthorized"})
			_, _ = c.Writer.Write(rawData)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

// SaveAuthSession 注册和登陆时都需要保存seesion信息
func SaveAuthSession(c *gin.Context, userID int64, username string) {
	session := sessions.Default(c)
	sessionID, _ := c.Get(CookieName)
	userInfo := &userInfo{
		UserID:   userID,
		UserName: username,
	}
	userInfoBytes, _ := json.Marshal(userInfo)
	session.Set(sessionID, string(userInfoBytes))
	_ = session.Save()
}

// ClearAuthSession 退出时清除session
func ClearAuthSession(c *gin.Context) {
	sessionID, exsit := c.Get(CookieName)
	if !exsit || sessionID.(string) == "" {
		return
	}
	session := sessions.Default(c)
	session.Delete(sessionID)
	_ = session.Save()
}

// GetSessionUserID ...
func GetSessionUserID(c *gin.Context) int64 {
	sessionID, exsit := c.Get(CookieName)
	if !exsit || sessionID.(string) == "" {
		return 0
	}
	session := sessions.Default(c)
	sessionValue := session.Get(sessionID)
	if sessionValue == nil {
		return 0
	}
	userInfo := &userInfo{}
	err := json.Unmarshal([]byte(sessionValue.(string)), &userInfo)
	if err != nil || userInfo == nil || userInfo.UserID <= 0 {
		return 0
	}
	return userInfo.UserID
}

// GetSessionUsername ...
func GetSessionUsername(c *gin.Context) string {
	sessionID, exsit := c.Get(CookieName)
	if !exsit || sessionID.(string) == "" {
		return ""
	}
	session := sessions.Default(c)
	sessionValue := session.Get(sessionID)
	if sessionValue == nil {
		return ""
	}
	userInfo := &userInfo{}
	err := json.Unmarshal([]byte(sessionValue.(string)), &userInfo)
	if err != nil || userInfo == nil || userInfo.UserID <= 0 {
		return ""
	}
	return userInfo.UserName
}
