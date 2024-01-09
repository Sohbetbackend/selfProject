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

const sqlBookFields = `b.id, b.category_id, b.author_id, b.name, b.page`
const sqlBookInsert = `insert into books`
const sqlBookUpdate = `update books b set id=id`
const sqlBookDelete = `delete from books b where id = ANY($1::int[])`
const sqlBookSelect = `select ` + sqlBookFields + ` from books b where b.id = ANY($1::int[])`
const sqlBookSelectMany = `select ` + sqlBookFields + `, count(*) over() as total from books b where b.id=b.id`

const sqlBookAuthor = `select ` + sqlAuthorFields + `, b.id from books b
	right join authors a on (a.id=b.author_id) where b.id = ANY($1::int[])`

const sqlBookCategory = `select ` + sqlCategoriesFields + `, b.id from books b
	right join categories ct on (ct.id=b.category_id) where b.id = ANY($1::int[])`

func scanBook(rows pgx.Row, m *models.Book, addColumns ...interface{}) (err error) {
	err = rows.Scan(parseColumnsForScan(m, addColumns...)...)
	return
}

func (d *PgxStore) BookFindById(ID string) (*models.Book, error) {
	row, err := d.BookFindByIds([]string{ID})
	if err != nil {
		return nil, err
	}
	if len(row) < 1 {
		return nil, pgx.ErrNoRows
	}
	return row[0], nil
}

func (d *PgxStore) BookFindByIds(Ids []string) ([]*models.Book, error) {
	books := []*models.Book{}
	err := d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		rows, err := conn.Query(context.Background(), sqlBookSelect, (Ids))
		for rows.Next() {
			book := models.Book{}
			err := scanBook(rows, &book)
			if err != nil {
				return err
			}
			books = append(books, &book)
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	return books, nil
}

func (d *PgxStore) BooksFindBy(f models.BookFilterRequest) (books []*models.Book, total int, err error) {
	args := []interface{}{}
	qs, args := BooksListBuildQuery(f, args)
	err = d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), qs, args...)
		for rows.Next() {
			book := models.Book{}
			err = scanBook(rows, &book, &total)
			if err != nil {
				return err
			}
			books = append(books, &book)
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, 0, err
	}
	return books, total, nil
}

func (d *PgxStore) BookCreate(model *models.Book) (*models.Book, error) {
	qs, args := BookCreateQuery(model)
	qs += " RETURNING id"
	err := d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		err = conn.QueryRow(context.Background(), qs, args...).Scan(&model.ID)
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	editModel, err := d.BookFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxStore) BookUpdate(model *models.Book) (*models.Book, error) {
	qs, args := BookUpdateQuery(model)
	err := d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Query(context.Background(), qs, args...)
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	editModel, err := d.BookFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxStore) BookDelete(items []*models.Book) ([]*models.Book, error) {
	ids := []uint{}
	for _, i := range items {
		ids = append(ids, i.ID)
	}
	err := d.runQuery(context.Background(), func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Query(context.Background(), sqlBookDelete, (ids))
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	return items, nil
}

func BookCreateQuery(m *models.Book) (string, []interface{}) {
	args := []interface{}{}
	cols := ""
	vals := ""
	q := BookAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		cols += ", " + k
		vals += ", $" + strconv.Itoa(len(args))
	}
	qs := sqlBookInsert + " (" + strings.Trim(cols, ", ") + ") VALUES (" + strings.Trim(vals, ", ") + ")"
	return qs, args
}

func BookUpdateQuery(m *models.Book) (string, []interface{}) {
	args := []interface{}{}
	sets := ""
	q := BookAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		sets += ", " + k + "=$" + strconv.Itoa(len(args))
	}
	args = append(args, m.ID)
	qs := strings.ReplaceAll(sqlBookUpdate, "set id=id", "set id=id "+sets+" ") + "where id=$" + strconv.Itoa(len(args))
	return qs, args
}

func BookAtomicQuery(m *models.Book) map[string]interface{} {
	q := map[string]interface{}{}
	q["name"] = m.Name
	if m.Page != nil {
		q["page"] = m.Page
	}
	if m.AuthorId != nil {
		q["author_id"] = m.AuthorId
	}
	if m.CategoryId != nil {
		q["category_id"] = m.CategoryId
	}
	return q
}

func BooksListBuildQuery(f models.BookFilterRequest, args []interface{}) (string, []interface{}) {
	var wheres string = ""

	if f.ID != nil && *f.ID != 0 {
		args = append(args, *f.ID)
		wheres += " and b.id=$" + strconv.Itoa(len(args))
	}
	if f.AuthorId != nil && *f.AuthorId != 0 {
		args = append(args, *f.AuthorId)
		wheres += " and b.author_id=$" + strconv.Itoa(len(args))
	}
	if f.CategoryId != nil && *f.CategoryId != 0 {
		args = append(args, *f.CategoryId)
		wheres += " and b.category_id=$" + strconv.Itoa(len(args))
	}
	wheres += " group by b.id "
	wheres += " order by b.id desc"
	qs := sqlBookSelectMany
	qs = strings.ReplaceAll(qs, "b.id=b.id", "b.id=b.id "+wheres+" ")
	return qs, args
}

func (d *PgxStore) BookLoadRelations(l *[]*models.Book) error {
	ids := []string{}
	for _, m := range *l {
		ids = append(ids, strconv.Itoa(int(m.ID)))
	}
	if len(ids) < 1 {
		return nil
	}

	if rs, err := d.BookLoadAuthor(ids); err == nil {
		for _, r := range rs {
			for _, m := range *l {
				if r.ID == m.ID {
					m.Author = r.Relation
				}
			}
		}
	} else {
		return err
	}

	if rs, err := d.BookLoadCategory(ids); err != nil {
		for _, r := range rs {
			for _, m := range *l {
				if r.ID == m.ID {
					m.Category = r.Relation
				}
			}
		}
	} else {
		return err
	}
	return nil
}

type BookLoadAuthorItem struct {
	ID       uint
	Relation *models.Author
}

func (d *PgxStore) BookLoadAuthor(ids []string) ([]BookLoadAuthorItem, error) {
	res := []BookLoadAuthorItem{}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), sqlBookAuthor, (ids))
		for rows.Next() {
			sub := models.Author{}
			pid := uint(0)
			err = scanAuthors(rows, &sub, &pid)
			if err != nil {
				return err
			}
			res = append(res, BookLoadAuthorItem{ID: pid, Relation: &sub})
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}

	return res, nil
}

type BookLoadCategoryItem struct {
	ID       uint
	Relation *models.Category
}

func (d *PgxStore) BookLoadCategory(ids []string) ([]BookLoadCategoryItem, error) {
	res := []BookLoadCategoryItem{}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), sqlBookCategory, (ids))
		for rows.Next() {
			sub := models.Category{}
			pid := uint(0)
			err = scanCategories(rows, &sub, &pid)
			if err != nil {
				return err
			}
			res = append(res, BookLoadCategoryItem{ID: pid, Relation: &sub})
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}

	return res, nil
}
