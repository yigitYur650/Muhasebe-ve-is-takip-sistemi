package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Employee struct {
	ID            string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email" gorm:"unique"`
	PasswordHash  string    `json:"-"`
	Role          string    `json:"role" gorm:"default:staff"`
	NfcTagID      string    `json:"nfc_tag_id"`
	HourlyRate    float64   `json:"hourly_rate"`
	WalletBalance float64   `json:"wallet_balance"`
	CreatedAt     time.Time `json:"created_at"`
}
type Shift struct {
	ID         string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EmployeeID string     `json:"employee_id"`
	StoreID    string     `json:"store_id"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    *time.Time //pointer başta boş olduğu anlamına geliyor
	PhoneUsage int        `json:"phone_usage_minutes" gorm:"column:phone_usage_minutes"`
}
type StartShiftDTO struct {
	NfcID   string `json:"nfc_tag_id"` //Telofondan gelen kart
	StoreID string `json:"store_id"`   //Hangi mağaza
}

// DB Veritabanı bağlantı değişkeni
var DB *gorm.DB

// Veritabanına bağlanma fonksiyonu
func connectDB() {
	// Docker ayarları: localhost, port 5432, şifre 12345
	dsn := "host=localhost user=postgres password=12345 dbname=personel_takip port=5432 sslmode=disable TimeZone=Europe/Istanbul"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı! Hata: ", err)
	}
	DB.AutoMigrate(&Employee{}, &Shift{})
	fmt.Println("🚀 Veritabanına başarıyla bağlanıldı!")
}

// CreateEmployee Çalışan oluşturma Fonksiyonu //Annemin ayakkabıları ayakkabıcıda genişletilecek
func CreateEmployee(c *fiber.Ctx) error {
	//Boş kart
	employee := new(Employee)
	//gelen veriyi buraya yaz
	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Veri Anlaşılamadı" + err.Error()})
	}
	// Veritabanına kaydet
	result := DB.Create(&employee)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Kaydedilmedi" + result.Error.Error()})
	}
	//Oluşan çalışanı dön
	return c.Status(201).JSON(employee)
}
func StartShift(c *fiber.Ctx) error {
	// gelen veriyi al
	payload := new(StartShiftDTO)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "kart numarası"})
	}
	// bu kart kimin
	var employee Employee
	// veritabanında NFC ID si eşleşen çalışanı bul
	if result := DB.Where("nfc_tag_id = ?", payload.NfcID).First(&employee); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tanımsız kart"})
	}
	// 3. Adam zaten içerde mi? (COUNT YÖNTEMİ - BULLETPROOF)
	var count int64
	DB.Model(&Shift{}).Where("employee_id = ? AND end_time IS NULL", employee.ID).Count(&count)

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Zaten açık bir mesain var! Önce çıkış yapmalısın."})
	}

	// 4. Mesaiyi başlat
	newShift := Shift{
		EmployeeID: employee.ID,
		StoreID:    payload.StoreID, // Mağaza ID'sini de kaydettik
		StartTime:  time.Now(),
	}
	DB.Create(&newShift)

	return c.Status(200).JSON(fiber.Map{
		"message": "Hoşgeldin " + employee.FullName + " Mesain başladı",
		"time":    newShift.StartTime,
	})
}
func EndShift(c *fiber.Ctx) error {
	//kart verisini al
	payload := new(StartShiftDTO)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Kart verisi okunmadı"})
	}
	// kart kimin
	var employee Employee
	if result := DB.Where("nfc_tag_id = ?", payload.NfcID).First(&employee); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Tanımsız kart"})
	}
	var activeShift Shift
	// başta * işaretli endtime sütunu NULL olan kaydı getir
	if result := DB.Where("employee_id = ? AND end_time IS NULL", employee.ID).First(&activeShift); result.Error != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Zaten dışarıdasın çıkış yapamazsın"})
	}
	//çıkış işlemi
	now := time.Now()
	activeShift.EndTime = &now //şimdiki zamanı attık
	//burada ne kullandığımızı anlamadım
	duration := now.Sub(activeShift.StartTime)
	minutes := int(duration.Minutes())
	//veritabanını güncelle
	DB.Save(&activeShift)
	return c.Status(200).JSON(fiber.Map{
		"message":        "güle güle " + employee.FullName,
		"end_time":       now,
		"worked_minutes": minutes,
	})
}

func main() {
	// 1. Veritabanına bağlan
	connectDB()
	// 2. Fiber uygulamasını başlat
	app := fiber.New()
	app.Post("/api/employees", CreateEmployee)
	// 3. Basit bir test rotası
	app.Post("/api/shifts/start", StartShift)
	app.Post("/api/shifts/end", EndShift)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Selam Kanka! Personel Takip Sistemi Çalışıyor 🚀")
	})
	// 4. Uygulamayı 3000 portundan yayınla
	log.Fatal(app.Listen(":3000"))
}
