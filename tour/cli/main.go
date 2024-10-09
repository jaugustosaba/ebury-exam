package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jaugustosaba/ebury-exam/tour"
)

func loadCSV(t *tour.Tour, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return t.LoadFromCSV(file)
}

func main() {
	t := tour.NewTour()
	for _, path := range os.Args[1:] {
		err := loadCSV(t, path)
		if err != nil {
			fmt.Printf("cannot read CSV file %s: %s", path, err.Error())
			os.Exit(1)
		}
	}

	for {
		fmt.Printf("please enter the route: ")
		var originDestiny string
		_, err := fmt.Scanf("%s", &originDestiny)
		if err != nil {
			break
		}
		cities := strings.Split(originDestiny, "-")
		if len(cities) != 2 {
			fmt.Printf("invalid input: '%s'\n", originDestiny)
			continue
		}
		originID, found := t.GetCityID(cities[0])
		if !found {
			fmt.Printf("unknown origin city: '%s'\n", cities[0])
			continue
		}
		destinyID, found := t.GetCityID(cities[1])
		if !found {
			fmt.Printf("unknown destiny city: '%s'\n", cities[1])
			continue
		}
		route, cost, err := t.ShortestRoute(originID, destinyID)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}
		routeStr := make([]string, len(route))
		for i := range route {
			routeStr[i] = t.GetCityName(route[i])
		}
		fmt.Printf("best route: %s > $%d\n", strings.Join(routeStr, " - "), cost)
	}
}
