package photos

type Environment string

const (
	EnvironmentProduction Environment = "production"
	EnvironmentLocal      Environment = "local"
)

var allEnvironments = []Environment{
	EnvironmentLocal, EnvironmentProduction,
}

func (e Environment) IsProduction() bool {
	return e == EnvironmentProduction
}

func (e Environment) Valid() bool {
	for _, ee := range allEnvironments {
		if e == ee {
			return true
		}
	}

	return false
}
