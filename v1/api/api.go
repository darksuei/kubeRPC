package api

type Request struct {
    Method string      `json:"method"`
    Params interface{} `json:"params"`
}

type Response struct {
    Result interface{} `json:"result"`
    Error  string      `json:"error,omitempty"`
}