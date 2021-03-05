package models

import (
	"encoding/json"
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"

	"golang.org/x/crypto/bcrypt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type User struct {
	ID        string    `rethinkdb:"id, omitempty" json:"id"`
	Nickname  string    `rethinkdb:"nickname" json:"nickname"`
	Email     string    `rethinkdb:"email" json:"email"`
	Password  string    `rethinkdb:"password" json:"password,omitempty"`
	CreatedAt time.Time `rethinkdb:"created_at" json:"created_at"`
	UpdatedAt time.Time `rethinkdb:"updated_at" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// MarshalJSON overridden in order to omit the password when outputting TO json, but not when RECEIVING json
func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Password string `json:"password,omitempty"`
		*Alias
	}{
		Password: "",
		Alias:    (*Alias)(&u),
	})
}

func (u *User) beforeNewUser() error {
	hashedPassord, err := Hash(u.Password)

	if err != nil {
		return err
	}
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	u.Password = string(hashedPassord)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}

		if u.Email == "" {
			return errors.New("Required Email")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil

	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}

		if u.Email == "" {
			return errors.New("Required Email")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}

		if u.Password == "" {
			return errors.New("Required Password")
		}

		if u.Email == "" {
			return errors.New("Required Email")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	}
}

// NewUser creates a new user
func (u *User) NewUser(dbname string, session *r.Session) (*User, error) {
	// check for conflict
	if !u.canCreateUser(dbname, session) {
		return nil, errors.New("Cannot create user; duplicate nickname or email address")
	}

	u.beforeNewUser()

	uuidResp, err := r.UUID(time.Now().String()).Run(session)

	var row interface{}

	err = uuidResp.One(&row)
	if err != nil {
		return nil, err
	}

	uuid, ok := row.(string)

	if !ok {
		return nil, errors.New("Failed to generate UUID")
	}

	u.ID = uuid
	resp, err := r.DB(dbname).Table("users").Insert(u).RunWrite(session)

	if err != nil || resp.Inserted < 1 {
		return nil, err
	}

	return u, nil
}

func (u *User) canCreateUser(dbname string, session *r.Session) bool {
	resp, err := r.DB(dbname).Table("users").Filter(r.Row.Field("nickname").Eq(u.Nickname).Or(r.Row.Field("email").Eq(u.Email))).Run(session)

	if err != nil {
		return false
	}

	existingUsers := &[]User{}

	if err = resp.All(existingUsers); err != nil {
		return false
	}

	return len(*existingUsers) == 0
}

func (u *User) CanUpdateUser(nicknameChanged, emailChanged bool, dbName string, session *r.Session) bool {
	var term r.Term

	if nicknameChanged && emailChanged {
		term = r.Row.Field("nickname").Eq(u.Nickname).Or(r.Row.Field("email").Eq(u.Email))
	} else if nicknameChanged && !emailChanged {
		term = r.Row.Field("nickname").Eq(u.Nickname)
	} else if emailChanged && !nicknameChanged {
		term = r.Row.Field("email").Eq(u.Email)
	} else {
		return true
	}

	resp, err := r.DB(dbName).Table("users").Filter(term).Run(session)

	if err != nil {
		log.Println(err)

		return false
	}

	existingUsers := &[]User{}

	if err = resp.All(existingUsers); err != nil {
		log.Println(err)
		return false
	}

	return len(*existingUsers) == 0
}

// FindAllUsers returns a list of all users in DB
func (u *User) FindAllUsers(dbname string, session *r.Session) (*[]User, error) {
	users := &[]User{}

	resp, err := r.DB(dbname).Table("users").Run(session)

	if err != nil {
		return nil, err
	}

	if err = resp.All(users); err != nil {
		return nil, err
	}

	return users, nil
}

// FindUserByID will get a user either by uuid or email
func (u *User) FindUserByID(dbName string, session *r.Session, id string) (*User, error) {
	if err := checkmail.ValidateFormat(id); err == nil {
		return u.FindUserByEmail(dbName, session, id)
	}

	resp, err := r.DB(dbName).Table("users").Get(id).Run(session)

	if err != nil {
		return nil, err
	}

	user := &User{}

	err = resp.One(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByEmail will get a user by email
func (u *User) FindUserByEmail(dbName string, session *r.Session, email string) (*User, error) {
	resp, err := r.DB(dbName).Table("users").Filter(r.Row.Field("email").Eq(email)).Run(session)

	if err != nil {
		return nil, err
	}

	user := &User{}

	err = resp.One(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateAUser updates a user
func (u *User) UpdateAUser(dbName string, session *r.Session, id string) (*User, error) {
	user, err := u.FindUserByID(dbName, session, id)

	if err != nil {
		return nil, errors.New("User doesn't exist")
	}

	nicknameChanged := u.Nickname != user.Nickname
	emailChanged := u.Email != user.Email

	user.Nickname = u.Nickname
	user.Email = u.Email
	user.UpdatedAt = time.Now()

	if !user.CanUpdateUser(nicknameChanged, emailChanged, dbName, session) {
		return nil, errors.New("duplicate nickname or email address")
	}

	resp, err := r.DB(dbName).Table("users").Get(user.ID).Update(user).RunWrite(session)

	if err != nil || (resp.Updated < 1 && resp.Replaced < 1) {
		return nil, err
	}

	return user, nil
}

// DeleteAUser Deletes a user
func (u *User) DeleteAUser(dbName string, session *r.Session, id string) (int64, error) {
	user, err := u.FindUserByID(dbName, session, id)

	if err != nil {
		return 0, errors.New("User doesn't exist")
	}

	resp, err := r.DB(dbName).Table("user").Get(user.ID).Delete().RunWrite(session)

	if err != nil || resp.Deleted < 1 {
		return 0, err
	}

	return 1, nil
}

func (u *User) GetLevelIds(dbname string, session *r.Session) ([]TileMap, error) {
	var levels []TileMap

	resp, err := r.DB(dbname).Table("levels").Filter(r.Row.Field("user_id").Eq(u.ID)).Pluck("id", "meta_data").Run(session)

	if err != nil {
		return nil, err
	}

	if err = resp.All(&levels); err != nil {
		return nil, err
	}

	return levels, nil
}
