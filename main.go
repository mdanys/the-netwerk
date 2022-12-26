package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/yudapc/go-rupiah"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Result = []Invest{}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

const charset = "0123456789"

type Keterangan struct {
	ID            uint
	JenisKelamin  string `json:"jenis_kelamin" form:"jenis_kelamin"`
	Usia          uint   `json:"usia" form:"usia"`
	Perokok       string `json:"perokok" form:"perokok"`
	Nominal       int    `json:"nominal" form:"nominal"`
	LamaInvestasi int    `json:"lama_investasi" form:"lama_investasi"`
}

type Invest struct {
	ID         uint    `json:"-" form:"-"`
	Awal       int     `json:"awal" form:"awal"`
	Bunga      int     `json:"bunga" form:"bunga"`
	Akhir      int     `json:"akhir" form:"akhir"`
	Persentase float64 `json:"-" form:"-"`
}

type User struct {
	ID                uint   `json:"-" form:"-"`
	Nama              string `json:"nama" form:"nama"`
	JenisKelamin      string `json:"jenis_kelamin" form:"jenis_kelamin"`
	Usia              uint   `json:"usia" form:"usia"`
	Email             string `gorm:"unique" json:"email" form:"email"`
	Perokok           string `json:"perokok" form:"perokok"`
	Nominal           int    `json:"nominal" form:"nominal"`
	LamaInvestasi     int    `json:"lama_investasi" form:"lama_investasi"`
	PeriodePembayaran string `json:"periode_pembayaran" form:"periode_pembayaran"`
	MetodeBayar       string `json:"metode_bayar" form:"metode_bayar"`
}

type Transaksi struct {
	ID                uint   `json:"-" form:"-"`
	TanggalTransaksi  string `json:"tgl_transaksi" form:"tgl_transaksi"`
	NoTransaksi       string `json:"no_transaksi" form:"no_transaksi"`
	Nama              string `json:"nama" form:"nama"`
	JenisKelamin      string `json:"jenis_kelamin" form:"jenis_kelamin"`
	Usia              uint   `json:"usia" form:"usia"`
	Nominal           int    `json:"nominal" form:"nominal"`
	LamaInvestasi     int    `json:"lama_investasi" form:"lama_investasi"`
	PeriodePembayaran string `json:"periode_pembayaran" form:"periode_pembayaran"`
	MetodeBayar       string `json:"metode_bayar" form:"metode_bayar"`
	TotalBayar        string `json:"total_bayar" form:"total_bayar"`
}

func Perhitungan() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input Keterangan
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot bind input data",
			})
		}

		var data Invest
		if input.Perokok == "Ya" && input.JenisKelamin == "Pria" {
			data.Persentase = 1
			if input.Usia > 0 && input.Usia <= 30 {
				data.Persentase += 1
			} else if input.Usia >= 31 && input.Usia <= 50 {
				data.Persentase += 0.5
			} else if input.Usia > 50 {
				data.Persentase += 0
			}
		} else if input.Perokok == "Tidak" && input.JenisKelamin == "Pria" {
			data.Persentase = 2
			if input.Usia > 0 && input.Usia <= 30 {
				data.Persentase += 1
			} else if input.Usia >= 31 && input.Usia <= 50 {
				data.Persentase += 0.5
			} else if input.Usia > 50 {
				data.Persentase += 0
			}
		} else if input.Perokok == "Ya" && input.JenisKelamin == "Wanita" {
			data.Persentase = 2
			if input.Usia > 0 && input.Usia <= 30 {
				data.Persentase += 1
			} else if input.Usia >= 31 && input.Usia <= 50 {
				data.Persentase += 0.5
			} else if input.Usia > 50 {
				data.Persentase += 0
			}
		} else if input.Perokok == "Tidak" && input.JenisKelamin == "Wanita" {
			data.Persentase = 3
			if input.Usia > 0 && input.Usia <= 30 {
				data.Persentase += 1
			} else if input.Usia >= 31 && input.Usia <= 50 {
				data.Persentase += 0.5
			} else if input.Usia > 50 {
				data.Persentase += 0
			}
		}

		data.Awal = input.Nominal

		for i := 1; i <= input.LamaInvestasi; i++ {
			data.Awal += data.Bunga
			data.Bunga = int(float64(data.Awal) * ((data.Persentase) / (100)))
			data.Akhir = data.Awal + data.Bunga
			Result = append(Result, data)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success",
			"status":  200,
			"data":    Result,
		})
	}
}

func Transaction(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input User
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot bind input data",
			})
		}

		var transaksi Transaksi

		transaksi.NoTransaksi = "TRX" + String(8)
		transaksi.TanggalTransaksi = time.Now().Format(time.RFC822)
		transaksi.Nama = input.Nama
		transaksi.JenisKelamin = input.JenisKelamin
		transaksi.Usia = input.Usia
		transaksi.Nominal = input.Nominal
		transaksi.LamaInvestasi = input.LamaInvestasi
		transaksi.PeriodePembayaran = input.PeriodePembayaran
		transaksi.MetodeBayar = input.MetodeBayar

		if input.PeriodePembayaran == "tahunan" {
			totalBayar := float64(input.Nominal - (input.Nominal / 12))
			formatRupiah := rupiah.FormatRupiah(totalBayar)
			transaksi.TotalBayar = formatRupiah
		} else {
			totalBayar := float64(input.Nominal)
			formatRupiah := rupiah.FormatRupiah(totalBayar)
			transaksi.TotalBayar = formatRupiah
		}

		if err := db.Create(&transaksi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "there is a problem on server",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success",
			"status":  200,
			"data":    transaksi,
		})
	}
}

func GetData(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var resQuery []Transaksi
		if err := db.Find(&resQuery).Error; err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "there is a problem on server",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success",
			"status":  200,
			"data":    resQuery,
		})
	}
}

func ConnectDB() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/thenetwerk?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("cannot initialize database", err.Error())
		return nil
	}

	return db
}

func autoGenerate(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return autoGenerate(length, charset)
}

func GenerateToken(id int, role string) string {
	claim := make(jwt.MapClaims)
	claim["authorized"] = true
	claim["role"] = role
	claim["id"] = id
	claim["expired"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	str, err := token.SignedString([]byte("Pr0J3ct"))
	if err != nil {
		log.Error("error on token signed string", err.Error())
		return "cannot generate token"
	}

	return str
}

func ExtractToken(c echo.Context) (id int, role string) {
	token := c.Get("user").(*jwt.Token)
	if token.Valid {
		claim := token.Claims.(jwt.MapClaims)
		return int(claim["id"].(float64)), string(claim["role"].(string))
	}

	return 0, ""
}


func main() {
	e := echo.New()
	db := ConnectDB()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Transaksi{})

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.POST("/", Perhitungan())

	e.POST("/invest", Transaction(db))

	e.GET("/invest", GetData(db))

	e.Start(":8000")
}
