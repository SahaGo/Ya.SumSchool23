package main

import (
	"Ya.SumSchool23/controllers"
	"Ya.SumSchool23/controllers/dto"
	controller_errors "Ya.SumSchool23/controllers/errors"
	"Ya.SumSchool23/repositories"
	"Ya.SumSchool23/services"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func main() {
	setupViper()

	db := initDb(getConnectionString())
	defer db.Close()

	migrateDb(db, "postgres", viper.GetString("migrations.url"))

	//repository
	courierRepository := repositories.NewCourierRepository(db)
	orderRepository := repositories.NewOrderRepository(db)

	//service
	courierService := services.NewCourierService(courierRepository)
	orderService := services.NewOrderService(orderRepository)

	//controller
	pingController := controllers.NewPingController()
	courierController := controllers.NewCourierController(courierService)
	orderController := controllers.NewOrderController(orderService)

	e := echo.New()
	e.Validator = controllers.NewCustomValidator()
	setupPingRoutes(pingController, e)
	setupCourierRoutes(courierController, e)
	setupOrdersRoutes(orderController, e)

	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Logger.Fatal(e.Start(":8080"))
}

func getConnectionString() string {
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	dbName := viper.GetString("db.dbname")
	sslMode := viper.GetString("db.sslmode")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, password, host, port, dbName, sslMode)
}

func setupViper() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		log.Fatal("ENVIRONMENT must be specified (ENVIRONMENT=DEV/PROD/etc)")
	}
	if env == "DEV" {
		viper.SetConfigName("development")
		viper.AddConfigPath("../configs/")
	} else if env == "PROD" {
		viper.SetConfigName("production")
		viper.AddConfigPath("/etc/app/configs/")
	} else {
		log.Fatal("Unknown environment, expected PROD or DEV")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}
}

func setupPingRoutes(c *controllers.PingController, e *echo.Echo) {
	e.GET("/ping", c.Ping)
}

func setupCourierRoutes(c *controllers.CourierController, e *echo.Echo) {
	e.GET("/couriers", c.GetCouriers)
	e.GET("/couriers/:courier_id", c.GetCourierById)
	e.GET("/couriers/meta-info/:courier_id", c.GetCourierMetaById)
	e.POST("/couriers", c.PostCouriers)
}

func setupOrdersRoutes(c *controllers.OrderController, e *echo.Echo) {
	e.GET("/orders", c.GetOrders)
	e.GET("/orders/:order_id", c.GetOrderById)
	e.POST("/orders", c.PostOrders)
	e.POST("/orders/complete", c.PostOrdersComplete)

}

func initDb(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	return db
}

func migrateDb(db *sql.DB, dbName string, sourceURL string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to get Postgres driver: %s", err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(sourceURL, dbName, driver)
	if err != nil {
		log.Fatalf("failed to initialize db migration: %s ", err.Error())
	}
	m.Up()
}

func customHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

	errType := controller_errors.GetType(err)
	switch errType {
	case controller_errors.BadRequest:
		_ = c.JSON(http.StatusBadRequest, dto.EmptyResponse{})
		return
	case controller_errors.NotFound:
		_ = c.JSON(http.StatusNotFound, dto.EmptyResponse{})
		return
	case controller_errors.TooManyRequests:
		_ = c.JSON(http.StatusTooManyRequests, dto.EmptyResponse{})
	case controller_errors.NotImplemented:
		_ = c.JSON(http.StatusNotImplemented, dto.EmptyResponse{})
	default:
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		_ = c.String(code, err.Error())
	}
}
