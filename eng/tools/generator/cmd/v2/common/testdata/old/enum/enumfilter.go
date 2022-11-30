package enum

type EnumRemove string

const (
	EnumRemoveA EnumRemove = "A"
	EnumRemoveB EnumRemove = "B"
	EnumRemoveC EnumRemove = "C"
)

func PossibleEnumRemoveValues() []EnumRemove {
	return []EnumRemove{
		EnumRemoveA,
		EnumRemoveB,
		EnumRemoveC,
	}
}

type EnumExist string

const (
	EnumExistA EnumExist = "A"
)

func PossibleEnumExistValues() []EnumExist {
	return []EnumExist{
		EnumExistA,
	}
}
