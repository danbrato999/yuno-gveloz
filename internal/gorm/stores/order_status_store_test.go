package stores_test

import (
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

var _ = Describe("OrderStatusStore", func() {
	var (
		testDB          *gorm.DB
		existingOrderID uint
		store           services.OrderStatusStore
	)

	BeforeEach(func() {
		var err error
		testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		Expect(testDB).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		err = testDB.AutoMigrate(&models.Order{}, &models.OrderStatus{})
		Expect(err).NotTo(HaveOccurred())

		store = stores.NewOrderStatusStore(testDB)

		testOrder := models.Order{
			Source: domain.OrderSourcePhone,
			Status: domain.OrderStatusPending,
			Time:   time.Now(),
		}

		err = testDB.Save(&testOrder).Error
		Expect(err).NotTo(HaveOccurred())

		existingStatus := models.OrderStatus{
			OrderID: testOrder.ID,
			Status:  testOrder.Status,
		}

		Expect(testDB.Save(&existingStatus).Error).ToNot(HaveOccurred())

		existingOrderID = testOrder.ID
	})

	Describe("AddCurrentStatus", func() {
		When("order exists", func() {
			It("adds new entry to db", func() {
				order := &domain.Order{
					ID:     existingOrderID,
					Status: domain.OrderStatusPreparing,
				}

				Expect(store.AddCurrentStatus(order)).ToNot(HaveOccurred())

				var total int64

				err := testDB.
					Model(&models.OrderStatus{}).
					Where("order_id = ?", existingOrderID).
					Count(&total).
					Error

				Expect(err).ToNot(HaveOccurred())
				Expect(total).To(BeNumerically("==", 2))
			})
		})
	})

	Describe("GetHistory", func() {
		It("returns the statuses of an order", func() {
			history, err := store.GetHistory(existingOrderID)
			Expect(err).ToNot(HaveOccurred())
			Expect(history).To(HaveLen(1))
			Expect(history[0].Status).To(Equal(domain.OrderStatusPending))
		})
	})
})
