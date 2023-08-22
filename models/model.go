
package main

import (
    "database/sql"
    "errors"
)


type art struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func (p *art) getArt(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (p *art) updateArt(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (p *art) deleteArt(db *sql.DB) error {
  return errors.New("Not implemented")
}

func (p *art) createArt(db *sql.DB) error {
  return errors.New("Not implemented")
}

func getArts(db *sql.DB, start, count int) ([]art, error) {
  return nil, errors.New("Not implemented")
}