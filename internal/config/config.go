package config

import "github.com/LearnShareApp/learn-share-backend/pkg/db/postgres"

type Config struct {
	postgres.DbConfig
}
