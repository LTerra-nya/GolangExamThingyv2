package main

//wdym identical
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"time"
)

const outputFormat = "{TrainID: %v, DepartureStationID: %v, ArrivalStationID: %v, Price: %v, ArrivalTime: %v, DepartureTime: %v}\n"

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

func (t *Train) UnmarshalJSON(data []byte) error {
	var objects map[string]*json.RawMessage //we put all json data into a map, keys represent the names of variables and *json.RawMessage is raw json data

	if err := json.Unmarshal(data, &objects); err != nil {
		return err
	}

	//then we normally unmarshal the data into all the Train struct variables except for ArrivalTime and Departure time.
	if err := json.Unmarshal(*objects["trainId"], &t.TrainID); err != nil {
		return err
	}

	if err := json.Unmarshal(*objects["departureStationId"], &t.DepartureStationID); err != nil {
		return err
	}

	if err := json.Unmarshal(*objects["arrivalStationId"], &t.ArrivalStationID); err != nil {
		return err
	}

	if err := json.Unmarshal(*objects["price"], &t.Price); err != nil {
		return err
	}

	var arr, dep string //we use string variables to get the lines from the json file

	if err := json.Unmarshal(*objects["arrivalTime"], &arr); err != nil {
		return err
	}

	if err := json.Unmarshal(*objects["departureTime"], &dep); err != nil {
		return err
	}

	var h, m, s int

	if _, err := fmt.Sscanf(arr, "%d:%d:%d", &h, &m, &s); err != nil { //then we convert them into the varaibles h(ours), m(inutes),and s(econds)
		return err
	}

	t.ArrivalTime = time.Date(0, time.January, 1, //and then we set the Arrival and Departure time!
		h, m, s, 0,
		time.UTC)

	if _, err := fmt.Sscanf(dep, "%d:%d:%d", &h, &m, &s); err != nil {
		return err
	}

	t.DepartureTime = time.Date(0, time.January, 1,
		h, m, s, 0,
		time.UTC)

	return nil
}
func main() {
	var departureStation string
	var arrivalStation string
	var criteria string

	fmt.Println("Enter departure station:") //departure
	fmt.Scanln(&departureStation)

	if departureStation == "" {
		log.Fatal(errors.New("empty departure station"))
	}

	if _, err := strconv.Atoi(departureStation); err != nil {
		log.Fatal(errors.New("bad departure station input"))
	}

	fmt.Println("Enter arrival station:") //arrival
	fmt.Scanln(&arrivalStation)

	if arrivalStation == "" {
		log.Fatal(errors.New("empty arrival station"))
	}

	if _, err := strconv.Atoi(arrivalStation); err != nil {
		log.Fatal(errors.New("bad arrival station input"))
	}

	fmt.Println("Enter criteria for search(price, arrival-time, departure-time):") //criteria
	fmt.Scanln(&criteria)

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range result {
		fmt.Printf(
			outputFormat, val.TrainID, val.DepartureStationID, val.ArrivalStationID, val.Price, val.ArrivalTime, val.DepartureTime)
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	byteValue, err := ioutil.ReadFile("data.json") //Opening file
	if err != nil {
		return nil, err
	}

	var trains Trains

	if err := json.Unmarshal(byteValue, &trains); err != nil { //Unmarshalling the file into the variable trains
		return nil, err
	}

	dep, _ := strconv.Atoi(departureStation)
	arr, _ := strconv.Atoi(arrivalStation)

	var needed Trains //trains we need to sort

	for _, train := range trains {
		if train.DepartureStationID == dep && train.ArrivalStationID == arr {
			needed = append(needed, train)
		}
	}
	if len(needed) == 0 { //if we have no matching trains, we return nothing

		return nil, nil
	}
	switch criteria {
	case "price": // sort by price, lowest to highest
		sort.Slice(needed, func(i, j int) bool { return needed[i].Price < needed[j].Price })

		return needed, nil

	case "arrival-time": //sort by arrival time, lowest to highest
		sort.Slice(needed, func(i, j int) bool { return needed[i].ArrivalTime.Before(needed[j].ArrivalTime) })

		return needed, nil

	case "departure-time": //sort by departure time, lowest to highest
		sort.Slice(needed, func(i, j int) bool { return needed[i].DepartureTime.Before(needed[j].DepartureTime) })

		return needed, nil

	default:
		return nil, errors.New("unsupported criteria")
	}
}
