package stores_test

import (
	"fmt"
	"time"

	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/domain/services"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/stores"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ = Describe("OrderStore", func() {
	var (
		testDB          *gorm.DB
		existingOrderID uint
		store           services.OrderStore
	)

	BeforeEach(func() {
		var err error
		testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		Expect(testDB).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		err = testDB.AutoMigrate(&models.Order{}, &models.OrderDish{}, &models.OrderPosition{})
		Expect(err).NotTo(HaveOccurred())

		store = stores.NewOrderStore(testDB)

		testOrder := models.Order{
			Dishes: []models.OrderDish{
				{Name: "Pizza"},
				{Name: "Pasta"},
			},
			Source: domain.OrderSourcePhone,
			Status: domain.OrderStatusPending,
			Time:   time.Now(),
		}

		err = testDB.Save(&testOrder).Error
		Expect(err).NotTo(HaveOccurred())

		existingOrderID = testOrder.ID
	})

	Describe("FindByID", func() {
		When("the order exists", func() {
			It("returns the order", func() {
				order, err := store.FindByID(existingOrderID)
				Expect(err).NotTo(HaveOccurred())
				Expect(order).NotTo(BeNil())
				Expect(order.ID).To(Equal(existingOrderID))
				Expect(order.Status).To(Equal(domain.OrderStatusPending))
				Expect(order.Source).To(Equal(domain.OrderSourcePhone))
				Expect(order.Dishes).To(HaveLen(2))
				Expect(order.Dishes[0].Name).To(Equal("Pizza"))
				Expect(order.Dishes[1].Name).To(Equal("Pasta"))
			})
		})

		When("the order does not exist", func() {
			It("returns nil and no error", func() {
				order, err := store.FindByID(999)
				Expect(err).NotTo(HaveOccurred())
				Expect(order).To(BeNil())
			})
		})
	})

	Describe("GetAll", func() {
		When("there are orders", func() {
			It("returns all orders", func() {
				orders, err := store.GetAll(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(orders).NotTo(BeEmpty())
				Expect(orders).To(HaveLen(1))
				Expect(orders[0].Dishes).NotTo(BeEmpty())
			})
		})

		When("there are a lot of orders", func() {
			const count = 40000
			BeforeEach(func() {
				for i := 0; i < count; i++ {
					order := models.Order{
						Status: domain.OrderStatusPending,
						Source: domain.OrderSourcePhone,
						Dishes: []models.OrderDish{
							{Name: fmt.Sprintf("dish-%d", i)},
						},
					}

					Expect(testDB.Save(&order).Error).To(Succeed())
				}
			})

			It("returns all orders", func() {
				orders, err := store.GetAll(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(orders).NotTo(BeEmpty())
				Expect(len(orders)).To(Equal(count + 1))
			})
		})

		When("filtering by status", func() {
			It("returns only matching orders", func() {
				filters := &domain.OrderFilters{AnyStatus: []domain.OrderStatus{domain.OrderStatusPending}}
				orders, err := store.GetAll(filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(orders).To(HaveLen(1))
				Expect(orders[0].Dishes).NotTo(BeEmpty())
			})
		})

		When("ordering by priority", func() {
			var positions []models.OrderPosition

			BeforeEach(func() {
				orders := make([]domain.Order, 3)

				for i := 0; i < 3; i++ {
					order := models.Order{
						Dishes: []models.OrderDish{
							{Name: fmt.Sprintf("dish-%d", i)},
						},
						Source: domain.OrderSourcePhone,
						Status: domain.OrderStatusPending,
						Time:   time.Now(),
					}

					Expect(testDB.Save(&order).Error).ToNot(HaveOccurred())
					orders[i] = domain.Order{ID: order.ID}
				}

				positions = []models.OrderPosition{
					{OrderID: orders[1].ID, Position: 1},
					{OrderID: orders[2].ID, Position: 2},
					{OrderID: existingOrderID, Position: 3},
					{OrderID: orders[0].ID, Position: 4},
				}

				Expect(testDB.Create(&positions).Error).ToNot(HaveOccurred())
			})

			It("should return orders in the right order", func() {
				result, err := store.GetAll(&domain.OrderFilters{PrioritySort: true})
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(HaveLen(4))

				for i, pos := range positions {
					Expect(result[i].ID).To(Equal(pos.OrderID))
				}
			})
		})

		When("no orders match the filter", func() {
			It("returns an empty list", func() {
				filters := &domain.OrderFilters{AnyStatus: []domain.OrderStatus{domain.OrderStatusDone}}
				orders, err := store.GetAll(filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(orders).To(BeEmpty())
			})
		})
	})

	Describe("Save", func() {
		When("saving a new order", func() {
			It("persists the order", func() {
				newOrder := domain.Order{
					NewOrder: domain.NewOrder{
						Dishes: []domain.Dish{{Name: "Burger"}},
						Source: "Web",
						Time:   time.Now(),
					},
					Status: domain.OrderStatusPending,
				}

				savedOrder, err := store.Save(newOrder)
				Expect(err).NotTo(HaveOccurred())
				Expect(savedOrder.ID).NotTo(BeZero())

				fetchedOrder, err := store.FindByID(savedOrder.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(fetchedOrder).NotTo(BeNil())
				Expect(fetchedOrder.Dishes).To(HaveLen(1))
			})
		})

		When("updating an existing order", func() {
			It("updates the order main attributes", func() {
				testOrder, err := store.FindByID(existingOrderID)
				Expect(err).NotTo(HaveOccurred())

				testOrder.Dishes = nil
				testOrder.Status = domain.OrderStatusDone

				updatedOrder, err := store.Save(*testOrder)
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedOrder.Status).To(Equal(domain.OrderStatusDone))

				var count int64
				err = testDB.Model(&models.OrderDish{}).Where("order_id = ?", existingOrderID).Count(&count).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(BeNumerically("==", 2))
			})

			It("updates the order dishes", func() {
				testOrder, err := store.FindByID(existingOrderID)
				Expect(err).NotTo(HaveOccurred())

				testOrder.Dishes = testOrder.Dishes[1:]
				updatedOrder, err := store.Save(*testOrder)
				Expect(err).NotTo(HaveOccurred())
				Expect(updatedOrder.Status).To(Equal(domain.OrderStatusPending))

				var count int64
				err = testDB.Model(&models.OrderDish{}).Where("order_id = ?", existingOrderID).Count(&count).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(BeNumerically("==", 1))
			})
		})
	})
})
