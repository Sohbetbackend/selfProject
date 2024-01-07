package pgx

import (
	"context"
	"strconv"
	"strings"

	"github.com/Sohbetbackend/selfProject/internal/models"
	"github.com/Sohbetbackend/selfProject/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const sqlAuthorFields = `a.id, a.first_name, a.last_name`
const sqlAuthorSelect = `select ` + sqlAuthorFields + ` from authors a where a.id = ANY($1::[]int)`
const sqlAuthorSelectMany = `select ` + sqlAuthorFields + ` count(*) over() as total from authors a
	where a.id=a.id`
const sqlAuthorInsert = `insert into authors`
const sqlAuthorUpdate = `update authors a set id=id`
const sqlAuthorDelete = `delete from authors a where id = ANY($1::int[])`

func scanAuthors(rows pgx.Row, m *models.Author, addColumns ...interface{}) (err error) {
	err = rows.Scan(parseColumnsForScan(m, addColumns...))
	return
}

func (d *PgxStore) AuthorFindById(ID string) (*models.Author, error) {
	row, err := d.AuthorFindByIds([]string{ID})
	if err != nil {
		return nil, err
	}
	if len(row) < 1 {
		return nil, pgx.ErrNoRows
	}
	return row[0], nil
}

func (d *PgxStore) AuthorFindByIds(Ids []string) ([]*models.Author, error) {
	authors := []*models.Author{}
	err := d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		rows, err := conn.Query(context.Background(), sqlAuthorSelect, (Ids))
		for rows.Next() {
			author := models.Author{}
			err := scanAuthors(rows, &author)
			if err != nil {
				return err
			}
			authors = append(authors, &author)
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	return authors, nil
}

func (d *PgxStore) AuthorsFindBy(f models.AuthorsFilterRequest) (authors []models.Author, total int, err error) {
	args := []interface{}{}
	qs, args := AuthorsListBuildQuery(f, args)
	err = d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		rows, err := conn.Query(context.Background(), qs, args...)
		for rows.Next() {
			author := models.Author{}
			err = scanAuthors(rows, &author, &total)
			if err != nil {
				return err
			}
			authors = append(authors, author)
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, 0, err
	}
	return authors, total, nil
}

func (d *PgxStore) AuthorsCreate(model *models.Author) (*models.Author, error) {
	qs, args := AuthorsCreateQuery(model)
	qs += " RETURNING id"
	err := d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Query(context.Background(), qs, args...)
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	editModel, err := d.AuthorFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxStore) AuthorsUpdate(model *models.Author) (*models.Author, error) {
	qs, args := AuthorsUpdateQuery(model)
	err := d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Query(context.Background(), qs, args...)
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	editModel, err := d.AuthorFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxStore) AuthorsDelete(items []*models.Author) ([]*models.Author, error) {
	ids := []uint{}
	for _, i := range items {
		ids = append(ids, i.ID)
	}
	err := d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Query(context.Background(), sqlAuthorDelete, (ids))
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	return items, nil
}

func AuthorsCreateQuery(m *models.Author) (string, []interface{}) {
	args := []interface{}{}
	cols := ""
	vals := ""
	q := AuthorsAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		cols += ", " + k
		vals += ", $" + strconv.Itoa(len(args))
	}
	qs := sqlAuthorInsert + " (" + strings.Trim(cols, ", ") + ") VALUES (" + strings.Trim(vals, ", ") + ")"
	return qs, args
}

func AuthorsUpdateQuery(m *models.Author) (string, []interface{}) {
	args := []interface{}{}
	sets := ""
	q := AuthorsAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		sets += ", " + k + "=$" + strconv.Itoa(len(args))
	}
	args = append(args, m.ID)
	qs := strings.ReplaceAll(sqlAuthorUpdate, "set id=id", "set id=id "+sets+" ") + "where id=$" + strconv.Itoa(len(args))
	return qs, args
}

func AuthorsAtomicQuery(m *models.Author) map[string]interface{} {
	q := map[string]interface{}{}
	q["first_name"] = m.FirstName
	q["last_name"] = m.LastName
	return q
}

func AuthorsListBuildQuery(f models.AuthorsFilterRequest, args []interface{}) (string, []interface{}) {
	var wheres string = ""

	if f.ID != nil && *f.ID != 0 {
		args = append(args, *f.ID)
		wheres += " and a.id=$" + strconv.Itoa(len(args))
	}
	wheres += " group by a.id "
	wheres += " order by a.id desc"
	qs := sqlAuthorSelectMany
	qs = strings.ReplaceAll(qs, "a.id=a.id", "a.id=a.id "+wheres+" ")
	return qs, args
}
