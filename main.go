package main

import (
	loginServer "backend/gen/http/login/server"
	signupServer "backend/gen/http/signup/server"
	uploadServer "backend/gen/http/upload/server"
	genlogin "backend/gen/login"
	gensignup "backend/gen/signup"
	genupload "backend/gen/upload"
	"net/http"
	"os"

	login "backend/login"
	signup "backend/signup"
	upload "backend/upload"

	goahttp "goa.design/goa/v3/http"
)

// func exitErrorf(msg string, args ...interface{}) {
// 	fmt.Fprintf(os.Stderr, msg+"\n", args...)
// 	os.Exit(1)
// }

// func aws3() {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(os.Getenv("BUCKETEER_AWS_REGION"))},
// 	)

// 	// Create S3 service client
// 	svc := s3.New(sess)
// 	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(os.Getenv("BUCKETEER_BUCKET_NAME"))})
// 	if err != nil {
// 		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
// 	}

// 	for _, item := range resp.Contents {
// 		fmt.Println("Name:         ", *item.Key)
// 		fmt.Println("Last modified:", *item.LastModified)
// 		fmt.Println("Size:         ", *item.Size)
// 		fmt.Println("Storage class:", *item.StorageClass)
// 		fmt.Println("")
// 	}
// }

func main() {
	// aws3()
	port := os.Getenv("PORT")

	sL := &login.Service{}                      //# Create Service
	loginEndpoints := genlogin.NewEndpoints(sL) // # Create endpoints
	sU := &upload.Service{}
	uploadEndpoints := genupload.NewEndpoints(sU)
	sS := &signup.Service{}
	signupEndpoints := gensignup.NewEndpoints(sS)

	mux := goahttp.NewMuxer()      //# Create HTTP muxer
	dec := goahttp.RequestDecoder  //# Set HTTP request decoder
	enc := goahttp.ResponseEncoder // # Set HTTP response encoder

	loginSvr := loginServer.New(loginEndpoints, mux, dec, enc, nil, nil) // # Create Goa HTTP server
	loginServer.Mount(mux, loginSvr)                                     //# Mount Goa server on mux

	uploadSvr := uploadServer.New(uploadEndpoints, mux, dec, enc, nil, nil) // # Create Goa HTTP server
	uploadServer.Mount(mux, uploadSvr)                                      //# Mount Goa server on mux

	signupSvr := signupServer.New(signupEndpoints, mux, dec, enc, nil, nil) // # Create Goa HTTP server
	signupServer.Mount(mux, signupSvr)                                      //# Mount Goa server on mux

	httpsvr := &http.Server{ // # Create Go HTTP server
		Addr: ":" + port, // # Configure server address (this is for heroku)
		//Addr:    "localhost:8080", // this is for localhost obviously
		Handler: mux, // # Set request handler
	}
	if err := httpsvr.ListenAndServe(); err != nil { // # Start HTTP server
		panic(err)
	}
}
