package apiclient

import (
	"strconv"

	"github.com/go-resty/resty/v2"
)

type ApiClient struct {
	client *resty.Client
}

func New(client *resty.Client) *ApiClient {
	return &ApiClient{client: client}
}

type ErrorResponse struct {
	ErrorValue string `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.ErrorValue
}

type CreateSubscriptionRequest struct {
	ChatID      int64
	AddressArID int64
}

func (api *ApiClient) CreateSubscription(chatID int64, addressArID int64) error {
	_, err := api.client.R().
		SetBody(CreateSubscriptionRequest{
			ChatID:      chatID,
			AddressArID: addressArID,
		}).
		Post("/subscription")
	if err != nil {
		return err
	}

	return nil
}

type GeocodeResponseAddress struct {
	ID            uint32
	Distance      float64
	AddressString string
}
type GeocodeResponse struct {
	Addresses []GeocodeResponseAddress
}

func (api *ApiClient) Geocode(
	lat float64,
	lng float64,
	accuracy float64,
) (*GeocodeResponse, error) {

	type Response struct {
		Result *GeocodeResponse
	}

	resp, err := api.client.R().
		SetResult(&Response{}).
		SetError(&ErrorResponse{}).
		SetPathParams(
			map[string]string{
				"lat":      strconv.FormatFloat(lat, 'f', -1, 64),
				"lng":      strconv.FormatFloat(lng, 'f', -1, 64),
				"accuracy": strconv.FormatFloat(accuracy, 'f', -1, 64),
			},
		).
		Get("/address/geocode/{lat}/{lng}/{accuracy}")
	if err != nil {
		return nil, err
	}
	if backendError := resp.Error(); backendError != nil {
		return nil, backendError.(*ErrorResponse)
	}
	result := resp.Result()
	return result.(*Response).Result, nil
}
