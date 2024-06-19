package responses

type Pattern struct {
	Data  any   `json:"data" extensions:"x-nullable"`
	Error error `json:"error" extensions:"x-nullable"`
}
