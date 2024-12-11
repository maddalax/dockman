package pages

import (
	"dockman/app"
	"dockman/app/ui"
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
)

func RegisterUser(ctx *h.RequestContext) *h.Partial {
	if !ctx.IsHttpPost() {
		return nil
	}

	if ctx.FormValue("password") != ctx.FormValue("password-confirm") {
		ctx.Response.WriteHeader(400)
		return ui.SwapFormError(ctx, "passwords do not match")
	}

	payload := &app.User{
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
	}

	_, err := app.UserCreate(
		ctx.ServiceLocator(),
		payload,
	)

	if err != nil {
		ctx.Response.WriteHeader(400)
		return ui.SwapFormError(ctx, err.Error())
	}

	session, err := app.UserLogin(ctx.ServiceLocator(), payload.Email, ctx.FormValue("password"))

	if err != nil {
		ctx.Response.WriteHeader(500)
		return ui.SwapFormError(ctx, "something went wrong")
	}

	session.Write(ctx)

	return h.RedirectPartial("/")
}

func LoginUser(ctx *h.RequestContext) *h.Partial {
	if !ctx.IsHttpPost() {
		return nil
	}

	payload := &app.User{
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
	}

	session, err := app.UserLogin(
		ctx.ServiceLocator(),
		payload.Email,
		payload.Password,
	)

	if err != nil {
		ctx.Response.WriteHeader(400)
		return ui.SwapFormError(ctx, err.Error())
	}

	session.Write(ctx)

	return h.RedirectPartial("/")
}

func Login(ctx *h.RequestContext) *h.Page {
	isRegister := !app.UserIsInitialSetup(ctx.ServiceLocator())

	return RootPage(
		ctx,
		h.Div(
			h.Class("flex flex-col items-center justify-center mx-auto min-h-screen bg-neutral-100 w-full"),
			h.Div(
				h.Class("bg-white p-8 rounded-lg shadow-lg"),
				h.H2F(
					h.Ternary(isRegister, fmt.Sprintf("Setup %s", app.AppName), fmt.Sprintf("Sign in to %s", app.AppName)),
					h.Class("text-3xl font-bold text-center mb-6"),
				),
				h.Form(
					h.TriggerChildren(),
					h.PostPartial(h.Ternary(isRegister, RegisterUser, LoginUser)),
					h.Attribute("hx-swap", "none"),
					h.Class("flex flex-col gap-4 max-w-md"),
					ui.Input(ui.InputProps{
						Id:       "username",
						Name:     "email",
						Label:    "Email Address",
						Type:     "email",
						Required: true,
						Children: []h.Ren{
							h.Attribute("autocomplete", "off"),
							h.MaxLength(50),
						},
					}),
					ui.Input(ui.InputProps{
						Id:       "password",
						Name:     "password",
						Label:    "Password",
						Type:     "password",
						Required: true,
						Children: []h.Ren{
							h.MinLength(6),
						},
						HelpText: h.Div(
							h.Pf("Do not lose your password. It cannot be recovered."),
							h.Pf("%s will have to be reset to recover access.", app.AppName),
						),
					}),
					h.If(
						isRegister,
						ui.Input(ui.InputProps{
							Id:       "password-confirm",
							Name:     "password-confirm",
							Label:    "Confirm Password",
							Type:     "password",
							Required: true,
							Children: []h.Ren{
								h.MinLength(6),
							},
						}),
					),
					// Error message
					ui.FormError(""),
					// Submit button at the bottom
					ui.SubmitButton(ui.ButtonProps{
						Text: h.Ternary(isRegister, "Create default admin account", "Sign In"),
					}),
				),
			),
		),
	)
}
