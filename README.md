# Tiny URL

Service for shortening URL by providing special aliases, which will redirect
to actual URL. Similar services: TinyURL (itself), byt.ly, goo.gl, qlink.me,
etc

Originally idea came from [course] of system design. Included parts:

- [x] Why do we need URL shortening?
    * Understood main idea of service
- [ ] Requirements and Goals of the System
    * Created minimal functionality of service
- [ ] Capacity Estimation and Constraints
- [ ] System APIs
    * Developed create URL operation
- [x] Database Design
    * Read about SQL vs NoSQL in order for choosing database type
    * Choose NoSQL database - MongoDB, because of weak relationship 
      between entities. Also, it is easier to scale
    * Copy original schema of database
- [ ] Basic System Design and Algorithm
    * Developed algorithm of encoding original url with md5 and base64
- [ ] Data Partitioning and Replication
- [ ] Cache
- [ ] Load Balancer (LB)
- [ ] Purging or DB cleanup
- [ ] Telemetry
- [ ] Security and Permissions

## Swagger

Run `swag init -g internal/app/app.go` for generating openapi documentation

## Build

Run `go build cmd/app/main.go` for building project

## Run

Run built binary with `./main`

## Format

Before any commit run `gofmt -s -w .` for formatting whole project

[course]: https://www.educative.io/courses/grokking-the-system-design-interview/m2ygV4E81AR
