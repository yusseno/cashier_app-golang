package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"time"
)

type SessionsRepository struct {
	db db.DB
}

func NewSessionsRepository(db db.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) ReadSessions() ([]model.Session, error) {
	records, err := u.db.Load("sessions")
	if err != nil {
		return nil, err
	}

	var listSessions []model.Session
	err = json.Unmarshal([]byte(records), &listSessions)
	if err != nil {
		return nil, err
	}

	return listSessions, nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}
	newSession := []model.Session{}
	// Select target token and delete from listSessions
	for _, element := range listSessions {
		// fmt.Print("ini listSessions :")
		if element.Token == tokenTarget {
			// fmt.Println("ini ada")
			continue
		}
		newSession = append(newSession, element)
	}
	// fmt.Println("ini token target : " + tokenTarget)
	// TODO: answer here

	jsonData, err := json.Marshal(newSession)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	dataSession, _ := u.ReadSessions()
	for _, element := range dataSession {
		if element.Token == session.Token {
			return nil
		}
	}
	dataSession = append(dataSession, session)
	// fmt.Println(listUser)
	jsonData, err := json.Marshal(dataSession)
	if err != nil {
		return err
	}
	u.db.Save("sessions", jsonData)
	return nil // TODO: replace this
}

func (u *SessionsRepository) CheckExpireToken(token string) (model.Session, error) {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return model.Session{}, err
	}

	for _, element := range listSessions {
		if element.Token == token {
			check := u.TokenExpired(element)
			if check == true {
				return model.Session{}, fmt.Errorf("Token is Expired!") // TODO: replace this
			}
			return element, nil
		}
	}
	return model.Session{}, nil
}

// }

func (u *SessionsRepository) ResetSessions() error {
	err := u.db.Reset("sessions", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) TokenExist(req string) (model.Session, error) {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return model.Session{}, err
	}
	for _, element := range listSessions {
		if element.Token == req {
			return element, nil
		}
	}
	return model.Session{}, fmt.Errorf("Token Not Found!")
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
