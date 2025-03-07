package prod

import "github.com/lib/pq"

type ProductRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images"`
	Price       int            `json:"price"`
}
