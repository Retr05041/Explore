package playerhandler

type Player struct {
    Name string
    Inventory []string
}

func NewPlayer(name string) (*Player) {
    tmpPlayer := new(Player)
    tmpPlayer.Name = name
    tmpPlayer.Inventory = append(tmpPlayer.Inventory, "TestItem1")
    tmpPlayer.Inventory = append(tmpPlayer.Inventory, "TestItem2")
    return tmpPlayer
}

func (p *Player) IsInInv(item string) (bool) {
    for _, token := range p.Inventory {
       if token == item { return true } 
    } 
    return false
}

func (p *Player) AddToInv(item string) {
    p.Inventory = append(p.Inventory, item)
}
