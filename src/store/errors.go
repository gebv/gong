package store

import "fmt"

var (
	CollNameBucket = "bucket"
	CollNameFile   = "file"
)

//dico errors
//config.toml
//[[errors]]
//name="not found"

//[[errors]]
//name="unknown"
//message="unknown error, see details in the log"
//comment="unknown error, see details in the log"

//[[errors]]
//name="not allowed"

//[[errors]]
//name="not supported"

//[[errors]]
//name="not valid_data"

//[[errors]]
//name="rejected"

//config.toml
//AUTOGENERATE.DICO>>>
//	The text in the section 'AUTOGENERATE.DICO' automatically generated, please do not edit it
//[DICO.VERSION]:	 0.0.2
//[DICO.COMMAND]:	  errors

var ErrNotFound = fmt.Errorf("not_found")

// ErrUnknown unknown error, see details in the log.
var ErrUnknown = fmt.Errorf("unknown error, see details in the log")

var ErrNotAllowed = fmt.Errorf("not_allowed")

var ErrNotSupported = fmt.Errorf("not_supported")

var ErrNotValidData = fmt.Errorf("not_valid_data")

var ErrRejected = fmt.Errorf("rejected")

//<<<AUTOGENERATE.DICO
