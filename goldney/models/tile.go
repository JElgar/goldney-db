package models

import(
    "fmt"
    errors "jameselgar.com/goldney/errors"
)

type TileStore interface {
  CreateTile (t *Tile) (*Tile, *errors.ApiError)
}

type Tile struct {
    Title         string    `json:"title"`
    Subtitle      string    `json:"subtitle"`
    Description   string    `json:"description"`
    Sections      []Section `json:"sections"`
    Email         string    `json:"email"`
}

func (db *DB) CreateTile (t *Tile) (*Tile, *errors.ApiError) {
  // Temporarily got rid of sections
  sqlStmt := `INSERT INTO tiles (title, subtitle, description, email) VALUES($1,$2, $3, $4);`

  res, insertErr := db.Exec(sqlStmt, t.Title, t.Subtitle, t.Description, t.Email)
  switch insertErr{
  case nil:
    fmt.Println("User inserted")
    fmt.Println(res)
    return t, nil
  default:
    return t, &errors.ApiError{insertErr, "Unknown Error during Insertion of User", 400}
        }

}
