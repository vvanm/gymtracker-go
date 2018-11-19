package food

import (
	"github.com/gin-gonic/gin"
	"github.com/ravendb/ravendb-go-client"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/util/helpers"
	"github.com/vvanm/gymtracker-go/util/jwt"
)

type Food struct {
	FoodID string `json:"foodID,omitempty"`
	//
	Name    string  `json:"name,omitempty"`
	Cal     float64 `json:"cal"`
	Carb    float64 `json:"carb"`
	Protein float64 `json:"protein"`
	Fat     float64 `json:"fat"`
	Amount  float64 `json:"amount"`
}

func Update(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)

	var f Food
	c.ShouldBind(&f)

	f.FoodID = c.Param("foodID")

	//create patch
	patch := ravendb.PatchRequest_forScript(
		`this.entries = this.entries.map(entry => {
				if(entry.foodID == args.payload.foodID){
					entry = args.payload
				}
				return entry
			})`,
	)

	//add values
	patch.SetValues(map[string]interface{}{
		"payload": f,
	},
	)

	patchOp := ravendb.NewPatchOperation(claims.ID+"/foods", nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, err.Error())
		return
	}

	c.JSON(200, gin.H{})

}

func Delete(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)

	patch := ravendb.PatchRequest_forScript(`this.entries = this.entries.filter(e => e.foodID != args.foodID)`)
	patch.SetValues(map[string]interface{}{
		"foodID": c.Param("foodID"),
	},
	)

	patchOp := ravendb.NewPatchOperation(claims.ID+"/foods", nil, patch, nil, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{})

}

func AddNew(c *gin.Context) {
	//get claims
	claims := jwt.ClaimsFromCookie(c)

	var f Food
	c.ShouldBind(&f)

	f.FoodID = helpers.UUID()

	//create payload
	payload := map[string]interface{}{
		"payload": f,
	}

	//create patches
	patch := ravendb.PatchRequest_forScript(`this.entries.push(args.payload)`)
	patchIfMissing := ravendb.PatchRequest_forScript(`this.entries = [args.payload];this["@metadata"] = {"@collection" : "foods"};`)

	//add values
	patch.SetValues(payload)
	patchIfMissing.SetValues(payload)

	//execute
	patchOp := ravendb.NewPatchOperation(claims.ID+"/foods", nil, patch, patchIfMissing, false)
	err := raven.Store.Operations().Send(patchOp)

	if err != nil {
		c.JSON(400, gin.H{})
		return
	}

	c.JSON(200, gin.H{"foodID": f.FoodID})
	return

}
