package models

type Stream struct {
	ClientID  string `json:"clientID" binding:"required"`
	Device    string `json:"device" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Geo       string `json:"geo" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}
