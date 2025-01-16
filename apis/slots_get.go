package apis

import (
	"encoding/json"
	"flatpeak-devtask/structs"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	isContinuous     bool = false // Default behavioural value
	timeSlotDuration int  = 30    // In minutes
)

func GetSlots(c *gin.Context) {
	validFrom := time.Now()
	if duration, ok := c.GetQuery("duration"); ok {
		duration, err := strconv.Atoi(duration)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		if duration < 30 {
			duration = 30
		}
		if continuous, ok := c.GetQuery("continuous"); ok {
			isContinuous, err = strconv.ParseBool(continuous)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
			}
		}
		response := []structs.Slot{}

		arr, err := getIntensityFromNow24h(validFrom.Format(time.RFC3339))
		if err != nil || arr == nil {
			c.IndentedJSON(http.StatusInternalServerError, nil)
		}
		data := []structs.Intensity{}
		if isContinuous {
			data = getBestContinuousSlot(arr, duration/timeSlotDuration)
		} else {
			data = getBestSlots(arr, duration/timeSlotDuration)
		}
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, nil)
		}
		for i := 0; i < len(data); i++ {
			response = append(response, structs.ConvertIntensityItemsToSlot(data[i]))
		}
		c.IndentedJSON(http.StatusOK, response)
	} else {
		c.IndentedJSON(http.StatusInternalServerError, nil)
	}
}

func getIntensityFromNow24h(validFrom string) ([]structs.Intensity, error) {
	url := fmt.Sprintf("https://api.carbonintensity.org.uk/intensity/%s/fw24h", validFrom)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var responseWrapper structs.ResponseWrapper
	if err := json.Unmarshal(body, &responseWrapper); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return responseWrapper.Data, nil
}

func getBestSlots(arr []structs.Intensity, duration int) []structs.Intensity {
	leastIntensityTimeSlot := []structs.Intensity{}
	copy := arr
	for d := 0; d < duration; d++ {
		lowestForecast, lowestForecastIndex := *arr[0].Intensity.Forecast, 0
		for i := 1; i < len(arr); i++ {
			if arr[i].Intensity.Forecast != nil {
				if *arr[i].Intensity.Forecast < lowestForecast {
					lowestForecast = *arr[i].Intensity.Forecast
					lowestForecastIndex = i
				}
			}
		}
		leastIntensityTimeSlot = append(leastIntensityTimeSlot, arr[lowestForecastIndex])
		copy = removeCopy(copy, lowestForecastIndex)
	}
	return leastIntensityTimeSlot
}

func getBestContinuousSlot(arr []structs.Intensity, duration int) []structs.Intensity {
	leastIntensityTimeSlot := []structs.Intensity{}
	lowestForecast, lowestForecastIndex := 0, 0
	for i := 0; i < len(arr); i++ {
		totalCarbonForecast := 0
		if arr[i].Intensity.Forecast != nil {
			totalCarbonForecast = *arr[i].Intensity.Forecast
		}
		for j := 0; j < duration; j++ {
			if arr[j].Intensity.Forecast != nil {
				totalCarbonForecast += *arr[j].Intensity.Forecast
			}
		}
		if i == 0 {
			lowestForecast = totalCarbonForecast
		}
		if totalCarbonForecast < lowestForecast {
			lowestForecast = totalCarbonForecast
			lowestForecastIndex = i
		}
	}
	for i := lowestForecastIndex; i < lowestForecastIndex+duration; i++ {
		leastIntensityTimeSlot = append(leastIntensityTimeSlot, arr[i])
	}
	return leastIntensityTimeSlot
}

func removeCopy(slice []structs.Intensity, i int) []structs.Intensity {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
