package models

import(
    "fmt"
    errors "jameselgar.com/goldney/errors"
)

type TileStore interface {
  AddTile (t *Tile) (*Tile, *errors.ApiError)
  GetActiveTiles () ([]Tile, *errors.ApiError)
  GetAllTiles () ([]Tile, *errors.ApiError)
  UpdateTile (t *Tile) (*Tile, *errors.ApiError)
  AddSection (s *Section) (int, *errors.ApiError)
  GetSections (id int) ([]Section, *errors.ApiError)
  UpdateSection (s *Section) (int, int, *errors.ApiError)
  SetActive (t *Tile) *errors.ApiError
  DeleteSection (id int) (*errors.ApiError)
}

type Tile struct {
    Title         string    `json:"title"`
    Subtitle      string    `json:"subtitle"`
    Description   string    `json:"description"`
    Sections      []Section `json:"sections"`
    Email         string    `json:"email"`
    Id            int       `json:"id"`
    Active        int       `json:"active"`
    DeleteSecs    []int     `json:"delSecs"`
}

func (db *DB) AddTile (t *Tile) (*Tile, *errors.ApiError) {
  fmt.Println("Adding tile")
  // Temporarily got rid of sections
  sqlStmt := `
      INSERT INTO tiles (title, subtitle, description, email) 
      VALUES($1,$2, $3, $4) 
      RETURNING id;`
  var id int
  insertErr := db.QueryRow(sqlStmt, t.Title, t.Subtitle, t.Description, t.Email).Scan(&id)
  switch insertErr{
  case nil:
    fmt.Println("Tile has been added to db")
  default:
    fmt.Println("There was an error")
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

func (db *DB) UpdateTile (t *Tile) (*Tile, *errors.ApiError) {
  fmt.Println("Updating Tiles")
//  sqlStmt := `
//      UPDATE tiles 
//      SET title = $1, 
//          subtitle = &2, 
//          description = $3, 
//          email = $4 
//      WHERE id = $5;`
    sqlStmt := `UPDATE tiles 
                SET title = $1,
                    subtitle = $2,
                    description = $3,
                    email = $4 
                WHERE id = $5;`
 
  fmt.Println("Tile id is: ")
  fmt.Println(t.Id)
  fmt.Println(t.Title)
  fmt.Println(t.Subtitle)
  fmt.Println(t.Description)
  fmt.Println(t.Email)
  //res, updateErr := db.Exec(sqlStmt, t.Title, t.Subtitle, t.Description, t.Email, t.Id)
  res, updateErr := db.Exec(sqlStmt, t.Title, t.Subtitle, t.Description, t.Email, int(t.Id))
  if updateErr != nil {
    panic(updateErr)
    return nil, &errors.ApiError{updateErr, "Error updating", 400}
  }
  count, updateErr := res.RowsAffected()
  if updateErr != nil {
    panic(updateErr)
    return nil, &errors.ApiError{updateErr, "Error updating", 400}
  }
  if count > 1 {
    return nil, &errors.ApiError{updateErr, "uhoh more than 1 with id", 400}
  }
  if count < 1 {
    return nil, &errors.ApiError{updateErr, "Could not find tile with given ID", 400}
  }
  fmt.Println("Updated Tile Data, now doing sections")
  for _, s := range t.Sections {
    fmt.Println("Adding section")
    s.Tile_id = t.Id
    _, c, err := db.UpdateSection(&s)
    if err != nil {
      panic(err)
      return t, &errors.ApiError{err, "Error updating Section", 400}
    }
    // If the section with that id doesnt already exist add it
    if c == 0 {
      _, err := db.AddSection(&s)
      if err != nil {
        panic(err)
        return t, &errors.ApiError{err, "Error inserting new Section", 400}
      }
    }
  }
  return t, nil
}

func (db *DB) GetAllTiles () ([]Tile, *errors.ApiError) {
  fmt.Println("Getting tiles")
  var tiles []Tile
  rows, err := db.Query("SELECT * FROM tiles")
    if err != nil {
      panic(err)
      return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
    }
    defer rows.Close()
    fmt.Println("Got tiles")
    for rows.Next() {
        fmt.Println("Rows")
        var tile Tile
        if err := rows.Scan(&tile.Id, &tile.Title, &tile.Subtitle, &tile.Description, &tile.Email, &tile.Active); err != nil {
          panic(err)
            return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
        }
       
        s, sectionErr := db.GetSections(tile.Id)
        if sectionErr != nil {
            panic(err)
            return nil, &errors.ApiError{err, "Error accessing sections for tile", 400}
        }
        tile.Sections = s
        tiles = append(tiles, tile)
        fmt.Println(tile.Title)
    }
    fmt.Println("Returned tiles")
    return tiles, nil
}

func (db *DB) GetActiveTiles () ([]Tile, *errors.ApiError) {
  fmt.Println("Getting tiles")
  var tiles []Tile
  rows, err := db.Query("SELECT * FROM tiles WHERE active=1")
    if err != nil {
      panic(err)
      return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
    }
    defer rows.Close()
    fmt.Println("Got tiles")
    for rows.Next() {
        fmt.Println("Rows")
        var tile Tile
        if err := rows.Scan(&tile.Id, &tile.Title, &tile.Subtitle, &tile.Description, &tile.Email, &tile.Active); err != nil {
          panic(err)
            return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
        }
       
        s, sectionErr := db.GetSections(tile.Id)
        if sectionErr != nil {
            panic(err)
            return nil, &errors.ApiError{err, "Error accessing sections for tile", 400}
        }
        tile.Sections = s
        tiles = append(tiles, tile)
        fmt.Println(tile.Title)
    }
    fmt.Println("Returned tiles")
    return tiles, nil
}


func (db *DB) SetActive (t *Tile) *errors.ApiError {
    sqlStmt := `UPDATE tiles 
                SET active = $1
                WHERE id = $2;`
  res, updateErr := db.Exec(sqlStmt, int(t.Active), int(t.Id))
  if updateErr != nil {
    panic(updateErr)
    return &errors.ApiError{updateErr, "Error upadting active", 400}
  }
  count, updateErr := res.RowsAffected()
  if updateErr != nil {
    panic(updateErr)
    return &errors.ApiError{updateErr, "Error updating active", 400}
  }
  if count > 1 {
    return &errors.ApiError{updateErr, "uhoh more than 1 given id", 400}
  }
  if count < 1 {
    return &errors.ApiError{updateErr, "Could not find tile with given ID", 400}
  }
  fmt.Println("Updated Tile activity")
  return nil
}
