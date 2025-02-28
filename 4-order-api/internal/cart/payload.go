package cart

import "github.com/lib/pq"

type OrderRequest struct {
	Products pq.StringArray `json:"products"`
}

