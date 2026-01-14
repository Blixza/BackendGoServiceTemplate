package town_domain

import "github.com/google/uuid"

type Model struct {
	ID              uuid.UUID
	Name            string
	Balance         int
	OwnerNickname   string
	XCoordOverworld int
	YCoordOverworld int
	ZCoordOverworld int
	XCoordNether    int
	YCoordNether    int
	ZCoordNether    int
}

func ListToString(towns []Model) string {
	result := ""
	for i := range towns {
		if i != 1 {
			result += "," + towns[i].Name
		} else {
			result += towns[i].Name
		}
	}

	return result
}
