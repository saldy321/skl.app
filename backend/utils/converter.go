package utils

import (
"strconv"
"github.com/google/uuid"
)


// StrToInt buat convert string dari FormValue jadi Integer
func StrToInt(s string) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        return 0 // Balikin 0 kalau gagal parse atau input kosong
    }
    return i
}

func ParseUUID(s string) uuid.UUID {
    id, err := uuid.Parse(s)
    if err != nil {
        return uuid.Nil // Balikin ID kosong kalau gagal
    }
    return id
}