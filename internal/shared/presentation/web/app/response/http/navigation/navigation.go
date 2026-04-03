package navigation

import (
	"net/http"

	webflash "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/flash"
)

type Navigator struct {
	flash *webflash.Adapter
}

func New(secureCookies bool) *Navigator {
	return &Navigator{flash: webflash.New(secureCookies)}
}

func (n *Navigator) SeeOther(w http.ResponseWriter, r *http.Request, location string) {
	http.Redirect(w, r, location, http.StatusSeeOther)
}

func (n *Navigator) ReadFlash(w http.ResponseWriter, r *http.Request, cookieName string) webflash.Message {
	return n.flash.Read(w, r, cookieName)
}

func (n *Navigator) ClearFlash(w http.ResponseWriter, cookieName string) {
	n.flash.Clear(w, cookieName)
}

func (n *Navigator) SeeOtherWithFlash(w http.ResponseWriter, r *http.Request, location, cookieName string, flash webflash.Message) {
	n.flash.Redirect(w, r, location, cookieName, flash)
}
