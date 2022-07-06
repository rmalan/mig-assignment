package models

type Activity struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserId   uint   `json:"user_id"`
	Date     string `json:"date"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Activity string `json:"activity"`
}
