package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

var _ = API("locallyimagined", func() {
	Title("Locally Imagined backend")
	Description("Serves all frontend requests")
	cors.Origin("http://localhost:3000", func() { // Define CORS policy, may be prefixed with "*" wildcard
		cors.Headers("*")                                       // One or more authorized headers, use "*" to authorize all
		cors.Methods("GET", "POST", "DELETE", "PUT", "OPTIONS") // One or more authorized HTTP methods
		cors.Expose("*")                                        // One or more headers exposed to clients
		cors.MaxAge(600)                                        // How long to cache a preflight request response
		cors.Credentials()                                      // Sets Access-Control-Allow-Credentials header
	})
	cors.Origin("http://localhost", func() { // Define CORS policy, may be prefixed with "*" wildcard
		cors.Headers("*")                                       // One or more authorized headers, use "*" to authorize all
		cors.Methods("GET", "POST", "DELETE", "PUT", "OPTIONS") // One or more authorized HTTP methods
		cors.Expose("*")                                        // One or more headers exposed to clients
		cors.MaxAge(600)                                        // How long to cache a preflight request response
		cors.Credentials()                                      // Sets Access-Control-Allow-Credentials header
	})
	cors.Origin("https://locallyimagined.netlify.app", func() { // Define CORS policy, may be prefixed with "*" wildcard
		cors.Headers("*")                                       // One or more authorized headers, use "*" to authorize all
		cors.Methods("GET", "POST", "DELETE", "PUT", "OPTIONS") // One or more authorized HTTP methods
		cors.Expose("*")                                        // One or more headers exposed to clients
		cors.MaxAge(600)                                        // How long to cache a preflight request response
		cors.Credentials()                                      // Sets Access-Control-Allow-Credentials header
	})
	Error("unauthorized", ErrorResult, "Credentials are invalid")
	Error("internal", ErrorResult, "Internal Error")
	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
		Response("internal", StatusInternalServerError)
	})
})

var _ = Service("login", func() {
	Error("unauthorized", ErrorResult, "Credentials are invalid")
	Method("Login", func() {
		Security(LoginBasicAuth)
		Payload(func() {
			Username("username", String, "Raw username")
			Password("password", String, "User password")
			Required("username", "password")
		})
		Result(func() {
			Attribute("LoginResponse", LoginResponse, "JWT and UserID")
		})
		HTTP(func() {
			POST("/login")
			Response(func() {
				Body("LoginResponse")
			})
		})
	})
})

var _ = Service("signup", func() {
	Error("unauthorized", ErrorResult, "Credentials are invalid")
	Method("Signup", func() {
		Security(SignupBasicAuth)
		Payload(func() {
			Username("username", String, "Raw username")
			Password("password", String, "User password")
			Attribute("user", NewUser)
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
	Error("unauthorized", ErrorResult, "Credentials are invalid")
	Error("internal", ErrorResult, "Internal Error")
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
			POST("/posts/create")
			Body("post")
			Response(func() {
				Body("Posted")
			})
		})
	})
	Method("delete_post", func() {
		Security(JWTAuth)
		Payload(func() {
			Token("token", String, "jwt used for auth")
			Attribute("postID", String, "Post to delete")
			Required("token", "postID")
		})
		HTTP(func() {
			DELETE("/posts/delete/{postID}")
		})
	})
	Method("edit_post", func() {
		Security(JWTAuth)
		Payload(func() {
			Token("token", String, "jwt used for auth")
			Attribute("postID", String, "Post ID")
			Attribute("title", String, "Post title")
			Attribute("description", String, "Post description")
			Attribute("price", String, "Post price")
			Attribute("content", Content, "Image content")
			Attribute("medium", String, "Art type")
			Attribute("sold", Boolean, "is sold")
			Attribute("deliverytype", String, "Delivery type")
			Attribute("imageID", String, "Image ID")
			Required("token", "postID")
		})
		Result(func() {
			Attribute("Posted", PostResponse)
		})
		// take these out of param
		HTTP(func() {
			PUT("/posts/edit/{postID}")
			Param("title")
			Param("description")
			Param("price")
			Param("medium")
			Param("sold")
			Param("deliverytype")
			Param("imageID")
			Body("content")
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
			GET("/posts/getpage/{page}")
			Response(func() {
				Body("Posts")
			})
		})
	})
	Method("get_artist_post_page", func() {
		Payload(func() {
			Attribute("userID", String, "User ID to get posts for")
			Attribute("page", Int, "Page to get posts for")
			Required("userID", "page")
		})
		Result(func() {
			Attribute("Posts", ArrayOf(PostResponse))
		})
		HTTP(func() {
			GET("/posts/artistposts/{page}")
			Param("userID")
			Response(func() {
				Body("Posts")
			})
		})
	})
	Method("get_post_page_filtered", func() {
		Payload(func() {
			Attribute("page", Int, "Page to get posts for")
			Attribute("keyword", String, "Search bar keyword to search for in title and description")
			Attribute("startDate", String, "Filter attribute to see posts after given date")
			Attribute("endDate", String, "Filter attribute to see posts before given date")
			Attribute("medium", String, "Filter attribute to see posts with given medium type")
			Required("page")
		})
		Result(func() {
			Attribute("Posts", ArrayOf(PostResponse))
		})
		HTTP(func() {
			GET("/posts/getpagefiltered/{page}")
			Param("keyword")
			Param("startDate")
			Param("endDate")
			Param("medium")
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
			GET("/posts/getimages/{postID}")
			Response(func() {
				Body("Images")
			})
		})
	})
})

var _ = Service("users", func() {
	Error("unauthorized", String, "Credentials are invalid")
	Method("update_bio", func() {
		Security(JWTAuth)
		Payload(func() {
			Token("token", String, "jwt used for auth")
			Attribute("bio", Bio, "New bio to be added")
			Required("token", "bio")
		})
		Result(func() {
			Attribute("user", User, "Updated user")
		})
		HTTP(func() {
			PUT("/users/updatebio")
			Body("bio")
			Response(func() {
				Body("user")
			})
		})
	})
	Method("update_profile_picture", func() {
		Security(JWTAuth)
		Payload(func() {
			Token("token", String, "jwt used for auth")
			Attribute("content", Content, "New Profile Photo")
			Required("token", "content")
		})
		Result(func() {
			Attribute("imageID", ProfilePhoto)
		})
		HTTP(func() {
			PUT("/users/updateprofilepic")
			Body("content")
			Response(func() {
				Body("imageID")
			})
		})
	})
	Method("get_user_info", func() {
		Payload(func() {
			Attribute("userID", String, "userid of user whose info to retrieve")
			Required("userID")
		})
		Result(func() {
			Attribute("user", User)
		})
		HTTP(func() {
			GET("/users/info")
			Param("userID")
			Response(func() {
				Body("user")
			})
		})
	})
})

// BasicAuth defines a security scheme using basic authentication. The scheme
// protects the "login" action used to create JWTs.
var LoginBasicAuth = BasicAuthSecurity("login", func() {
	Description("Basic authentication used to authenticate security principal during login")
})

// BasicAuth defines a security scheme using basic authentication. The scheme
// protects the "signup" action used to create JWTs.
var SignupBasicAuth = BasicAuthSecurity("signup", func() {
	Description("Basic authentication used to authenticate security principal during signup")
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	Description(`Secures endpoint by requiring a valid JWT token retrieved via the login endpoint.`)
})

var LoginResponse = Type("LoginResponse", func() {
	Description("Response from logging in")
	Attribute("jwt", String, "jwt used for future authentication")
	Attribute("userID", String, "users ID")
})

var Content = Type("Content", func() {
	Description("Image Content")
	Attribute("content", String, "raw image content")
})

var Bio = Type("Bio", func() {
	Description("Updated Bio")
	Attribute("bio", String, "New Bio")
})

var NewUser = Type("NewUser", func() {
	Description("Describes a user at signup")
	Attribute("firstName", String, "First name")
	Attribute("lastName", String, "Last name")
	Attribute("phone", String, "Phone number")
	Attribute("email", String, "Email")
	Required("firstName", "lastName", "phone", "email")
})

var User = Type("User", func() {
	Description("Describes a user")
	Attribute("firstName", String, "First name")
	Attribute("lastName", String, "Last name")
	Attribute("phone", String, "Phone number")
	Attribute("email", String, "Email")
	Attribute("bio", String, "Bio")
	Attribute("profpicID", String, "Prof Pic UUID")
	Required("firstName", "lastName", "phone", "email", "bio", "profpicID")
})

var ProfilePhoto = Type("Profile Photo", func() {
	Description("Profile Photo uuid")
	Attribute("imageID", String, "Image ID")
})

// we probably dont need this, change createpost to return postresponse
var Post = Type("Post", func() {
	Description("Describes a post payload")
	Attribute("title", String, "Post title")
	Attribute("description", String, "Post description")
	Attribute("price", String, "Post price")
	Attribute("content", ArrayOf(String), "Post content")
	Attribute("medium", String, "Art type")
	Attribute("deliverytype", String, "Delivery type")
	Required("title", "description", "price", "content", "medium", "deliverytype")
})

var PostResponse = Type("PostResponse", func() {
	Description("Describes a post response")
	Attribute("title", String, "Post title")
	Attribute("description", String, "Post description")
	Attribute("price", String, "Post price")
	Attribute("imageIDs", ArrayOf(String), "Image ID")
	Attribute("postID", String, "Post ID")
	Attribute("medium", String, "Art type")
	Attribute("uploadDate", String, "Upload Date")
	Attribute("sold", Boolean, "is sold")
	Attribute("deliverytype", String, "Delivery type")
	Attribute("userID", String, "User id associated with post")
	Attribute("profpicID", String, "prof pic id")
	Attribute("username", String, "Username associated with post")
	Required("title", "description", "price", "imageIDs", "postID", "medium", "uploadDate", "sold", "deliverytype", "username", "userID", "profpicID")
})
