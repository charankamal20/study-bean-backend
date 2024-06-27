package models



// album represents data about a record album.
type User struct {
    ID     string  `json:"id"`
    Username  string  `json:"username"`
    Password string  `json:"password"`
}


// albums slice to seed record album data.
var Users = []User{
    {ID: "1", Username: "Blue Train", Password: "John Coltrane"},
    {ID: "2", Username: "Jeru", Password: "Gerry Mulligan"},
    {ID: "3", Username: "Sarah Vaughan and Clifford Brown", Password: "Sarah Vaughan"},
}
