package diary

import (
	"github.com/gin-gonic/gin"

	"github.com/vvanm/gymtracker-go/food"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/util/jwt"
	"log"
	"strconv"
	"time"
)

type Diary struct {
	ID string
	//
	Owner      string `json:"owner,omitempty"`
	Date       int64  `json:"date,omitempty"`
	DateString string `json:"dateString,omitempty"`
	//
	Foods []food.Food `json:"foods"`
	//
	RoutineID string           `json:"routineID"`
	Exercises []Diary_Exercise `json:"exercises"`
}

type Diary_Exercise struct {
	ExerciseID string               `json:"exerciseID,omitempty"`
	Sets       []Diary_Exercise_Set `json:"sets"`
}

type Diary_Exercise_Set struct {
	Reps   int     `json:"reps,omitempty"`
	Weight float64 `json:"weight,omitempty"`
}

func GetDiary(c *gin.Context) {
	claims := jwt.ClaimsFromCookie(c)
	date := c.Param("date")
	diaryKey := claims.ID + "/diaries/" + date

	var d *Diary

	session, err := raven.Store.OpenSession()
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	err = session.Load(&d, diaryKey)

	if d != nil {
		c.JSON(200, d)
		return
	}

	epoInt, _ := strconv.ParseInt(date, 10, 64)

	newDiary := &Diary{
		Date:       epoInt,
		DateString: time.Unix(epoInt/1000, 0).Format(time.RFC1123Z),
		Owner:      claims.ID,
		RoutineID:  "",
		Foods:      []food.Food{},
		Exercises:  []Diary_Exercise{},
	}

	_ = session.StoreWithID(newDiary, diaryKey)

	_ = session.SaveChanges()

	c.JSON(200, newDiary)

}
