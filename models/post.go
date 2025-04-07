package models

import "time"


type Post struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    UserID      uint64    `json:"user_id"`
    CreatedAt   time.Time `json:"created_at"`
}