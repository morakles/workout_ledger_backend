package param_workout_exercise

type AddWorkoutExerciseDTO struct {
	WorkoutID     int64
	ExerciseID    int64
	ExerciseOrder int
}

type WorkoutExerciseDTO struct {
	WorkoutID     int64  `json:"workout_id"`
	ExerciseID    int64  `json:"exercise_id"`
	ExerciseOrder int    `json:"exercise_order"`
	ExerciseName  string `json:"exercise_name,omitempty"`
}
