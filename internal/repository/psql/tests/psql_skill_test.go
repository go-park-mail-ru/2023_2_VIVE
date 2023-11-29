package psql

import (
	"HnH/internal/repository/psql"
	"HnH/pkg/testHelper"
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var testAddSkillsByVacID = []struct {
	vacancyID int
	skills    []string
}{
	{
		vacancyID: 1,
		skills:    []string{"python", "c++"},
	},
	{
		vacancyID: 2,
		skills:    []string{"js", "golang"},
	},
}

func TestAddSkillsByVacID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testAddSkillsByVacID {
		for i, skill := range testCase.skills {
			mock.
				ExpectQuery(testHelper.InsertQuery).
				WithArgs(skill).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(i + 1),
				)

			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.vacancyID, i+1).
				WillReturnResult(driver.RowsAffected(1))

		}
		actual := repo.AddSkillsByVacID(ctxWithLogger, testCase.vacancyID, testCase.skills)
		if actual != nil {
			t.Errorf("unexpected err: %v", actual)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return

		}
	}
}

var testAddSkillsByVacIDOnExistance = []struct {
	vacancyID int
	skills    []string
}{
	{
		vacancyID: 1,
		skills:    []string{"python", "c++"},
	},
	{
		vacancyID: 2,
		skills:    []string{"js", "golang"},
	},
}

func TestAddSkillsByVacIDOnExistance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testAddSkillsByVacIDOnExistance {
		for i, skill := range testCase.skills {
			mock.
				ExpectQuery(testHelper.InsertQuery).
				WithArgs(skill).
				WillReturnError(sql.ErrNoRows)

			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(skill).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(i + 1),
				)

			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.vacancyID, i+1).
				WillReturnResult(driver.RowsAffected(1))

		}
		actual := repo.AddSkillsByVacID(ctxWithLogger, testCase.vacancyID, testCase.skills)
		if actual != nil {
			t.Errorf("unexpected err: %v", actual)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return

		}
	}
}

var testAddSkillsByCvID = []struct {
	cvID   int
	skills []string
}{
	{
		cvID:   1,
		skills: []string{"python", "c++"},
	},
	{
		cvID:   2,
		skills: []string{"js", "golang"},
	},
}

func TestAddSkillsByCvID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testAddSkillsByCvID {
		for i, skill := range testCase.skills {
			mock.
				ExpectQuery(testHelper.InsertQuery).
				WithArgs(skill).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(i + 1),
				)

			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.cvID, i+1).
				WillReturnResult(driver.RowsAffected(1))

		}
		actual := repo.AddSkillsByCvID(ctxWithLogger, testCase.cvID, testCase.skills)
		if actual != nil {
			t.Errorf("unexpected err: %v", actual)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return

		}
	}
}

var testAddSkillsByCvIDOnExistance = []struct {
	cvID   int
	skills []string
}{
	{
		cvID:   1,
		skills: []string{"python", "c++"},
	},
	{
		cvID:   2,
		skills: []string{"js", "golang"},
	},
}

func TestAddSkillsByCvIDOnExistance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testAddSkillsByCvIDOnExistance {
		for i, skill := range testCase.skills {
			mock.
				ExpectQuery(testHelper.InsertQuery).
				WithArgs(skill).
				WillReturnError(sql.ErrNoRows)

			mock.
				ExpectQuery(testHelper.SelectQuery).
				WithArgs(skill).
				WillReturnRows(
					sqlmock.NewRows([]string{"id"}).
						AddRow(i + 1),
				)

			mock.
				ExpectExec(testHelper.InsertQuery).
				WithArgs(testCase.cvID, i+1).
				WillReturnResult(driver.RowsAffected(1))

		}
		actual := repo.AddSkillsByVacID(ctxWithLogger, testCase.cvID, testCase.skills)
		if actual != nil {
			t.Errorf("unexpected err: %v", actual)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return

		}
	}
}

var testGetSkillsByVacID = []struct {
	vacancyID      int
	expectedSkills []string
}{
	{
		vacancyID:      1,
		expectedSkills: []string{"python", "c++"},
	},
	{
		vacancyID:      2,
		expectedSkills: []string{"js", "golang"},
	},
}

func TestGetSkillsByVacID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testGetSkillsByVacID {
		rows := sqlmock.NewRows([]string{"skill"})
		for _, skill := range testCase.expectedSkills {
			rows.AddRow(skill)
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.vacancyID).
			WillReturnRows(rows)

		actualSkills, getErr := repo.GetSkillsByVacID(ctxWithLogger, testCase.vacancyID)
		if getErr != nil {
			t.Errorf("unexpected err: %v", getErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actualSkills, testCase.expectedSkills) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedSkills, actualSkills))
		}
	}
}

var testGetSkillsByCvID = []struct {
	cvID           int
	expectedSkills []string
}{
	{
		cvID:           1,
		expectedSkills: []string{"python", "c++"},
	},
	{
		cvID:           2,
		expectedSkills: []string{"js", "golang"},
	},
}

func TestGetSkillsByCvID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := psql.NewPsqlSkillRepository(db)
	for _, testCase := range testGetSkillsByCvID {
		rows := sqlmock.NewRows([]string{"skill"})
		for _, skill := range testCase.expectedSkills {
			rows.AddRow(skill)
		}
		mock.
			ExpectQuery(testHelper.SelectQuery).
			WithArgs(testCase.cvID).
			WillReturnRows(rows)

		actualSkills, getErr := repo.GetSkillsByCvID(ctxWithLogger, testCase.cvID)
		if getErr != nil {
			t.Errorf("unexpected err: %v", getErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if !reflect.DeepEqual(actualSkills, testCase.expectedSkills) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedSkills, actualSkills))
		}
	}
}
