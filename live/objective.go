package live

const threshold float64 = 0.05 // 0.05 km = 50 meters

type Objective struct {
	Location Location
}

func (o *Objective) HasBeenMet(party *Party) (hasBeenMet bool) {
	for userID := range party.Users {
		c := party.Users[userID]
		hasBeenMet = hasBeenMet || o.isThere(c)
	}
	return 
}

func (o *Objective) isThere(c *Connection) (isAtNextScene bool) {
	distanceFromScene := CalculateDistance(o.Location, c.Location)
	isAtNextScene = distanceFromScene < threshold
	return
}