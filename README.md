# Series Tracker Backend
## Description
The Series Tracker is a simple web application designed to help users manage their backlog of TV series. It allows users to track the status of their favorite series, 
monitor their progress, and organize their viewing experience. Developed as part of an educational exercise, we were given an existing frontend found [here](https://github.com/denn1s/series-tracker) and
required to implement the necessary API endpoints to handle its functionality. This project implements the following technologies for the API:

- [Go](https://go.dev) as the programming language
- [Echo](https://echo.labstack.com) as a lightweight HTTP library for Go
- [Air](https://github.com/air-verse/air) for hot-reloading during development
- [Docker](https://www.docker.com) for containerization
- [Docker Compose](https://docs.docker.com/compose/) for managing different containers & volumes
- [PostgreSQL](https://www.postgresql.org) for persistent data storage

In this README.md you'll find instructions on how to run the project, endpoint documentation (how to interact with the API) and project documentation (structure & general information for further development)

### Example
![image](https://github.com/user-attachments/assets/317797e6-e387-4812-b782-4df3c148d57f)

## Running this Project
### Dependencies
This project requires you to have [Docker](https://www.docker.com) & [Docker Compose](https://docs.docker.com/compose/)

### Commands
Head to the root of the directory and execute the following commands to run the backend in port 8080

```
docker compose up --build
```

After having run the command once, you can again pull up the container in the future without building using

```
docker compose up
```
Keep in mind the project is set up to run the database in port 5233.
### Running the Frontend
You're welcome to visit the existing frontend [repository](https://github.com/denn1s/series-tracker) and run it locally, keep in mind backend has CORS configured for port 80 as it's the default nginx port.

This project however includes the unmodified frontend repository with an additional Dockerfile, this lets you simply run the frontend more easily. To do so, navigate to the frontend folder and execute the following commands:
```
docker build -t frontend .
```
To build the container, you can name it however you'd like by replacing 'frontend'. After that, run
```
docker run -p 80:80 frontend
```
Where the name for running must match what you've named the container when building. After it's built it can be run by simply running the 'docker run' command again at any time.

## Documentation
### Swagger
This project includes documentation via Swagger, after running the project you can visit the route '/swagger/' or simply click [here](http://localhost:8080/swagger/index.html) to view the documentation.
### Postman
I've also included a Postman collection found [here](https://documenter.getpostman.com/view/45312303/2sB2qdgzj).

## About this Project
This section mainly goes over the structure of the API itself & general documentation for further development. Here I tackle some of the main parts of the project with brief descriptions of what they do so you don't go blindly looking through the project.

### Structure
This project uses a layered architecture, separating the API into the following layers:

- Repository: This layer is in charge of the main operations & communication with the database, this includes deleting, updating, creating and retrieving information from the database. There's no business logic involved within this layer.
  
- Service: This layer the 'core' of the application, it enforces business rules / checks & uses repository functions to manage data.
  
- Handlers: This layer is in charge of managing incoming HTTP requests, extracting information and mapping accordingly to the service functions.
  
### Error Handling
For error handling, the specific errors are declared at the top level of the internal package & a custom HTTP error handler within the handlers package takes care of logging the error & sending only basic information back as a response.

### Hot Reloading
This API uses [Air](https://github.com/air-verse/air) for hot reloading, this means that when the container is running it will automatically reflect any changes made in the backend to simplify development.
