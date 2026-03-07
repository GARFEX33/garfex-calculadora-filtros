// internal/pdf/domain/empresa.go
package domain

// EmpresaPresentacion contiene los datos de la empresa para mostrar en la memoria de cálculo.
type EmpresaPresentacion struct {
	ID              string
	NombreCompleto  string
	LogoPath        string
	Direccion       string
	Telefono        string
	Email           string
	ColorPrimario   string
	ColorSecundario string
}

// empresasCatalogo es el catálogo estático de empresas disponibles para presentación.
// Los datos son placeholder para MVP — se actualizarán con datos reales antes de producción.
var empresasCatalogo = map[string]EmpresaPresentacion{
	"garfex": {
		ID:              "garfex",
		NombreCompleto:  "GARFEX",
		LogoPath:        "garfex.png",
		Direccion:       "",
		Telefono:        "+52 55 1193-0515",
		Email:           "jcgarcia@garfex.mx",
		ColorPrimario:   "#7C0000", /* rojoGarfex */
		ColorSecundario: "#F4CF00", /* amarilloGarfex */
	},
	"summaa": {
		ID:              "summaa",
		NombreCompleto:  "GRUPO SUMMAA ENERGIA",
		LogoPath:        "summaa.png",
		Direccion:       "CALLE TEZIUTLAN #43, COL. SAN LUCAS, DELEGACION COYOACAN, CP 04030, CDMX",
		Telefono:        "+52 55 5243-9127/28",
		Email:           "ventas@summaa.com",
		ColorPrimario:   "#004A99",
		ColorSecundario: "#1B75BB",
	},
	"siemens": {
		ID:              "siemens",
		NombreCompleto:  "Siemens S.A. de C.V.",
		LogoPath:        "siemens.png",
		Direccion:       "Lago Alberto 319, Anáhuac I Secc, Miguel Hidalgo, CDMX, C.P. 11320",
		Telefono:        "+52 55 5229-3600",
		Email:           "contacto@siemens.com.mx",
		ColorPrimario:   "#009999",
		ColorSecundario: "#000000",
	},
}

// BuscarEmpresaPorID busca una empresa en el catálogo estático por su ID.
// Retorna la empresa y true si existe; retorna zero value y false si no existe.
func BuscarEmpresaPorID(id string) (EmpresaPresentacion, bool) {
	empresa, ok := empresasCatalogo[id]
	return empresa, ok
}
