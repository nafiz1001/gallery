package arts

type Art struct {
	Genres   []string `json:"genres"`
	Id       int      `json:"id"`
	Picture  string   `json:"picture"`
	Price    int      `json:"price"`
	Quantity int      `json:"quantity"`
	Title    string   `json:"title"`
}
