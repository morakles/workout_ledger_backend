package param_workout_exercise

type AddWorkoutExerciseDTO struct {
	WorkoutID          int64
	ExerciseID         int64
	ExerciseOrder      int
	DefaultSets        *int
	DefaultRestSeconds *int
	DefaultReps        *int
}

type WorkoutExerciseDTO struct {
	WorkoutID          int64  `json:"workout_id"`
	ExerciseID         int64  `json:"exercise_id"`
	ExerciseOrder      int    `json:"exercise_order"`
	ExerciseName       string `json:"exercise_name,omitempty"`
	DefaultSets        *int   `json:"default_sets,omitempty"`
	DefaultRestSeconds *int   `json:"default_rest_seconds,omitempty"`
	DefaultReps        *int   `json:"default_reps,omitempty"`
}
