package models

import(
  errors "jameselgar.com/goldney/errors"
  "fmt"
)

type Section struct {
    Title         string    `json:"title"`
    Description   string    `json:"description"`
    Tile_id       int       `json:"tile_id"`
    Type          string    `json:"type"`
}

//func CreateSection

func (db *DB) AddSection (s *Section) (int, *errors.ApiError) {
    sqlStmt := `INSERT INTO sections (title, description, tile_id) VALUES ($1,$2, $3) RETURNING id;`
  
    res, insertErr := db.Exec(sqlStmt, s.Title, s.Description, s.Tile_id)
    switch insertErr{
    case nil:
      var id int
      fmt.Println("The id is::")
      fmt.Println(res)
//      fmt.Printf("Section has been added to db %d", res)
      return id, nil
    default:
      return -1, &errors.ApiError{insertErr, "Unknown Error during Insertion of User", 400}
          }
}

func (db *DB) GetSections (id int) ([]Section, *errors.ApiError) {
  var sections []Section
  rows, err := db.Query("SELECT title, description, type FROM sections WHERE tile_id=$1", id)
    if err != nil {
          panic(err)
      return nil, &errors.ApiError{err, "Error whilst accessing sections from database", 400}
    }
    defer rows.Close()
    for rows.Next() {
        fmt.Println("Sections")
        var section Section
        if err := rows.Scan(&section.Title, &section.Description, &section.Type); err != nil {
          panic(err)
            return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
        }
       
        sections = append(sections, section)
        fmt.Println(section.Title)
    }
    return sections, nil
}
