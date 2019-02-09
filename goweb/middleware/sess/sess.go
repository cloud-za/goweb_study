package sess

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func New(secret string, name string, maxAge int) gin.HandlerFunc {
	store := memstore.NewStore([]byte(secret))
	store.Options(sessions.Options{
		MaxAge:   maxAge,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	return sessions.Sessions(name, store)
}

func ValidateSession(priv int16) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("privilege")
		if v == nil {
			c.AbortWithStatusJSON(401, gin.H{
				"status": -1,
				"msg":    "invalid access",
			})
			return
		} else {
			sessPriv := v.(int16)
			if sessPriv < priv {
				c.AbortWithStatusJSON(401, gin.H{
					"status": -1,
					"msg":    "invalid access",
				})
				return
			}
			c.Set("session_username", session.Get("username"))
			c.Set("session_displayname", session.Get("displayname"))
			c.Set("session_uid", session.Get("uid"))
			c.Set("session_locale", session.Get("locale"))
			c.Set("session_privilege", sessPriv)
		}
	}
}
