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
		Result(String)
		HTTP(func() {
			GET("/login/{username}/{password}")
		})
	})
	Method("Signup", func() {
		Payload(func() {
			Attribute("username", String, "Raw username")
			Attribute("password", String, "User password")
		})
		Result(String)
		HTTP(func() {
			GET("/signup/{username}/{password}")
		})
	})
})

// var JWTAuth = JWTSecurity("jwt", func() {
// 	Description(`Secures endpoint by requiring a valid JWT token retrieved via the login service.`)
// })
