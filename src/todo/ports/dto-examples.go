package todo_ports

import (
	"time"

	gonull "github.com/LukaGiorgadze/gonull/v2"
)

// these "Data Transfert Objects" are not ports/adapters, but are used to transfer data between different tiers of the application;

// examples

type Dbmapping1DTO struct {
	Id uint32
	// numbers
	Anint           *int32
	Abigint         *int64
	Anintunsighed   *uint32
	Abigintunsigned *uint64
	Adecimal        *float64
	Afloat          *float32
	Adouble         *float64
	// booleans
	Aboolean *bool
	// strings
	Avarchar  *string
	Adatetime *time.Time
	Ablob     *[]byte
}

type Dbmapping2DTO struct {
	Id uint32
	// numbers
	Anint           gonull.Nullable[int32]
	Abigint         gonull.Nullable[int64]
	Anintunsighed   gonull.Nullable[uint32]
	Abigintunsigned gonull.Nullable[uint64]
	Adecimal        gonull.Nullable[float64]
	Afloat          gonull.Nullable[float32]
	Adouble         gonull.Nullable[float64]
	// booleans
	Aboolean gonull.Nullable[bool]
	// strings
	Avarchar  gonull.Nullable[string]
	Adatetime gonull.Nullable[time.Time]
	Ablob     gonull.Nullable[[]byte]
}
