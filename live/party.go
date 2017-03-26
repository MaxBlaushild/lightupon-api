package live

type Party struct {
	Passcode string
	Users map[uint]*Connection
	Objectives []Objective
	CurrentObjectiveIndex int
}

func (p *Party) CurrentObjective() (currentObjective Objective) {
	return p.Objectives[p.CurrentObjectiveIndex]
}

func (p *Party) Exists() (exists bool) {
	exists = (len(p.Passcode) > 0)
	return
}

func (p *Party) CurrentObjectiveMet() (hasBeenMet bool) {
	currentObjective := p.CurrentObjective()
	hasBeenMet = currentObjective.HasBeenMet(p)
	return
}

func (p *Party) Push(response Response) {
	response.NextSceneAvailable = p.CurrentObjectiveMet()
	response.CurrentSceneOrder = p.CurrentObjectiveIndex
	for userID := range p.Users {
		select {
		case p.Users[userID].Send <- response:
		default:
		}
	}
}