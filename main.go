package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Result = []Invest{}

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

func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.POST("/", Perhitungan())

	e.Start(":8000")
}
