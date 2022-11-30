package enum

type EnumAdd string

const (
	EnumAddA EnumAdd = "A"
	EnumAddB EnumAdd = "B"
)

func PossibleEnumAddValues() []EnumAdd {
	return []EnumAdd{
		EnumAddA,
		EnumAddB,
	}
}

type EnumExist string

const (
	EnumExistA EnumExist = "A"
	EnumExistB EnumExist = "B"
)

func PossibleEnumExistValues() []EnumExist {
	return []EnumExist{
		EnumExistA,
		EnumExistB,
	}
}
