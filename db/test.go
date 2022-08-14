package db

import (
	"context"
	"githab.com/kbats183/argotech/backend/models"
)

func (db Database) GetTestByID(testID int) (models.Test, error) {
	test := models.Test{Questions: make([]models.TestQuestion, 0)}
	rows, err := db.Conn.Query("SELECT * FROM test_questions WHERE test_id=$1 ORDER BY id;", testID)
	if err != nil {
		return test, err
	}
	for rows.Next() {
		var question models.TestQuestion
		err := question.ScanRow(rows)
		if err != nil {
			return test, err
		}
		test.Questions = append(test.Questions, question)
	}
	return test, nil
}

func (db Database) GetTestAnswers(testID int, userID models.UserID) (models.TestAnswerList, error) {
	answers := make(models.TestAnswerList, 0)
	rows, err := db.Conn.Query("SELECT ta.* FROM test_answers ta INNER JOIN test_questions tq on ta.question_id = tq.id WHERE user_id=$1 AND test_id=$2 ORDER BY question_id;", userID, testID)
	if err != nil {
		return answers, err
	}
	for rows.Next() {
		var answer models.TestAnswer
		err := answer.ScanRow(rows)
		if err != nil {
			return answers, err
		}
		answers = append(answers, answer)
	}
	return answers, nil
}

func (db Database) AddTestAnswer(userID models.UserID, answer models.TestAnswerUserData) error {
	_, err := db.Conn.Exec("INSERT INTO test_answers(user_id, question_id, answer) VALUES ($1, $2, $3) ON CONFLICT (user_id, question_id) DO UPDATE SET answer = EXCLUDED.answer;",
		userID, answer.QuestionID, answer.Answer)
	return err
}

func (db Database) GetTestAnswersCount(testID int, userID models.UserID) (models.TestAnswerCount, error) {
	var count models.TestAnswerCount
	err := db.Conn.QueryRow("SELECT count(*) FROM test_answers INNER JOIN test_questions tq on tq.id = test_answers.question_id WHERE test_id = $1 AND user_id = $2;",
		testID, userID).Scan(&count)
	return count, err
}

func (db Database) AddTestAnswers(userID models.UserID, answers models.TestAnswerUserDataList) error {
	tx, err := db.Conn.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, answer := range answers {
		_, err := tx.Exec("INSERT INTO test_answers(user_id, question_id, answer) VALUES ($1, $2, $3) ON CONFLICT (user_id, question_id) DO UPDATE SET answer = EXCLUDED.answer;",
			userID, answer.QuestionID, answer.Answer)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}
	return tx.Commit()
}

func (db Database) GetProfessionByTest(userID models.UserID, testID int) (models.ProfessionWithRatingList, error) {
	list := make(models.ProfessionWithRatingList, 0)
	query := `WITH scores AS (SELECT ttp.profession_id, a.answer * ttp.correlation as score
                FROM test_answers a
                         INNER JOIN test_to_profession ttp on a.question_id = ttp.question_id
                         INNER JOIN test_questions q on q.id = a.question_id
                WHERE a.user_id = $1
                  AND q.test_id = $2)
SELECT p.*, exists(SELECT 1 FROM favourite_professions fp WHERE p.id = fp.profession_id AND fp.user_id = $1) as "is_favourite"
FROM professions p
         INNER JOIN scores sc on p.id = sc.profession_id
GROUP BY p.id, is_favourite
ORDER BY sum(sc.score) DESC
LIMIT 5;`
	rows, err := db.Conn.Query(query, userID, testID)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var profession models.ProfessionWithRating
		err := profession.ScanRow(rows)
		if err != nil {
			return list, err
		}
		list = append(list, profession)
	}
	return list, nil
}
