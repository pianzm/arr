package model

import (
	"encoding/json"
)

const (
	ArticleCollection = "articles"
)

// QueryParameters for search
type QueryParameters struct {
	Page     int64
	Limit    int64
	Offset   int64
	StrPage  string `json:"page" query:"page"`
	StrLimit string `json:"limit" query:"limit"`
	Email    string `json:"email" query:"email"`
}

type StatusRequest struct {
	RequestID string `json:"request_id"`
}

func (u *QueryParameters) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *QueryParameters) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, u); err != nil {
		return err
	}
	return nil
}

func (u *QueryParameters) String() string {
	return "Q Param: " + u.Email
}

// redis pattern:
type QueueStatus struct {
	RequestID string          `json:"req_id"`
	Completed bool            `json:"completed"`
	Parameter QueryParameters `json:"params"`
	FilePath  string          `json:"path"`
}

func (u *QueueStatus) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *QueueStatus) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, u); err != nil {
		return err
	}
	return nil
}

type Member struct {
	ID        string `json:"id" faker:"uuid_digit"`
	FirstName string `json:"first_name" faker:"first_name"`
	LastName  string `json:"last_name" faker:"last_name"`
	Email     string `json:"email" faker:"email,unique"`
}
