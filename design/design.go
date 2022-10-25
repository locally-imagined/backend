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
	Method("Add", func() {
		Payload(func() {
			Token("token", String, "JWT token used for auth")
			Attribute("a", Int, "First operand")
			Attribute("b", Int, "Second operand")
		})
		Result(Int)
		HTTP(func() {
			GET("/add/{a}/{b}")
		})
	})
})
