package api

import (
	"fmt"

	"github.com/jaugustosaba/ebury-exam/tour"
)

var (
	StatusOk    = "ok"
	StatusError = "error"
)

type Response[T any] struct {
	Status   string `json:"status"`
	Reason   string `json:"reason,omitempty"`
	Response T      `json:"response,omitempty"`
}

type AddRouteRequest struct {
	Origin  string `json:"origin"`
	Destiny string `json:"destiny"`
	Cost    int    `json:"cost"`
}

type AddRouteOutput struct {
}

type AddRouteResponse = Response[AddRouteOutput]

type ShortestRouteRequest struct {
	Origin  string `json:"origin"`
	Destiny string `json:"destiny"`
}

type ShortestRouteOutput struct {
	ShortestRoute []string `json:"shortestRoute"`
	Cost          int      `json:"cost"`
}

type ShortestRouteResponse = Response[*ShortestRouteOutput]

type Service struct {
	tour *tour.Tour
}

func NewService(tour *tour.Tour) *Service {
	return &Service{
		tour: tour,
	}
}

func (aService *Service) AddRoute(request *AddRouteRequest) *AddRouteResponse {
	if len(request.Origin) == 0 || len(request.Destiny) == 0 {
		return &AddRouteResponse{
			Status: StatusError,
			Reason: "city name cannot be empty",
		}
	}

	originID := aService.tour.AddCity(request.Origin)
	destinyID := aService.tour.AddCity(request.Destiny)
	aService.tour.AddRoute(originID, destinyID, request.Cost)
	return &AddRouteResponse{
		Status:   StatusOk,
		Response: AddRouteOutput{},
	}
}

func (aService *Service) ShortestRoute(request *ShortestRouteRequest) *ShortestRouteResponse {
	originID, found := aService.tour.GetCityID(request.Origin)
	if !found {
		return &ShortestRouteResponse{
			Status: StatusError,
			Reason: fmt.Sprintf("unknown origin city: `%s`", request.Origin),
		}
	}

	destinyID, found := aService.tour.GetCityID(request.Destiny)
	if !found {
		return &ShortestRouteResponse{
			Status: StatusError,
			Reason: fmt.Sprintf("unknown destiny city: `%s`", request.Destiny),
		}
	}

	route, cost, err := aService.tour.ShortestRoute(originID, destinyID)
	if err != nil {
		return &ShortestRouteResponse{
			Status: StatusError,
			Reason: err.Error(),
		}
	}

	routeStr := make([]string, len(route))
	for i := range route {
		routeStr[i] = aService.tour.GetCityName(route[i])
	}

	return &ShortestRouteResponse{
		Status: StatusOk,
		Response: &ShortestRouteOutput{
			ShortestRoute: routeStr,
			Cost:          cost,
		},
	}
}
