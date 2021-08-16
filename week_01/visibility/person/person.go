package person

var (
	Public  = 1
	private = 1
)

type Person struct {
	ID     int // public
	Name   string
	secret string // private
}

func (p *Person) SetSecret(newSecret string) {
	// struct method, setter
	p.secret = newSecret
}
