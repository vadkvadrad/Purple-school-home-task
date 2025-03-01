package cart

import "github.com/lib/pq"

type OrderRequest struct {
	Products pq.Int64Array `json:"products"`
}
