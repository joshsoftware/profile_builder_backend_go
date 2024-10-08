package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
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

// Constants for table names
var (
	ProfileTable = "profiles"
	SkillsTable  = "skills"
)

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
	GetProfileIDByEmail(ctx context.Context, email string, tx pgx.Tx) (int, error)
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

	insertQuery, args, err := psql.Insert(ProfileTable).
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
func (profileStore *ProfileStore) ListProfiles(ctx context.Context, tx pgx.Tx) ([]specs.ListProfiles, error) {
	queryBuilder := psql.Select(constants.ListProfilesColumns...).From("profiles p").OrderBy("p.created_at DESC")
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		zap.S().Error("Error executing query: ", err)
		return nil, err
	}
	defer rows.Close()

	var profiles []specs.ListProfiles
	for rows.Next() {
		var profile specs.ListProfiles
		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Email,
			&profile.YearsOfExperience,
			&profile.PrimarySkills,
			&profile.IsCurrentEmployee,
			&profile.IsActive,
			&profile.JoshJoiningDate,
			&profile.CreatedAt,
			&profile.UpdatedAt,
			&profile.IsProfileComplete,
		)
		if err != nil {
			zap.S().Error("Error scanning row: ", err)
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// ListSkills returns a list of all skills in the Database that are currently available
func (profileStore *ProfileStore) ListSkills(ctx context.Context, tx pgx.Tx) (values specs.ListSkills, err error) {
	sql, args, err := psql.Select("name").From(SkillsTable).ToSql()
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
	query, args, err := psql.Select(constants.ResponseProfileColumns...).
		Column(sq.Expr(`(CASE 
				WHEN EXISTS (
					SELECT 1 FROM invitations 
					WHERE profile_id = ? AND is_profile_complete = 0
				) THEN 'YES' 
				ELSE 'NO' 
			END) AS is_invited`, profileID)).From(ProfileTable).Where(sq.Eq{"id": profileID}).ToSql()

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
		if err := rows.Scan(&value.ProfileID, &value.Name, &value.Email, &value.Gender, &value.Mobile, &value.Designation, &value.Description, &value.Title, &value.YearsOfExperience, &value.PrimarySkills, &value.SecondarySkills, &value.GithubLink, &value.LinkedinLink, &value.CareerObjectives, &value.JoshJoiningDate, &value.IsInvited); err != nil {
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
	updateQuery, args, err := psql.Update(ProfileTable).
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

	selectEmailQuery, args, err := psql.Select("email").From(ProfileTable).Where(sq.Eq{"id": profileID}).ToSql()
	if err != nil {
		zap.S().With("profile_id", profileID).Error("Error generating email select query: ", err)
	}

	var email string
	err = tx.QueryRow(ctx, selectEmailQuery, args...).Scan(&email)
	if err != nil {
		zap.S().Error("Error fetching email: ", err)
	}

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

	if email == "" {
		zap.S().Info("No email found for profile ID: ", profileID)
		return nil
	}

	selectRoleQuery, args, err := psql.Select("role").From(userTable).Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		zap.S().With("email", email).Error("Error generating role select query: ", err)
		return err
	}

	var role string
	err = tx.QueryRow(ctx, selectRoleQuery, args...).Scan(&role)
	if err == pgx.ErrNoRows {
		zap.S().Info("No user found for email: ", email, " skipping user deletion.")
	} else if err != nil {
		zap.S().Error("Error fetching role: ", err)
		return err
	}

	if role == constants.Employee {
		deleteQuery, args, err := psql.Delete("users").Where(sq.Eq{"email": email}).ToSql()
		if err != nil {
			zap.S().With("email", email).Error("Error generating user delete query: ", err)
			return err
		}

		result, err := tx.Exec(ctx, deleteQuery, args...)
		if err != nil {
			zap.S().With("query", deleteQuery, "args", args).Error("Error executing delete user query", zap.Error(err))
			return err
		}

		if result.RowsAffected() == 0 {
			zap.S().Info("No user found for email : ", email)
		}
	} else {
		zap.S().Info("User with email: ", email, " is an admin,, skipping delete.")
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
	updateQuery := psql.Update(ProfileTable)

	if updateRequest.IsCurrentEmployee != nil {
		updateQuery = updateQuery.Set("is_current_employee", *updateRequest.IsCurrentEmployee)
	}

	if updateRequest.IsActive != nil {
		updateQuery = updateQuery.Set("is_active", *updateRequest.IsActive)
	}

	updateQuery = updateQuery.Set("updated_at", updateRequest.UpdatedAt).Where(sq.Eq{"id": profileID})
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

	dbName := os.Getenv("BACKUP_DBNAME")
	dbUser := os.Getenv("BACKUP_USER")
	dbPass := os.Getenv("BACKUP_PASSWORD")
	backupFile := fmt.Sprintf("%s/%s_backup_%s.sql", backupDir, dbName, time.Now().Format("20060102150405"))

	os.Setenv("PGPASSWORD", dbPass)
	defer os.Unsetenv("PGPASSWORD")

	cmd := exec.Command("pg_dump", "-U", dbUser, "-f", backupFile, dbName)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to run pg_dump: %v, output: %s", err, output)
	}

	zap.S().Infow("Database backed up successfully", "fileName", backupFile)
}

// GetProfileIDByEmail returns the profile ID for a given email.
func (profileStore *ProfileStore) GetProfileIDByEmail(ctx context.Context, email string, tx pgx.Tx) (int, error) {
	query := psql.Select("id").From("profiles").Where(sq.Eq{"email": email})
	sql, args, err := query.ToSql()
	if err != nil {
		zap.S().Error("Error generating select query: ", err)
		return 0, err
	}

	var profileID int
	err = tx.QueryRow(ctx, sql, args...).Scan(&profileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, errors.ErrNoRecordFound
		}
		zap.S().Error("Error executing select query: ", err)
		return 0, err
	}

	return profileID, nil
}
