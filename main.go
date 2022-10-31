package main

import (
	"backend/gen/auth"
	"backend/gen/http/auth/server"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	login "backend/auth"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	goahttp "goa.design/goa/v3/http"
)

type locallyImaginedClaims struct {
	Name  string
	Email string
	jwt.RegisteredClaims
}

func MakeToken(email string) (string, error) {
	uuid := uuid.New()
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       uuid.String()}

	payload := &locallyImaginedClaims{
		Name:             email,
		Email:            email,
		RegisteredClaims: claims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedJWT, err := token.SignedString([]byte("test"))
	if err != nil {
		return "", fmt.Errorf("unable to sign token with zendesk secret: %v", err)
	}
	return signedJWT, nil
}

func DecodeToken(tokenString string) *jwt.Token {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("test"), nil
	})
	// ... error handling
	if err != nil {
		log.Println("NOOIOIO")
	}
	// do something with decoded claims
	log.Println(token)
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
	}
	return token
}

func main() {
	port := os.Getenv("PORT")
	s := &login.Service{}                                 //# Create Service
	endpoints := auth.NewEndpoints(s)                     // # Create endpoints
	mux := goahttp.NewMuxer()                             //# Create HTTP muxer
	dec := goahttp.RequestDecoder                         //# Set HTTP request decoder
	enc := goahttp.ResponseEncoder                        // # Set HTTP response encoder
	svr := server.New(endpoints, mux, dec, enc, nil, nil) // # Create Goa HTTP server
	server.Mount(mux, svr)                                //# Mount Goa server on mux
	httpsvr := &http.Server{                              // # Create Go HTTP server
		Addr: ":" + port, // # Configure server address (this is for heroku)
		//Addr:    "localhost:8080", // this is for localhost obviously
		Handler: mux, // # Set request handler
	}
	if err := httpsvr.ListenAndServe(); err != nil { // # Start HTTP server
		panic(err)
	}
}
