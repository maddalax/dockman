package pages

import "github.com/maddalax/htmgo/framework/h"

func LogoutPage(ctx *h.RequestContext) *h.Page {
	ctx.Response.Header().Set("Set-Cookie", "session_id=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT")
	ctx.Redirect("/login", 302)
	return h.EmptyPage()
}
