package barnd

import (
	"car_rental/genprotos/brand"
	"car_rental/storage"
	"context"

	"github.com/google/uuid"
)

// BrandService is a struct that implements the server interface
type BrandService struct {
	Stg storage.StorageI
	brand.UnimplementedBrandServiceServer
}

// BrandService ...
func (a *BrandService) CreateBrand(ctx context.Context, req *brand.CreateBrandRequest) (*brand.Brand, error) {
	id := uuid.New()
	res, err := a.Stg.CreateBrand(id.String(), req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (a *BrandService) GetBrandList(ctx context.Context, req *brand.GetBrandListRequest) (*brand.GetBrandListResponse, error) {
	res, err := a.Stg.GetBrandList(req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (a *BrandService) GetBrandByID(ctx context.Context, req *brand.GetBrandByIDRequest) (*brand.Brand, error) {
	res, err := a.Stg.GetBrandByID(req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (a *BrandService) UpdateBrand(ctx context.Context, req *brand.UpdateBrandRequest) (*brand.Brand, error) {
	id := req.Id
	err := a.Stg.UpdateBrand(id, req)
	if err != nil {
		return &brand.Brand{}, err
	}
	r := &brand.GetBrandByIDRequest{BrandId: id}
	res, err := a.Stg.GetBrandByID(r)
	if err != nil {
		return &brand.Brand{}, err
	}
	return res, nil
}

func (a *BrandService) DeleteBrand(ctx context.Context, req *brand.DeleteBrandRequest) (*brand.Brand, error) {
	res, err := a.Stg.GetBrandByID((*brand.GetBrandByIDRequest)(req))
	if err != nil {
		return &brand.Brand{}, err
	}
	err = a.Stg.DeleteBrand(req)
	if err != nil {
		return &brand.Brand{}, err
	}
	return res, nil
}
