package models

import (
    "time"
)

type Query struct {
    Query       string
    Result      string
    Timestamp   time.Time
}    