package server

//dico stringer
//config.toml
// fieldprefix = "Code"

// [[strings]]
// name = "not found"

// [[strings]]
// name = "unknown"

// [[strings]]
// name = "success"

// [[strings]]
// name = "error"

// [[strings]]
// name = "invalid data"

// [[strings]]
// name = "existing"

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  stringer

var CodeNotFound = "not_found"
var CodeUnknown = "unknown"
var CodeSuccess = "success"
var CodeError = "error"
var CodeInvalidData = "invalid_data"
var CodeExisting = "existing"

//<<<AUTOGENERATE.DICO

// NewFail
// Если один параметр, то он уточняющим статусом
// Если два - первый статус, второй сообщение
// По умолчанию статус CodeError
func F(args ...string) ResponseDTO {
	if len(args) == 1 {
		return ResponseDTO{args[0], "", nil}
	}

	if len(args) == 2 {
		return ResponseDTO{args[0], args[1], nil}
	}

	return ResponseDTO{CodeError, "", nil}
}

// NewSuccess
// Если один параметр, то он уточняющим статусом
// Если два - первый статус, второй сообщение
// По умолчанию статус CodeSuccess
func OK(data interface{}, args ...string) ResponseDTO {
	if len(args) == 1 {
		return ResponseDTO{args[0], "", data}
	}

	if len(args) == 2 {
		return ResponseDTO{args[0], args[1], data}
	}

	return ResponseDTO{CodeSuccess, "", data}
}

type ResponseDTO struct {
	Code    string      `json:",omitempty"`
	Message string      `json:",omitempty"`
	Data    interface{} `json:",omitempty"`
}
