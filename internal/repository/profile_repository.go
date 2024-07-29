package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	errors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// ProfileStore implements the ProfileStorer interface.
type ProfileStore struct {
	db *pgxpool.Pool
}

// ProfileStorer defines methods to interact with user profile data.
type ProfileStorer interface {
	CreateProfile(ctx context.Context, pd ProfileRepo, tx pgx.Tx) (int, error)
	ListProfiles(ctx context.Context, tx pgx.Tx) (values []specs.ListProfiles, err error)
	GetProfile(ctx context.Context, profileID int, tx pgx.Tx) (value specs.ResponseProfile, err error)
	UpdateProfile(ctx context.Context, profileID int, pd UpdateProfileRepo, tx pgx.Tx) (int, error)
	UpdateSequence(ctx context.Context, us UpdateSequenceRequest, tx pgx.Tx) (ID int, err error)
	DeleteProfile(ctx context.Context, profileID int, tx pgx.Tx) (err error)
	CountRecords(ctx context.Context, ProfileID int, ComponentName string, tx pgx.Tx) (Count int, err error)
	UpdateProfileStatus(ctx context.Context, profileID int, updateRequest UpdateProfileStatusRepo, tx pgx.Tx) error
	ListSkills(ctx context.Context, tx pgx.Tx) (values specs.ListSkills, err error)
	BeginTransaction(ctx context.Context) (tx pgx.Tx, err error)
	HandleTransaction(ctx context.Context, tx pgx.Tx, incomingErr error) (err error)
	BackupAllProfiles(backupDir string)
}

// NewProfileRepo creates a new instance of ProfileRepo.
func NewProfileRepo(db *pgxpool.Pool) ProfileStorer {
	return &ProfileStore{
		db: db,
	}
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// CreateProfile inserts a new user profile into the database.
func (profileStore *ProfileStore) CreateProfile(ctx context.Context, pd ProfileRepo, tx pgx.Tx) (int, error) {

	values := []interface{}{
		pd.Name, pd.Email, pd.Gender, pd.Mobile, pd.Designation, pd.Description, pd.Title,
		pd.YearsOfExperience, pd.PrimarySkills, pd.SecondarySkills, pd.JoshJoiningDate, pd.GithubLink, pd.LinkedinLink, pd.CareerObjectives, 1, 1, pd.CreatedAt, pd.UpdatedAt, pd.CreatedByID, pd.UpdatedByID,
	}

	insertQuery, args, err := psql.Insert("profiles").
		Columns(constants.CreateUserColumns...).
		Values(values...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		zap.S().Error("Error generating profile insert query: ", err)
		return 0, err
	}

	if tx == nil {
		zap.S().Error("Transaction is nil")
		return 0, fmt.Errorf("internal error: transaction is nil")
	}

	var profileID int
	err = tx.QueryRow(ctx, insertQuery, args...).Scan(&profileID)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return 0, errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing insert query: ", err)
		return 0, err
	}

	return profileID, nil
}

// ListProfiles returns a list of all profiles in the Database that are currently available
func (profileStore *ProfileStore) ListProfiles(ctx context.Context, tx pgx.Tx) (values []specs.ListProfiles, err error) {
	sql, args, err := psql.Select(constants.ListProfilesColumns...).From("profiles p").LeftJoin("invitations i ON p.id = i.profile_id").OrderBy("p.created_at DESC").ToSql()
	if err != nil {
		zap.S().Error("Error generating list project select query: ", err)
		return []specs.ListProfiles{}, err
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("error executing list project insert query:", err)
		return []specs.ListProfiles{}, err
	}

	for rows.Next() {
		var value specs.ListProfiles
		err := rows.Scan(&value.ID, &value.Name, &value.Email, &value.YearsOfExperience, &value.PrimarySkills, &value.IsCurrentEmployee, &value.IsActive, &value.ProfileComplete)
		if err != nil {
			zap.S().Error("error scanning row:", err)
			return []specs.ListProfiles{}, err
		}
		values = append(values, value)
	}
	defer rows.Close()
	return values, nil
}

// ListSkills returns a list of all skills in the Database that are currently available
func (profileStore *ProfileStore) ListSkills(ctx context.Context, tx pgx.Tx) (values specs.ListSkills, err error) {
	sql, args, err := psql.Select("name").From("skills").ToSql()
	if err != nil {
		zap.S().Error("Error generating list skills select query: ", err)
		return specs.ListSkills{}, err
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("error executing list skills insert query:", err)
		return specs.ListSkills{}, err
	}

	for rows.Next() {
		var value string
		rows.Scan(&value)
		values.Name = append(values.Name, value)
	}
	defer rows.Close()

	return values, nil
}

// GetProfile returns a details profile in the Database that are currently available for perticular ID
func (profileStore *ProfileStore) GetProfile(ctx context.Context, profileID int, tx pgx.Tx) (value specs.ResponseProfile, err error) {
	query, args, err := psql.Select(constants.ResponseProfileColumns...).From("profiles").
		Where(sq.Eq{"id": profileID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating list project select query: ", err)
		return specs.ResponseProfile{}, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		zap.S().Error("Error executing get profile query: ", err)
		return specs.ResponseProfile{}, err
	}

	if rows.Next() {
		if err := rows.Scan(&value.ProfileID, &value.Name, &value.Email, &value.Gender, &value.Mobile, &value.Designation, &value.Description, &value.Title, &value.YearsOfExperience, &value.PrimarySkills, &value.SecondarySkills, &value.JoshJoiningDate, &value.GithubLink, &value.LinkedinLink, &value.CareerObjectives); err != nil {
			zap.S().Error("Error scanning row: ", err)
			return specs.ResponseProfile{}, err
		}
	} else {
		zap.S().Info("No profile found for profileID: ", profileID)
		return specs.ResponseProfile{}, errors.ErrNoRecordFound
	}
	defer rows.Close()

	return value, nil
}

// UpdateProfile updates an existing user profile in the database.
func (profileStore *ProfileStore) UpdateProfile(ctx context.Context, profileID int, pd UpdateProfileRepo, tx pgx.Tx) (int, error) {

	updateQuery, args, err := psql.Update("profiles").
		SetMap(map[string]interface{}{
			"name": pd.Name, "email": pd.Email,
			"gender": pd.Gender, "mobile": pd.Mobile,
			"designation": pd.Designation, "description": pd.Description,
			"title": pd.Title, "years_of_experience": pd.YearsOfExperience,
			"josh_joining_date": pd.JoshJoiningDate, "primary_skills": pd.PrimarySkills, "secondary_skills": pd.SecondarySkills, "github_link": pd.GithubLink, "linkedin_link": pd.LinkedinLink, "updated_at": pd.UpdatedAt, "updated_by_id": pd.UpdatedByID,
		}).
		Where(sq.Eq{"id": profileID}).ToSql()
	if err != nil {
		zap.S().Error("Error generating profile update query: ", err)
		return 0, err
	}

	res, err := tx.Exec(ctx, updateQuery, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return 0, errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return 0, errors.ErrInvalidProfile
		}
		zap.S().Error("Error executing update query: ", err)
		return 0, err
	}

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : profile")
		return 0, errors.ErrInvalidRequestData
	}

	return profileID, nil
}

// DeleteProfile delete an existing user profile in the database.
func (profileStore *ProfileStore) DeleteProfile(ctx context.Context, profileID int, tx pgx.Tx) (err error) {
	deleteQuery, args, err := psql.Delete("profiles").Where(sq.Eq{"id": profileID}).ToSql()
	if err != nil {
		zap.S().With("profile_id", profileID).Error("Error generating profile delete query: ", err)
		return err
	}

	result, err := tx.Exec(ctx, deleteQuery, args...)
	if err != nil {
		zap.S().With("query", deleteQuery, "args", args).Error("Error executing delete profile query", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrNoData
	}
	return nil
}

// CountRecords counts an existing user records for particular component in the database.
func (profileStore *ProfileStore) CountRecords(ctx context.Context, ProfileID int, ComponentName string, tx pgx.Tx) (Count int, err error) {
	countQuery, args, err := psql.Select("count(*)").
		From(ComponentName).
		Where(sq.Eq{"profile_id": ProfileID}).
		ToSql()
	if err != nil {
		zap.S().Error("Error generating count query: ", err)
		return 0, err
	}

	var currentCount int
	err = tx.QueryRow(ctx, countQuery, args...).Scan(&currentCount)
	if err != nil {
		zap.S().Error("Error fetching current count: ", err)
		return 0, err
	}

	return currentCount, nil
}

// UpdateSequence updates an existing component's priorities in the database.
func (profileStore *ProfileStore) UpdateSequence(ctx context.Context, us UpdateSequenceRequest, tx pgx.Tx) (int, error) {

	for compID, priority := range us.ComponentPriorities {
		updateQuery, args, err := psql.Update(us.ComponentName).
			Set("priorities", priority).
			Set("updated_at", us.UpdatedAt).
			Set("updated_by_id", us.UpdatedByID).
			Where(sq.Eq{"id": compID, "profile_id": us.ProfileID}).
			ToSql()
		if err != nil {
			zap.S().Error("Error constructing update query: ", err)
			return 0, err
		}

		res, err := tx.Exec(ctx, updateQuery, args...)
		if err != nil {
			if helpers.IsDuplicateKeyError(err) {
				return 0, errors.ErrDuplicateKey
			}
			if helpers.IsInvalidProfileError(err) {
				return 0, errors.ErrInvalidProfile
			}
			zap.S().Error("Error executing update query: ", err)
			return 0, err
		}

		if res.RowsAffected() == 0 {
			zap.S().Infof("No rows affected for update component sequence : %s", us.ComponentName)
		}
	}

	return us.ProfileID, nil
}

// UpdateProfileStatus updates an existing profile's status in the database.
func (profileStore *ProfileStore) UpdateProfileStatus(ctx context.Context, profileID int, updateRequest UpdateProfileStatusRepo, tx pgx.Tx) error {
	updateQuery := psql.Update("profiles")

	if updateRequest.IsCurrentEmployee != nil {
		updateQuery = updateQuery.Set("is_current_employee", *updateRequest.IsCurrentEmployee)
	}

	if updateRequest.IsActive != nil {
		updateQuery = updateQuery.Set("is_active", *updateRequest.IsActive)
	}

	updateQuery = updateQuery.Where(sq.Eq{"id": profileID})
	query, args, err := updateQuery.ToSql()
	if err != nil {
		zap.S().Error("Error generating update profile status query: ", err)
		return err
	}

	res, err := tx.Exec(ctx, query, args...)
	if err != nil {
		if helpers.IsDuplicateKeyError(err) {
			return errors.ErrDuplicateKey
		}
		if helpers.IsInvalidProfileError(err) {
			return errors.ErrInvalidProfile
		}
		zap.S().Error("Error while executing update query: ", err)
		return err
	}

	if res.RowsAffected() == 0 {
		zap.S().Warn("invalid request for update : profile status")
		return errors.ErrInvalidRequestData
	}
	return nil
}

// BeginTransaction used to begin transaction while each task
func (profileStore *ProfileStore) BeginTransaction(ctx context.Context) (tx pgx.Tx, err error) {
	tx, err = profileStore.db.BeginTx(ctx, pgx.TxOptions{})
	return
}

// HandleTransaction used to handle transaction while each task
func (profileStore *ProfileStore) HandleTransaction(ctx context.Context, tx pgx.Tx, incomingErr error) (err error) {
	if incomingErr != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return incomingErr
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// BackupAllProfiles used to backing up all the data at every midnight by CronJOb
func (profileStore *ProfileStore) BackupAllProfiles(backupDir string) {
	zap.S().Info("Starting backing up data...")
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		zap.S().Errorw("Failed to create backup directory", "error", err)
		return
	}

	fileName := fmt.Sprintf("%s%s-profile_builder.sql", backupDir, time.Now().Format("20060102"))

	file, err := os.Create(fileName)
	if err != nil {
		zap.S().Errorw("Failed to create backup file", "error", err)
		return
	}
	defer file.Close()

	for _, table := range constants.BackupTables {
		if err := dumpTable(profileStore, file, table); err != nil {
			zap.S().Errorw("Failed to dump table", "table", table, "error", err)
			return
		}
	}

	zap.S().Infow("Database backed up successfully", "fileName", fileName, "time", time.Now())
}

func dumpTable(profileStore *ProfileStore, file *os.File, table string) error {
	rows, err := profileStore.db.Query(context.Background(), fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return err
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()
	columnNames := make([]string, len(columns))
	for i, col := range columns {
		columnNames[i] = string(col.Name)
	}

	if _, err := fmt.Fprintf(file, "-- Dumping data for table %s\n", table); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(file, "COPY %s (%s) FROM stdin;\n", table, helpers.JoinValues(columnNames, ", ")); err != nil {
		return err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(file, "%s\n", helpers.JoinValues(values, "\t")); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(file, "\\."); err != nil {
		return err
	}

	return rows.Err()
}
