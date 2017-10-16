Go Server
=========

HTTP Server Router with middleware

## Example

```go
import "github.com/euskadi31/go-server"

router := server.NewRouter(&server.Configuration{
    Host: "127.0.0.1",
    Port: 1337,
})

router.EnableMetrics()
router.EnableCors()
router.EnableHealthCheck()

router.AddHealthCheck("my-health-check", NewMyHealthCheck())

router.Use(MyMiddleWare())

router.AddController(MyController())


panic(router.ListenAndServe())

```


## License

go-server is licensed under [the MIT license](LICENSE.md).
