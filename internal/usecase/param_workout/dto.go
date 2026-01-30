package param_workout

type CreateParamWorkoutDTO struct {
	Name                   string
	NumberOfSets           int
	RestBetweenSetsSeconds int
}

type ParamWorkoutDTO struct {
	ID                     int64
	Name                   string
	NumberOfSets           int
	RestBetweenSetsSeconds int
}
