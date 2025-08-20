package postgres

import (
	"context"
	"fmt"

	"github.com/et0/avito-tech-internship-spring-2025/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const URI = "postgres://%s:%s@%s:%s/%s"

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(cfgDB *config.Database) (*Postgres, error) {
	conn := fmt.Sprintf(URI, cfgDB.Username, cfgDB.Password, cfgDB.Host, cfgDB.Port, cfgDB.Basename)
	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	//
	cfg.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	return &Postgres{Pool: pool}, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
