package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/a.vandam/product-search-challenge/src/domain/entities"
	"gitlab.com/a.vandam/product-search-challenge/src/logger"
)

func TestGetAProduct(t *testing.T) {
	testCase := getProdByIDTestCase{

		testName: "retrieve one product, not a palindrome",
		id:       123,
		existingProductsInPortMock: map[int]entities.ProductInfo{
			123: {
				ID:                 123,
				Title:              "a random product",
				FullPrice:          1000,
				FinalPrice:         1000,
				PriceModifications: 0.0,
			},
		},
		expectedProd: entities.ProductInfo{
			ID:                 123,
			Title:              "a random product",
			FullPrice:          1000,
			FinalPrice:         1000,
			PriceModifications: 0.0,
		},
		expectedErr: "",
	}

	t.Run(testCase.testName, testCase.testAndAssert)
}

func TestGetProductWithPalindromeId(t *testing.T) {
	testCase := getProdByIDTestCase{

		testName: "retrieve a product with palindrome id",
		id:       181,
		existingProductsInPortMock: map[int]entities.ProductInfo{
			181: {
				ID:          181,
				Title:       "a palindromic(?) product",
				FullPrice:   1000,
				Description: "palindrome",
			},
		},
		expectedProd: entities.ProductInfo{
			ID:                 181,
			Title:              "a palindromic(?) product",
			FullPrice:          1000,
			FinalPrice:         500,
			PriceModifications: -0.5,
			Description:        "palindrome",
		},
		expectedErr: "",
	}

	t.Run(testCase.testName, testCase.testAndAssert)
}

func TestNoProductFound(t *testing.T) {
	testCase := getProdByIDTestCase{

		testName: "retrieve no product as id doesn't match any",
		id:       55,
		existingProductsInPortMock: map[int]entities.ProductInfo{
			181: {
				ID:          181,
				Title:       "a palindromic(?) product",
				FullPrice:   1000,
				Description: "palindrome",
			},
		},
		expectedProd: entities.ProductInfo{},
		expectedErr:  "no products found with id: 55",
	}

	t.Run(testCase.testName, testCase.testAndAssert)
}

// End test cases

type getProdByIDTestCase struct {
	testName                   string
	id                         int
	existingProductsInPortMock map[int]entities.ProductInfo
	errorPortInMock            error
	expectedProd               entities.ProductInfo
	expectedErr                string
}

func (testCase getProdByIDTestCase) testAndAssert(t *testing.T) {
	t.Logf("testing function")

	/**---------------------- FUNCTION UNDER TEST -----------------------**/
	/*Dependencies*/
	mockedPort := mockedGetProductByIDPort{
		products: testCase.existingProductsInPortMock,
		err:      testCase.errorPortInMock,
	}
	loggerFactory := logger.LogFactory{LogLevel: "INFO"}
	log := loggerFactory.CreateLog("")
	svc := GetProductByIDServiceDefinition{
		Port: mockedPort,
		Log:  log,
	}
	/*END Dependencies*/
	testCtx := context.Background()
	product, err := svc.GetProductByID(testCtx, testCase.id)
	/**---------------------- END FUNCTION UNDER TEST -----------------------**/

	if !assert.Equal(t, testCase.expectedProd, product, "difference in value expected (%v) and obtained (%v)", testCase.expectedProd, product) {
		t.Fail()
	}
	if testCase.expectedErr != "" && err == nil {
		t.Logf("test failed as the function did not return an expected error: %v vs %v", err, testCase.expectedErr)
		t.FailNow()
	}
	if testCase.expectedErr == "" && err != nil {
		t.Logf("test failed as the function returned an error when it shouldn't: %v", err)
		t.FailNow()
	}
	if testCase.expectedErr != "" && err != nil {
		//comparing errors
		if !assert.EqualErrorf(t, err, testCase.expectedErr, "function returned an unexpected error: expected: %v vs found: %v", testCase.expectedErr, err) {
			t.FailNow()
		}
	}

	t.Logf("OK!!!! - test case:  %v  - OK!!!!", testCase.testName)
}

/*Mocking*/
type mockedGetProductByIDPort struct {
	products map[int]entities.ProductInfo
	err      error
}

func (mock mockedGetProductByIDPort) GetProductByID(ctx context.Context, id int) (entities.ProductInfo, error) {
	return mock.products[id], mock.err
}
