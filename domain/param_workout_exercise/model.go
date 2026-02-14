package param_workout_exercise

type WorkoutExercise struct {
	WorkoutID          int64
	ExerciseID         int64
	ExerciseOrder      int
	ExerciseName       string
	DefaultSets        *int
	DefaultRestSeconds *int
	DefaultReps        *int
}
