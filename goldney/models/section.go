package models

import(
  errors "jameselgar.com/goldney/errors"
  "fmt"
  "database/sql"
)

type Section struct {
    Title         string    `json:"title"`
    Description   string    `json:"description"`
    Tile_id       int       `json:"tile_id"`
    Type          string    `json:"type"`
    ImageName     string    `json:"image_name"`
    ImageLink     string    `json:"image_link"`
    Id            int       `json:"id"`
}

//func CreateSection

func (db *DB) AddSection (s *Section) (int, *errors.ApiError) {
    var res sql.Result
    var insertErr error
    fmt.Print("Adding a section of type: ")
    fmt.Println(s.Type)
    if (s.Type == "text") {
      sqlStmt := `INSERT INTO sections (title, description, tile_id) VALUES ($1,$2, $3) RETURNING id;`
      res, insertErr = db.Exec(sqlStmt, s.Title, s.Description, s.Tile_id)
    } else if (s.Type == "image") {
      sqlStmt := `INSERT INTO sections (tile_id, type, image_name, image_link) VALUES ($1,$2, $3, $4) RETURNING id;`
      res, insertErr = db.Exec(sqlStmt, s.Tile_id, s.Type, s.ImageName, s.ImageLink)
    }
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

func (db *DB) UpdateSection (s *Section) (int, int, *errors.ApiError) {
    var res sql.Result
    var insertErr error
    if (s.Type == "text") {
      sqlStmt := `UPDATE sections 
                  SET title = $1, description = $2
                  WHERE id=$3;`
      res, insertErr = db.Exec(sqlStmt, s.Title, s.Description, s.Tile_id)
    } else if (s.Type == "image") {
      sqlStmt := `UPDATE sections 
                  SET image_name = $1, image_link = $2 
                  WHERE id=$3;`
      res, insertErr = db.Exec(sqlStmt, s.ImageName, s.ImageLink, s.Id)
    }
    switch insertErr{
    case nil:
      var id int
      fmt.Println("The id is::")
      fmt.Println(res)
//      fmt.Printf("Section has been added to db %d", res)
      count, updateErr := res.RowsAffected()
      if updateErr != nil {
        return -1, 0, &errors.ApiError{updateErr, "Error updating", 400}
      }
      return id, int(count), nil
    default:
      return -1, 0, &errors.ApiError{insertErr, "Unknown Error during Insertion of User", 400}
          }
}

func (db *DB) GetSections (id int) ([]Section, *errors.ApiError) {
  var sections []Section
  rows, err := db.Query("SELECT title, description, type, image_link, image_name, id FROM sections WHERE tile_id=$1", id)
    if err != nil {
          panic(err)
      return nil, &errors.ApiError{err, "Error whilst accessing sections from database", 400}
    }
    defer rows.Close()
    for rows.Next() {
        fmt.Println("Sections")
        var section Section
        if err := rows.Scan(&section.Title, &section.Description, &section.Type, &section.ImageLink, &section.ImageName, &section.Id); err != nil {
          panic(err)
            return nil, &errors.ApiError{err, "Error whilst accessing tiles from database", 400}
        }
       
        sections = append(sections, section)
        fmt.Println(section.Title)
    }
    return sections, nil
}

func (db *DB) DeleteSection (id int) (*errors.ApiError) {
  sqlStmt := `DELETE sections 
              WHERE id=$1;`
  res, delErr := db.Exec(sqlStmt, id)
  
  count, delErr := res.RowsAffected()
  if delErr != nil {
    panic(delErr)
    return &errors.ApiError{delErr, "Error deleteing", 400}
  }
  if count > 1 {
    return &errors.ApiError{delErr, "uhoh more than 1 with id", 400}
  }
  if count < 1 {
    return &errors.ApiError{delErr, "Could not find tile with given ID", 400}
  }
  fmt.Println("Updated Tile Data, now doing sections")
  return nil
}
