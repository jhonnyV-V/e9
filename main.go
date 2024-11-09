package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DiagramResponse struct {
	SpecificVolumeLiquid float64 `json:"specific_volume_liquid"`
	SpecificVolumeVapor  float64 `json:"specific_volume_vapor"`
}

const R = 8.3145

func getAproxTemp(pressure float64) float64 {
	//got this by trying linear regression
	return float64((27.6382 + (47.2362 * pressure)) + 273.15)
}

func getAproxSaturatedLiquidMole(pressure float64) float64 {
	//got this by trying linear regression
	return float64((0.545 * pressure) - 0.0065)
}

func getAproxSaturatedGasMole(pressure float64) float64 {
	//got this by trying linear regression
	return float64(598.1093 - (pressure * 59.2664))
}

func getAproxVolume(p, t, n float64) float64 {
	return (t * R * n) / (p * 1e6)
}

func roundFloat(val float64) float64 {
	ratio := math.Pow(10, 5)
	return math.Round(val*ratio) / ratio
}

func main() {
	e := echo.New()

	e.GET("/phase-change-diagram", func(c echo.Context) error {
		rawPressure := c.QueryParam("pressure")
		if rawPressure == "" {
			return c.String(http.StatusBadRequest, "missing pressure query param")
		}
		v, err := strconv.ParseFloat(rawPressure, 32)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid pressure value")
		}
		pressure := float64(v)
		fmt.Printf("pressure: %fmpa\n", pressure)

		temp := getAproxTemp(v)
		fmt.Printf("temp: %f 󰔄 +273.15\n", temp)

		moleLiquid := getAproxSaturatedLiquidMole(v)
		fmt.Printf("mole liquid number: %f\n", moleLiquid)

		resultLiquid := getAproxVolume(pressure, temp, moleLiquid)
		fmt.Printf("result liquid: %f\n", resultLiquid)
		resultLiquid = roundFloat(resultLiquid)
		fmt.Printf("rounded result liquid: %f\n", resultLiquid)

		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")

		moleGas := getAproxSaturatedGasMole(v)
		fmt.Printf("mole gas number: %f\n", moleGas)

		resultGas := getAproxVolume(pressure, temp, moleGas)
		fmt.Printf("result gas: %f\n", resultGas)
		resultGas = roundFloat(resultGas)
		fmt.Printf("rounded result gas: %f\n", resultGas)

		response := DiagramResponse{
			SpecificVolumeLiquid: resultLiquid,
			SpecificVolumeVapor:  resultGas,
		}

		return c.JSON(http.StatusOK, response)
	})

	e.Start(":8080")
}
