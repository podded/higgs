package higgs

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type (

	ESIPosition struct {
		X float64 `json:"x" bson:"x"`
		Y float64 `json:"y" bson:"y"`
		Z float64 `json:"z" bson:"z"`
	}

	ESIRegion struct {
		Constellations []int  `json:"constellations" bson:"constellations"`
		Description    string `json:"description" bson:"description"`
		Name           string `json:"name" bson:"name"`
		RegionID       int    `json:"region_id" bson:"_id"`
	}

	ESIConstellation struct {
		ConstellationID int         `json:"constellation_id" bson:"_id"`
		Name            string      `json:"name" bson:"name"`
		Systems         []int       `json:"systems" bson:"systems"`
		Postion         ESIPosition `json:"position" bson:"position"`
		RegionID        int         `json:"region_id" bson:"region_id"`
	}

	ESISystem struct {
		ConstellationID int                `json:"constellation_id" bson:"constellation_id"`
		Name            string             `json:"name" bson:"name"`
		Planets         []ESISystemPlanets `json:"planets" bson:"planets"`
		Position        ESIPosition        `json:"position" bson:"position"`
		SecurityClass   string             `json:"security_class" bson:"security_class"`
		SecurityStatus  float64            `json:"security_status" bson:"security_status"`
		StarID          int                `json:"star_id" bson:"star_id"`
		Stargates       []int              `json:"stargates" bson:"stargates"`
		Stations        []int              `json:"stations" bson:"stations"`
		SystemID        int                `json:"system_id" bson:"_id"`
	}

	ESISystemPlanets struct {
		PlanetID      int   `json:"planet_id" bson:"planet_id"`
		Moons         []int `json:"moons" bson:"moons,omitempty"`
		AsteroidBelts []int `json:"asteroid_belts" bson:"asteroid_belts,omitempty"`
	}

	ESIStar struct {
		Age           int64   `json:"age" bson:"age"`
		Luminosity    float64 `json:"luminosity" bson:"luminosity"`
		Name          string  `json:"name" bson:"name"`
		Radius        int64   `json:"radius" bson:"radius"`
		SolarSystemID int     `json:"solar_system_id" bson:"solar_system_id"`
		SpectralClass string  `json:"spectral_class" bson:"spectral_class"`
		Temperature   int     `json:"temperature" bson:"temperature"`
		TypeID        int     `json:"type_id" bson:"type_id"`
		StarID        int     `json:"star_id,omitempty" bson:"_id"`
	}

	ESIPlanet struct {
		Name     string      `json:"name" bson:"name"`
		PlanetID int32       `json:"planet_id" bson:"_id"`
		Position ESIPosition `json:"position" bson:"position"`
		SystemID int32       `json:"system_id" bson:"system_id"`
		TypeID   int32       `json:"type_id" bson:"type_id"`
	}

	ESIMoon struct {
		MoonID   int32       `json:"moon_id" bson:"_id"`
		Name     string      `json:"name" bson:"name"`
		Position ESIPosition `json:"position" bson:"position"`
		SystemID int32       `json:"system_id" bson:"system_id"`
	}

	ESIAsteroidBelt struct {
		BeltID   int32       `json:"belt_id,omitempty" bson:"_id"`
		Name     string      `json:"name" bson:"name"`
		Position ESIPosition `json:"position" bson:"position"`
		SystemID int32       `json:"system_id" bson:"system_id"`
	}

	ESIStargate struct {
		Destination ESIStargateDestination `json:"destination" bson:"destination"`
		Name        string                 `json:"name" bson:"name"`
		Position    ESIPosition            `json:"position" bson:"position"`
		StargateID  int32                  `json:"stargate_id" bson:"_id"`
		SystemID    int32                  `json:"system_id" bson:"system_id"`
		TypeID      int32                  `json:"type_id" bson:"type_id"`
	}

	ESIStargateDestination struct {
		StargateID int32 `json:"stargate_id" bson:"stargate_id"`
		SystemID   int32 `json:"system_id" bson:"system_id"`
	}

	ESIStation struct {
		MaxDockableShipVolume  float64     `json:"max_dockable_ship_volume" bson:"max_dockable_ship_volume"`
		Name                   string      `json:"name" bson:"name"`
		OfficeRentalCost       float64     `json:"office_rental_cost" bson:"office_rental_cost"`
		Owner                  int32       `json:"owner" bson:"owner"`
		Position               ESIPosition `json:"position" bson:"position"`
		RaceID                 int32       `json:"race_id" bson:"race_id"`
		ReprocessingEfficiency float32     `json:"reprocessing_efficiency" bson:"reprocessing_efficiency"`
		Services               []string    `json:"services" bson:"services"`
		StationID              int32       `json:"station_id" bson:"_id"`
		SystemID               int32       `json:"system_id" bson:"system_id"`
		TypeID                 int32       `json:"type_id" bson:"type_id"`
	}

	ESIType struct {
		Capacity        float64              `json:"capacity,omitempty" bson:"capacity,omitempty"`
		Description     string               `json:"description" bson:"description"`
		DogmaAttributes []TypeDogmaAttribute `json:"dogma_attributes,omitempty" bson:"dogma_attributes,omitempty"`
		DogmaEffects    []TypeDogmaEffect    `json:"dogma_effects,omitempty" bson:"dogma_effects,omitempty"`
		GraphicID       int32                `json:"graphic_id,omitempty" bson:"graphic_id,omitempty"`
		GroupID         int32                `json:"group_id" bson:"group_id"`
		IconID          int32                `json:"icon_id,omitempty" bson:"icon_id,omitempty"`
		MarketGroupID   int32                `json:"market_group_id,omitempty" bson:"market_group_id,omitempty"`
		Mass            float64              `json:"mass,omitempty" bson:"mass,omitempty"`
		Name            string               `json:"name" bson:"name"`
		PackagedVolume  float64              `json:"packaged_volume,omitempty" bson:"packaged_volume,omitempty"`
		PortionSize     int32                `json:"portion_size,omitempty" bson:"portion_size,omitempty"`
		Published       bool                 `json:"published" bson:"published"`
		Radius          float64              `json:"radius,omitempty" bson:"radius,omitempty"`
		TypeID          int32                `json:"type_id" bson:"_id"`
		Volume          float64              `json:"volume,omitempty" bson:"volume,omitempty"`
	}

	TypeDogmaAttribute struct {
		AttributeID int32   `json:"attribute_id" bson:"attribute_id"`
		Value       float64 `json:"value" bson:"value"`
	}

	TypeDogmaEffect struct {
		EffectID  int32 `json:"effect_id" bson:"effect_id"`
		IsDefault bool  `json:"is_default" bson:"is_default"`
	}

	ESIGroup struct {
		CategoryID int32   `json:"category_id" bson:"category_id"`
		GroupID    int32   `json:"group_id" bson:"_id"`
		Name       string  `json:"name" bson:"name"`
		Published  bool    `json:"published" bson:"published"`
		Types      []int32 `json:"types" bson:"types"`
	}

	ESICategory struct {
		CategoryID int32   `json:"category_id" bson:"_id"`
		Groups     []int32 `json:"groups" bson:"groups"`
		Name       string  `json:"name" bson:"name"`
		Published  bool    `json:"published" bson:"published"`
	}
)


func (db *DB) DeleteStaticData() error {

	// I know this is bad but I really dont care about errors here for now

	collection := db.Database.Database(db.DBName).Collection("regions")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("constellations")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("solarsystems")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("stars")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("planets")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("moons")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("asteroid_belts")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("stargates")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("stations")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("categories")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("groups")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("types")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("regions")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("ancestries")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("bloodlines")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database(db.DBName).Collection("factions")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	return nil

}

func (db *DB) InsertRegion(region ESIRegion) error {

	collection := db.Database.Database(db.DBName).Collection("regions")

	_, err := collection.InsertOne(context.TODO(), region)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve region")
	}

	return nil

}

func (db *DB) InsertConstellation(cons ESIConstellation) error {

	collection := db.Database.Database(db.DBName).Collection("constellations")

	_, err := collection.InsertOne(context.TODO(), cons)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve constellation")
	}

	return nil

}

func (db *DB) InsertSystem(system ESISystem) error {

	collection := db.Database.Database(db.DBName).Collection("solarsystems")

	_, err := collection.InsertOne(context.TODO(), system)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve system")
	}

	return nil

}

func (db *DB) InsertStar(star ESIStar) error {

	collection := db.Database.Database(db.DBName).Collection("stars")

	_, err := collection.InsertOne(context.TODO(), star)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve star")
	}

	return nil

}

func (db *DB) InsertPlanet(planet ESIPlanet) error {

	collection := db.Database.Database(db.DBName).Collection("planets")

	_, err := collection.InsertOne(context.TODO(), planet)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve planet")
	}

	return nil

}

func (db *DB) InsertMoon(moon ESIMoon) error {

	collection := db.Database.Database(db.DBName).Collection("moons")

	_, err := collection.InsertOne(context.TODO(), moon)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve moon")
	}

	return nil

}

func (db *DB) InsertAsteroidBelt(belt ESIAsteroidBelt) error {

	collection := db.Database.Database(db.DBName).Collection("asteroid_belts")

	_, err := collection.InsertOne(context.TODO(), belt)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve asteroid belt")
	}

	return nil

}

func (db *DB) InsertStargate(gate ESIStargate) error {

	collection := db.Database.Database(db.DBName).Collection("stargates")

	_, err := collection.InsertOne(context.TODO(), gate)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve stargate")
	}

	return nil

}

func (db *DB) InsertStation(station ESIStation) error {

	collection := db.Database.Database(db.DBName).Collection("stations")

	_, err := collection.InsertOne(context.TODO(), station)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve station")
	}

	return nil
}

func (db *DB) InsertType(typeESI ESIType) error {

	collection := db.Database.Database(db.DBName).Collection("types")

	_, err := collection.InsertOne(context.TODO(), typeESI)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve type")
	}

	return nil
}

func (db *DB) InsertGroup(group ESIGroup) error {

	collection := db.Database.Database(db.DBName).Collection("groups")

	_, err := collection.InsertOne(context.TODO(), group)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve group")
	}

	return nil
}

func (db *DB) InsertCategory(category ESICategory) error {

	collection := db.Database.Database(db.DBName).Collection("categories")

	_, err := collection.InsertOne(context.TODO(), category)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve category")
	}

	return nil
}

func (db *DB) GetSystems() (systems []ESISystem, err error) {
	collection := db.Database.Database(db.DBName).Collection("solarsystems")

	ctx := context.Background()

	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return systems, errors.Wrap(err, "error retrieving existing systems")
	}

	defer c.Close(ctx)

	for c.Next(ctx) {

		var id ESISystem

		err := c.Decode(&id)
		if err != nil {
			return systems, errors.Wrap(err, "Failed to morp system into struct")
		}

		systems = append(systems, id)
	}

	return systems, nil
}



