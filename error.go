package openrouter

type apiError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Param   string `json:"param"`
		Type    string `json:"type"`
	} `json:"error"`
}
