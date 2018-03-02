package cal

import (
	"revelforce/internal/platform/domain/models"
	"strings"
)

const (
	dateFormat     = "20060102T150405"
	dateLayout     = "20060102"
	dateTimeLayout = "20060102T150405"
)

func getAddress(t *models.Trip) string {
	address := ""

	if len(t.Venues) > 0 {
		for _, v := range t.Venues {
			if v.Primary {
				address = v.Name + ", " + v.Address + ", " + v.City + ", " + v.State + ", " + v.Zip
			}
		}

		if address == "" {
			v := t.Venues[0]
			address = v.Name + ", " + v.Address + ", " + v.City + ", " + v.State + ", " + v.Zip
		}
	}

	return address
}

func stripSpaces(s string) string {
	return strings.Replace(s, " ", "+", -1)
}
