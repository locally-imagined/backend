package main

import (
	"net/http"

	goahttp "goa.design/goa/v3/http"

	"github.com/locally-imagined/goa/gen/calc"
	"github.com/locally-imagined/goa/gen/http/calc/server"
	helpers "github.com/locally-imagined/goa/service"
)

func main() {
	//port := os.Getenv("PORT")
	s := &helpers.Service{}                               //# Create Service
	endpoints := calc.NewEndpoints(s)                     // # Create endpoints
	mux := goahttp.NewMuxer()                             //# Create HTTP muxer
	dec := goahttp.RequestDecoder                         //# Set HTTP request decoder
	enc := goahttp.ResponseEncoder                        // # Set HTTP response encoder
	svr := server.New(endpoints, mux, dec, enc, nil, nil) // # Create Goa HTTP server
	server.Mount(mux, svr)                                //# Mount Goa server on mux
	httpsvr := &http.Server{                              // # Create Go HTTP server
		//Addr:    ":" + port, // # Configure server address (this is for heroku)
		Addr:    "localhost:8080", // this is for localhost obviously
		Handler: mux,              // # Set request handler
	}
	if err := httpsvr.ListenAndServe(); err != nil { // # Start HTTP server
		panic(err)
	}
}
