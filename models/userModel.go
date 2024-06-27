package models



// album represents data about a record album.
type User struct {
    ID     string  `json:"id"`
    Email  string  `json:"email"`
    Password string  `json:"password"`
    Username string `json:"username"`
}


// albums slice to seed record album data.
var Users = []User{
    {ID: "1", Email: "Blue Train", Password: "John Coltrane", Username: "blue_train"},
    {ID: "2", Email: "Jeru", Password: "Gerry Mulligan", Username: "jeru"},
    {ID: "3", Email: "Sarah Vaughan and Clifford Brown", Password: "Sarah Vaughan", Username: "sarah_vaughan"},
}
