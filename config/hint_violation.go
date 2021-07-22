package config

const (
	HintApd    = "Sopir tidak menggunakan APD"
	HintOut    = "Sopir turun dari armada"
	HintSign   = "Sopir melanggar rambu-rambu / marka jalan"
	HintKernet = "Membawa kernet atau tumpangan orang"
	HintPark   = "Parkir tidak pada tempatnya"
	HintTrapic = "Melawan arus lalu lintas"
	HintDerm   = "Melintas di daerah dermaga, bukan truck lossing"
	HintCont   = "Kontainer belum ditutup rapat"
	HintDoor   = "Posisi pintu kontainer salah"
	HintWeapon = "Membawa senjata tajam"
	HintDrug   = "Membawa atau menggunakan obat terlarang"
	HintSpeed  = "Melebihi batas maksimal kecepatan"
)

func GetHintAvailable() []string {
	return []string{
		HintApd,
		HintSign,
		HintOut,
		HintKernet,
		HintPark,
		HintTrapic,
		HintDerm,
		HintCont,
		HintDoor,
		HintWeapon,
		HintDrug,
		HintSpeed,
	}
}
