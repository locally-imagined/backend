# Backend Repo for Locally Imagined

## Tools Used:
+ GO
+ GOA V3

## Platforms Used:
+ Heroku
+ Amazon S3
+ PostgreSQL

## Services:
+ Login
+ Signup
+ Postings
+ Users

## Endpoints:
+ Login
+ Signup
+ CreatePost - Creates a post for a user and stores information DB and image content in S3
+ EditPost - Edits a users post
+ DeletePost - Deletes a users post from DB and S3
+ GetPostPage - Returns 25 most recent posts
+ GetArtistPostPage - Returns 25 most recent posts for artist
+ GetPostPageFiltered - Returns 25 most recent posts that match the filter
+ GetImagesForPost - Returns ordered image ID's for a post
+ GetUserInfo - Returns all relevant information about a user
+ UpdateProfilePicture - Updates a users profile picture
+ UpdateBio - Updates a users bio

## To Create a Service:
+ Follow the style of Postings
+ Mount service endpoints on mux in main.go

## To Create an Endpoint:
+ Add to design.go under correct service
+ Implement endpoint in "service".go and client.go
