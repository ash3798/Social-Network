# Social-Network
Social Network is a Golang based backend for the app where registered users can have wall of comments , subcomments and reactions

### Authentication mechanism
Authentication is done using JWT(JSON WEB-TOKEN). Only the registered users with valid token can access the API's otherwise content will not be accessible.

### Database Used 
* Postgres

## Prerequisite
* You should have Docker installed on your system. (Docker desktop for windows)

## Getting started

### Setup the Database for your application

1. Pull the postgres image
```shell
docker pull postgres
```
2. Run the docker container with this image 
```shell
docker run --name postgres-db -p 5432:5432 -e POSTGRES_PASSWORD=<password> -d postgres
```
> - Here POSTGRES_PASSWORD is environment variable, this is the password for the "postgres" user created by default
> - Postgres is running on port 5432 inside container, port mapping is done here to allow access from outside the container also.

3. Once container is up and running , you can use database visualizer tool like "DBeaver" to check if postgres is accessible.
> - Hostname will be the IP of your machine ,  Port will same as mentioned in docker run command
> - Username will be "postgres" if you are using default user.

### Setup the application

1. Clone the repository into your system.
```
git clone https://github.com/ash3798/Social-Network.git
```
2. Navigate to Docker folder inside the repository .Here Dockerfile will be there which we will use to create container
3. Before building the docker file , clone the repository again inside Docker folder , using the same command as above.
> This workaround is needed to avoid entering the github access token in init.sh file which will be required to clone the repository inside container. Here in our case repository will get copied from machine while building image .
4. After cloning repository in Docker folder , run the docker build command using Dockerfile present.
```
docker build . -t socialnetwork:latest
```
5. You can check the created image 
```
docker images

REPOSITORY      TAG       IMAGE ID       CREATED         SIZE
socialnetwork   latest    a209a80aa593   7 seconds ago   880MB
```
6. Run the container using the image-id or name .
```bash
docker run --name Social-Network-app -e HOSTNAME=<host-ip> -e DATABASE_PORT=<db-port> -e DATABASE_PASSWORD=<password> -p 9999:9999 socialnetwork
```
> - HOSTNAME is the IP of the machine where postgres is running.
> - Application runs on the port 9999 by default. It can be modified using "APP_PORT" environment variable , but then port also should be changed accordingly in -p argument of docker run command.
> - Database username and Database name have been set to "postgres" by default in application
> - If you are using custom user and custom database , you can set username with "DATABASE_USERNAME" and database name with "DATABASE_NAME" environment variables.
> - Additional environment variables that can be tweaked have been listed below.


#### Note : If you see error like this on running docker run 
```
standard_init_linux.go:228: exec user process caused: no such file or directory
```
> - Open init.sh in Notepad++ -> Edit -> EOL Conversions -> select Unix(LF)
> - Save the file and build the docker image again
> - Docker Run this newly created image now

7. Once container is up and running , application will be accessible on port number used in docker run command.
8. Docker logs of container can be used to view the application logs.
```shell
docker logs -f Social-Network-app
```

### Environment Variables
* APP_PORT : Used to set port on which application runs
* HOSTNAME : Used to set the hostname of machine 
* DATABASE_PORT : Used to set port on which database is accessible
* DATABASE_USERNAME : Used to set username of database. Default : "postgres"
* DATABASE_PASSWORD : Used to set password of database
* DATABASE_NAME : Used to set name of database. Default : "postgres"
* ACCESS_SECRET : Used to set secret key to be used in generation of JWT Token
* TOKEN_EXPIRE_TIME_SEC : Used to set time in seconds after which token will expire. Default : "900"seconds

## USAGE TIPS :

### API available 
1) **Create User**
```json
POST   /createuser

Body :   {
          "username" : "testusername",
          "name" : "testname",
          "password" : "pass123"
         }
```

2. **Login**
```json
POST   /login

Body :  {
          "username" : "testusername",
          "password" : "pass123"
        }
```
> In response to successful login you will get a JWT token , which you can use for accessing all other API's .Without valid token access will be denied.
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjQ3MTc0MDMsInVzZXJuYW1lIjoiYXNoIn0.vBO-kDDume2jvQTS7RhhlmyAvAew7bxLwxO5cpbFZFk"
}
```
> To use it , add the token string to the header of request like this :
```
{"Authorization" : "Bearer <token>" }
```

3. **Create Comment**
```json
POST  /comment

Body :  {
          "comment_text" : "test comment",
          "receiver_username" : "testusername"
        }
```
> receiver_username is the name of user one whose wall you are creating comment
> Response of this api call gives back ID of the created comment which can be used for other actions/API calls.

4. **Create Subcomment**
```json
POST  /subcomment

Body :  {
          "comment_text" : "test comment",
          "parent_comment_id" : 2
        }
```
> parent_comment_id is the ID of comment on which you are creating subcomment

5. **Create reaction**
```json
POST  /reaction

Body :  {
          "comment_id" : 2,
          "reaction" : "dislike"
        }
```
> - comment_id is ID of comment on which you wanted to react
> - currently only 3 reactions are supported : **[  like  ,  dislike  ,  +1  ]** . If something else is passed , api will reject it.
> - Reactions are string type and are case insensitive.

6. **Delete comment**
```json
DELETE  /comment?id=<comment-id>
```
> - Id of the comment to be deleted has to be passed as query param in URL itself
> - Delete is idempotent

7. **Generate Wall**
```json
GET   /wall
```
> Response will be in form of array of comments along with the count of reactions made on user's wall
```json
[
    {
        "comment_id": 2,
        "comment_text": "ash commented on nit",
        "sender_username": "ash",
        "timestamp": 1624710842,
        "reactions": {
            "dislike": 5
        }
    },
    {
        "comment_id": 3,
        "comment_text": "ash commented on nit again",
        "sender_username": "ash",
        "timestamp": 1624710870,
        "reactions": {}
    },
    {
        "comment_id": 4,
        "comment_text": "ash subcommented",
        "sender_username": "ash",
        "timestamp": 1624710925,
        "reactions": {
            "dislike": 1
        }
    }
]
```
