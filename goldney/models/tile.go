package models

import(
    "fmt"
    errors "jameselgar.com/goldney/errors"
)

type TileStore interface {
  AddTile (t *Tile) (*Tile, *errors.ApiError)
  GetTiles () ([]Tile, *errors.ApiError)
  AddSection (s *Section) (int, *errors.ApiError)
  GetSections (id int) ([]Section, *errors.ApiError)
}

type Tile struct {
    Title         string    `json:"title"`
    Subtitle      string    `json:"subtitle"`
    Description   string    `json:"description"`
    Sections      []Section `json:"sections"`
    Email         string    `json:"email"`
}

func (db *DB) AddTile (t *Tile) (*Tile, *errors.ApiError) {
  fmt.Println("Adding tile")
  // Temporarily got rid of sections
  sqlStmt := `INSERT INTO tiles (title, subtitle, description, email) VALUES($1,$2, $3, $4) RETURNING id;`
  var id int
  insertErr := db.QueryRow(sqlStmt, t.Title, t.Subtitle, t.Description, t.Email).Scan(&id)
  switch insertErr{
  case nil:
    fmt.Println("Tile has been added to db")
  default:
    return t, &errors.ApiError{insertErr, "Unknown Error during Insertion of User", 400}
  }
  
  fmt.Println("Adding section")

  for _, s := range t.Sections {
    fmt.Println("Adding section")
    s.Tile_id = id
    a, err := db.AddSection(&s)
    if err != nil {
      return t, &errors.ApiError{err, "Error inserting Section", 400}
    }
    fmt.Println(a)
  }
  return t, nil
}

func (db *DB) GetTiles () ([]Tile, *errors.ApiError) {
  var tiles []Tile
  rows, err := db.Query("SELECT * FROM tiles")
    if err != nil {
      return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
    }
    defer rows.Close()
    for rows.Next() {
        fmt.Println("Rows")
        var tile Tile
        var id int
        if err := rows.Scan(&id, &tile.Title, &tile.Subtitle, &tile.Description, &tile.Email); err != nil {
          panic(err)
            return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
        }
       
        s, sectionErr := db.GetSections(id)
        if sectionErr != nil {
            return nil, &errors.ApiError{err, "Error accessing sections for tile", 400}
        }
        tile.Sections = s
        tiles = append(tiles, tile)
        fmt.Println(tile.Title)
    }
    return tiles, nil
}
