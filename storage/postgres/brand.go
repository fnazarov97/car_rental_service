package postgres

import (
	"car_rental/genprotos/brand"
	"database/sql"
	"errors"
)

// AddBrand ...
func (p Postgres) CreateBrand(id string, req *brand.CreateBrandRequest) (res *brand.Brand, err error) {

	_, err = p.DB.Exec(`Insert into brand("id", "name", "discription", "year", "created_at") 
							VALUES($1, $2, $3, $4, now())`, id, req.Name, req.Discription, req.Year)
	if err != nil {
		return nil, errors.New("create error")
	}
	Id := &brand.GetBrandByIDRequest{
		BrandId: id,
	}
	res, err = p.GetBrandByID(Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetBrandByID ...
func (p Postgres) GetBrandByID(req *brand.GetBrandByIDRequest) (*brand.Brand, error) {
	result := &brand.Brand{}
	var (
		updated_at sql.NullString
		deleted_at sql.NullString
	)
	row := p.DB.QueryRow(`SELECT id, name, discription, year, created_at, updated_at, deleted_at FROM brand WHERE id = $1`, req.BrandId)
	err := row.Scan(&result.BrandId, &result.Name, &result.Discription, &result.Year, &result.CreatedAt, &updated_at, &deleted_at)
	if updated_at.Valid {
		result.UpdatedAt = updated_at.String
	}
	if deleted_at.Valid {
		return nil, errors.New("brand not found")
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetBrandList ...
func (p Postgres) GetBrandList(req *brand.GetBrandListRequest) (*brand.GetBrandListResponse, error) {
	resp := &brand.GetBrandListResponse{
		Brands: []*brand.Brand{},
	}
	rows, err := p.DB.Queryx(`SELECT
	"id", "name", "discription", "year", "created_at", "updated_at"
	FROM "brand" WHERE "deleted_at" IS NULL AND ("name" || "discription" ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3
	`, req.Search, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			a         brand.Brand
			update_at sql.NullString
		)

		err := rows.Scan(
			&a.BrandId,
			&a.Name,
			&a.Discription,
			&a.Year,
			&a.CreatedAt,
			&update_at,
		)
		if err != nil {
			return nil, err
		}
		if update_at.Valid {
			a.UpdatedAt = update_at.String
		}
		resp.Brands = append(resp.Brands, &a)
	}

	return resp, nil
}

// UpdateBrand ...
func (p Postgres) UpdateBrand(id string, req *brand.UpdateBrandRequest) error {
	res, err := p.DB.NamedExec(`
	UPDATE 
		"brand"  
	SET 
		"name"=:n, "discription"=:d, "year"=:y, "updated_at"=now() WHERE deleted_at IS NULL AND id=:id`, map[string]interface{}{
		"id": id,
		"n":  req.Name,
		"d":  req.Discription,
		"y":  req.Year,
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
	return errors.New("brand not found")
}

// DeleteBrand ...
func (p Postgres) DeleteBrand(req *brand.DeleteBrandRequest) error {
	res, err := p.DB.Exec(`UPDATE "brand"  SET deleted_at=now() WHERE "id"=$1 AND "deleted_at" IS NULL`, req.BrandId)
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

	return errors.New("brand had been deleted already")
}
