package models

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Password  string `json:"password"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    IzinUsaha string `json:"izin_usaha"`
    Verified  bool   `json:"verified"`
}
