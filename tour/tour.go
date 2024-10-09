package tour

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

var (
	ErrColumnCount = errors.New("invalid number of columns on CSV")
	ErrNoPath      = errors.New("no path to destiny")
)

// City IDs are positive integers
type CityID int

// A tour is a graph where nodes are cities and paths are routes
type Tour struct {
	cities      []string                  // ID -> City
	lookup      map[string]CityID         // City -> ID
	connections map[CityID]map[CityID]int // City -> City -> Cost
}

// Creates a new tour
func NewTour() *Tour {
	return &Tour{
		cities:      []string{},
		lookup:      map[string]CityID{},
		connections: map[CityID]map[CityID]int{},
	}
}

// Gets a city ID by name
func (aTour *Tour) GetCityID(name string) (CityID, bool) {
	id, found := aTour.lookup[name]
	if !found {
		id = CityID(-1)
	}
	return id, found
}

// Asks if a city exists in a tour
func (aTour *Tour) HasCity(name string) bool {
	_, found := aTour.lookup[name]
	return found
}

// Gets city name by city ID
func (aTour *Tour) GetCityName(cityID CityID) string {
	if cityID < 0 || int(cityID) >= len(aTour.cities) {
		panic(fmt.Errorf("invalid city ID"))
	}
	return aTour.cities[cityID]
}

// Adds a city to a tour
func (aTour *Tour) AddCity(name string) CityID {
	cityID, found := aTour.lookup[name]
	if !found {
		cityID = CityID(len(aTour.cities))
		aTour.cities = append(aTour.cities, name)
		aTour.lookup[name] = cityID
	}
	return cityID
}

func (aTour *Tour) addEdge(originID, destinyID CityID, cost int) {
	destinyMap, found := aTour.connections[originID]
	if !found {
		destinyMap = map[CityID]int{}
		aTour.connections[originID] = destinyMap
	}
	destinyMap[destinyID] = cost
}

// Adds a route between two cities and its cost
func (aTour *Tour) AddRoute(originID, destinyID CityID, cost int) {
	aTour.addEdge(originID, destinyID, cost)
	aTour.addEdge(destinyID, originID, cost)
}

// Gets the cost between two cities if a direct route exists
func (aTour *Tour) Cost(originID, destinyID CityID) (int, bool) {
	destinyMap := aTour.connections[originID]
	cost, found := destinyMap[destinyID]
	return cost, found
}

// Loads routes from a CSV
func (aTour *Tour) LoadFromCSV(reader io.Reader) error {
	csv := bufio.NewReader(reader)
	for eof := false; !eof; {
		line, err := csv.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				eof = true
			} else {
				return err
			}
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		columns := strings.Split(line, ",")
		if len(columns) != 3 {
			return ErrColumnCount
		}
		origin := strings.TrimSpace(columns[0])
		destiny := strings.TrimSpace(columns[1])
		cost, err := strconv.Atoi(strings.TrimSpace(columns[2]))
		if err != nil {
			return err
		}
		originID := aTour.AddCity(origin)
		destinyID := aTour.AddCity(destiny)
		aTour.AddRoute(originID, destinyID, cost)

	}
	return nil
}

// cost queue element
type costInfo struct {
	Index int
	City  CityID
	Cost  int
}

// cost queue used as a priority queue
type costQueue []*costInfo

func (aQueue costQueue) Len() int {
	return len(aQueue)
}

func (aQueue costQueue) Less(i, j int) bool {
	return aQueue[i].Cost < aQueue[j].Cost
}

func (aQueue costQueue) Swap(i, j int) {
	temp := aQueue[i]
	aQueue[i] = aQueue[j]
	aQueue[i].Index = i
	aQueue[j] = temp
	aQueue[j].Index = j
}

func (aQueue *costQueue) Push(x any) {
	item := x.(*costInfo)
	item.Index = len(*aQueue)
	*aQueue = append(*aQueue, item)
}

func (aQueue *costQueue) Pop() any {
	qi := (*aQueue)[len(*aQueue)-1]
	*aQueue = (*aQueue)[0 : len(*aQueue)-1]
	return qi
}

// Computes shortest route between two cities using Djikstra's algorithm
func (aTour *Tour) ShortestRoute(originID, destinyID CityID) ([]CityID, int, error) {
	// initializes priority queue
	queue := make(costQueue, 0, len(aTour.cities))
	heap.Init(&queue)

	// previous city
	previous := make([]CityID, len(aTour.cities))
	// current cost for a given city ID
	costByCity := make([]*costInfo, len(aTour.cities))

	// cost initialization (infinity)
	for i := range costByCity {
		previous[i] = -1 // no previous city

		cost := math.MaxInt
		if i == int(originID) {
			cost = 0
		}
		costByCity[i] = &costInfo{
			City: CityID(i),
			Cost: cost,
		}
		heap.Push(&queue, costByCity[i])
	}

	// discover minimal routes for each reachable city
	for len(queue) > 0 {
		info := heap.Pop(&queue).(*costInfo)
		for nextID, cost := range aTour.connections[info.City] {
			if info.Cost+cost < costByCity[nextID].Cost {
				previous[nextID] = info.City
				costByCity[nextID].Cost = info.Cost + cost
				heap.Fix(&queue, costByCity[nextID].Index)
			}
		}
	}

	// if we cannot reach destiny (no route possible)
	if previous[destinyID] < 0 {
		return nil, 0, ErrNoPath
	}

	// path in reverse order
	path := make([]CityID, 0, len(aTour.cities))
	for i := destinyID; i != originID; {
		path = append(path, i)
		i = previous[i]
	}
	path = append(path, originID)

	// reverse list
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - i - 1
		tmp := path[j]
		path[j] = path[i]
		path[i] = tmp
	}

	return path, costByCity[destinyID].Cost, nil

}
