package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/pianzm/arr/src/member/v1/model"
)

type MemberRepoPostgres struct {
	pgdb *sql.DB
}

func NewMemberRepo(db *sql.DB) *MemberRepoPostgres {
	return &MemberRepoPostgres{
		pgdb: db,
	}
}

func (repo MemberRepoPostgres) Insert(ctx context.Context, members []*model.Member) error {
	tx, err := repo.pgdb.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("members", "first_name", "last_name", "email", "created_at"))
	if err != nil {
		tx.Rollback()
		return err
	}
	createdTime := time.Now()
	for _, param := range members {
		if _, err := stmt.Exec(param.FirstName, param.LastName, param.Email, createdTime); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err = stmt.Exec(); err != nil {
		tx.Rollback()
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
