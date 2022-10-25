package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("auth", func() {
	Method("Login", func() {
		Payload(func() {
			Attribute("username", String, "Raw username")
			Attribute("password", String, "Hashed user password")
		})
		Result(String)
		HTTP(func() {
			GET("/login/{username}/{password}")
		})
	})
	// Method("Add", func() {
	// 	Security(JWTAuth)
	// 	Payload(func() {
	// 		Token("token", String, "JWT token used for auth")
	// 		Attribute("a", Int, "First operand")
	// 		Attribute("b", Int, "Second operand")
	// 	})
	// 	Result(Int)
	// 	HTTP(func() {
	// 		GET("/add/{a}/{b}")
	// 	})
	// })
})

// var JWTAuth = JWTSecurity("jwt", func() {
// 	Description(`Secures endpoint by requiring a valid JWT token retrieved via the login service.`)
// })
