package routing

import (
	"os"
	"strings"

	"github.com/arsmn/fiber-swagger/example/docs"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	model "url-shortener/model"
)

var (
	defaultSkipper = func(c *fiber.Ctx) bool {
		return strings.Contains(c.Path(), "swagger")
	}
)

// FiberSkipper skip next
type fiberSkipper func(c *fiber.Ctx) bool

// FiberMiddleware init
type FiberMiddleware struct {
	Skipper fiberSkipper
	Config  interface{}
}

// InitFiber init http fiber
func InitFiber() *FiberMiddleware {
	preforkEnv := os.Getenv("PREFORK")
	m := &FiberMiddleware{
		Skipper: defaultSkipper,
		Config: fiber.New(fiber.Config{
			Prefork:       preforkEnv == "TRUE",
			CaseSensitive: false,
			StrictRouting: false,
			ErrorHandler: func(f *fiber.Ctx, err error) error {
				code := fiber.StatusInternalServerError
				if e, isOk := err.(*fiber.Error); isOk {
					code = e.Code
				}
				logrus.WithFields(
					logrus.Fields{
						"code": code,
					}).Error(err)
				return fiberError(code, f, err)
			},
		}),
	}

	return m
}

// InitFiberMiddleware fiber use
func (m *FiberMiddleware) InitFiberMiddleware() (*fiber.App, fiber.Router) {
	service := os.Getenv("APP_SERVICE")
	basePath := "/"

	f := m.Config.(*fiber.App)

	f.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	f.Use(recover.New())
	f.Use(m.fiberLogResponse)
	f.Use(m.fiberLogRequest)

	router := f.Group(basePath)

	if os.Getenv("APP_ENV") != "production" {
		docs.SwaggerInfo.Title = "Swagger " + service + " service"
		docs.SwaggerInfo.Host = os.Getenv("APP_HOST")
		docs.SwaggerInfo.BasePath = os.Getenv("SWAGGER_BASE_PATH")
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
		router.Get("/swagger/*", swagger.HandlerDefault)
	}
	router.Get("/healthz", healthcheck)

	return f, router
}

func healthcheck(c *fiber.Ctx) error {
	return c.JSON(model.NewBaseResponse(0, true))
}

func fiberError(code int, f *fiber.Ctx, err error) error {
	switch code {
	case 1:
		return f.Status(fiber.StatusBadRequest).JSON(model.NewBaseErrorResponse(code, err.Error()))
	case 2:
		return f.Status(fiber.StatusNotFound).JSON(model.NewBaseErrorResponse(code, err.Error()))
	case 3:
		return f.Status(fiber.StatusForbidden).JSON(model.NewBaseErrorResponse(code, "invalid token"))

	default:
		return f.Status(fiber.StatusBadRequest).JSON(model.NewBaseErrorResponse(code, err.Error()))
	}
}

func (m *FiberMiddleware) fiberLogRequest(c *fiber.Ctx) error {
	if m.Skipper(c) {
		return c.Next()
	}
	body := string(c.Request().Body())
	method := string(c.Request().Header.Method())

	newSpanID := strings.ReplaceAll(uuid.New().String(), "-", "")
	spanID := c.Get("span_id", newSpanID[len(newSpanID)-20:])

	c.Locals("span_id", spanID)

	logrus.WithFields(
		logrus.Fields{
			"spanID": spanID,
			"method": method,
			"path":   c.Path(),
			"body":   string(body),
		}).Infof("Request")
	return c.Next()
}

func (m *FiberMiddleware) fiberLogResponse(c *fiber.Ctx) error {
	if m.Skipper(c) {
		return c.Next()
	}
	if err := c.Next(); err != nil {
		return err
	}

	body := string(c.Response().Body())
	method := string(c.Request().Header.Method())
	spanID := c.Locals("span_id").(string)
	logrus.WithFields(
		logrus.Fields{
			"spanID": spanID,
			"method": method,
			"path":   c.Path(),
			"body":   string(body),
		}).Infof("Response")
	return nil
}
