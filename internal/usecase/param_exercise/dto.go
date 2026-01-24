package param_exercise

type CreateParamExerciseDTO struct {
	Name    string
	IconUrl string
}

type CreateParamExerciseResultDTO struct {
	ID int64
}

type ParamExerciseDTO struct {
	ID      int64
	Name    string
	IconUrl string
}
