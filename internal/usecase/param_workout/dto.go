package param_workout

type CreateParamWorkoutDTO struct {
	UserID                 int64
	Name                   string
	NumberOfSets           int
	RestBetweenSetsSeconds int
}

type ParamWorkoutDTO struct {
	ID                     int64
	UserID                 int64
	Name                   string
	NumberOfSets           int
	RestBetweenSetsSeconds int
}
