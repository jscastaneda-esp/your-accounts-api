package project

import (
	"api-your-accounts/project/domain"
	mocksShared "api-your-accounts/shared/domain/persistent/mocks"
	"api-your-accounts/shared/domain/test_utils"
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestSuite struct {
	suite.Suite
	name       string
	userId     uint
	typeBudget domain.ProjectType
	mock       sqlmock.Sqlmock
	mockTX     *mocksShared.Transaction
	repository domain.ProjectRepository
}

func (suite *TestSuite) SetupSuite() {
	suite.name = "Test"
	suite.userId = 1
	suite.typeBudget = domain.Budget

	require := require.New(suite.T())

	var (
		db  *sql.DB
		err error
	)

	db, suite.mock, err = sqlmock.New()
	require.NoError(err)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(err)

	suite.repository = NewRepository(DB)
}

func (suite *TestSuite) SetupTest() {
	suite.mockTX = mocksShared.NewTransaction(suite.T())
}

func (suite *TestSuite) TearDownTest() {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *TestSuite) TestWithTransactionSuccessNew() {
	require := require.New(suite.T())

	suite.mockTX.On("Get").Return(&gorm.DB{})

	repo := suite.repository.WithTransaction(suite.mockTX)

	require.NotNil(repo)
	require.NotEqual(suite.repository, repo)
}

func (suite *TestSuite) TestWithTransactionSuccessExists() {
	require := require.New(suite.T())

	suite.mockTX.On("Get").Return(&sql.DB{})

	repo := suite.repository.WithTransaction(suite.mockTX)

	require.NotNil(repo)
	require.Equal(suite.repository, repo)
}

func (suite *TestSuite) TestCreateSuccess() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "projects" ("updated_at","name","user_id","type") 
		VALUES ($1,$2,$3,$4) 
		RETURNING "id","created_at"
		`)).
		WithArgs(test_utils.AnyTime{}, suite.name, suite.userId, suite.typeBudget).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(999))
	suite.mock.ExpectCommit()
	project := &domain.Project{
		Name:   suite.name,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}

	res, err := suite.repository.Create(context.Background(), project)

	require.NoError(err)
	require.NotNil(res)
	require.Equal(uint(999), res.ID)
	require.Equal(project.Name, res.Name)
	require.Equal(project.UserId, res.UserId)
	require.Equal(project.Type, res.Type)
}

func (suite *TestSuite) TestCreateError() {
	require := require.New(suite.T())

	suite.mock.ExpectBegin()
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO "projects" ("updated_at","name","user_id","type") 
		VALUES ($1,$2,$3,$4) 
		RETURNING "id","created_at"
		`)).
		WithArgs(test_utils.AnyTime{}, suite.name, suite.userId, suite.typeBudget).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()
	project := &domain.Project{
		Name:   suite.name,
		UserId: suite.userId,
		Type:   suite.typeBudget,
	}

	res, err := suite.repository.Create(context.Background(), project)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.Nil(res)
}

func (suite *TestSuite) TestFindByUserIdSuccess() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "projects"
		WHERE "projects"."user_id" = $1
		ORDER BY created_at desc LIMIT 10
		`)).
		WithArgs(suite.userId).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "user_id", "type"}).
			AddRow(999, time.Now(), time.Now(), suite.name, suite.userId, suite.typeBudget).
			AddRow(1000, time.Now(), time.Now(), suite.name, suite.userId, suite.typeBudget),
		)

	projects, err := suite.repository.FindByUserId(context.Background(), suite.userId)

	require.NoError(err)
	require.NotNil(projects)
	require.Len(projects, 2)
	require.Equal(uint(999), projects[0].ID)
	require.Equal(suite.name, projects[0].Name)
	require.Equal(suite.userId, projects[0].UserId)
	require.Equal(suite.typeBudget, projects[0].Type)
	require.Equal(uint(1000), projects[1].ID)
	require.Equal(suite.name, projects[1].Name)
	require.Equal(suite.userId, projects[1].UserId)
	require.Equal(suite.typeBudget, projects[1].Type)
}

func (suite *TestSuite) TestFindByUserIdError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "projects"
		WHERE "projects"."user_id" = $1
		ORDER BY created_at desc LIMIT 10
		`)).
		WithArgs(suite.userId).
		WillReturnError(gorm.ErrRecordNotFound)

	projects, err := suite.repository.FindByUserId(context.Background(), suite.userId)

	require.EqualError(gorm.ErrRecordNotFound, err.Error())
	require.Empty(projects)
}

func (suite *TestSuite) TestExistsByNameAndUserIdAndTypeSuccessTrue() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(1) 
		FROM projects 
		WHERE name = $1
		AND user_id = $2
		AND type = $3
		`)).
		WithArgs(suite.name, suite.userId, suite.typeBudget).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := suite.repository.ExistsByNameAndUserIdAndType(context.Background(), suite.name, suite.userId, suite.typeBudget)

	require.NoError(err)
	require.True(exists)
}

func (suite *TestSuite) TestExistsByNameAndUserIdAndTypeSuccessFalse() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(1) 
		FROM projects 
		WHERE name = $1
		AND user_id = $2
		AND type = $3
		`)).
		WithArgs(suite.name, suite.userId, suite.typeBudget).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	exists, err := suite.repository.ExistsByNameAndUserIdAndType(context.Background(), suite.name, suite.userId, suite.typeBudget)

	require.NoError(err)
	require.False(exists)
}

func (suite *TestSuite) TestExistsByNameAndUserIdAndTypeError() {
	require := require.New(suite.T())

	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT COUNT(1) 
		FROM projects 
		WHERE name = $1
		AND user_id = $2
		AND type = $3
		`)).
		WithArgs(suite.name, suite.userId, suite.typeBudget).
		WillReturnError(gorm.ErrInvalidField)

	exists, err := suite.repository.ExistsByNameAndUserIdAndType(context.Background(), suite.name, suite.userId, suite.typeBudget)

	require.EqualError(gorm.ErrInvalidField, err.Error())
	require.False(exists)
}

func (suite *TestSuite) TestDeleteSuccess() {
	require := require.New(suite.T())
	projectExpected := &domain.Project{
		ID:        999,
		Name:      suite.name,
		UserId:    suite.userId,
		Type:      suite.typeBudget,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "projects"
		WHERE "projects"."id" = $1
		ORDER BY "projects"."id" LIMIT 1
		`)).
		WithArgs(projectExpected.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "user_id", "type"}).
			AddRow(projectExpected.ID, projectExpected.CreatedAt, projectExpected.UpdatedAt, projectExpected.Name, projectExpected.UserId, projectExpected.Type),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`
		DELETE FROM "projects"
		WHERE "projects"."id" = $1
		`)).
		WithArgs(projectExpected.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()

	err := suite.repository.Delete(context.Background(), projectExpected.ID)

	require.NoError(err)
}

func (suite *TestSuite) TestDeleteErrorFind() {
	require := require.New(suite.T())
	projectExpected := &domain.Project{
		ID:        999,
		Name:      suite.name,
		UserId:    suite.userId,
		Type:      suite.typeBudget,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "projects"
		WHERE "projects"."id" = $1
		ORDER BY "projects"."id" LIMIT 1
		`)).
		WithArgs(projectExpected.ID).
		WillReturnError(gorm.ErrInvalidField)

	err := suite.repository.Delete(context.Background(), projectExpected.ID)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func (suite *TestSuite) TestUpdateErrorSave() {
	require := require.New(suite.T())
	projectExpected := &domain.Project{
		ID:        999,
		Name:      suite.name,
		UserId:    suite.userId,
		Type:      suite.typeBudget,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mock.
		ExpectQuery(regexp.QuoteMeta(`
		SELECT *
		FROM "projects"
		WHERE "projects"."id" = $1
		ORDER BY "projects"."id" LIMIT 1
		`)).
		WithArgs(projectExpected.ID).
		WillReturnRows(sqlmock.
			NewRows([]string{"id", "created_at", "updated_at", "name", "user_id", "type"}).
			AddRow(projectExpected.ID, projectExpected.CreatedAt, projectExpected.UpdatedAt, projectExpected.Name, projectExpected.UserId, projectExpected.Type),
		)
	suite.mock.ExpectBegin()
	suite.mock.
		ExpectExec(regexp.QuoteMeta(`
		DELETE FROM "projects"
		WHERE "projects"."id" = $1
		`)).
		WithArgs(projectExpected.ID).
		WillReturnError(gorm.ErrInvalidField)
	suite.mock.ExpectRollback()

	err := suite.repository.Delete(context.Background(), projectExpected.ID)

	require.EqualError(gorm.ErrInvalidField, err.Error())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
