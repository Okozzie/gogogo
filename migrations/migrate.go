package main

import (
	"github.com/okozzie/gogogo/initializers"
	"github.com/okozzie/gogogo/models"
)

func init() {
	//Connect to the DB
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
}

func main() {

	//Migrate the models
	initializers.DB.AutoMigrate(&models.Ship{}, &models.Armament{}, &models.ShipArmament{})

	//Seed the database with armaments
	var armamentSeeder = []models.Armament{
		{Name: "Turbo Lasers"},
		{Name: "Ion Cannons"},
		{Name: "Laser Beams"},
	}
	initializers.DB.Create(&armamentSeeder)

	//Grab the 1st and 2nd seeded Armaments -- Turbo Lasers and Ion Cannons
	var armaments []models.Armament
	initializers.DB.Where("id in (?)", []int{1, 2}).Find(&armaments)

	ship := &models.Ship{
		Name:   "New Ship",
		Class:  "Class 1",
		Crew:   50,
		Image:  "image.jpg",
		Value:  100.0,
		Status: "active",
	}

	//Create the initial Ship
	initializers.DB.Create(ship)

	//Give that ship armaments by inserting into the Pivot Table
	for _, armament := range armaments {
		shipArmament := &models.ShipArmament{
			ShipID:     ship.ID,
			ArmamentID: armament.ID,
			Quantity:   20,
		}
		initializers.DB.Create(shipArmament)
	}

}
