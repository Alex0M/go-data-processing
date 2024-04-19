package db

import (
	"context"
	"out/internal/models"
)

func (d *Database) AddStream(s *models.Stream, c context.Context) error {
	_, err := d.Client.NewInsert().Model(s).Exec(c)
	if err != nil {
		return err
	}

	return nil
}
