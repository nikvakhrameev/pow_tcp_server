package pow

const powDifficulty = 7

type DifficultyStorage struct{}

func NewDifficultyStorage() DifficultyStorage {
	return DifficultyStorage{}
}

func (d DifficultyStorage) GetDifficulty() int {
	return powDifficulty
}
