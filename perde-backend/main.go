package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- 1. VERİTABANI MODELLERİ ---

// Müşteri Tablosu
type Customer struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	Orders  []Order `json:"orders" gorm:"foreignKey:CustomerID"` // Bir müşterinin siparişleri
}

// Sipariş Tablosu (Ana Fiş)
type Order struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	CustomerID  uint        `json:"customer_id"`                     // Hangi müşterinin?
	TotalAmount float64     `json:"total_amount"`                    // Toplam Tutar
	Note        string      `json:"note"`                            // Ek notlar
	CreatedAt   time.Time   `json:"created_at"`                      // Sipariş Tarihi
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"` // Siparişin içindeki perdeler
}

// Sipariş Kalemleri (Perdeler)
type OrderItem struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	OrderID uint    `json:"order_id"`
	Room    string  `json:"room"`   // Oda (Salon, Mutfak)
	Type    string  `json:"type"`   // Tül, Stor
	Width   float64 `json:"width"`  // En
	Height  float64 `json:"height"` // Boy
	Pile    float64 `json:"pile"`   // Pile Sıklığı
	Price   float64 `json:"price"`  // O perdenin fiyatı
}

var DB *gorm.DB

func ConnectDB() {
	var err error
	// SENİN AYARLARIN: Port 5433, Şifre 12345
	dsn := "host=localhost user=postgres password=12345 dbname=perde_db port=5433 sslmode=disable"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Veritabanı Hatası:", err)
	}
	fmt.Println("🚀 Veritabanı Bağlantısı Başarılı!")

	// Tabloları Otomatik Oluştur (İlişkilerle Beraber)
	DB.AutoMigrate(&Customer{}, &Order{}, &OrderItem{})
	fmt.Println("✅ Tablolar Hazırlandı (Customers, Orders, OrderItems)")
}

func main() {
	ConnectDB()
	app := fiber.New()
	app.Use(cors.New())

	// --- API ROTALARI ---

	// 1. Müşterileri Getir
	app.Get("/api/customers", func(c *fiber.Ctx) error {
		var customers []Customer
		DB.Find(&customers)
		return c.JSON(customers)
	})

	// 2. Yeni Müşteri Ekle
	app.Post("/api/customers", func(c *fiber.Ctx) error {
		customer := new(Customer)
		if err := c.BodyParser(customer); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		DB.Create(&customer)
		return c.JSON(customer)
	})

	// 3. SİPARİŞ KAYDETME (Yeni Özellik)
	app.Post("/api/orders", func(c *fiber.Ctx) error {
		order := new(Order)

		// Frontend'den gelen veriyi oku
		if err := c.BodyParser(order); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		// Sipariş tarihini şimdi olarak ayarla
		order.CreatedAt = time.Now()

		// Veritabanına kaydet (GORM kalemleri de otomatik kaydeder)
		result := DB.Create(&order)
		if result.Error != nil {
			return c.Status(500).SendString("Sipariş kaydedilemedi")
		}

		// Müşterinin bakiyesini (borcunu) güncelle
		var customer Customer
		if err := DB.First(&customer, order.CustomerID).Error; err == nil {
			customer.Balance += order.TotalAmount
			DB.Save(&customer)
		}

		return c.JSON(order)
	})
	// --- BU KODU MAIN FONKSIYONUNUN İÇİNE, DİĞER ROTALARIN ALTINA EKLE ---

	// 4. Müşteri Detayını Getir (Sipariş Geçmişiyle Birlikte)
	app.Get("/api/customers/:id", func(c *fiber.Ctx) error {
		id := c.Params("id") // URL'den ID'yi al (Örn: /api/customers/5)
		var customer Customer

		// GORM'un Sihirli Kısmı: Preload
		// Önce "Orders"ları yükle, sonra o Order'ların içindeki "Items"ları yükle.
		result := DB.Preload("Orders.Items").First(&customer, id)

		if result.Error != nil {
			return c.Status(404).SendString("Müşteri bulunamadı")
		}

		return c.JSON(customer)
	})

	log.Fatal(app.Listen(":3000"))
}
