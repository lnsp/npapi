package npapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type jsonResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

func fetch(data interface{}, b string, params ...interface{}) error {
	components := append([]interface{}{apiAddress}, params...)
	url := fmt.Sprintf(b, components...)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	response := jsonResponse{false, data}
	if err := decoder.Decode(&response); err != nil {
		return err
	}
	if !response.Status {
		return errors.New("request failed")
	}
	return nil
}

func endpoint(b string, params ...interface{}) string {
	components := append([]interface{}{apiAddress}, params...)
	return fmt.Sprintf(b, components...)
}

func parseStringMapToFloat(data map[string]string) (map[string]float64, error) {
	result := make(map[string]float64)
	for i, s := range data {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		result[i] = f
	}
	return result, nil
}
