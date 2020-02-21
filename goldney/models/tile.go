package models

import(
    "fmt"
    errors "jameselgar.com/goldney/errors"
)

type TileStore interface {
  CreateTile (t *Tile) (*Tile, *errors.ApiError)
  GetTiles () ([]Tile, *errors.ApiError)
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

func (db *DB) GetTiles () ([]Tile, *errors.ApiError) {
  var tiles []Tile
  rows, err := db.Query("SELECT * FROM tiles")
    if err != nil {
      return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
    }
    defer rows.Close()
    println(rows)
    for rows.Next() {
        var tile Tile
        if err := rows.Scan(tile); err != nil {
            return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
        }
        tiles = append(tiles, tile)
        fmt.Println(tile)
    }
    return tiles, nil
}
