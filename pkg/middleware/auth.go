package middleware

import (
	"context"
	"encoding/json"
	"github.com/TinyRogue/lembook-serv/internal/models"
	"github.com/TinyRogue/lembook-serv/pkg/jwt"
	"log"
	"net/http"
)

type ContextKey string

const ContextUserKey ContextKey = "user"
const ContextReqWriterKey ContextKey = "reqWriter"
const ContextJWT ContextKey = "jwt"

type JWTHandler struct {
	jwt string
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("auth")
		if cookie != nil {
			ctx := context.WithValue(r.Context(), ContextJWT, &JWTHandler{jwt: cookie.Value})
			r = r.WithContext(ctx)
		}

		ctx := context.WithValue(r.Context(), ContextReqWriterKey, &w)
		r = r.WithContext(ctx)

		// Unauthenticated
		if cookie == nil {
			next.ServeHTTP(w, r)
			return
		}

		uid, err := jwt.ParseToken(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := models.FindUserBy(r.Context(), "uid", uid)
		// Not exists
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		marshalledUser, err := json.Marshal(user)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userCtx := context.WithValue(r.Context(), ContextUserKey, &marshalledUser)
		r = r.WithContext(userCtx)

		next.ServeHTTP(w, r)
	})
}

func FindUserByCtx(ctx context.Context) *models.User {
	raw, _ := ctx.Value(ContextUserKey).([]byte)
	var user models.User
	err := json.Unmarshal(raw, &user)
	if err != nil {
		log.Println("Couldn't unmarshall user due to ", err)
		return nil
	}
	return &user
}

func GetResWriter(ctx context.Context) *http.ResponseWriter {
	resWriter, ok := ctx.Value(ContextReqWriterKey).(*http.ResponseWriter)
	if !ok {
		log.Fatalln("FATALITY")
	}
	return resWriter
}
