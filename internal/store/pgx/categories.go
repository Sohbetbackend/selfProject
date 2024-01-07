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

const sqlCategoriesFields = `ct.id, ct.name`
const sqlCategoriesSelect = `select ` + sqlCategoriesFields + ` from categories ct where ct.id = ANY($1::int[])`
const sqlCategoriesSelectMany = `select ` + sqlCategoriesFields + `, count(*) over()  as total from categories ct
	where ct.id=ct.id limit $1 offset $2`
const sqlCategoriesInsert = `insert into categories`
const sqlCategoriesUpdate = `update categories ct set id=id`
const sqlCategoriesDelete = `delete from categories ct where id = ANY($1::int[])`

func scanCategories(rows pgx.Rows, m *models.Category, addColumns ...interface{}) (err error) {
	err = rows.Scan(parseColumnsForScan(m, addColumns...)...)
	return
}

func (d *PgxStore) CategoryFindById(ID string) (*models.Category, error) {
	row, err := d.CategoryFindByIds([]string{ID})
	if err != nil {
		return nil, err
	}
	if len(row) < 1 {
		return nil, pgx.ErrNoRows
	}
	return row[0], nil
}

func (d *PgxStore) CategoryFindByIds(Ids []string) ([]*models.Category, error) {
	categories := []*models.Category{}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), sqlCategoriesSelect, (Ids))
		for rows.Next() {
			category := models.Category{}
			err := scanCategories(rows, &category)
			if err != nil {
				utils.LoggerDesc("Scan error").Error(err)
				return err
			}
			categories = append(categories, &category)
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	return categories, nil
}

func (d *PgxStore) CategoryFindBy(f models.CategoryFilterRequest) (categories []models.Category, total int, err error) {
	args := []interface{}{f.Limit, f.Offset}
	qs, args := CategoryListBuildQuery(f, args)
	err = d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), qs, args...)
		for rows.Next() {
			category := models.Category{}
			err = scanCategories(rows, &category, &total)
			if err != nil {
				utils.LoggerDesc("Scan error").Error(err)
				return err
			}
			categories = append(categories, category)
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, 0, err
	}
	return categories, total, nil
}

func (d *PgxStore) CategoryCreate(model *models.Category) (*models.Category, error) {
	qs, args := CategoryCreateQuery(model)
	qs += " RETURNING id"
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		err = tx.QueryRow(context.Background(), qs, args...).Scan(&model.ID)
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	editModel, err := d.CategoryFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxStore) CategoryUpdate(model *models.Category) (*models.Category, error) {
	qs, args := CategoryUpdateQuery(model)
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		_, err = tx.Query(context.Background(), qs, args...)
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	editModel, err := d.CategoryFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxStore) CategoryDelete(items []*models.Category) ([]*models.Category, error) {
	ids := []uint{}
	for _, i := range items {
		ids = append(ids, i.ID)
	}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		_, err = tx.Query(context.Background(), sqlCategoriesDelete, (ids))
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	return items, nil
}

func CategoryCreateQuery(m *models.Category) (string, []interface{}) {
	args := []interface{}{}
	cols := ""
	vals := ""
	q := CategoryAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		cols += ", " + k
		vals += ", $" + strconv.Itoa(len(args))
	}
	qs := sqlCategoriesInsert + " (" + strings.Trim(cols, ", ") + ") VALUES (" + strings.Trim(vals, ", ") + ")"
	return qs, args
}

func CategoryUpdateQuery(m *models.Category) (string, []interface{}) {
	args := []interface{}{}
	sets := ""
	q := CategoryAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		sets += ", " + k + "=$" + strconv.Itoa(len(args))
	}
	args = append(args, m.ID)
	qs := strings.ReplaceAll(sqlCategoriesUpdate, "set id=id", "set id=id "+sets+" ") + "where id=$" + strconv.Itoa(len(args))
	return qs, args
}

func CategoryAtomicQuery(m *models.Category) map[string]interface{} {
	q := map[string]interface{}{}
	q["name"] = m.Name
	return q
}

func CategoryListBuildQuery(f models.CategoryFilterRequest, args []interface{}) (string, []interface{}) {
	var wheres string = ""

	if f.ID != nil && *f.ID != 0 {
		args = append(args, *f.ID)
		wheres += " and ct.id=$" + strconv.Itoa(len(args))
	}
	wheres += " order by ct.id desc"
	qs := sqlCategoriesSelectMany
	qs = strings.ReplaceAll(qs, "ct.id=ct.id", "ct.id=ct.id "+wheres+" ")
	return qs, args
}
