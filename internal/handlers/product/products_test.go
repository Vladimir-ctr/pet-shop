package product

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	testifymock "github.com/stretchr/testify/mock"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	mock := new(ProductsMock)
	mock.On("GetAllProducts", testifymock.Anything).Return([]models.Product{
		{ID: 1, Name: "Cat Food", Price: 200, Stock: 50},
	}, nil)

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), mock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	mock := new(ProductsMock)
	mock.On("GetAllProducts", testifymock.Anything).Return(nil, errors.New("DB error"))

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта (200 OK)
	mock := new(ProductsMock)
	mock.On("CreateProduct", testifymock.Anything, testifymock.Anything).Return(1, nil)

	payload := `{"ID": 1, "Name": "Cat Food", "Price": 200, "Stock": 50}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	handler.CreateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта с невалидным JSON (400 Bad Request)
	mock := new(ProductsMock)
	mock.On("CreateProduct", testifymock.Anything, testifymock.Anything).Return(1, nil)

	payload := `{1}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", w.Code)
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта при ошибке сервиса (500 Internal Server Error)
	mock := new(ProductsMock)
	mock.On("CreateProduct", testifymock.Anything, testifymock.Anything).Return(0, errors.New("DB error"))

	payload := `{"ID": 1, "Name": "Cat Food", "Price": 200, "Stock": 50}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	handler.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected 500, got %d", w.Code)
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	mock := new(ProductsMock)
	mock.On("UpdateProduct", testifymock.Anything, testifymock.Anything).Return(nil)

	payload := `{"Name":"Tinker Food","Price":300,"Stock":10}`
	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	//для того чтобы chi.RouteContext был доступен в контексте запроса, нужно создать его вручную
	//для этого используем chi.NewRouteContext() и добавляем параметр id
	//затем с помощью context.WithValue добавляем chi.RouteContext в контекст
	//это нужно для того чтобы в хендлере можно было получить id продукта из URL
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)

	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта с невалидным JSON (400 Bad Request)
	mock := new(ProductsMock)
	mock.On("UpdateProduct", testifymock.Anything, testifymock.Anything).Return(nil)

	payload := `{1}`
	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	//для того чтобы chi.RouteContext был доступен в контексте запроса, нужно создать его вручную
	//для этого используем chi.NewRouteContext() и добавляем параметр id
	//затем с помощью context.WithValue добавляем chi.RouteContext в контекст
	//это нужно для того чтобы в хендлере можно было получить id продукта из URL
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)

	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта при ошибке сервиса (500 Internal Server Error)
	mock := new(ProductsMock)
	mock.On("UpdateProduct", testifymock.Anything, testifymock.Anything).Return(errors.New("DB error"))

	payload := `{"Name":"Tinker Food","Price":300,"Stock":10}`
	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	//для того чтобы chi.RouteContext был доступен в контексте запроса, нужно создать его вручную
	//для этого используем chi.NewRouteContext() и добавляем параметр id
	//затем с помощью context.WithValue добавляем chi.RouteContext в контекст
	//это нужно для того чтобы в хендлере можно было получить id продукта из URL
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)

	handler.UpdateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected 500, got %d", w.Code)
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта (200 OK)
	mock := new(ProductsMock)
	mock.On("DeleteProduct", testifymock.Anything, testifymock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	//делаем так же как и в UpdateProduct, чтобы chi.RouteContext был доступен в контексте запроса
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта с пустым id (400 Bad Request)
	mock := new(ProductsMock)
	mock.On("DeleteProduct", testifymock.Anything, testifymock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	//Удаляем проброс в контекст chi.RouteContext, чтобы id был пустым

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта при ошибке сервиса (500 Internal Server Error)
	mock := new(ProductsMock)
	mock.On("DeleteProduct", testifymock.Anything, testifymock.Anything).Return(errors.New("DB error"))

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	//делаем так же как и в UpdateProduct, чтобы chi.RouteContext был доступен в контексте запроса
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected 500, got %d", w.Code)
	}
}
