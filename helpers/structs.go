package helpers

// Structure for storing a method
type Method struct {
	Name        string   `json:"name"`
	Params      []string `json:"params"`
	Description string   `json:"description"`
}

// Structure for storing a service
type Service struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	ServiceName string   `json:"service_name"`
	Methods     []Method `json:"methods"`
}
