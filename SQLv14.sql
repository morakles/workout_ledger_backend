CREATE TABLE IF NOT EXISTS "exercises_in_workout" (
	"workout_id" BIGINT NOT NULL,
	"exercise_id" BIGINT NOT NULL,
	"exercise_order" INT NOT NULL,
	PRIMARY KEY ("workout_id", "exercise_id"),
	UNIQUE ("workout_id", "exercise_order"),
	FOREIGN KEY ("workout_id") REFERENCES "param_workout"("id") ON DELETE CASCADE,
	FOREIGN KEY ("exercise_id") REFERENCES "param_exercise"("id") ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS "idx_exercises_in_workout_workout_id"
ON "exercises_in_workout" ("workout_id");

CREATE INDEX IF NOT EXISTS "idx_exercises_in_workout_exercise_id"
ON "exercises_in_workout" ("exercise_id");
