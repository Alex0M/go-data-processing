package db

import (
	"out/internal/models"

	"github.com/gocql/gocql"
)

func (c *CassandraDB) AddStream(s *models.Stream) error {
	err := c.Session.Query(`INSERT INTO streams (timestamp, id, client_id, device, content, geo) VALUES (?, ?, ?, ?, ?, ?)`,
		s.Timestamp, gocql.TimeUUID(), s.ClientID, s.Device, s.Content, s.Geo).Exec()
	if err != nil {
		return err
	}

	return nil
}
