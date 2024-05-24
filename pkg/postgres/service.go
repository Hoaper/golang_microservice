package postgres

import (
	"EffectiveMobile/pkg/postgres/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetCar(regNum string) (*models.Car, error) {
	var car models.Car
	err := s.db.Where("reg_num = ?", regNum).First(&car).Error
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (s *Service) CreatePeople(name string, surname string) (*models.People, error) {
	people := &models.People{Name: name, Surname: surname}
	err := s.db.Create(people).Error
	if err != nil {
		return nil, err
	}
	return people, nil
}

func (s *Service) GetCars(filter models.Car, page int, pageSize int) ([]models.Car, error) {
	var cars []models.Car
	err := s.db.Where(&filter).Offset((page - 1) * pageSize).Limit(pageSize).Find(&cars).Error
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (s *Service) DeleteCarByRegNum(regNum string) error {
	return s.db.Where("reg_num = ?", regNum).Delete(&models.Car{}).Error
}

func (s *Service) UpdateCar(id uint, updatedCar models.Car) error {
	return s.db.Model(&models.Car{}).Where("id = ?", id).Updates(updatedCar).Error
}

func (s *Service) AddCars(regNums []string, peopleID uint) error {
	var people models.People
	if err := s.db.First(&people, peopleID).Error; err != nil {
		return err
	}

	for _, regNum := range regNums {
		car := models.Car{RegNum: regNum, Owner: people}
		if err := s.db.Create(&car).Error; err != nil {
			return err
		}
	}
	return nil
}
