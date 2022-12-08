package postgres

import (
	"car_rental/genprotos/brand"
	"car_rental/genprotos/car"
	"errors"
	"time"
)

// CreateCar ...
func (p Postgres) CreateCar(id string, req *car.CreateCarRequest) error {
	Id := &brand.GetBrandByIDRequest{
		BrandId: req.BrandId,
	}
	_, err := p.GetBrandByID(Id)
	if err != nil {
		return errors.New("brand not found to create car")
	}

	_, err = p.DB.Exec(`
	INSERT INTO
	 "car"("id", "model", "color", "year", "mileage", "brand_id", created_at) 
	VALUES($1, $2, $3, $4, $5, $6, now())
	`, id, req.Model, req.Color, req.Year, req.Mileage, req.BrandId)
	if err != nil {
		return err
	}
	return nil
}

// GetCarByID ...
func (p Postgres) GetCarByID(id string) (*car.GetCarByIDResponse, error) {
	res := &car.GetCarByIDResponse{
		Brand: &car.GetCarByIDResponse_Brand{},
	}
	var deletedAt *time.Time
	var updatedAt, brandUpdatedAt *string
	err := p.DB.QueryRow(`SELECT 
	c.id,
	c.model,
	c.color,
	c.year,
	c.mileage,
	c.created_at,
	c.updated_at,
	b.id,
	b.name,
	b.discription,
	b.created_at,
	b.updated_at
    FROM car AS c JOIN brand AS b ON c.brand_id = b.id WHERE c.id = $1`, id).Scan(
		&res.CarId,
		&res.Model,
		&res.Color,
		&res.Year,
		&res.Mileage,
		&res.CreatedAt,
		&updatedAt,
		&res.Brand.BrandId,
		&res.Brand.Name,
		&res.Brand.Discription,
		&res.Brand.CreatedAt,
		&brandUpdatedAt,
	)
	if err != nil {
		return res, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if brandUpdatedAt != nil {
		res.Brand.UpdatedAt = *brandUpdatedAt
	}

	if deletedAt != nil {
		return res, errors.New("car not found")
	}

	return res, err
}

// GetCarList ...
func (p Postgres) GetCarList(offset, limit int, search string) (*car.GetCarListResponse, error) {
	resp := &car.GetCarListResponse{
		Cars: make([]*car.Car, 0),
	}

	rows, err := p.DB.Queryx(`
	SELECT
		"id", "model", "color", "year", "mileage", "brand_id", "created_at", "updated_at"
	FROM 
		"car" WHERE deleted_at IS NULL AND (model ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3
	`, search, limit, offset)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		a := &car.Car{}

		var updatedAt *string

		err := rows.Scan(
			&a.CarId, &a.Model, &a.Color, &a.Year, &a.Mileage, &a.BrandId, &a.CreatedAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedAt != nil {
			a.UpdatedAt = *updatedAt
		}

		resp.Cars = append(resp.Cars, a)
	}

	return resp, err
}

// UpdateCar ...
func (p Postgres) UpdateCar(id string, entity *car.UpdateCarRequest) error {

	res, err := p.DB.NamedExec(`
	UPDATE 
		"car"  
	SET 
		"model"=:mo, "color"=:c, "year"=:y, "mileage"=:mi, "brand_id"=:b, updated_at=now() WHERE deleted_at IS NULL AND id=:id`, map[string]interface{}{
		"id": id,
		"mo": entity.Model,
		"c":  entity.Color,
		"y":  entity.Year,
		"mi": entity.Mileage,
		"b":  entity.BrandId,
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

	return errors.New("car not found")
}

// DeleteCar ...
func (p Postgres) DeleteCar(id string) error {
	res, err := p.DB.Exec("UPDATE car SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL", id)
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

	return errors.New("car not found")
}
