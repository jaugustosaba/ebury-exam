package tour

import (
	"strings"
	"testing"
)

var source = `GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5	
SCL,ORL,20`

var routes = []struct {
	origin, destiny string
	cost            int
}{
	{"GRU", "BRC", 10},
	{"BRC", "SCL", 5},
	{"GRU", "CDG", 75},
	{"GRU", "SCL", 20},
	{"GRU", "ORL", 56},
	{"ORL", "CDG", 5},
	{"SCL", "ORL", 20},
}

func TestAddCity(t *testing.T) {
	tour := NewTour()
	ID := tour.AddCity("A")
	if tour.GetCityName(ID) != "A" {
		t.Errorf("city name lookup by id failed")
	}
	if id, found := tour.GetCityID("A"); !found || (id != ID) {
		t.Errorf("city id lookup by name failed")
	}
}

func TestAddRoute(t *testing.T) {
	tour := NewTour()
	A := tour.AddCity("A")
	B := tour.AddCity("B")
	tour.AddRoute(A, B, 10)
	if cost, found := tour.Cost(A, B); !found || (cost != 10) {
		t.Fatalf("invalid cost from A->B")
	}
	if cost, found := tour.Cost(B, A); !found || (cost != 10) {
		t.Fatalf("invalid cost from B->A")
	}
}

func TestLoadFromCSV(t *testing.T) {
	tour := NewTour()
	reader := strings.NewReader(source)
	err := tour.LoadFromCSV(reader)
	if err != nil {
		t.Fatalf("unexpected error when loading CSV: %s", err.Error())
	}
	cities := []string{"GRU", "BRC", "SCL", "CDG", "ORL"}
	if len(tour.cities) != len(cities) {
		t.Fatalf("expecting %d cities found %d", len(tour.cities), len(cities))
	}
	for i := range routes {
		A, _ := tour.GetCityID(routes[i].origin)
		B, _ := tour.GetCityID(routes[i].destiny)
		cost, found := tour.Cost(A, B)
		if !found {
			t.Errorf("no route found for %s -> %s", routes[i].origin, routes[i].destiny)
		}
		if cost != routes[i].cost {
			t.Errorf("invalid cost for %s -> %s: %d", routes[i].origin, routes[i].destiny, cost)
		}
	}
}

func TestShortestPath(t *testing.T) {
	tour := NewTour()
	for i := range routes {
		origin := tour.AddCity(routes[i].origin)
		destiny := tour.AddCity(routes[i].destiny)
		tour.AddRoute(origin, destiny, routes[i].cost)
	}
	GRU, _ := tour.GetCityID("GRU")
	CDG, _ := tour.GetCityID("CDG")
	path, cost, err := tour.ShortestRoute(GRU, CDG)
	if err != nil {
		t.Fatalf("unexpected error when computing shortest path: %s", err.Error())
	}
	expectedPath := "GRU,BRC,SCL,ORL,CDG"
	pathStr := make([]string, len(path))
	for i := range path {
		pathStr[i] = tour.GetCityName(path[i])
	}
	str := strings.Join(pathStr, ",")
	if str != expectedPath {
		t.Errorf("invalid shortest path: %s ; expecting: %s", str, expectedPath)
	}
	if cost != 40 {
		t.Errorf("expecting cost 40 but found: %d", cost)
	}
}
