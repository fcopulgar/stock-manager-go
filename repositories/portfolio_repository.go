package repositories

import "github.com/fcopulgar/stock-manager-go/models"

type PortfolioRepository interface {
	GetAll() ([]models.Portfolio, error)
	GetByID(id int) (*models.Portfolio, error)
	Save(portfolio *models.Portfolio) error
	Update(portfolio *models.Portfolio) error
	Delete(id int) error
}
