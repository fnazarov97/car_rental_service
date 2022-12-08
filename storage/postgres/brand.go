package postgres

import (
	"car_rental/genprotos/brand_service"
	"database/sql"
	"errors"
)

// AddBrand ...
func (p Postgres) AddBrand(id string, req *brand_service.CreateBrandReq) (res *brand_service.CreateBrandRes, err error) {
	req.ID = id
	_, err = p.DB.Exec(`Insert into author(id, fullname, created_at) 
							VALUES($1,$2,now())`, req.ID, req.Fullname)
	if err != nil {
		return res, err
	}
	res = &author.CreateBrandRes{}
	return res, nil
}

// GetBrandByID ...
func (p Postgres) GetBrandByID(req *brand_service.Id) (*brand_service.GetBrandByIdRes, error) {
	result := &author.GetBrandByIdRes{
		Articles: make([]*brand_service.Article, 0),
	}
	var (
		updated_at sql.NullString
		deleted_at sql.NullString
	)
	row := p.DB.QueryRow("SELECT id, created_at, updated_at, deleted_at, fullname FROM author WHERE id=$1", req.Id)
	err := row.Scan(&result.Id, &result.CreatedAt, &updated_at, &deleted_at, &result.Fullname)
	if updated_at.Valid {
		result.UpdatedAt = updated_at.String
	}
	if deleted_at.Valid {
		result.DeletedAt = deleted_at.String
	}
	if err != nil {
		return result, err
	}
	ars, err := p.GetArticlesByBrandID(req)
	result.Articles = ars.Articles
	if err != nil {
		return result, err
	}
	return result, nil
}

// GetArticlesByBrandID ...
func (p Postgres) GetArticlesByBrandID(req *brand_service.Id) (*brand_service.GetArticles, error) {
	resp := &author.GetArticles{
		Articles: make([]*brand_service.Article, 0),
	}
	rows, err := p.DB.Queryx(`SELECT 
									 id, 
									 title, 
									 body, 
									 author_id, 
									 created_at,
									 updated_at,
									 deleted_at
							FROM article
							WHERE author_id = $1 `, req.Id)
	if err != nil {
		return resp, err
	}
	for rows.Next() {
		var (
			update_at  sql.NullString
			deleted_at sql.NullString
		)
		row := author.Article{
			Content: &author.Post{},
		}
		err := rows.Scan(&row.Id, &row.Content.Title, &row.Content.Body,
			&row.BrandId, &row.CreatedAt, &update_at, &deleted_at)
		if err != nil {
			return resp, err
		}
		if update_at.Valid {
			row.UpdatedAt = update_at.String
		}
		if deleted_at.Valid {
			row.DeletedAt = deleted_at.String
		}
		resp.Articles = append(resp.Articles, &row)
	}
	return resp, nil
}

// GetBrandList ...
func (p Postgres) GetBrandList(req *brand_service.GetBrandListReq) (*brand_service.GetBrands, error) {
	resp := &author.GetBrands{
		Brands: make([]*brand_service.Brand, 0),
	}
	rows, err := p.DB.Queryx(`SELECT
	id,
	fullname,
	created_at,
	updated_at
	FROM author WHERE deleted_at IS NULL AND (fullname ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3
	`, req.Search, int(req.Limit), int(req.Offset))
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var (
			a         author.Brand
			update_at sql.NullString
		)

		err := rows.Scan(
			&a.Id,
			&a.Fullname,
			&a.CreatedAt,
			&update_at,
		)
		if err != nil {
			return resp, err
		}
		if update_at.Valid {
			a.UpdatedAt = update_at.String
		}
		resp.Brands = append(resp.Brands, &a)
	}

	return resp, nil
}

// UpdateBrand ...
func (p Postgres) UpdateBrand(req *brand_service.UpdateBrandReq) error {
	res, err := p.DB.NamedExec("UPDATE author  SET fullname=:f, updated_at=now() WHERE deleted_at IS NULL AND id=:id", map[string]interface{}{
		"id": req.Id,
		"f":  req.Fullname,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}
	return errors.New("author not found")
}

// DeleteBrand ...
func (p Postgres) DeleteBrand(req *brand_service.Id) error {
	res, err := p.DB.Exec("UPDATE author  SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL", req.Id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("author had been deleted already")
}
