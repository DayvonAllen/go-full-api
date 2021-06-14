## Dependencies
- MongoDB
- Apache Zookeeper/ Apache Kafka
- Go programming language
---

## Setup
1. In the setup file there is a docker compose file.
2. `cd` into that setup file and then run:
    - Linux: `sudo docker-compose up` or `sudo docker-compose up -d`(for a detached start up)
    - Windows: `docker-compose up` or `docker-compose up -d`(for a detached start up)
3. Start MongoDB on port 27017
    - If using docker run the following command:
      - Linux: `sudo docker run -p 27017:27017 mongo`
      - Windows: `docker run -p 27017:27017 mongo`
4. Run the app:
   - Execute this command from the root folder in the terminal `go run main.go`
   - Root Folder is `go-full-api`
---

## Routes
- Get All users:
  - `GET:http://localhost:8080/users`(protected, needs token)
- Login:
  - `POST:http://localhost:8080/auth/login`
  - JSON: `{
    "email": "jdoedddd25455@gmail.com",
    "password": "password"
}`
- Register:
  - `POST:http://localhost:8080/users`
  - JSON: `{
        "username": "jdoe1744",
        "email": "jdoedddd25455@gmail.com",
        "password": "password"
}`
- Reset Password Query:
  - `POST:http://localhost:8080/auth/reset`
  - JSON: `{
    "email": "jdoedddd25455@gmail.com"
}`
- Reset Password:
  - `PUT:http://localhost:8080/auth/reset/<token should be in the console of the go app, place here>`
- Verify Account:
  - `PUT:http://localhost:8080/auth/account/<Token is in MongoDB user collection, place here>`
- Get User's account: (protected, needs token)
  - `GET:http://localhost:8080/users/account`
- Flag user: (protected, needs token):
  - `POST:http://localhost:8080/users/flag/<username of person to flag>`
- Update profile visibility: (protected, needs token)
  - `PUT:http://localhost:8080/users/profile-visibility`
  - JSON: `{
    "profileIsViewable": false
}`
- Update Message Acceptance(whether you want to receive messages or not): (protected, needs token)
  - `PUT:http://localhost:8080/users/message-acceptance`
  - JSON: `{
    "acceptMessages": false
}`
- Block user: (protected, needs token)
  - `PUT:http://localhost:8080/users/block/<username of user you want to block>`
- Get all blocked users: (protected, needs token)   
  - `GET:http://localhost:8080/users/blocked`
- Unblock user: (protected, needs token)
  - `PUT:http://localhost:8080/users/unblock/<username of user you want to unblock>`
- Delete current user account: (protected, needs token)
  - `DELETE:http://localhost:8080/users/delete`