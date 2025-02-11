package services_test

import (
	"errors"
	"fmt"
	"sync"

	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/domain/services/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("OrderService", func() {
	var (
		mockOrderStore    *mocks.MockOrderStore
		mockStatusStore   *mocks.MockOrderStatusStore
		mockPriorityQueue *mocks.MockPriorityQueue
		orderService      services.OrderService
	)

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		mockOrderStore = mocks.NewMockOrderStore(mockCtrl)
		mockPriorityQueue = mocks.NewMockPriorityQueue(mockCtrl)
		mockStatusStore = mocks.NewMockOrderStatusStore(mockCtrl)
		orderService = services.NewOrderService(mockOrderStore, mockPriorityQueue, mockStatusStore)
	})

	Context("CreateOrder", func() {
		It("should create an order successfully", func() {
			newOrder := domain.NewOrder{
				Dishes: []domain.Dish{{Name: "Pizza"}},
			}

			savedOrder := &domain.Order{
				NewOrder: newOrder,
				Status:   domain.OrderStatusPending,
				ID:       1,
			}

			mockOrderStore.EXPECT().Save(gomock.Any()).Return(savedOrder, nil)

			var wg sync.WaitGroup
			wg.Add(2)
			mockStatusStore.EXPECT().AddCurrentStatus(savedOrder).Do(func(o *domain.Order) {
				wg.Done()
			})
			mockPriorityQueue.EXPECT().Add(savedOrder).Do(func(o *domain.Order) {
				wg.Done()
			})

			order, err := orderService.CreateOrder(newOrder)

			Expect(err).To(Succeed())
			Expect(order).NotTo(BeNil())
			Expect(order.Status).To(Equal(domain.OrderStatusPending))

			wg.Wait()
		})

		It("should return an error if saving fails", func() {
			newOrder := domain.NewOrder{
				Dishes: []domain.Dish{{Name: "Pizza"}},
			}

			mockOrderStore.EXPECT().Save(gomock.Any()).Return(nil, errors.New("save error"))

			order, err := orderService.CreateOrder(newOrder)

			Expect(order).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("save error"))
		})
	})

	Context("FindByID", func() {
		It("should return an order with status history", func() {
			order := &domain.Order{ID: 1, Status: domain.OrderStatusPending}
			history := []domain.OrderStatusHistory{
				{
					Status: domain.OrderStatusPending,
				},
				{
					Status: domain.OrderStatusPreparing,
				},
			}

			mockOrderStore.EXPECT().FindByID(uint(1)).Return(order, nil)
			mockStatusStore.EXPECT().GetHistory(uint(1)).Return(history, nil)

			result, err := orderService.FindByID(1)

			Expect(err).To(Succeed())
			Expect(result).NotTo(BeNil())
			Expect(result.StatusHistory).To(Equal(history))
		})

		It("should return an error if order not found", func() {
			mockOrderStore.EXPECT().FindByID(uint(1)).Return(nil, nil)

			result, err := orderService.FindByID(1)

			Expect(result).To(BeNil())
			Expect(err).To(Equal(domain.ErrOrderNotFound))
		})

		It("should return an error if store fails", func() {
			testErr := fmt.Errorf("random error")
			mockOrderStore.EXPECT().FindByID(uint(1)).Return(nil, testErr)

			result, err := orderService.FindByID(1)

			Expect(result).To(BeNil())
			Expect(err).To(Equal(testErr))
		})
	})

	Context("UpdateStatus", func() {
		It("should update order status successfully", func() {
			order := &domain.Order{ID: 1, Status: domain.OrderStatusPending}

			mockOrderStore.EXPECT().FindByID(uint(1)).Return(order, nil)
			mockOrderStore.EXPECT().Save(gomock.Any()).Return(order, nil)

			var wg sync.WaitGroup
			wg.Add(1)
			mockStatusStore.EXPECT().AddCurrentStatus(order).Do(func(o *domain.Order) {
				wg.Done()
			})

			updatedOrder, err := orderService.UpdateStatus(1, domain.OrderStatusPreparing)

			Expect(err).To(BeNil())
			Expect(updatedOrder).NotTo(BeNil())
			Expect(updatedOrder.Status).To(Equal(domain.OrderStatusPreparing))

			wg.Wait()
		})

		It("should return an error if order does not exist", func() {
			mockOrderStore.EXPECT().FindByID(uint(1)).Return(nil, nil)

			order, err := orderService.UpdateStatus(1, domain.OrderStatusPreparing)

			Expect(order).To(BeNil())
			Expect(err).To(Equal(domain.ErrOrderNotFound))
		})

		It("should return an error if order is completed", func() {
			order := &domain.Order{ID: 1, Status: domain.OrderStatusDone}
			mockOrderStore.EXPECT().FindByID(uint(1)).Return(order, nil)

			result, err := orderService.UpdateStatus(1, domain.OrderStatusPreparing)

			Expect(result).To(BeNil())
			Expect(err).To(Equal(domain.ErrCompleteOrderUpdate))

		})

		It("should return an error if status transition is invalid", func() {
			order := &domain.Order{ID: 1, Status: domain.OrderStatusReady}
			mockOrderStore.EXPECT().FindByID(uint(1)).Return(order, nil)

			result, err := orderService.UpdateStatus(1, domain.OrderStatusPreparing)

			Expect(result).To(BeNil())
			Expect(err).To(Equal(domain.ErrInvalidOrderUpdate))
		})

		It("should remove order from queue when cancelled", func() {
			order := &domain.Order{ID: 1, Status: domain.OrderStatusPending}

			mockOrderStore.EXPECT().FindByID(uint(1)).Return(order, nil)
			mockOrderStore.EXPECT().Save(gomock.Any()).Return(order, nil)

			var wg sync.WaitGroup
			wg.Add(2)
			mockStatusStore.EXPECT().AddCurrentStatus(order).Do(func(o *domain.Order) {
				wg.Done()
			})
			mockPriorityQueue.EXPECT().Remove(order.ID).Do(func(id uint) {
				wg.Done()
			})

			updatedOrder, err := orderService.UpdateStatus(1, domain.OrderStatusCancelled)

			Expect(err).To(BeNil())
			Expect(updatedOrder).NotTo(BeNil())
			Expect(updatedOrder.Status).To(Equal(domain.OrderStatusCancelled))

			wg.Wait()
		})
	})

	Context("UpdateContent", func() {
		const fakeID uint = 123
		When("a new set of dishes is provided", func() {
			var dishes []domain.Dish
			BeforeEach(func() {
				dishes = []domain.Dish{{Name: "Lasagna"}}
			})

			It("should update when the order is in a valid state", func() {
				fakeOrder := &domain.Order{
					ID:     fakeID,
					Status: domain.OrderStatusPreparing,
				}
				mockOrderStore.EXPECT().FindByID(fakeID).Return(fakeOrder, nil)

				updatedOrder := &domain.Order{
					ID:     fakeID,
					Status: domain.OrderStatusPreparing,
					NewOrder: domain.NewOrder{
						Dishes: dishes,
					},
				}

				mockOrderStore.EXPECT().Save(gomock.Any()).Return(updatedOrder, nil)

				result, err := orderService.UpdateDishes(fakeID, dishes)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).ToNot(BeNil())
				Expect(result.Dishes).To(HaveLen(1))
			})

			DescribeTable("and the order is in an invalid state", func(status domain.OrderStatus) {
				fakeOrder := &domain.Order{
					ID:     fakeID,
					Status: status,
				}
				mockOrderStore.EXPECT().FindByID(fakeID).Return(fakeOrder, nil)

				result, err := orderService.UpdateDishes(fakeID, dishes)
				Expect(result).To(BeNil())
				Expect(err).To(Equal(domain.ErrInvalidOrderUpdate))
			},
				Entry("should error for ready", domain.OrderStatusReady),
				Entry("should error for done", domain.OrderStatusDone),
				Entry("should error for cancelled", domain.OrderStatusCancelled),
			)

			It("should error when the order doesn't exist", func() {
				mockOrderStore.EXPECT().FindByID(fakeID).Return(nil, nil)

				result, err := orderService.UpdateDishes(fakeID, dishes)
				Expect(result).To(BeNil())
				Expect(err).To(Equal(domain.ErrOrderNotFound))
			})
		})
		When("an empty set of dishes is provided", func() {
			It("should return an error", func() {
				result, err := orderService.UpdateDishes(123, nil)
				Expect(result).To(BeNil())
				Expect(err).To(Equal(domain.ErrInvalidOrderUpdate))
			})
		})
	})
})
