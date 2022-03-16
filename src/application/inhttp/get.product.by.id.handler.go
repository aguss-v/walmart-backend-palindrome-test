package inhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/a.vandam/product-search-challenge/src/domain/entities"
	"gitlab.com/a.vandam/product-search-challenge/src/domain/services"
	"gitlab.com/a.vandam/product-search-challenge/src/logger"
)

type GetProductByIdHandlerDependencies struct {
	Svc services.GetProductByIdServiceContract
	Log logger.LogContract
}

func CreateGetProdByIdHandlerFunc(dep GetProductByIdHandlerDependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		pathParamReceived := req.Context().Value(ProductIdPathParamCtxKey{}).([]string)
		idInPath, _ := strconv.Atoi(pathParamReceived[0])
		dep.Log.Info("received GET request for id: %v", idInPath)
		if idInPath == DefaultErrorInt || idInPath == 0 {
			errMsg := fmt.Sprintf("invalid product id path parameter sent: received: %v", idInPath)
			dep.Log.Error(errMsg)
			http.Error(rw, wrapErrAsJson(errors.New(errMsg)), http.StatusBadRequest)
			return
		}

		//Useful for time outs
		reqContext := req.Context()
		dep.Log.Debug("created context to obtain products")
		product, err := dep.Svc.GetProductById(idInPath, reqContext)
		if err != nil {
			dep.Log.Error("received an error while fetching product by id: %v", err)
			http.Error(rw, wrapErrAsJson(err), http.StatusInternalServerError)
			return
		}
		dep.Log.Info("mapping product to response")
		dep.Log.Debug("product being mapped: %+v", product)
		response, err := mapProductToJsonResponse(&product)
		if err != nil {
			dep.Log.Error("received an error while generating response: %v", err)
			http.Error(rw, wrapErrAsJson(err), http.StatusInternalServerError)
			return
		}
		dep.Log.Debug("response to write: %+v", string(response))
		rw.Write(response)
		dep.Log.Info("response has been sent back")
	})

}

func mapProductToJsonResponse(prod *entities.ProductInfo) ([]byte, error) {
	responseDTO := embeddingJsonResponse{
		ErrMsg: "",
		Resources: getProductJsonResponse{
			Id:                 prod.Id,
			Title:              prod.Title,
			Description:        prod.Description,
			ImageURL:           prod.ImageURL,
			FullPrice:          prod.FullPrice,
			FinalPrice:         prod.FinalPrice,
			PriceModifications: prod.PriceModifications,
		},
	}
	responseBody, err := json.Marshal(responseDTO)

	return responseBody, err
}

type embeddingJsonResponse struct {
	ErrMsg    string                 `json:"error"`
	Resources getProductJsonResponse `json:"resources"`
}
type getProductJsonResponse struct {
	Id                 int     `json:"id"`
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	ImageURL           string  `json:"imageURL"`
	FullPrice          float32 `json:"fullPrice"`
	FinalPrice         float32 `json:"finalPrice"`
	PriceModifications float32 `json:"priceModifications"`
}