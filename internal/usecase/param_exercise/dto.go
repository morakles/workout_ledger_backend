package param_exercise

type CreateParamExerciseDTO struct {
	Name    string
	IconUrl string
}

type CreateParamExerciseResultDTO struct {
	ID int64
}

type UpdateParamExerciseDTO struct {
	ID      int64
	Name    string
	IconUrl string
}

type ParamExerciseDTO struct {
	ID      int64
	Name    string
	IconUrl string
}
