package stores_test

import (
	"time"

	"github.com/danbrato999/yuno-gveloz/domain"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/models"
	"github.com/danbrato999/yuno-gveloz/internal/gorm/stores"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ = Describe("OrderPriorityStore", func() {
	var (
		testDB     *gorm.DB
		orderQueue []models.Order
		store      stores.OrderPositionStore
	)

	BeforeEach(func() {
		var err error
		testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		Expect(testDB).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		err = testDB.AutoMigrate(&models.Order{}, &models.OrderPosition{})
		Expect(err).NotTo(HaveOccurred())

		store = *stores.NewOrderPositionStore(testDB)

		times := []int{-10, -5, -2}
		orderQueue = make([]models.Order, len(times))

		for i, v := range times {
			testOrder := models.Order{
				Source: domain.OrderSourcePhone,
				Status: domain.OrderStatusPending,
				Time:   time.Now().Add(time.Duration(v) * time.Minute),
			}
			Expect(testDB.Save(&testOrder).Error).NotTo(HaveOccurred())

			position := models.OrderPosition{
				OrderID:  testOrder.ID,
				Position: uint(i + 1),
			}
			Expect(testDB.Save(&position).Error).NotTo(HaveOccurred())

			orderQueue[i] = testOrder
		}
	})

	Describe("Add", func() {
		When("adding a new order to the queue", func() {
			It("should set new order to latest position", func() {
				newOrder := models.Order{
					Source: domain.OrderSourceInPerson,
					Status: domain.OrderStatusPending,
					Time:   time.Now(),
				}
				Expect(testDB.Save(&newOrder).Error).NotTo(HaveOccurred())

				err := store.Add(&domain.Order{ID: newOrder.ID})
				Expect(err).ToNot(HaveOccurred())

				var positions []models.OrderPosition
				err = testDB.Order("position").Find(&positions).Error
				Expect(err).ToNot(HaveOccurred())
				Expect(positions).To(Equal([]models.OrderPosition{
					{OrderID: orderQueue[0].ID, Position: 1},
					{OrderID: orderQueue[1].ID, Position: 2},
					{OrderID: orderQueue[2].ID, Position: 3},
					{OrderID: newOrder.ID, Position: 4},
				}))
			})
		})

		When("there are no orders queued", func() {
			BeforeEach(func() {
				Expect(testDB.Exec("DELETE FROM order_positions").Error).ToNot(HaveOccurred())
			})

			It("should add the order at position 1", func() {
				Expect(store.Add(&domain.Order{ID: orderQueue[0].ID})).To(Succeed())

				var positions []models.OrderPosition
				err := testDB.Order("position").Find(&positions).Error
				Expect(err).ToNot(HaveOccurred())
				Expect(positions).To(Equal([]models.OrderPosition{
					{OrderID: orderQueue[0].ID, Position: 1},
				}))
			})
		})

		When("adding an existing order to the queue", func() {
			It("should return an error", func() {
				err := store.Add(&domain.Order{ID: orderQueue[0].ID})
				Expect(err).To(Equal(domain.ErrIncorrectOrderQueueing))
			})
		})
	})

	Describe("Remove", func() {
		It("should properly remove an order at the beginning", func() {
			err := store.Remove(orderQueue[0].ID)
			Expect(err).ToNot(HaveOccurred())

			var positions []models.OrderPosition
			err = testDB.Order("position").Find(&positions).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(positions).To(Equal([]models.OrderPosition{
				{OrderID: orderQueue[1].ID, Position: 1},
				{OrderID: orderQueue[2].ID, Position: 2},
			}))
		})

		It("should properly remove an order at the bottom", func() {
			err := store.Remove(orderQueue[2].ID)
			Expect(err).ToNot(HaveOccurred())

			var positions []models.OrderPosition
			err = testDB.Order("position").Find(&positions).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(positions).To(Equal([]models.OrderPosition{
				{OrderID: orderQueue[0].ID, Position: 1},
				{OrderID: orderQueue[1].ID, Position: 2},
			}))
		})

		It("should do nothing when removing an non existing order", func() {
			err := store.Remove(uint(666))
			Expect(err).ToNot(HaveOccurred())

			var positions []models.OrderPosition
			err = testDB.Order("position").Find(&positions).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(positions).To(Equal([]models.OrderPosition{
				{OrderID: orderQueue[0].ID, Position: 1},
				{OrderID: orderQueue[1].ID, Position: 2},
				{OrderID: orderQueue[2].ID, Position: 3},
			}))
		})
	})

	Describe("ShuffleAfter", func() {
		It("should move an order forward in the queue", func() {
			err := store.ShuffleAfter(orderQueue[0].ID, orderQueue[1].ID)
			Expect(err).ToNot(HaveOccurred())

			var positions []models.OrderPosition
			err = testDB.Order("position").Find(&positions).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(positions).To(Equal([]models.OrderPosition{
				{OrderID: orderQueue[1].ID, Position: 1},
				{OrderID: orderQueue[0].ID, Position: 2},
				{OrderID: orderQueue[2].ID, Position: 3},
			}))
		})

		It("should move an order backward in the queue", func() {
			err := store.ShuffleAfter(orderQueue[2].ID, orderQueue[0].ID)
			Expect(err).ToNot(HaveOccurred())

			var positions []models.OrderPosition
			err = testDB.Order("position").Find(&positions).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(positions).To(Equal([]models.OrderPosition{
				{OrderID: orderQueue[0].ID, Position: 1},
				{OrderID: orderQueue[2].ID, Position: 2},
				{OrderID: orderQueue[1].ID, Position: 3},
			}))
		})

		It("should not change order if already in correct place", func() {
			err := store.ShuffleAfter(orderQueue[1].ID, orderQueue[0].ID)
			Expect(err).ToNot(HaveOccurred())

			var positions []models.OrderPosition
			err = testDB.Order("position").Find(&positions).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(positions).To(Equal([]models.OrderPosition{
				{OrderID: orderQueue[0].ID, Position: 1},
				{OrderID: orderQueue[1].ID, Position: 2},
				{OrderID: orderQueue[2].ID, Position: 3},
			}))
		})

		It("should do nothing if one of the orders does not exist", func() {
			err := store.ShuffleAfter(uint(999), orderQueue[0].ID)
			Expect(err).ToNot(HaveOccurred())

			var positions []models.OrderPosition
			err = testDB.Order("position").Find(&positions).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(positions).To(Equal([]models.OrderPosition{
				{OrderID: orderQueue[0].ID, Position: 1},
				{OrderID: orderQueue[1].ID, Position: 2},
				{OrderID: orderQueue[2].ID, Position: 3},
			}))
		})

		It("should do nothing if both orders do not exist", func() {
			err := store.ShuffleAfter(uint(999), uint(888))
			Expect(err).ToNot(HaveOccurred())

			var positions []models.OrderPosition
			err = testDB.Order("position").Find(&positions).Error
			Expect(err).ToNot(HaveOccurred())
			Expect(positions).To(Equal([]models.OrderPosition{
				{OrderID: orderQueue[0].ID, Position: 1},
				{OrderID: orderQueue[1].ID, Position: 2},
				{OrderID: orderQueue[2].ID, Position: 3},
			}))
		})
	})
})
