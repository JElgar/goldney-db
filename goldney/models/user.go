package models

import(
    "fmt"
    errors "jameselgar.com/goldney/errors"
)

type Authenticate interface {
  Login (u *User) (*errors.ApiError)
}

type User struct {
  Username string `json:"username"`
  Password string `json:"password"`
}


func (db *DB) Login (u *User) (*errors.ApiError) {
  testPassword, err := HashSaltPwd([]byte("test"))
  if err != nil {
    panic(err);
  }
  var admin User = User{"admin", testPassword};
  isSame, err := ComparePassword(admin.Password, []byte(u.Password))
  if err != nil {
    panic (err);
  }
  if !isSame {
    fmt.Println("Incorrect password")
    return &errors.ApiError{nil, "Password does not match", 401}
  }
  fmt.Println("Passwords match")
  return nil
}
