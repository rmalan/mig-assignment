package models

type Attendance struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserId   uint   `json:"user_id"`
	Date     string `json:"date"`
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
}
