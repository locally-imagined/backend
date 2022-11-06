package design

import (
	. "goa.design/goa/v3/dsl"
)

// BasicAuth defines a security scheme using basic authentication. The scheme
// protects the "signin" action used to create JWTs.
var BasicAuth = BasicAuthSecurity("basic", func() {
	Description("Basic authentication used to authenticate security principal during signin")
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	Description(`Secures endpoint by requiring a valid JWT token retrieved via the signin endpoint.`)
})

var _ = Service("auth", func() {
	Error("unauthorized", String, "Credentials are invalid")
	Method("Login", func() {
		Security(BasicAuth)
		Payload(func() {
			Username("username", String, "Raw username")
			Password("password", String, "User password")
			Required("username", "password")
		})
		Result(func() {
			Attribute("jwt", String)
			Attribute("Access-Control-Allow-Origin")
		})
		HTTP(func() {
			POST("/login")
			Response(func() {
				Header("Access-Control-Allow-Origin")
				Body("jwt")
			})
		})
	})
	Method("Signup", func() {
		Payload(func() {
			Attribute("username", String, "Raw username")
			Attribute("password", String, "User password")
		})
		Result(func() {
			Attribute("jwt", String)
			Attribute("Access-Control-Allow-Origin")
		})
		HTTP(func() {
			GET("/signup/{username}/{password}")
			Response(func() {
				Header("Access-Control-Allow-Origin")
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
			Attribute("success", Boolean)
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
