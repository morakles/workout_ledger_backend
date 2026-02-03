package param_workout

type ParamWorkout struct {
	ID                     int64
	UserID                 int64
	Name                   string
	NumberOfSets           int
	RestBetweenSetsSeconds int
}
