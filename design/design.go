package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("calc", func() {
	Method("Multiply", func() {
		// Security(JWTAuth, func() {
		// 	Scope("multiply")
		// })
		Payload(func() {
			Attribute("a", Int, "First operand")
			Attribute("b", Int, "Second operand")
		})
		Result(Int)
		HTTP(func() {
			GET("/multiply/{a}/{b}")
		})
	})
	Method("Add", func() {
		Payload(func() {
			Attribute("a", Int, "First operand")
			Attribute("b", Int, "Second operand")
		})
		Result(Int)
		HTTP(func() {
			GET("/add/{a}/{b}")
		})
	})
})

// var JWTAuth = JWTSecurity("jwt", func() {
// 	Description(`Secures endpoint by requiring a valid JWT token.`)
// 	Scope("multiply", "Enable building shutdown")
// })

// func PayloadWithToken(dsl func()) {
// 	Payload(func() {
// 		Token("token", String, func() {
// 			Description("JWT used for authentication")
// 			Pattern(`^Bearer [A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`)
// 		})
// 		Required("token")
// 		dsl()
// 	})
// }
