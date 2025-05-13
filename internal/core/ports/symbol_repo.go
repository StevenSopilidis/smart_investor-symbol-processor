package ports

import "github.com/stevensopi/smart_investor/symbol_processor/internal/core/domain"

type ISymbolRepo interface {
	Put(data domain.SymbolData) error
	Get(ticker string, ammount int) ([]domain.SymbolData, error)
}
