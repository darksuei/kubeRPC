package helpers

type Method struct {
	Name        string   `json:"name"`
	Params      []string `json:"params"`
	Description string   `json:"description"`
}

type Service struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	ServiceName string   `json:"service_name"`
	Methods     []Method `json:"methods"`
}

type ServicePtr struct {
	Host        *string   `json:"host"`
	Port        *int      `json:"port"`
	ServiceName *string   `json:"service_name"`
	Methods     *[]Method `json:"methods"`
}
