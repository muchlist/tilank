package config

const (
	TypeAPD        = "APD"
	TypeProccedure = "PROSEDUR"
	TypeCrime      = "PRILAKU"
	TypeBehavior   = "KRIMINAL"
	TypeOther      = "LAINNYA"
)

func GetTypeAvailable() []string {
	return []string{TypeAPD, TypeProccedure, TypeCrime, TypeBehavior, TypeOther}
}
