package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/okozzie/gogogo/initializers"
	"github.com/okozzie/gogogo/models"
)

func Index(c *gin.Context) {

	/*
		Retrieve all Ships
		Filtering can be applied through URL params
	*/

	var ships []models.Ship
	chain := initializers.DB.Preload("Armaments")
	//initializers.DB.Preload("Armaments").Find(&ships)

	if c.Query("name") != "" {
		//Filter requests by name
		chain = chain.Where("name = ?", c.Query("name"))
	}

	if c.Query("class") != "" {
		//Filter requests by class
		chain = chain.Where("class = ?", c.Query("class"))
	}

	if c.Query("status") != "" {
		//Filter requests by status
		chain = chain.Where("status = ?", c.Query("status"))
	}

	//Retrieve records
	chain.Find(&ships)

	/*
		Pivot table contains an additional Quantity column so that has to be retrieved separately

		This loops through any armaments a ship may have and finds the corresponding Quantity
	*/
	for i, ship := range ships {
		for j, armament := range ship.Armaments {

			var tmpArmament models.ShipArmament
			initializers.DB.Table("ship_armaments").Where("armament_id = ?", armament.ID).Where("ship_id = ?", ship.ID).First(&tmpArmament)
			ships[i].Armaments[j].Quantity = tmpArmament.Quantity
		}
	}

	c.JSON(http.StatusOK, gin.H{"ships": ships})

}

func Show(c *gin.Context) {

	var ship models.Ship
	id := c.Param("id") //get ID from URL

	//Check if the Ship exists
	err := initializers.DB.Where("id = ?", id).Preload("Armaments").First(&ship).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	/*
		Pivot table contains an additional Quantity column so that has to be retrieved separately

		This loops through any armaments a ship may have and finds the corresponding Quantity
	*/
	for i, armament := range ship.Armaments {
		var tmpArmament models.ShipArmament
		initializers.DB.Table("ship_armaments").Where("armament_id = ?", armament.ID).Where("ship_id = ?", ship.ID).First(&tmpArmament)
		ship.Armaments[i].Quantity = tmpArmament.Quantity
	}

	c.JSON(http.StatusOK, ship)
}

func Store(c *gin.Context) {
	var ship models.Ship

	//Validate the payload based on the Ship struct
	if err := c.BindJSON(&ship); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "invalid payload"})
		return
	}

	//Retrieve all the armament names passed, so the models can be retrieved later on by name
	var armamentNames []string
	for _, armament := range ship.Armaments {
		armamentNames = append(armamentNames, armament.Name)
	}

	//Retrieve all the armaments passed based on their names
	var armaments []models.Armament
	initializers.DB.Where("name in (?)", armamentNames).Find(&armaments)

	//Create the new ship
	initializers.DB.Omit("Armaments").Create(&ship)

	//Insert into the Pivot table
	for i, armament := range armaments {
		shipArmament := &models.ShipArmament{
			ShipID:     ship.ID,
			ArmamentID: armament.ID,
			Quantity:   ship.Armaments[i].Quantity,
		}
		initializers.DB.Create(shipArmament)
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})

}

func Update(c *gin.Context) {
	var ship models.Ship
	id := c.Param("id")

	err := initializers.DB.Where("id = ?", id).Preload("Armaments").First(&ship).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	if err := c.BindJSON(&ship); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "invalid payload"})
		return
	}

	//Retrieve all the armament names passed, so the models can be retrieved later on
	var armamentNames []string
	for _, armament := range ship.Armaments {
		armamentNames = append(armamentNames, armament.Name)
	}

	//Retrieve all the armaments passed based on their names
	var armaments []models.Armament
	initializers.DB.Where("name in (?)", armamentNames).Find(&armaments)

	//Delete the associated records in pivot table
	initializers.DB.Table("ship_armaments").Where("ship_id = ?", id).Unscoped().Delete(&models.ShipArmament{})

	//Insert Armament associations into the Pivot table
	for i, armament := range armaments {
		shipArmament := &models.ShipArmament{
			ShipID:     ship.ID,
			ArmamentID: armament.ID,
			Quantity:   ship.Armaments[i].Quantity,
		}
		initializers.DB.Create(shipArmament)
	}

	//Update the Ship
	initializers.DB.Omit("Armaments").Model(&ship).Where("id = ?", id).Updates(&ship)

	c.JSON(http.StatusOK, gin.H{"message": "success"})

}

func Delete(c *gin.Context) {
	var ship models.Ship
	id := c.Param("id")

	err := initializers.DB.Where("id = ?", id).Preload("Armaments").First(&ship).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	//First delete the associated records in pivot table
	initializers.DB.Table("ship_armaments").Where("ship_id = ?", id).Unscoped().Delete(&models.ShipArmament{})

	//Now delete the Ship
	initializers.DB.Delete(&ship)

	c.JSON(200, gin.H{"message": "successfully deleted"})

}
