package storage

import (
	"car_rental/genprotos/brand"
	"car_rental/genprotos/car"
)

type StorageI interface {
	CreateBrand(id string, req *brand.CreateBrandRequest) (res *brand.Brand, err error)
	GetBrandByID(req *brand.GetBrandByIDRequest) (*brand.Brand, error)
	GetBrandList(req *brand.GetBrandListRequest) (*brand.GetBrandListResponse, error)
	UpdateBrand(id string, req *brand.UpdateBrandRequest) error
	DeleteBrand(req *brand.DeleteBrandRequest) error

	CreateCar(id string, req *car.CreateCarRequest) error
	GetCarByID(id string) (*car.GetCarByIDResponse, error)
	GetCarList(offset, limit int, search string) (*car.GetCarListResponse, error)
	UpdateCar(id string, entity *car.UpdateCarRequest) error
	DeleteCar(id string) error
}
