# Introduction
Welcome to the repository for an example of a microservice architecture in Go with gRPC.
The main purpose of the repository is to demonstrate the capabilities of microservices written in Go.

# Installation 

To install, you need to clone the current repository. You also need to have `Docker` and `Docker-compose`.
After cloning, you need to add a .env file to the project root or add the variables to the system:
```
DATABASE_URL=postgresql://YOUR_LOGIN:1111@postgres-db:5432/auth_db?sslmode=disable
JWT_SECRET=YOUT_SECRET
RDB_ADDR=redis-db:6379
RDB_PASSWORD=YOUR_RDB_PASSWORD
RDB_DB=0
AUTH_SERVER_ADDR=auth-server:50051
```

Then, in the root of the project, enter `docker compose up --build -d`. After that, the project will be assembled.
**Don’t forget to add an administrator to the database for the test! At the same time, the administrator’s password must be hashed.**

# Using
Routes:
```
  / - home page, using method GET
  /login - login request, using method POST with json content-type, format: {"login":"login", "password":"password"}
  /admin - using method GET, needed for login verification 
  /admin/logout  - using method POST,  using for logout account
```

Examples of curl requests for verification:
- `curl "http://localhost:8080/"` - base home-page request
- `curl -c cookies.txt -X POST -i "http://localhost:8080/login" \
                                   -H "Content-Type: application/json" \
                                   -d '{"login":"login", "password":"password"}'` - login request with save cookies in file
- `curl -b cookies.txt "http://localhost:8080/admin"` - request with cookies to admin-page
- `curl -b cookies.txt -X POST -i "http://localhost:8080/admin/logout"` - logout request


**Attention! All requests to the `/admin...` address must be made with a cookie!**
