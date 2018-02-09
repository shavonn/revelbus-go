package session

import (
	"time"

	"github.com/alexedwards/scs"
	"github.com/spf13/viper"
)

var session *scs.Manager

func GetSession() *scs.Manager {
	if session == nil {
		return createSession()
	}
	return session
}

func createSession() *scs.Manager {
	sessionManager := scs.NewCookieManager(viper.GetString("secret"))
	sessionManager.Lifetime(12 * time.Hour)
	sessionManager.Persist(true)
	return sessionManager
}
