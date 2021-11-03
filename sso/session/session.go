package session

import (
	"log"
	"net/http"
	"net/url"
	"sso/sso/settings"

	"gopkg.in/boj/redistore.v1"

	"encoding/gob"

	"github.com/gorilla/sessions"
)

var store *redistore.RediStore

func Init(cfg *settings.SessionConfig) (err error) {

	gob.Register(url.Values{})

	store, err = redistore.NewRediStore(10, "tcp", ":6379", "", []byte(cfg.HashKey))
	if err != nil {
		log.Printf("session初始化失败：%s\n", err.Error())
		return
	}
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 20,
		HttpOnly: true,
		Secure:   true,
	}

	return
}

func Get(r *http.Request, name string) (val interface{}, err error) {

	session, err := store.Get(r, settings.Conf.SessionConfig.SessionId)
	if err != nil {
		return
	}

	val = session.Values[name]

	return
}

func Set(w http.ResponseWriter, r *http.Request, name string, val interface{}) (err error) {

	session, err := store.Get(r, settings.Conf.SessionConfig.SessionId)
	if err != nil {
		return
	}

	session.Values[name] = val

	err = session.Save(r, w)

	return
}

func Delete(w http.ResponseWriter, r *http.Request, name string) (err error) {

	session, err := store.Get(r, settings.Conf.SessionConfig.SessionId)
	if err != nil {
		return
	}
	delete(session.Values, name)
	err = session.Save(r, w)
	return
}
