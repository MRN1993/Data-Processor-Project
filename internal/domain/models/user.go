package models

import "time"

type User struct {
	ID     int 
	Quota  int  
	monthly_data_limit int
	request_limit_per_minute int
	used_data  int
	request_count int
	last_request_time  time.Time
}
