package httpResponse

type Success struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type Error struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Errors  interface{} `json:"errors,omitempty"`
}