package design

import (
	. "goa.design/goa/v3/dsl"
)

// BasicAuth defines a security scheme using basic authentication. The scheme
// protects the "signin" action used to create JWTs.
var BasicAuth = BasicAuthSecurity("basic", func() {
	Description("Basic authentication used to authenticate security principal during signin")
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
	Method("upload_photo", func() {
		Payload(func() {
			Attribute("Authorization")
			Attribute("content", Bytes, "photo content")
			Required("Authorization")
		})
		Result(func() {
			Attribute("success", Boolean)
			Attribute("Access-Control-Allow-Origin")
		})
		HTTP(func() {
			Header("Authorization")
			GET("/upload/{content}")
			Response(func() {
				Header("Access-Control-Allow-Origin")
				Body("success")
			})
		})
	})
})

// var JWTAuth = JWTSecurity("jwt", func() {
// 	Description(`Secures endpoint by requiring a valid JWT token retrieved via the login service.`)
// })
