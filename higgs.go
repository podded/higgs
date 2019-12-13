package higgs

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Eve ID Ranges - https://gist.github.com/a-tal/5ff5199fdbeb745b77cb633b7f4400bb


func DeleteStaticData(config Configuration) error {
	client, err := newClient(config)

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return err
	}

	return client.Store.DeleteStaticData()
}

func PopulateStaticData(config Configuration) error {

	client, err := newClient(config)

	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return err
	}


	err = client.Store.DeleteStaticData()
	if err != nil{
		return errors.Wrap(err, "Failed to delete existing static data")
	}

	err = populateUniverse(client)
	if err != nil {
		return errors.Wrap(err, "Failed to populate the universe")
	}

	err = populateTypes(client)
	if err != nil {
		return errors.Wrap(err, "Failed to populate types")
	}

	err = populateGroups(client)
	if err != nil {
		return errors.Wrap(err, "Failed to populate groups")
	}

	err = populateCategories(client)
	if err != nil {
		return errors.Wrap(err, "Failed to populate categories")
	}

	return nil

}

func populateUniverse(client *Client) error {

	client.Log.Println("WARNING!!!!")
	client.Log.Println("This command will take a very long time to run!!! I mean that!!!")
	client.Log.Println("It is also not fault tolerant. If you see any errors you need to run it again!!!")

	// Just in case it has bulk errors straight away
	time.Sleep(30 * time.Second)

	err := populateRegions(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import region information")
	}

	err = populateConstellations(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import constellation information")
	}

	err = populateSystems(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import system information")
	}

	err = populateStars(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import stars")
	}

	err = populatePlanets(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import planets")
	}

	err = populateMoons(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import moons")
	}

	err = populateAsteroidBelts(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import asteroid belts")
	}

	err = populateStargates(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import stargates")
	}

	err = populateStations(client)
	if err != nil {
		return errors.Wrap(err, "Failed to import stations")
	}

	return nil
}

func populateRegions(client *Client) error {

	var waitgroup sync.WaitGroup

	// First Step is to populate the region list. Will do a goroutine each, there isnt that many
	const urlRegion = "https://esi.evetech.net/latest/universe/regions/?datasource=tranquility"
	regionsBody, err := client.MakeESIGet(urlRegion)
	var regions []int
	err = json.Unmarshal(regionsBody, &regions)
	if err != nil {
		return err
	}

	client.Log.Printf("Have to get %v regions", len(regions))

	const urlRegionSpecifc = "https://esi.evetech.net/latest/universe/regions/%v/?datasource=tranquility"
	for _, rr := range regions {
		r := rr
		waitgroup.Add(1)
		go func() {
			regionURL := fmt.Sprintf(urlRegionSpecifc, r)
			regionBody, err := client.MakeESIGet(regionURL)
			if err != nil {
				client.Log.Printf("Failed to query region %v", r)
				waitgroup.Done()
				return
			}

			region := ESIRegion{}
			err = json.Unmarshal(regionBody, &region)
			if err != nil {
				client.Log.Printf("Failed to decode region %v from esi; %v", r, string(regionBody))
				waitgroup.Done()
				return
			}
			// client.Log.Printf("Adding Region - %v\n", region.Name)
			err = client.Store.InsertRegion(region)
			if err != nil {
				log.Fatalln(errors.Wrap(err, "Failed to insert region"))
			}
			waitgroup.Done()
		}()

	}

	// Only waiting because I want to be in order for my OCD :P
	waitgroup.Wait()

	return nil
}

func populateConstellations(client *Client) error {

	var waitgroup sync.WaitGroup

	// Now grab all the constellations. Lets batch these out and do 50 goroutines... Dont want to go too fast...
	const urlConstellation = "https://esi.evetech.net/latest/universe/constellations/?datasource=tranquility"
	constellationsBody, err := client.MakeESIGet(urlConstellation)
	var constellations []int
	err = json.Unmarshal(constellationsBody, &constellations)
	if err != nil {
		return err
	}

	client.Log.Printf("Have to get %v constellations", len(constellations))

	const urlConstellationSpecifc = "https://esi.evetech.net/latest/universe/constellations/%v/?datasource=tranquility"

	var batches [][]int

	batchSize := (len(constellations) / client.MaxRoutines) + 1

	for batchSize < len(constellations) {
		constellations, batches = constellations[batchSize:], append(batches, constellations[0:batchSize:batchSize])
	}

	batches = append(batches, constellations)

	for _, b := range batches {
		batch := b
		waitgroup.Add(1)
		go func() {
			for _, r := range batch {
				constellationURL := fmt.Sprintf(urlConstellationSpecifc, r)
				constellationBody, err := client.MakeESIGet(constellationURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download constellation %v - got; %v\n", r, string(constellationBody))

					return
				}

				constellation := ESIConstellation{}
				err = json.Unmarshal(constellationBody, &constellation)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode constellation %v - got; %v\n", r, string(constellationBody))
					return
				}

				// client.Log.Printf("Adding Constellation - %v\n", constellation.Name)
				err = client.Store.InsertConstellation(constellation)
				if err != nil {
					log.Fatalln(errors.Wrap(err, "Failed to insert constellation"))
				}
			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateSystems(client *Client) error {

	var waitgroup sync.WaitGroup

	// Now grab all the systems. Lets definetly batch these out and do 50 goroutines... Dont want to go too fast...
	const urlSystems = "https://esi.evetech.net/latest/universe/systems/?datasource=tranquility"
	systemsBody, err := client.MakeESIGet(urlSystems)
	var systems []int
	err = json.Unmarshal(systemsBody, &systems)
	if err != nil {
		return err
	}

	client.Log.Printf("Have to get %v systems", len(systems))

	const urlSystemSpecifc = "https://esi.evetech.net/latest/universe/systems/%v/?datasource=tranquility"

	batches := [][]int{}

	batchSize := (len(systems) / client.MaxRoutines) + 1

	for batchSize < len(systems) {
		systems, batches = systems[batchSize:], append(batches, systems[0:batchSize:batchSize])
	}

	batches = append(batches, systems)

	for _, b := range batches {
		batch := b
		waitgroup.Add(1)
		go func() {
			for _, r := range batch {
				systemURL := fmt.Sprintf(urlSystemSpecifc, r)
				systemBody, err := client.MakeESIGet(systemURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download System %v - got; %v; %v\n", r, string(systemBody), err)

					return
				}

				system := ESISystem{}
				err = json.Unmarshal(systemBody, &system)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode System %v - got; %v; %v\n", r, string(systemBody), err)
					client.Log.Println(err)
					return
				}

				// client.Log.Printf("Adding System - %v\n", system.Name)
				err = client.Store.InsertSystem(system)

				if err != nil {
					client.Log.Printf("dup - %v\n", system.SystemID)
				}

			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateStars(client *Client) error {

	var waitgroup sync.WaitGroup

	// All of the following will use parts from the systems objects so lets get all the systems here
	systemlist, err := client.Store.GetSystems()
	if err != nil {
		return err
	}

	// Now that we have got all the systems in place we can get all the information from them regarding stars
	var starIDComplete []int
	for _, sys := range systemlist {
		starIDComplete = append(starIDComplete, sys.StarID)
	}

	// Prevent duplicates
	starIDs := uniqueIDs(starIDComplete)

	client.Log.Printf("Have to get %v stars", len(starIDs))

	const urlStar = "https://esi.evetech.net/v1/universe/stars/%v/?datasource=tranquility"

	batches := [][]int{}

	batchSize := (len(starIDs) / client.MaxRoutines) + 1

	for batchSize < len(starIDs) {
		starIDs, batches = starIDs[batchSize:], append(batches, starIDs[0:batchSize:batchSize])
	}

	batches = append(batches, starIDs)

	for _, b := range batches {
		batch := b
		waitgroup.
			Add(1)
		go func() {
			for _, r := range batch {
				if r == 0 {
					// There are 250 of these.......
					continue
				}
				starURL := fmt.Sprintf(urlStar, r)
				starBody, err := client.MakeESIGet(starURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download star %v - got; %v; %v\n", r, string(starBody), err)

					return
				}

				star := ESIStar{}
				err = json.Unmarshal(starBody, &star)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode star %v - got; %v; %v\n", r, string(starBody), err)
					client.Log.Println(err)
					return
				}

				star.StarID = r

				// client.Log.Printf("Adding star - %v\n", star.Name)
				err = client.Store.InsertStar(star)

				if err != nil {
					client.Log.Printf("dup - %v\n", star.StarID)
				}

			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populatePlanets(client *Client) error {

	// Declare this early as we are going to use it a lot....
	var waitgroup sync.WaitGroup

	// All of the following will use parts from the systems objects so lets get all the systems here
	systemlist, err := client.Store.GetSystems()
	if err != nil {
		return err
	}

	// So the suns are done.. Now get ids for all planets, moons and asteroid belts...
	var planetList []int
	for _, sys := range systemlist {
		if sys.Planets != nil {
			for _, planet := range sys.Planets {
				planetList = append(planetList, planet.PlanetID)
			}
		}
	}
	// Now lets scrape the planets!!

	client.Log.Printf("Have to get %v planets", len(planetList))

	const urlPlanets = "https://esi.evetech.net/v1/universe/planets/%v/?datasource=tranquility"

	batches := [][]int{}

	batchSize := (len(planetList) / client.MaxRoutines) + 1

	for batchSize < len(planetList) {
		planetList, batches = planetList[batchSize:], append(batches, planetList[0:batchSize:batchSize])
	}

	batches = append(batches, planetList)

	for _, b := range batches {
		batch := b
		waitgroup.
			Add(1)
		go func() {
			for _, r := range batch {
				if r == 0 {
					// There are 250 of these.......
					continue
				}
				planetURL := fmt.Sprintf(urlPlanets, r)
				planetBody, err := client.MakeESIGet(planetURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download planet %v - got; %v; %v\n", r, string(planetBody), err)

					return
				}

				planetData := ESIPlanet{}
				err = json.Unmarshal(planetBody, &planetData)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode planet %v - got; %v; %v\n", r, string(planetBody), err)
					client.Log.Println(err)
					return
				}

				// client.Log.Printf("Adding planet - %v\n", planetData.Name)
				err = client.Store.InsertPlanet(planetData)

				if err != nil {
					client.Log.Printf("dup - %v; %v\n", planetData.PlanetID, err)
				}

			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateMoons(client *Client) error {

	systemlist, err := client.Store.GetSystems()
	if err != nil {
		return err
	}

	// Declare this early as we are going to use it a lot....
	var waitgroup sync.WaitGroup

	var moonList []int
	for _, sys := range systemlist {
		if sys.Planets != nil {
			for _, planet := range sys.Planets {
				if planet.Moons != nil {
					for _, moon := range planet.Moons {
						moonList = append(moonList, moon)
					}
				}
			}
		}
	}

	client.Log.Printf("Have to get %v moons", len(moonList))

	// THATS NO MOON!!!

	const urlMoon = "https://esi.evetech.net/v1/universe/moons/%v/?datasource=tranquility"

	batches := [][]int{}

	batchSize := (len(moonList) / client.MaxRoutines) + 1

	for batchSize < len(moonList) {
		moonList, batches = moonList[batchSize:], append(batches, moonList[0:batchSize:batchSize])
	}

	batches = append(batches, moonList)

	for _, b := range batches {
		batch := b
		waitgroup.
			Add(1)
		go func() {
			for _, r := range batch {
				if r == 0 {
					// There are 250 of these.......
					continue
				}
				moonURL := fmt.Sprintf(urlMoon, r)
				moonBody, err := client.MakeESIGet(moonURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download moon %v - got; %v; %v\n", r, string(moonBody), err)

					return
				}

				moon := ESIMoon{}
				err = json.Unmarshal(moonBody, &moon)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode planet %v - got; %v; %v\n", r, string(moonBody), err)
					client.Log.Println(err)
					return
				}

				// client.Log.Printf("Adding moon - %v\n", moon.Name)
				err = client.Store.InsertMoon(moon)

				if err != nil {
					client.Log.Printf("dup - %v; %v\n", moon.MoonID, err)
				}

			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateAsteroidBelts(client *Client) error {
	systemlist, err := client.Store.GetSystems()
	if err != nil {
		return err
	}

	// Declare this early as we are going to use it a lot....
	var waitgroup sync.WaitGroup

	var beltList []int
	for _, sys := range systemlist {
		if sys.Planets != nil {
			for _, planet := range sys.Planets {
				if planet.AsteroidBelts != nil {
					for _, belt := range planet.AsteroidBelts {
						beltList = append(beltList, belt)
					}
				}
			}
		}
	}

	client.Log.Printf("Have to get %v asteroid belts", len(beltList))

	const urlBelt = "https://esi.evetech.net/v1/universe/asteroid_belts/%v/?datasource=tranquility"

	batches := [][]int{}

	batchSize := (len(beltList) / client.MaxRoutines) + 1

	for batchSize < len(beltList) {
		beltList, batches = beltList[batchSize:], append(batches, beltList[0:batchSize:batchSize])
	}

	batches = append(batches, beltList)

	for _, b := range batches {
		batch := b
		waitgroup.
			Add(1)
		go func() {
			for _, r := range batch {
				if r == 0 {
					// There are 250 of these.......
					continue
				}
				beltURL := fmt.Sprintf(urlBelt, r)
				beltBody, err := client.MakeESIGet(beltURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download belt %v - got; %v; %v\n", r, string(beltBody), err)

					return
				}

				belt := ESIAsteroidBelt{}
				err = json.Unmarshal(beltBody, &belt)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode belt %v - got; %v; %v\n", r, string(beltBody), err)
					client.Log.Println(err)
					return
				}

				belt.BeltID = int32(r)
				// client.Log.Printf("Adding moon - %v\n", moon.Name)
				err = client.Store.InsertAsteroidBelt(belt)

				if err != nil {
					client.Log.Printf("dup - %v; %v\n", belt.BeltID, err)
				}

			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateStargates(client *Client) error {
	systemlist, err := client.Store.GetSystems()
	if err != nil {
		return err
	}

	// Declare this early as we are going to use it a lot....
	var waitgroup sync.WaitGroup

	var gateList []int
	for _, sys := range systemlist {
		if sys.Stargates != nil {
			for _, gate := range sys.Stargates {
				gateList = append(gateList, gate)
			}
		}
	}

	client.Log.Printf("Have to get %v stargates", len(gateList))

	const urlGate = "https://esi.evetech.net/v1/universe/stargates/%v/?datasource=tranquility"

	batches := [][]int{}

	batchSize := (len(gateList) / client.MaxRoutines) + 1

	for batchSize < len(gateList) {
		gateList, batches = gateList[batchSize:], append(batches, gateList[0:batchSize:batchSize])
	}

	batches = append(batches, gateList)

	for _, b := range batches {
		batch := b
		waitgroup.
			Add(1)
		go func() {
			for _, r := range batch {
				if r == 0 {
					continue
				}
				gateURL := fmt.Sprintf(urlGate, r)
				gateBody, err := client.MakeESIGet(gateURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download gate %v - got; %v; %v\n", r, string(gateBody), err)

					return
				}

				gate := ESIStargate{}
				err = json.Unmarshal(gateBody, &gate)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode gate %v - got; %v; %v\n", r, string(gateBody), err)
					client.Log.Println(err)
					return
				}

				// client.Log.Printf("Adding moon - %v\n", moon.Name)
				err = client.Store.InsertStargate(gate)

				if err != nil {
					client.Log.Printf("dup - %v; %v\n", gate.StargateID, err)
				}

			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateStations(client *Client) error {
	systemlist, err := client.Store.GetSystems()
	if err != nil {
		return err
	}

	// Declare this early as we are going to use it a lot....
	var waitgroup sync.WaitGroup

	var stationList []int
	for _, sys := range systemlist {
		if sys.Stations != nil {
			for _, station := range sys.Stations {
				stationList = append(stationList, station)
			}
		}
	}

	client.Log.Printf("Have to get %v stations", len(stationList))

	const urlStations = "https://esi.evetech.net/v2/universe/stations/%v/?datasource=tranquility"

	batches := [][]int{}

	batchSize := (len(stationList) / client.MaxRoutines) + 1

	for batchSize < len(stationList) {
		stationList, batches = stationList[batchSize:], append(batches, stationList[0:batchSize:batchSize])
	}

	batches = append(batches, stationList)

	for _, b := range batches {
		batch := b
		waitgroup.
			Add(1)
		go func() {
			for _, r := range batch {
				if r == 0 {
					continue
				}
				stationURL := fmt.Sprintf(urlStations, r)
				stationBody, err := client.MakeESIGet(stationURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download station %v - got; %v; %v\n", r, string(stationBody), err)

					return
				}

				station := ESIStation{}
				err = json.Unmarshal(stationBody, &station)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode station %v - got; %v; %v\n", r, string(stationBody), err)
					client.Log.Println(err)
					return
				}

				err = client.Store.InsertStation(station)

				if err != nil {
					client.Log.Printf("dup - %v; %v\n", station.StationID, err)
				}

			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateTypes(client *Client) error {
	var waitgroup sync.WaitGroup

	// Now grab all the types
	const urlTypes = "https://esi.evetech.net/v1/universe/types/?datasource=tranquility&page=%v"
	var types []int

	page := 1
	for {
		url := fmt.Sprintf(urlTypes, page)
		typesBody, err := client.MakeESIGet(url)
		if string(typesBody) == "[]" {
			break
		}
		var t []int
		err = json.Unmarshal(typesBody, &t)
		if err != nil {
			return err
		}

		// client.Log.Printf("Getting types page %v\n", page)
		types = append(types, t...)

		page++
	}

	client.Log.Printf("Have to get %v types from ESI", len(types))

	const urlTypeSpecifc = "https://esi.evetech.net/v3/universe/types/%v/?datasource=tranquility"

	var batches [][]int

	// Because there are so many typeids to fetch, going to double the number of goroutines
	batchSize := (len(types) / (client.MaxRoutines * 2)) + 1

	for batchSize < len(types) {
		types, batches = types[batchSize:], append(batches, types[0:batchSize:batchSize])
	}

	batches = append(batches, types)

	for _, b := range batches {
		batch := b
		waitgroup.Add(1)
		go func() {
			for _, r := range batch {
				typeURL := fmt.Sprintf(urlTypeSpecifc, r)
				typeBody, err := client.MakeESIGet(typeURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download type %v - got; %v\n", r, string(typeBody))

					return
				}

				typeESI := ESIType{}
				err = json.Unmarshal(typeBody, &typeESI)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode type %v - got; %v\n", r, string(typeBody))
					return
				}

				// client.Log.Printf("Adding type - %v - %v\n", typeESI.TypeID, typeESI.Name)
				client.Store.InsertType(typeESI)
			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateGroups(client *Client) error {
	var waitgroup sync.WaitGroup

	// Now grab all the types
	const urlGroups = "https://esi.evetech.net/v1/universe/groups/?datasource=tranquility&page=%v"
	var groups []int

	page := 1
	for {
		url := fmt.Sprintf(urlGroups, page)
		groupBody, err := client.MakeESIGet(url)
		if string(groupBody) == "[]" {
			break
		}
		var t []int
		err = json.Unmarshal(groupBody, &t)
		if err != nil {
			return err
		}

		// client.Log.Printf("Getting groups page %v\n", page)
		groups = append(groups, t...)

		page++
	}

	client.Log.Printf("Have to get %v groups from ESI", len(groups))

	const urlGroupSpecifc = "https://esi.evetech.net/v1/universe/groups/%v/?datasource=tranquility"

	var batches [][]int

	// Because there are so many typeids to fetch, going to double the number of goroutines
	batchSize := (len(groups) / (client.MaxRoutines * 2)) + 1

	for batchSize < len(groups) {
		groups, batches = groups[batchSize:], append(batches, groups[0:batchSize:batchSize])
	}

	batches = append(batches, groups)

	for _, b := range batches {
		batch := b
		waitgroup.Add(1)
		go func() {
			for _, r := range batch {
				groupURL := fmt.Sprintf(urlGroupSpecifc, r)
				groupBody, err := client.MakeESIGet(groupURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download group %v - got; %v\n", r, string(groupBody))

					return
				}

				group := ESIGroup{}
				err = json.Unmarshal(groupBody, &group)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode group %v - got; %v\n", r, string(groupBody))
					return
				}

				// client.Log.Printf("Adding group - %v - %v\n", group.GroupID, group.Name)
				client.Store.InsertGroup(group)
			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()

	return nil
}

func populateCategories(client *Client) error {
	var waitgroup sync.WaitGroup

	const urlCategories = "https://esi.evetech.net/v1/universe/categories/?datasource=tranquility"
	categoriesBody, err := client.MakeESIGet(urlCategories)
	var categories []int
	err = json.Unmarshal(categoriesBody, &categories)
	if err != nil {
		return err
	}

	client.Log.Printf("Have to get %v categories from ESI", len(categories))

	const urlGroupSpecifc = "https://esi.evetech.net/v1/universe/categories/%v/?datasource=tranquility"

	var batches [][]int

	// Because there are so many typeids to fetch, going to double the number of goroutines
	batchSize := (len(categories) / (client.MaxRoutines * 2)) + 1

	for batchSize < len(categories) {
		categories, batches = categories[batchSize:], append(batches, categories[0:batchSize:batchSize])
	}

	batches = append(batches, categories)

	for _, b := range batches {
		batch := b
		waitgroup.Add(1)
		go func() {
			for _, r := range batch {
				categoryURL := fmt.Sprintf(urlGroupSpecifc, r)
				categoryBody, err := client.MakeESIGet(categoryURL)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to download category %v - got; %v\n", r, string(categoryBody))

					return
				}

				category := ESICategory{}
				err = json.Unmarshal(categoryBody, &category)
				if err != nil {
					waitgroup.Done()
					client.Log.Printf("Failed to decode category %v - got; %v\n", r, string(categoryBody))
					return
				}

				// client.Log.Printf("Adding category - %v - %v\n", category.CategoryID, category.Name)
				client.Store.InsertCategory(category)
			}
			waitgroup.Done()
		}()
	}

	waitgroup.Wait()
	return nil
}

func uniqueIDs(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
