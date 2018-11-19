package user

import (
	"errors"
	"github.com/vvanm/gymtracker-go/exercise"
	"github.com/vvanm/gymtracker-go/food"
	"github.com/vvanm/gymtracker-go/raven"
	"github.com/vvanm/gymtracker-go/routine"
	"github.com/vvanm/gymtracker-go/util/auth"
)

type User struct {
	ID       string
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	//
	Foods     []food.Food         `json:"foods"`
	Routines  []routine.Routine   `json:"routines"`
	Exercises []exercise.Exercise `json:"exercises"`
}

func GetUser(q string) (u *User, err error) {
	session, err := raven.Store.OpenSession()
	if err != nil {
		return
	}
	defer session.Close()

	var users []*User
	err = session.Advanced().RawQuery(q).ToList(&users)

	if err != nil {
		return
	}

	if len(users) == 1 {
		u = users[0]
	}

	return

}

func (u *User) Create() error {
	if u.Email == "" || u.Password == "" {
		return errors.New("Empty inputs")
	}

	//open session
	session, err := raven.Store.OpenSession()
	if err != nil {
		return err
	}
	defer session.Close()

	//verify unique email
	var users []*User
	err = session.Advanced().RawQuery("from Users where email == '" + u.Email + "'").ToList(&users)
	if err != nil {
		return err
	}

	if len(users) > 0 {
		return errors.New("Email in use")
	}

	//Process the password
	u.Password = auth.SecurePassword(u.Password)

	//add user
	err = session.StoreWithID(u, "users|")
	if err != nil {
		return err
	}

	//push to raven
	err = session.SaveChanges()
	if err != nil {
		return err
	}

	return nil

}
