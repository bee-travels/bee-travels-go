package server

import (
	"encoding/json"
	"fmt"
	"github.com/bee-travels/bee-travels-go/services/destination-v2/wrappers/database"
	instana "github.com/instana/go-sensor"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
)

type RequestContext = iris.Context
type RequestHandler = func(ctx RequestContext)
type RouterInitializer = func(router PathRouter, pool database.Pool, sensor *instana.Sensor)

func Response(ctx iris.Context, code int, response interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.StatusCode(code)

	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	if code != http.StatusOK {
		if span, ok := instana.SpanFromContext(ctx.Request().Context()); ok {
			span.SetTag("reason", string(data))
		}
	}
	ctx.Write(data)
}

func Start(serviceName string, init RouterInitializer) error {
	tracer := instana.NewTracerWithOptions(
		&instana.Options{
			Service:           serviceName,
			EnableAutoProfile: true,
		},
	)

	sensor := instana.NewSensorWithTracer(tracer)

	pool, err := database.NewDatabasePool(sensor)
	if err != nil {
		return errors.Errorf("Database connection not available: %v", err)
	}

	app := iris.New()

	// Add Instana tracer to all calls
	app.WrapRouter(func(w http.ResponseWriter, req *http.Request, router http.HandlerFunc) {
		adapter := instana.TracingHandlerFunc(sensor, "", router)
		adapter.ServeHTTP(w, req)
	})

	// Initialize the default readiness probe for k8s
	app.Get("/ready", func(ctx *context.Context) {
		ctx.StatusCode(http.StatusOK)
		ctx.Text("ok")
	})

	// Initialize user registered handlers
	init(&innerRouter{party: app}, pool, sensor)

	port := ":9001"
	if addr, ok := os.LookupEnv("PORT"); ok {
		port = addr
	}
	address := fmt.Sprintf(":%s", port)

	fmt.Println("Starting webserver...")
	err = app.Listen(address)
	if err != nil {
		return errors.Errorf("Couldn't start authentication server: %v", err)
	}
	return nil
}

func newCors(methods ...string) RequestHandler {
	found := false
	for _, method := range methods {
		if method == iris.MethodOptions {
			found = true
			break
		}
	}
	if !found {
		methods = append(methods, iris.MethodOptions)
	}
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: methods,
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Authorization",
			"X-INSTANA-T",
			"X-INSTANA-S",
			"X-INSTANA-L",
		},
	})
}

type PathRouter interface {
	Path(path string, init func(router PathRouter))
	Get(handler RequestHandler)
	Post(handler RequestHandler)
	Delete(handler RequestHandler)
	Put(handler RequestHandler)
	Patch(handler RequestHandler)
	Head(handler RequestHandler)
}

type handlerRegistration struct {
	handler RequestHandler
}

type innerRouter struct {
	party   iris.Party
	sensor  *instana.Sensor
	methods map[string]*handlerRegistration
}

func (ir innerRouter) Get(handler RequestHandler) {
	if ir.methods[http.MethodGet] != nil {
		panic(fmt.Errorf("GET already registered"))
	}
	ir.methods[http.MethodGet] = &handlerRegistration{
		handler: handler,
	}
}

func (ir innerRouter) Post(handler RequestHandler) {
	if ir.methods[http.MethodPost] != nil {
		panic(fmt.Errorf("POST already registered"))
	}
	ir.methods[http.MethodPost] = &handlerRegistration{
		handler: handler,
	}
}

func (ir innerRouter) Delete(handler RequestHandler) {
	if ir.methods[http.MethodDelete] != nil {
		panic(fmt.Errorf("DELETE already registered"))
	}
	ir.methods[http.MethodDelete] = &handlerRegistration{
		handler: handler,
	}
}

func (ir innerRouter) Put(handler RequestHandler) {
	if ir.methods[http.MethodPut] != nil {
		panic(fmt.Errorf("PUT already registered"))
	}
	ir.methods[http.MethodPut] = &handlerRegistration{
		handler: handler,
	}
}

func (ir innerRouter) Patch(handler RequestHandler) {
	if ir.methods[http.MethodPatch] != nil {
		panic(fmt.Errorf("PATCH already registered"))
	}
	ir.methods[http.MethodPatch] = &handlerRegistration{
		handler: handler,
	}
}

func (ir innerRouter) Head(handler RequestHandler) {
	if ir.methods[http.MethodHead] != nil {
		panic(fmt.Errorf("HEAD already registered"))
	}
	ir.methods[http.MethodHead] = &handlerRegistration{
		handler: handler,
	}
}

func (ir innerRouter) Path(path string, init func(router PathRouter)) {
	ir.party.PartyFunc(path, func(inner iris.Party) {
		iir := innerRouter{
			party:   inner,
			methods: make(map[string]*handlerRegistration),
		}

		pathTemplate := inner.GetRelPath()
		for strings.Contains(pathTemplate, ":") {
			index := strings.Index(pathTemplate, ":")
			end := index + strings.Index(pathTemplate[index:], "}")
			pathTemplate = pathTemplate[:index] + pathTemplate[end:]
		}

		pathTemplateAdapter := func(next RequestHandler) RequestHandler {
			return func(ctx RequestContext) {
				if span, ok := instana.SpanFromContext(ctx.Request().Context()); ok {
					span.SetTag("http.path_tpl", pathTemplate)
				}
				next(ctx)
			}
		}

		// Initialize inner paths
		init(iir)

		methods := make([]string, 0)
		for k, _ := range iir.methods {
			methods = append(methods, k)
		}

		iir.party.Reset()

		cors := newCors(methods...)
		iir.party.Use(cors)
		for k, registration := range iir.methods {
			handler := pathTemplateAdapter(registration.handler)

			switch k {
			case iris.MethodGet:
				iir.party.Get("", handler)
			case iris.MethodPost:
				iir.party.Post("", handler)
			case iris.MethodPut:
				iir.party.Put("", handler)
			case iris.MethodPatch:
				iir.party.Patch("", handler)
			case iris.MethodDelete:
				iir.party.Delete("", handler)
			case iris.MethodHead:
				iir.party.Head("", handler)
			}
		}
		iir.party.Options("", pathTemplateAdapter(cors))
	})
}
