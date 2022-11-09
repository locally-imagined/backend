package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("auth", func() {
	Method("Login", func() {
		Payload(func() {
			Attribute("username", String, "Raw username")
			Attribute("password", String, "User password")
		})
		Result(func() {
			Attribute("jwt", String)
			Attribute("Access-Control-Allow-Origin")
		})
		HTTP(func() {
			GET("/login/{username}/{password}")
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
			POST("/upload/{content}")
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
