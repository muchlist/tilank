package config

const (
	TypeAPD        = "APD"
	TypeSign       = "RAMBU-RAMBU"
	TypeProccedure = "PROSEDUR"
	TypeCrime      = "PRILAKU"
	TypeBehavior   = "KRIMINAL"
	TypeOther      = "LAINNYA"
)

func GetTypeAvailable() []string {
	return []string{TypeAPD, TypeSign, TypeProccedure, TypeCrime, TypeBehavior, TypeOther}
}
