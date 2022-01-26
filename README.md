# HTTP REST Movie omdb API


## Description
This is an example of implementation of Clean Architecture in Go (Golang) projects.

Rule of Clean Architecture by Uncle Bob
 * Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
 * Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
 * Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
 * Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
 * Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has  4 Domain layer :
 * Models Layer
 * Repository Layer
 * Usecase Layer  
 * Delivery Layer

It may different already, but the concept still the same in application level, also you can see the change log from v1 to current version in Master.

### How To Run This Project
> Make Sure you have run the movies.sql in your mysql


Since the project already use Go Module, I recommend to put the source code in any folder but GOPATH.

#### Run the Testing

```bash
$ make test
```

#### Run the Applications
Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/bxcodec/go-clean-arch.git

#move to project
$ cd go-clean-arch

# Build the docker image first
$ make docker

# Run the application
$ make run

# check if the containers are running
$ docker ps


# Stop
$ make stop
```

# Execute the call
## Fetch movies
```
localhost:9090/movies?searchword=Batman&paginatioon=4
Params:
- searchword : the title of the movie
- paginatioon : page number of the result
```
## Get Single Movie
```
localhost:9090/movies/{:id}
Params:
- id : imdb movie id
```

