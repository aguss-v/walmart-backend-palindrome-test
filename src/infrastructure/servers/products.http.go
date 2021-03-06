package servers

import (
	"net/http"

	"gitlab.com/a.vandam/product-search-challenge/src/logger"
)

/*ProductsHTTPServer contains configurations for starting a http server*/
type ProductsHTTPServer struct {
	RouterFunc http.HandlerFunc
	Host       string
	Port       string
	Log        logger.LogContract
}

/*Start should be invoked after struct initialization. This methods starts the HTTP server, with it's main handlerFunc and port declared*/
func (svr *ProductsHTTPServer) Start() error {
	svr.Log.Info("starting http server on host:%v and port: %v", svr.Host, svr.Port)
	http.HandleFunc("/api", svr.RouterFunc)
	svr.Log.Info("main handler has been registered")
	hostPortPath := svr.Host + ":" + svr.Port
	svr.Log.Info("starting listen and server")
	err := http.ListenAndServe(hostPortPath, svr.RouterFunc)
	if err != nil {
		svr.Log.Error("error while starting http server: %v", err)
		return err
	}
	return nil
}
