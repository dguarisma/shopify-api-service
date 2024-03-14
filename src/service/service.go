package service

import (
	"desarrollosmoyan/lambda/src/controller"
	"desarrollosmoyan/lambda/src/response"
	"errors"
	"net/http"
)

var (
	ErrEmptyBody      = errors.New("todas las peticiones necesitan body excepto 'GET'")
	ErrNotJsonOrArray = errors.New("tormato incorrecto, debe ser tipo json/array")
	ErrNotJson        = errors.New("formato incorrecto, debe ser tipo json")
)

type Service struct {
	filters IFilters
	repo    controller.CRUD
}

func NewService(filter IFilters, repo controller.CRUD) *Service {
	return &Service{
		filters: filter,
		repo:    repo,
	}
}

func (s *Service) Get(queryParams map[string]string) response.Result {
	if s.filters.IsQueryParamsEmpty(queryParams) {
		return s.repo.GetAll()
	}

	if resPath, ok := s.filters.IsSearchByMany(queryParams); ok {
		return s.repo.GetManyBy(resPath)
	}

	if id, IsSearchById := s.filters.IsSearchById(queryParams); IsSearchById {
		return s.repo.GetByID(id)
	}

	resPath := s.filters.AdaptSearch(queryParams)
	return s.repo.GetManyBy(resPath)
}

func (s *Service) Insert(msg string) response.Result {
	body := s.filters.StringToArrByte(msg)

	if s.filters.IsEmpty(body) {
		return response.NewFailResult(ErrEmptyBody, http.StatusBadRequest)
	}

	if s.filters.IsArray(body) {
		return s.repo.InsertMany(body)
	}

	if s.filters.IsObject(body) {
		return s.repo.InsertOne(body)
	}

	return response.NewFailResult(ErrNotJsonOrArray, http.StatusBadRequest)
}

func (s *Service) Update(msg string) response.Result {
	body := s.filters.StringToArrByte(msg)

	if s.filters.IsEmpty(body) {
		return response.NewFailResult(ErrEmptyBody, http.StatusBadRequest)
	}

	if s.filters.IsObject(body) {
		return s.repo.Update(body)
	}
	return response.NewFailResult(ErrNotJson, http.StatusBadRequest)
}

func (s *Service) Delete(msg string) response.Result {
	body := s.filters.StringToArrByte(msg)
	if s.filters.IsEmpty(body) {
		return response.NewFailResult(ErrEmptyBody, http.StatusBadRequest)
	}

	if s.filters.IsObject(body) {
		return s.repo.DeleteById(body)
	}
	return response.NewFailResult(ErrNotJson, http.StatusBadRequest)
}
