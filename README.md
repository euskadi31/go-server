Go Server
=========

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
