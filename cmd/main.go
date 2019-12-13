package main

import (
	"log"

	"github.com/podded/higgs"

	"github.com/spf13/viper"

	"github.com/pkg/profile"
)

func main() {

	// Start out by reading in our config file

	defer profile.Start().Stop()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var config higgs.Configuration

	//Set some reasonable defaults

	config = higgs.Configuration{
		Database: higgs.DatabaseConfig{
			URI:        "mongodb://localhost:27017",
			Database:   "podded",
		},
		Web: higgs.HttpConfig{
			UserAgent:  "Crypta-Eve/Podded install (BUT I AM BAD AND HAVENT CHANGED DEFAULT UA)",
			TimeoutSec: 30,
		},
		App: higgs.AppConfig{
			MaxRoutines: 20,
		},
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config, does it exist? err: %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error interpreting the config, is it valid? err: %s", err)
	}

	if err := higgs.PopulateStaticData(config); err != nil {
		log.Fatalf("Error deleting static data. err: %s", err)
	}



}
