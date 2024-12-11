package app

import (
	"dockman/app/util/json2"
	"errors"
	"github.com/google/uuid"
	"github.com/maddalax/htmgo/framework/h"
	"github.com/maddalax/htmgo/framework/service"
	"github.com/nats-io/nats.go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	Token     string    `json:"token"`
	UserId    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (session *Session) Write(ctx *h.RequestContext) {
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    session.Token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
		Path:     "/",
	}
	ctx.SetCookie(&cookie)
}

func UserIsInitialSetup(locator *service.Locator) bool {
	anyUsers, err := UserGetByPredicate(locator, func(u *User) bool {
		return true
	})

	if err != nil {
		return true
	}

	if anyUsers != nil {
		return true
	}

	return false
}

func UserCreate(locator *service.Locator, user *User) (*User, error) {

	if UserIsInitialSetup(locator) {
		return nil, errors.New("registration is disabled")
	}

	anyUsers, err := UserGetByPredicate(locator, func(u *User) bool {
		return true
	})

	if err != nil {
		return nil, err
	}

	if anyUsers != nil {
		return nil, errors.New("registration is disabled")
	}

	client := service.Get[KvClient](locator)

	user.Email = strings.TrimSpace(user.Email)
	user.Email = strings.ToLower(user.Email)

	u, err := UserGetByPredicate(locator, func(u *User) bool {
		return strings.ToLower(u.Email) == strings.ToLower(user.Email)
	})

	if u != nil {
		return nil, errors.New("user already exists by that email")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPass)
	user.Id = uuid.NewString()

	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "users",
	})

	if err != nil {
		return nil, err
	}

	serialized, err := json2.Serialize(user)

	if err != nil {
		return nil, err
	}

	_, err = bucket.Create(user.Id, serialized)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func CurrentUser(ctx *h.RequestContext) *User {
	if ctx.Get("user") != nil {
		return ctx.Get("user").(*User)
	}
	return nil
}

func UserGetByPredicate(locator *service.Locator, predicate func(user *User) bool) (*User, error) {
	client := service.Get[KvClient](locator)
	users, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "users",
	})
	if err != nil {
		return nil, err
	}

	keys, err := users.ListKeys()

	if err != nil {
		return nil, err
	}

	for key := range keys.Keys() {
		raw, err := users.Get(key)
		if err != nil {
			continue
		}
		user, err := json2.Deserialize[User](raw.Value())
		if err != nil {
			continue
		}
		if predicate(user) {
			return user, nil
		}
	}

	return nil, nil
}

const sessionDuration = 24 * time.Hour

// UserLogin verifies user credentials and generates a session token
func UserLogin(locator *service.Locator, email, password string) (*Session, error) {
	client := service.Get[KvClient](locator)

	// Retrieve the user by email
	user, err := UserGetByPredicate(locator, func(u *User) bool {
		return u.Email == strings.ToLower(strings.TrimSpace(email))
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Create a session token
	token := uuid.NewString()
	session := &Session{
		Token:     token,
		UserId:    user.Id,
		ExpiresAt: time.Now().Add(sessionDuration),
	}

	// Store the session in the "sessions" bucket
	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "sessions",
	})
	if err != nil {
		return nil, err
	}

	serialized, err := json2.Serialize(session)
	if err != nil {
		return nil, err
	}

	_, err = bucket.Create(token, serialized)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSession checks the validity of a session token
func ValidateSession(ctx *h.RequestContext) (*User, error) {
	sessionTokenCookie, err := ctx.Request.Cookie("session_id")

	if err != nil {
		return nil, errors.New("authorization token not provided")
	}

	sessionToken := sessionTokenCookie.Value

	if sessionToken == "" {
		return nil, errors.New("authorization token not provided")
	}

	client := service.Get[KvClient](ctx.ServiceLocator())

	// Retrieve the session from the "sessions" bucket
	bucket, err := client.GetOrCreateBucket(&nats.KeyValueConfig{
		Bucket: "sessions",
	})

	if err != nil {
		return nil, err
	}

	raw, err := bucket.Get(sessionToken)

	if errors.Is(err, nats.ErrKeyNotFound) {
		return nil, errors.New("invalid or expired session token")
	}

	if err != nil {
		return nil, err
	}

	session, err := json2.Deserialize[Session](raw.Value())

	if err != nil {
		return nil, err
	}

	// Check if the session is expired
	if time.Now().After(session.ExpiresAt) {
		_ = bucket.Delete(sessionToken) // Clean up expired session
		return nil, errors.New("session token expired")
	}

	// Retrieve the user associated with the session
	user, err := UserGetByPredicate(ctx.ServiceLocator(), func(u *User) bool {
		return u.Id == session.UserId
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found for session")
	}

	return user, nil
}
