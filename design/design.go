package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
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

var User = Type("User", func() {
	Description("Describes a user")
	Attribute("firstName", String, "First name")
	Attribute("lastName", String, "Last name")
	Attribute("phone", String, "Phone number")
	Attribute("email", String, "Email")
	Required("firstName", "lastName", "phone", "email")
})

var Post = Type("Post", func() {
	Description("Describes a post")
	Attribute("title", String, "Post title")
	Attribute("description", String, "Post description")
	Attribute("price", String, "Post price")
	Attribute("content", Bytes, "Post content")
	Required("title", "description", "price", "content")
})

var PostResponse = Type("PostResponse", func() {
	Description("Describes a post")
	Attribute("title", String, "Post title")
	Attribute("description", String, "Post description")
	Attribute("price", String, "Post price")
	Attribute("imageID", String, "Image ID")
	Attribute("postID", String, "Post ID")
	Attribute("uploadDate", String, "Upload Date")
	Required("title", "description", "price", "imageID", "postID", "uploadDate")
})

var _ = Service("login", func() {
	cors.Origin("http://localhost:3000", func() { // Define CORS policy, may be prefixed with "*" wildcard
		cors.Headers("*")                      // One or more authorized headers, use "*" to authorize all
		cors.Methods("GET", "POST", "OPTIONS") // One or more authorized HTTP methods
		cors.Expose("*")                       // One or more headers exposed to clients
		cors.MaxAge(600)                       // How long to cache a preflight request response
		cors.Credentials()                     // Sets Access-Control-Allow-Credentials header
	})
	cors.Origin("http://localhost:3001", func() { // Define CORS policy, may be prefixed with "*" wildcard
		cors.Headers("*")                      // One or more authorized headers, use "*" to authorize all
		cors.Methods("GET", "POST", "OPTIONS") // One or more authorized HTTP methods
		cors.Expose("*")                       // One or more headers exposed to clients
		cors.MaxAge(600)                       // How long to cache a preflight request response
		cors.Credentials()                     // Sets Access-Control-Allow-Credentials header
	})
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
		})
		HTTP(func() {
			POST("/login")
			Response(func() {
				Body("jwt")
			})
		})
	})
})

var _ = Service("signup", func() {
	cors.Origin("http://localhost:3000", func() { // Define CORS policy, may be prefixed with "*" wildcard
		cors.Headers("*")                      // One or more authorized headers, use "*" to authorize all
		cors.Methods("GET", "POST", "OPTIONS") // One or more authorized HTTP methods
		cors.Expose("*")                       // One or more headers exposed to clients
		cors.MaxAge(600)                       // How long to cache a preflight request response
		cors.Credentials()                     // Sets Access-Control-Allow-Credentials header
	})
	cors.Origin("http://localhost:3001", func() { // Define CORS policy, may be prefixed with "*" wildcard
		cors.Headers("*")                      // One or more authorized headers, use "*" to authorize all
		cors.Methods("GET", "POST", "OPTIONS") // One or more authorized HTTP methods
		cors.Expose("*")                       // One or more headers exposed to clients
		cors.MaxAge(600)                       // How long to cache a preflight request response
		cors.Credentials()                     // Sets Access-Control-Allow-Credentials header
	})
	Method("Signup", func() {
		Security(SignupBasicAuth)
		Payload(func() {
			Username("username", String, "Raw username")
			Password("password", String, "User password")
			Attribute("user", User)
			Required("username", "password", "user")
		})
		Result(func() {
			Attribute("jwt", String)
		})
		HTTP(func() {
			POST("/signup")
			Body("user")
			Response(func() {
				Body("jwt")
			})
		})
	})
})

var _ = Service("postings", func() {
	Error("unauthorized", String, "Credentials are invalid")
	Method("create_post", func() {
		Security(JWTAuth)
		Payload(func() {
			Token("token", String, "jwt used for auth")
			Attribute("post", Post, "Post info")
			Required("token", "post")
		})
		Result(func() {
			Attribute("Posted", PostResponse)
		})
		HTTP(func() {
			POST("/create")
			Body("post")
			Response(func() {
				Body("Posted")
			})
		})
	})
	Method("get_post_page", func() {
		Payload(func() {
			Attribute("page", Int, "Page to get posts for")
			Required("page")
		})
		Result(func() {
			Attribute("Posts", ArrayOf(PostResponse))
		})
		HTTP(func() {
			GET("/posts/{page}")
			Response(func() {
				Body("Posts")
			})
		})
	})
	Method("get_images_for_post", func() {
		Payload(func() {
			Attribute("postID", String, "Post to get images for")
			Required("postID")
		})
		Result(func() {
			Attribute("Images", ArrayOf(String))
		})
		HTTP(func() {
			GET("/posts/{postID}")
			Response(func() {
				Body("Images")
			})
		})
	})
})
