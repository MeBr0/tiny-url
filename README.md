# Tiny URL

![develop](https://github.com/mebr0/tiny-url/actions/workflows/develop.yml/badge.svg)

Service for shortening URL by providing special aliases, which will redirect
to actual URL. Similar services: TinyURL (itself), byt.ly, goo.gl, qlink.me,
etc

Originally idea came from [course] of system design. Included parts:

- [x] Why do we need URL shortening?
    * Understood main idea of service
- [ ] Requirements and Goals of the System
    * Created minimal functionality of service
- [ ] Capacity Estimation and Constraints
- [x] System APIs
    * Developed create and delete URL operations
    * Limit active URL count
- [x] Database Design
    * Read about SQL vs NoSQL in order for choosing database type
    * Choose NoSQL database - MongoDB, because of weak relationship 
      between entities. Also, it is easier to scale
    * Copy original schema of database
- [x] Basic System Design and Algorithm
    * Developed algorithm of encoding original url with md5 and base64
- [ ] Data Partitioning and Replication
- [x] Cache
    * Implement caching with Redis
    * Developed cache invalidation policy
- [ ] Load Balancer (LB)
- [ ] Purging or DB cleanup
    * Deleting expired URLs whenever it was accessed
- [ ] Telemetry
- [ ] Security and Permissions

## Variables

Use these variables to run project

```dotenv
MONGO_URI=mongodb://localhost:27017
MONGO_USER=<username>
MONGO_PASSWORD=<password>
MONGO_NAME=<db>

REDIS_URI=localhost:6379
REDIS_PASSWORD=<password>>
REDIS_DB=<db>
REDIS_TTL=10s

AUTH_ACCESS_TOKEN_TTL=5m
AUTH_PASSWORD_SALT=<salt>
AUTH_JWT_KEY=<key>

URL_ALIAS_LENGTH=8
URL_DEFAULT_EXPIRATION=30
URL_COUNT_LIMIT=3
```

## Commands

`go generate` - generate mock classes _(in package)_

`make fmt` - format whole project with gofmt _(do it before any commit)_

`make swag` - generate openapi documentation

`make cover` - run unit tests and show coverage report

`make build` - build project

`make run` - build and run project

## Docker

Use dockerfiles in `build` directory for building images and running containers

Use `build/Dockerfile` for building images on unix systems. 
Use `build/Dockerfile.multi` for building images on non-unix systems

[course]: https://www.educative.io/courses/grokking-the-system-design-interview/m2ygV4E81AR
