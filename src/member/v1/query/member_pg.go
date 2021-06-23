package query

import (
	"context"
	"database/sql"

	"github.com/pianzm/arr/src/member/v1/model"
)

type MemberQueryPostgres struct {
	readDB *sql.DB
}

func NewMemberQuery(db *sql.DB) *MemberQueryPostgres {
	return &MemberQueryPostgres{
		readDB: db,
	}
}

func (q MemberQueryPostgres) GetAll(ctx context.Context) ([][]string, error) {
	rawQuery := `SELECT "first_name", "last_name", "email" FROM members ORDER BY id DESC`
	rows, err := q.readDB.Query(rawQuery)
	if err != nil {
		return nil, err
	}
	results := [][]string{}
	for rows.Next() {
		row := model.Member{}
		if err := rows.Scan(&row.FirstName, &row.LastName, &row.Email); err != nil {
			return nil, err
		}
		m := []string{row.FirstName, row.LastName, row.Email}
		results = append(results, m)
	}
	return results, nil
}

func (q MemberQueryPostgres) GetTotal(ctx context.Context) (int64, error) {
	var total int64
	totalQuery := `SELECT COUNT(id) FROM members`

	if err := q.readDB.QueryRow(totalQuery).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
