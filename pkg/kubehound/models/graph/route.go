package graph

type Route struct {
	StoreID      string `json:"storeID" mapstructure:"storeID"`
	App          string `json:"app" mapstructure:"app"`
	Team         string `json:"team" mapstructure:"team"`
	Service      string `json:"service" mapstructure:"service"` // TODO[TG] this is probably wrong
	RunID        string `json:"runID" mapstructure:"runID"`
	Cluster      string `json:"cluster" mapstructure:"cluster"`
	IsNamespaced bool   `json:"isNamespaced" mapstructure:"isNamespaced"`
	Namespace    string `json:"namespace" mapstructure:"namespace"`
	Name         string `json:"name" mapstructure:"name"`
	// TODO[TG] add more
}
