Go Server ![Last release](https://img.shields.io/github/release/euskadi31/go-server.svg)
=========

[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-server)](https://goreportcard.com/report/github.com/euskadi31/go-server)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://img.shields.io/travis/euskadi31/go-server/master.svg)](https://travis-ci.org/euskadi31/go-server) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-server/master.svg)](https://coveralls.io/github/euskadi31/go-server?branch=master) |
| develop | [![Build Status](https://img.shields.io/travis/euskadi31/go-server/develop.svg)](https://travis-ci.org/euskadi31/go-server) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-server/develop.svg)](https://coveralls.io/github/euskadi31/go-server?branch=develop) |

HTTP Server Router with middleware

## Example

```go
import "github.com/euskadi31/go-server"

router := server.NewRouter()

router.EnableMetrics()
router.EnableCors()
router.EnableHealthCheck()

router.AddHealthCheck("my-health-check", NewMyHealthCheck())

router.Use(MyMiddleWare())

router.AddController(MyController())

panic(http.ListenAndServe(":1337", router))

```


## License

go-server is licensed under [the MIT license](LICENSE.md).
