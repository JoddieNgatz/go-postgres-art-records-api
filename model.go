
package main

import (
    "database/sql"
)


type art struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func (p *art) getArt(db *sql.DB) error {
  return db.QueryRow("SELECT name, price FROM arts WHERE id=$1",
        p.ID).Scan(&p.Name, &p.Price)
}

func (p *art) updateArt(db *sql.DB) error {
  _, err :=
        db.Exec("UPDATE arts SET name=$1, price=$2 WHERE id=$3",
            p.Name, p.Price, p.ID)

    return err
}

func (p *art) deleteArt(db *sql.DB) error {
  _, err := db.Exec("DELETE FROM arts WHERE id=$1", p.ID)

  return err
}

func (p *art) createArt(db *sql.DB) error {
  err := db.QueryRow(
    "INSERT INTO arts(name, price) VALUES($1, $2) RETURNING id",
    p.Name, p.Price).Scan(&p.ID)

if err != nil {
    return err
}

return nil
}

func getArts(db *sql.DB, start, count int) ([]art, error) {
  rows, err := db.Query(
    "SELECT id, name,  price FROM arts LIMIT $1 OFFSET $2",
    count, start)

if err != nil {
    return nil, err
}

defer rows.Close()

arts := []art{}

for rows.Next() {
    var p art
    if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
        return nil, err
    }
    arts = append(arts, p)
}

return arts, nil
}