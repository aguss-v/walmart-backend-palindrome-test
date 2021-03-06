package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"gitlab.com/a.vandam/product-search-challenge/src/application/inhttp"
	"gitlab.com/a.vandam/product-search-challenge/src/application/outbound/mongodb"
	"gitlab.com/a.vandam/product-search-challenge/src/domain/services"
	"gitlab.com/a.vandam/product-search-challenge/src/infrastructure/configs"
	"gitlab.com/a.vandam/product-search-challenge/src/infrastructure/servers"
	"gitlab.com/a.vandam/product-search-challenge/src/logger"
)

func init() {

}

var log logger.LogContract

func main() {

	fmt.Println("Starting the application...")
	logConfigs, configErr := configs.GetLogConfigs()
	logFactory := logger.LogFactory{LogLevel: logConfigs.LogLevel}
	log = logFactory.CreateLog("main")
	if configErr != nil {
		log.Error("failed to start application ")
		os.Exit(1)
	}
	mongoDbConnector, dbConnErr := startDBOrFail(logFactory.CreateLog(""))
	defer mongoDbConnector.Disconnect(context.Background())
	if dbConnErr != nil {
		log.Error("impossible to start DB connector, exiting application: %v", dbConnErr)
		os.Exit(1)
	}

	log.Info("Starting product adapter")
	productsAdapter := mongodb.ProductsAdapter{
		DBConnector: mongoDbConnector,
		Log:         logFactory.CreateLog(""),
	}
	//handlers
	getProdByIDHandlerFunc := startGetProductByIDHandler(&productsAdapter, &logFactory)
	getProdByTextHandlerFunc := startGetProductByTextHandler(&productsAdapter, &logFactory)

	//routes
	log.Info("creating products router")
	productsRouter := inhttp.ProductsRouter{
		Log: logFactory.CreateLog(""),
	}
	productsRouter = productsRouter.AddRoute(http.MethodGet, "/api/products/([0-9]+)", *getProdByIDHandlerFunc)
	productsRouter = productsRouter.AddRoute(http.MethodGet, "/api/products/search", *getProdByTextHandlerFunc)
	log.Debug("products router registered routes: %+v", productsRouter.RegisteredRoutes())

	//server
	log.Info("creating main http server")
	httpServer := servers.ProductsHTTPServer{
		RouterFunc: productsRouter.Serve,
		Host:       "",
		Port:       "8080",
		Log:        logFactory.CreateLog(""),
	}
	err := httpServer.Start()
	if err != nil {
		log.Error("error while starting server: %v", err)
		os.Exit(1)
	}

}
func startDBOrFail(dbLog logger.LogContract) (*mongodb.MongoConnector, error) {
	dbConfigs, dbConfigErr := configs.GetProductsDatabaseConfigs()
	if dbConfigErr != nil {
		errMsg := fmt.Sprintf("failed to retrieve db configurations: %v ", dbConfigErr)
		log.Error(errMsg)
		return &mongodb.MongoConnector{}, errors.New(errMsg)
	}
	mongoConn := &mongodb.MongoConnector{Log: dbLog}
	mongoConn, dbConnErr := mongoConn.Connect(context.Background(), dbConfigs)
	if dbConnErr != nil {
		errMsg := fmt.Sprintf("failed to connect  to db: %v ", dbConnErr)
		log.Error(errMsg)
		return &mongodb.MongoConnector{}, errors.New(errMsg)
	}
	pingErr := mongoConn.Ping(context.Background())
	if pingErr != nil {
		errMsg := fmt.Sprintf("failed to ping  db: %v ", dbConnErr)
		log.Error(errMsg)
		return &mongodb.MongoConnector{}, errors.New(errMsg)
	}
	return mongoConn, nil
}

func startGetProductByIDHandler(productsAdapter *mongodb.ProductsAdapter, logFactory *logger.LogFactory) *http.HandlerFunc {
	getProductByIDServiceDef := services.GetProductByIDServiceDefinition{
		Port: productsAdapter,
		Log:  logFactory.CreateLog(""),
	}
	getProdByIDHandlerFunc := inhttp.CreateGetProdByIDHandlerFunc(inhttp.GetProductByIDHandlerDependencies{
		Svc: getProductByIDServiceDef,
		Log: logFactory.CreateLog(""),
	})
	log.Info("get prod handler func started")
	return &getProdByIDHandlerFunc
}

func startGetProductByTextHandler(productsAdapter *mongodb.ProductsAdapter, logFactory *logger.LogFactory) *http.HandlerFunc {
	getProductByTextServiceDef := services.GetProductsByTextServiceDefinition{
		Port: productsAdapter,
		Log:  logFactory.CreateLog(""),
	}
	getProdByTextHandlerFunc := inhttp.CreateGetProductByFieldHandlerFunc(inhttp.GetProductsByFieldDependencies{
		Svc: getProductByTextServiceDef,
		Log: logFactory.CreateLog(""),
	})
	log.Info("get prod by field handler func started")
	return &getProdByTextHandlerFunc
}

// func routes() {

// 	[]route{
// 		newRoute(http.MethodGet, "/api/products/([0-9]+))", CreateGetProdByIdHandlerFunc()),
// 		//newRoute(http.MethodGet, "/api/products/search)", apiUpdateWidget),
// 	}
// }
