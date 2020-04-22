package models

import(
    "fmt"
    errors "jameselgar.com/goldney/errors"
    secret "jameselgar.com/goldney/secret"
)
type Authenticate interface {
  Login (u *User) (*errors.ApiError)
}

type User struct {
  Username string `json:"username"`
  Password string `json:"password"`
}


func (db *DB) Login (u *User) (*errors.ApiError) {
  //testPassword, err := HashSaltPwd([]byte("test"))
  testPassword := secret.AdminPassword
  var admin User = User{"admin", testPassword};
  isSame, err := ComparePassword(admin.Password, []byte(u.Password))
  if err != nil {
    panic (err);
  }
  if !isSame || u.Username != admin.Username {
    fmt.Println("Incorrect credentials")
    return &errors.ApiError{nil, "Incorrect Credentials", 401}
  }
  fmt.Println("Passwords match")
  return nil
}
