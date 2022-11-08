package design

import (
	. "goa.design/goa/v3/dsl"
)

// BasicAuth defines a security scheme using basic authentication. The scheme
// protects the "login" action used to create JWTs.
var LoginBasicAuth = BasicAuthSecurity("login", func() {
	Description("Basic authentication used to authenticate security principal during signin")
})

// BasicAuth defines a security scheme using basic authentication. The scheme
// protects the "signup" action used to create JWTs.
var SignupBasicAuth = BasicAuthSecurity("signup", func() {
	Description("Basic authentication used to authenticate security principal during signin")
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	Description(`Secures endpoint by requiring a valid JWT token retrieved via the signin endpoint.`)
})

var _ = Service("login", func() {
	Error("unauthorized", String, "Credentials are invalid")
	Method("Login", func() {
		Security(LoginBasicAuth)
		Payload(func() {
			Username("username", String, "Raw username")
			Password("password", String, "User password")
			Required("username", "password")
		})
		Result(func() {
			Attribute("jwt", String)
			Attribute("Access-Control-Allow-Headers")
			Attribute("Access-Control-Allow-Methods")
			Attribute("Access-Control-Allow-Origin")
			Attribute("Access-Control-Allow-Credentials")
		})
		HTTP(func() {
			POST("/login")
			Response(func() {
				Header("Access-Control-Allow-Headers")
				Header("Access-Control-Allow-Methods")
				Header("Access-Control-Allow-Origin")
				Header("Access-Control-Allow-Credentials")
				Body("jwt")
			})
		})
	})
})

var _ = Service("signup", func() {
	Method("Signup", func() {
		Security(SignupBasicAuth)
		Payload(func() {
			Username("username", String, "Raw username")
			Password("password", String, "User password")
			Required("username", "password")
		})
		Result(func() {
			Attribute("jwt", String)
			Attribute("Access-Control-Allow-Origin")
			Attribute("Access-Control-Allow-Credentials")
		})
		HTTP(func() {
			POST("/signup")
			Response(func() {
				Header("Access-Control-Allow-Origin")
				Header("Access-Control-Allow-Credentials")
				Body("jwt")
			})
		})
	})
})

var _ = Service("upload", func() {
	Error("unauthorized", String, "Credentials are invalid")
	Method("upload_photo", func() {
		Security(JWTAuth)
		Payload(func() {
			Token("token", String, "jwt used for auth")
			Attribute("content", Bytes, "photo content")
			Required("token")
		})
		Result(func() {
			Attribute("success", String)
			Attribute("Access-Control-Allow-Origin")
		})
		HTTP(func() {
			GET("/upload/{content}")
			Response(func() {
				Header("Access-Control-Allow-Origin")
				Body("success")
			})
		})
	})
})
