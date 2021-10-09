# Appointy-Shortlist-Assignment
Instagram API using GoLang and MongoDB


Clone the repository in your preferred directory
open your terminal

# Run the server on a separate terminal: 
```bash
go run main.go
```

# To check the running server through your browser:
```bash
curl localhost:8081
```
# To create a user:
```bash
curl localhost:8081/users -X POST -d '{"name": "sample-name", "email": "sample-email", "password": "sample-password"}'
```
The password posted by the user will be encrypted and stored on the mongodb database for concrete security

# To get user details by id:
```bash
curl localhost:8081/users/?id=<user_id>
```
# To create a post:
```bash
curl localhost:8081/posts -X POST -d '{"postowner": <id-of-user>, "caption": "sample-caption", "imageURL": "http://www.sample.com", "timestamp": ""}'
```
The timestamp field is initialized using the time.Now() library function in goLang

# To get post details by id:
```bash
curl localhost:8081/posts/?id=<post_id>
```
# List all posts of a user:
```bash
curl localhost:8081/posts/users/?id=<user_id>
```
Note: The id is passed as a parameter in all the GET Requests

