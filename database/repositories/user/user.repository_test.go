package user

import (
	"skripsi-sso/database/entities"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository *userRepo
	user       entities.User
	role       entities.Role
	permission entities.Permission
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	mDB := []*gorm.DB{s.DB, s.DB, s.DB}
	s.repository = NewUserRepo(mDB)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_FindUserAuthForDevice() {
	email := "test@ainosi.com"
	userId := uint(1)
	roleIds := []int64{1, 2}
	roleNames := []string{"ADMIN", "VIEWER"}
	permissions := []string{"master.user", "master.role"}
	logos := []string{"/logo.png", "/logo.png"}
	address := []string{"Jl. Yogyakarta", "Jl. Yogyakarta"}
	groupNames := []string{"xyz", "zyz"}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "`+s.user.TableName()+`"
			WHERE (email = $1 and type = $2)`)).
		WithArgs(email, "device").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(userId, email))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT roles.name FROM "` + s.role.TableName() + `"
			LEFT JOIN model_has_roles ON model_has_roles.role_id = roles.id
			WHERE (model_has_roles.model_id IN ($1))`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(roleNames[0]).
			AddRow(roleNames[1]))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT roles.id FROM "` + s.role.TableName() + `"
			LEFT JOIN model_has_roles ON model_has_roles.role_id = roles.id
			WHERE (model_has_roles.model_id IN ($1))`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(roleIds[0]).
			AddRow(roleIds[1]))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT permissions.name FROM "`+s.permission.TableName()+`"
			INNER JOIN role_has_permissions ON role_has_permissions.permission_id = permissions.id
			WHERE (role_has_permissions.role_id in ($1,$2)`)).
		WithArgs(roleIds[0], roleIds[1]).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).
			AddRow(permissions[0]).
			AddRow(permissions[1]))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT master_group.group_logo FROM "role_has_group"
			INNER JOIN master_group ON master_group.group_id = role_has_group.group_id
			WHERE (role_has_group.role_id in ($1,$2))`)).
		WithArgs(roleIds[0], roleIds[1]).
		WillReturnRows(sqlmock.NewRows([]string{"group_logo"}).
			AddRow(logos[0]).
			AddRow(logos[1]))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT master_group.group_extras->>'address' as address FROM "role_has_group"
			INNER JOIN master_group ON master_group.group_id = role_has_group.group_id
			WHERE (role_has_group.role_id in ($1,$2))`)).
		WithArgs(roleIds[0], roleIds[1]).
		WillReturnRows(sqlmock.NewRows([]string{"address"}).
			AddRow(address[0]).
			AddRow(address[1]))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT master_group.group_name FROM "role_has_group"
			INNER JOIN master_group ON master_group.group_id = role_has_group.group_id
			WHERE (role_has_group.role_id in ($1,$2))`)).
		WithArgs(roleIds[0], roleIds[1]).
		WillReturnRows(sqlmock.NewRows([]string{"group_name"}).
			AddRow(groupNames[0]).
			AddRow(groupNames[1]))

	res, err := s.repository.FindUserAuthForDevice(email)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
	require.Equal(s.T(), userId, res.UserId)
	require.Equal(s.T(), roleNames, res.Roles)
	require.Equal(s.T(), permissions, res.Permissions)
	require.Equal(s.T(), logos[0], res.Logo)
	require.Equal(s.T(), address[0], res.Address)
	require.Equal(s.T(), groupNames[0], res.GroupName)
}

func (s *Suite) Test_FindUserByEmailForCss() {
	email := "test@ainosi.com"
	userId := uint(1)
	procodes := []string{"D001", "D002"}
	expiredAt, _ := time.Parse(time.RFC3339, "2019-06-07 00:12:12")

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, email, password, expired_at FROM "`+s.user.TableName()+`"
			WHERE (email = $1 and type = $2)`)).
		WithArgs(email, "css").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "expired_at"}).
			AddRow(userId, email, "asdf", expiredAt))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT master_procode.code FROM "user_has_procode"
			LEFT JOIN master_procode ON master_procode.id = user_has_procode.procode_id 
			WHERE (user_has_procode.user_id in ($1))`)).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"code"}).
			AddRow(procodes[0]).
			AddRow(procodes[1]))

	res, err := s.repository.FindUserByEmailForCss(email)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
	require.Equal(s.T(), userId, res.UserId)
	require.Equal(s.T(), email, res.Email)
	require.Equal(s.T(), expiredAt, res.ExpiredAt)
	require.Equal(s.T(), procodes, res.ProcodePermissions)
}
