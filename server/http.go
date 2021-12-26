package server

type Http struct {
	Address           string `json:"address"`
	ReadTimeOutInSec  int    `json:"read_time_out_in_sec"`
	WriteTimeOutInSec int    `json:"write_time_out_in_sec"`
	IdleTimeOutInSec  int    `json:"idle_time_out_in_sec"`
}
