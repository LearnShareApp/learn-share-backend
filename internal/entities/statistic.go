package entities

type StudentStatistic struct {
	CountOfFinishedLesson     int `db:"count_of_finished_lesson"`
	CountOfVerificationLesson int `db:"count_of_verification_lesson"`
	CountOfWaitingLesson      int `db:"count_of_waiting_lesson"`
	CountOfTeachers           int `db:"count_of_teachers"`
}

type TeacherStatistic struct {
	CountOfFinishedLesson int `db:"count_of_finished_lesson"`
	CountOfStudents       int `db:"count_of_students"`
}
