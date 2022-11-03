package main

import (
	"backend/gen/auth"
	authServer "backend/gen/http/auth/server"
	uploadServer "backend/gen/http/upload/server"
	"backend/gen/upload"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	login "backend/auth"
	uploads "backend/upload"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func aws3() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("BUCKETEER_BUCKET_NAME")),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
}

func main() {
	// aws3()
	port := os.Getenv("PORT")
	sL := &login.Service{}                 //# Create Service
	authEndpoints := auth.NewEndpoints(sL) // # Create endpoints
	sU := &uploads.Service{}
	upEndpoints := upload.NewEndpoints(sU)
	mux := goahttp.NewMuxer()                                           //# Create HTTP muxer
	dec := goahttp.RequestDecoder                                       //# Set HTTP request decoder
	enc := goahttp.ResponseEncoder                                      // # Set HTTP response encoder
	authSvr := authServer.New(authEndpoints, mux, dec, enc, nil, nil)   // # Create Goa HTTP server
	authServer.Mount(mux, authSvr)                                      //# Mount Goa server on mux
	uploadSvr := uploadServer.New(upEndpoints, mux, dec, enc, nil, nil) // # Create Goa HTTP server
	uploadServer.Mount(mux, uploadSvr)                                  //# Mount Goa server on mux
	httpsvr := &http.Server{                                            // # Create Go HTTP server
		Addr: ":" + port, // # Configure server address (this is for heroku)
		//Addr:    "localhost:8080", // this is for localhost obviously
		Handler: mux, // # Set request handler
	}
	if err := httpsvr.ListenAndServe(); err != nil { // # Start HTTP server
		panic(err)
	}
}
