package gin_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/domain/services/mocks"
	internalGin "github.com/danbrato999/yuno-gveloz/internal/gin"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

const baseAPIUri = "/api/v1/orders"

var _ = Describe("OrdersHandler", func() {
	var (
		ctrl        *gomock.Controller
		mockService *mocks.MockOrderService
		router      *gin.Engine
		recorder    *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockService = mocks.NewMockOrderService(ctrl)
		recorder = httptest.NewRecorder()
		router = internalGin.GetServer(mockService)
	})

	Describe("Create Order", func() {
		var validNewOrder domain.NewOrder

		BeforeEach(func() {
			validNewOrder = domain.NewOrder{
				Time:   time.Now(),
				Source: domain.OrderSourcePhone,
				Dishes: []domain.Dish{
					{
						Name: "Pasta",
					},
				},
			}
		})

		When("the request is valid", func() {
			It("should return 201 Created", func() {
				order := &domain.Order{ID: 1, NewOrder: validNewOrder, Status: domain.OrderStatusPending}

				mockService.EXPECT().CreateOrder(gomock.Any()).Return(order, nil)

				body, _ := json.Marshal(validNewOrder)
				req, _ := http.NewRequest(http.MethodPost, baseAPIUri, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusCreated))
				Expect(recorder.Body.String()).To(ContainSubstring(`"id":1`))
			})
		})

		DescribeTable("request is incorrect", func(request domain.NewOrder) {
			body, _ := json.Marshal(request)
			req, _ := http.NewRequest(http.MethodPost, baseAPIUri, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		},
			Entry("when no time is provided", domain.NewOrder{Source: domain.OrderSourcePhone, Dishes: []domain.Dish{{Name: "Pizza"}}}),
			Entry("when no dishes are provided", domain.NewOrder{Time: time.Now(), Source: domain.OrderSourcePhone}),
			Entry("when no source is provided", domain.NewOrder{Time: time.Now(), Dishes: []domain.Dish{{Name: "Pizza"}}}),
			Entry("when invalid source is provided", domain.NewOrder{Source: "test", Dishes: []domain.Dish{{Name: "Pizza"}}, Time: time.Now()}),
		)

		When("service fails", func() {
			It("should return 500 Internal Server Error", func() {
				mockService.EXPECT().CreateOrder(gomock.Any()).Return(nil, errors.New("error"))

				body, _ := json.Marshal(validNewOrder)
				req, _ := http.NewRequest(http.MethodPost, baseAPIUri, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Find Order", func() {
		When("order exists", func() {
			It("should return 200 OK", func() {
				order := &domain.OrderWithStatusHistory{
					Order: domain.Order{ID: 1, Status: domain.OrderStatusPending},
				}

				mockService.EXPECT().FindByID(uint(1)).Return(order, nil)

				req, _ := http.NewRequest(http.MethodGet, baseAPIUri+"/1", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(ContainSubstring(`"id":1`))
			})
		})

		When("order is not found", func() {
			It("should return 404 Not Found", func() {
				mockService.EXPECT().FindByID(uint(1)).Return(nil, domain.ErrOrderNotFound)

				req, _ := http.NewRequest(http.MethodGet, baseAPIUri+"/1", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("List Orders", func() {
		When("active orders are requested", func() {
			It("should return 200 OK with filtered orders", func() {
				orders := []domain.Order{{ID: 1, Status: domain.OrderStatusPending}}
				mockService.EXPECT().FindMany(gomock.Len(1)).Return(orders, nil)

				req, _ := http.NewRequest(http.MethodGet, baseAPIUri+"?active=true", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(ContainSubstring(`"id":1`))
			})
		})

		When("all orders are requested", func() {
			It("should return 200 OK with filtered orders", func() {
				orders := []domain.Order{{ID: 1, Status: domain.OrderStatusPending}}
				mockService.EXPECT().FindMany(gomock.Len(0)).Return(orders, nil)

				req, _ := http.NewRequest(http.MethodGet, baseAPIUri, nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(ContainSubstring(`"id":1`))
			})
		})

		When("service fails", func() {
			It("should return 500 Internal Server Error", func() {
				mockService.EXPECT().FindMany(gomock.Any()).Return(nil, errors.New("error"))

				req, _ := http.NewRequest(http.MethodGet, baseAPIUri, nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Update Order Status", func() {
		When("order status is updated successfully", func() {
			It("should return 200 OK", func() {
				order := &domain.Order{ID: 1, Status: domain.OrderStatusDone}
				mockService.EXPECT().UpdateStatus(uint(1), domain.OrderStatusDone).Return(order, nil)

				req, _ := http.NewRequest(http.MethodPut, baseAPIUri+"/1/status/done", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(recorder.Body.String()).To(ContainSubstring(`"status":"done"`))
			})
		})

		When("order update fails due to invalid status", func() {
			It("should return 400 Bad Request", func() {
				mockService.EXPECT().UpdateStatus(uint(1), domain.OrderStatusDone).Return(nil, domain.ErrInvalidStatusUpdate)

				req, _ := http.NewRequest(http.MethodPut, baseAPIUri+"/1/status/done", nil)
				router.ServeHTTP(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
