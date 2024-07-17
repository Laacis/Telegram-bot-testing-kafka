package order_generation_service

type Product struct {
	Id              int     `json:"id"`
	ProductKey      string  `json:"productKey"`
	Name            string  `json:"name"`
	Manufacturer    string  `json:"manufacturer"`
	ThermalCategory string  `json:"thermalCategory"` //TODO make an enum later?!
	BuyPrice        float64 `json:"buyPrice"`
	UnitsPerPallet  int     `json:"unitsPerPallet"`
	UnitWeightKg    int     `json:"unitWeightKg"`
}
