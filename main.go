package main

import (
	loginServer "backend/gen/http/login/server"
	postingsServer "backend/gen/http/postings/server"
	signupServer "backend/gen/http/signup/server"
	genlogin "backend/gen/login"
	genpostings "backend/gen/postings"
	gensignup "backend/gen/signup"
	"net/http"
	"os"

	login "backend/login"
	postings "backend/postings"
	signup "backend/signup"

	goahttp "goa.design/goa/v3/http"
)

func main() {
	port := os.Getenv("PORT")
	sL := &login.Service{}                      //# Create Service
	loginEndpoints := genlogin.NewEndpoints(sL) // # Create endpoints
	sP := &postings.Service{}
	postingsEndpoints := genpostings.NewEndpoints(sP)
	sS := &signup.Service{}
	signupEndpoints := gensignup.NewEndpoints(sS)

	mux := goahttp.NewMuxer()      //# Create HTTP muxer
	dec := goahttp.RequestDecoder  //# Set HTTP request decoder
	enc := goahttp.ResponseEncoder // # Set HTTP response encoder

	loginSvr := loginServer.New(loginEndpoints, mux, dec, enc, nil, nil) // # Create Goa HTTP server
	loginServer.Mount(mux, loginSvr)                                     //# Mount Goa server on mux

	postingsSvr := postingsServer.New(postingsEndpoints, mux, dec, enc, nil, nil) // # Create Goa HTTP server
	postingsServer.Mount(mux, postingsSvr)                                        //# Mount Goa server on mux

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
