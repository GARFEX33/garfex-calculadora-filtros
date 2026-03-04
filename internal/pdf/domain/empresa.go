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
		NombreCompleto:  "Garfex",
		LogoPath:        "assets/logos/garfex.png",
		Direccion:       "Av. Insurgentes Sur 1234, Col. Del Valle, CDMX, C.P. 03100",
		Telefono:        "+52 55 1193-0515",
		Email:           "jcgarcia@garfex.mx",
		ColorPrimario:   "#1a3a5c",
		ColorSecundario: "#2e7d32",
	},
	"summa": {
		ID:              "summa",
		NombreCompleto:  "Summa Ingeniería Eléctrica S.A. de C.V.",
		LogoPath:        "assets/logos/summa.png",
		Direccion:       "Blvd. Manuel Ávila Camacho 800, Lomas de Chapultepec, CDMX, C.P. 11000",
		Telefono:        "+52 55 9876-5432",
		Email:           "contacto@summa.mx",
		ColorPrimario:   "#b71c1c",
		ColorSecundario: "#fbc02d",
	},
	"siemens": {
		ID:              "siemens",
		NombreCompleto:  "Siemens S.A. de C.V.",
		LogoPath:        "assets/logos/siemens.png",
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
