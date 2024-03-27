package domain

type Carrie struct {
	Id           int    `json:"id"`
	Cid          string `json:"cid" validate:"required"`
	Company_name string `json:"company_name" validate:"required"`
	Address      string `json:"address" validate:"required"`
	Telephone    string `json:"telephone" validate:"required"`
	Locality_id  string `json:"locality_id" validate:"required"`
}

type CarrieLocality struct {
	Locality_id   string `json:"locality_id" validate:"required"`
	Locality_name string `json:"local_name" validate:"required"`
	Cant_carries  int    `json:"carries_count" validate:"required"`
}
