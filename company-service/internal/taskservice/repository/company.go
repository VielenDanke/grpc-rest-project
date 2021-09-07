package repository

import "context"

type CompanyRepository interface {
	GetByBin(ctx context.Context, bin string) (string, error)
}
