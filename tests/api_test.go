package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
)

var apiUrl string

func init() {
	apiUrl = os.Getenv("API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:8080"
	}
}

func TestPingHandler(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/ping", apiUrl))
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read HTTP body")

	require.Equal(t, "pong", string(body), "Wrong ping response")
}

func TestPostOrdersAndPostOrdersComplete(t *testing.T) {
	r := bytes.NewReader([]byte(`{"orders": [{"weight": 1, "regions": 2, "delivery_hours": ["13:14-15:16"], "cost": 5}]}`))
	resp, err := http.Post(fmt.Sprintf("%s/orders", apiUrl), "application/json", r)
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read HTTP body")

	response := make([]OrderDto, 0)
	err = json.Unmarshal(body, &response)

	require.Equal(t, 1, len(response), "must create only one order")
	require.Equal(t, float64(1), response[0].Weight)
	require.Equal(t, int64(2), response[0].Regions)
	require.Equal(t, int64(5), response[0].Cost)
	require.Equal(t, []string{"13:14-15:16"}, response[0].DeliveryHours)
}

func TestPostOrdersComplete(t *testing.T) {
	postOrderRequest := bytes.NewReader([]byte(`{"orders": [{"weight": 6, "regions": 3, "delivery_hours": ["16:16-17:17"], "cost": 10}]}`))
	postOrderResponse, err := http.Post(fmt.Sprintf("%s/orders", apiUrl), "application/json", postOrderRequest)
	if err != nil {
		t.Error("cannot execute post order request")
	}
	defer postOrderResponse.Body.Close()
	body, err := io.ReadAll(postOrderResponse.Body)
	if err != nil {
		t.Error("cannot read post order response")
	}
	response := make([]OrderDto, 0)
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Error("cannot unmarshal post order response")
	}

	id := strconv.FormatInt(response[0].OrderId, 10)
	str := fmt.Sprintf(`{"complete_info": [{"courier_id": 110 ,"order_id": %s, "complete_time": "2023-05-11T18:58:12.340Z"}]}`, id)
	rComplete := bytes.NewReader([]byte(str))

	respComplete, err := http.Post(fmt.Sprintf("%s/orders/complete", apiUrl), "application/json", rComplete)

	require.NoError(t, err, "HTTP error")
	defer respComplete.Body.Close()

	require.Equal(t, http.StatusOK, respComplete.StatusCode, "HTTP status code")

	bodyComplete, err := io.ReadAll(respComplete.Body)
	require.NoError(t, err, "failed to read HTTP body")

	responseComplete := make([]OrderDto, 0)
	err = json.Unmarshal(bodyComplete, &responseComplete)

	require.Equal(t, 1, len(responseComplete), "must create only one order")
	require.Equal(t, float64(6), responseComplete[0].Weight)
	require.Equal(t, int64(3), responseComplete[0].Regions)
	require.Equal(t, int64(10), responseComplete[0].Cost)
	require.Equal(t, []string{"16:16-17:17"}, responseComplete[0].DeliveryHours)
	require.Equal(t, "2023-05-11T18:58:12.340Z", *responseComplete[0].CompletedTime)
}

func TestPostCouriers(t *testing.T) {
	r := bytes.NewReader([]byte(`{"couriers":[{"courier_type": "AUTO","regions": [5], "working_hours": ["16:18-20:21"]}]}`))
	resp, err := http.Post(fmt.Sprintf("%s/couriers", apiUrl), "application/json", r)
	require.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read HTTP body")

	response := new(PostCouriersResponse)
	err = json.Unmarshal(body, &response)

	couriers := response.Couriers
	require.Equal(t, 1, len(couriers), "must create only one order")
	require.Equal(t, "AUTO", couriers[0].CourierType)
	require.Equal(t, []int64{5}, couriers[0].Regions)
	require.Equal(t, []string{"16:18-20:21"}, couriers[0].WorkingHours)
}

type PostCouriersResponse struct {
	Couriers []CourierDto `json:"couriers"`
}

type CourierDto struct {
	CourierId    int64    `json:"courier_id"`
	CourierType  string   `json:"courier_type"`
	Regions      []int64  `json:"regions"`
	WorkingHours []string `json:"working_hours"`
}

type OrderDto struct {
	OrderId       int64    `json:"order_id"`
	Weight        float64  `json:"weight"`
	Regions       int64    `json:"regions"`
	DeliveryHours []string `json:"delivery_hours"`
	Cost          int64    `json:"cost"`
	CompletedTime *string  `json:"completed_time"`
}
