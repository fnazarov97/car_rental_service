package car

import (
	"car_rental/genprotos/car"
	"car_rental/storage"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CarService is a struct that implements the server interface
type CarService struct {
	Stg storage.StorageI
	car.UnimplementedCarServiceServer
}

// CreateCar ...
func (a *CarService) CreateCar(c context.Context, req *car.CreateCarRequest) (*car.Car, error) {
	id := uuid.New()
	err := a.Stg.CreateCar(id.String(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.CreateCar: %s", err.Error())
	}
	res, err := a.Stg.GetCarByID(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.Stg.GetCarByID: %s", err.Error())
	}
	return &car.Car{
		CarId:     res.CarId,
		Model:     res.Model,
		Color:     res.Color,
		Year:      res.Year,
		Mileage:   res.Mileage,
		BrandId:   res.Brand.BrandId,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

// GetCarByID ...
func (a *CarService) GetCarByID(c context.Context, req *car.GetCarByIDRequest) (*car.GetCarByIDResponse, error) {
	res, err := a.Stg.GetCarByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.GetCarByID: %s", err.Error())
	}
	return res, nil
}

// GetCarList ...
func (a *CarService) GetCarList(c context.Context, req *car.GetCarListRequest) (*car.GetCarListResponse, error) {
	res, err := a.Stg.GetCarList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.GetCarList: %s", err.Error())
	}
	return res, nil
}

// UpdateCar ...
func (a *CarService) UpdateCar(c context.Context, req *car.UpdateCarRequest) (*car.Car, error) {
	err := a.Stg.UpdateCar(req.Id, req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.UpdateCar: %s", err.Error())
	}
	res, e := a.Stg.GetCarByID(req.Id)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.UpdateCar: %s", e.Error())
	}
	return &car.Car{
		CarId:     res.CarId,
		Model:     res.Model,
		Color:     res.Color,
		Year:      res.Year,
		Mileage:   res.Mileage,
		BrandId:   res.Brand.BrandId,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

// DeleteCar ...
func (a *CarService) DeleteCar(c context.Context, req *car.DeleteCarRequest) (*car.Car, error) {
	res, e := a.Stg.GetCarByID(req.Id)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.UpdateCar: %s", e.Error())
	}
	err := a.Stg.DeleteCar(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "a.Stg.DeleteCar: %s", err.Error())
	}

	return &car.Car{
		CarId:     res.CarId,
		Model:     res.Model,
		Color:     res.Color,
		Year:      res.Year,
		Mileage:   res.Mileage,
		BrandId:   res.Brand.BrandId,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}
