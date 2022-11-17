package enumfilter_test

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
