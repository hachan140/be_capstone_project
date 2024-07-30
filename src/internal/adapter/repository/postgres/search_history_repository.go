package postgres

import (
	"be-capstone-project/src/internal/adapter/repository/postgres/model"
	"be-capstone-project/src/internal/core/storage"
	"time"
)

type SearchHistoryRepository struct {
	storage *storage.Database
}

type ISearchHistoryRepository interface {
	GetAllSearchHistoryPersonalize(userID uint, input string) ([]string, error)
	GetExistedHistoryOfUser(userID uint, keyword string) (*model.SearchHistory, error)
	UpdateExistedKeyword(searchHistoryID uint) error
	SaveSearchHistory(searchHistory *model.SearchHistory) error
}

func NewSearchHistoryRepository(storage *storage.Database) ISearchHistoryRepository {
	return &SearchHistoryRepository{storage: storage}
}

func (s *SearchHistoryRepository) GetAllSearchHistoryPersonalize(userID uint, input string) ([]string, error) {
	var keywords []string
	err := s.storage.Raw("select keywords from search_histories where (user_id = ? or user_id = 0) and keywords like ? order by updated_at desc", userID, "%"+input+"%").Scan(&keywords).Error
	if err != nil {
		return nil, err
	}
	return keywords, nil
}

func (s *SearchHistoryRepository) GetExistedHistoryOfUser(userID uint, keyword string) (*model.SearchHistory, error) {
	var searchHistory *model.SearchHistory
	err := s.storage.Raw("select * from search_histories where user_id = ? and LOWER(keywords) = ?", userID, keyword).Scan(&searchHistory).Error
	if err != nil {
		return nil, err
	}
	return searchHistory, nil
}

func (s *SearchHistoryRepository) UpdateExistedKeyword(searchHistoryID uint) error {
	err := s.storage.Exec("update search_histories set updated_at = ? where id = ?", time.Now(), searchHistoryID).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *SearchHistoryRepository) SaveSearchHistory(searchHistory *model.SearchHistory) error {
	err := s.storage.Model(searchHistory).Create(&searchHistory).Error
	if err != nil {
		return err
	}
	return nil
}
