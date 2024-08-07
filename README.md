# Creating a Profile Builder Application with Admin Control and Some Advanced Features.

<p>
<b>Requirements - </b> The system revolves around a single role: Admin.

<b> Admin Privileges: </b>
Create, view, update user profiles.
Control all profile-related operations.

General Features: Admin can create standardized templates for user profiles. Only admins who manage all user profile activities and print the pdf of created profiles.

</p>

## Setup

This Project uses Postgres DB to handle database queries.
There are few records already seeded into database and whatever updations you make on database, it will persist even after you close the application. You can run the CleanUp command to start fresh.

Firstly, run the following command to get the Project on local system

```bash
git clone github.com/joshsoftware/profile_builder_backend_go.git
```

1. Run following command to download all dependencies

```bash
go mod download
or
go mod tidy
```

2. Run following command to start Application

```bash
make run
```

2. Run following command to run unit test cases

```bash
make test
```

3. Run following command to check test coverage

```bash
make test-cover

#you can also check code test coverage on top. Click on codeccov badge to check more about test coverage
```

4. Run following command to erase database to start fresh

```bash
make cleanDB
```

---

At Next, run the following command in terminal for DB related activities

```bash
sudo -u postgres psql
```

1. Run following command to Create a Database for Application

```bash
CREATE DATABASE profile_builder;
\c profile_builder;
```

## Commands to Make migrations of DB

2.Run following command to install migrate dependency -

```bash
curl -L -o migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz
```

3.Run following command to extract that dependency -

```bash
tar -xvzf migrate.tar.gz

for MacOs
brew install golang-migrate
```
4.Run following command to extract that dependency(Optional) -

```bash
sudo mv migrate /usr/local/bin/
```
5.Finally, Run following command to apply all migrations -

```bash
make migrate-up
```

## Postman Collection

[here](postman_collection.json)

## Project Structure

```
josh@josh:~/JOSH/Profile Builder Backend$ tree
.
├── cmd
│   └── main.go
├── coverage.out
├── Dockerfile
├── Dockerfile.dev
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── handler
│   │   │   ├── achievement_handler.go
│   │   │   ├── certificate_handler.go
│   │   │   ├── decoder.go
│   │   │   ├── education_handler.go
│   │   │   ├── experience_handler.go
│   │   │   ├── profile_handler.go
│   │   │   ├── project_handler.go
│   │   │   └── user_login_handler.go
│   │   ├── router.go
│   │   └── tests
│   │       ├── achievement_handler_test.go
│   │       ├── certificate_handler_test.go
│   │       ├── education_handler_test.go
│   │       ├── experience_handler_test.go
│   │       ├── profile_handler_test.go
│   │       ├── project_handler_test.go
│   │       └── user_login_handler_test.go
│   ├── app
│   │   └── service
│   │       ├── achievement_service.go
│   │       ├── certificate_service.go
│   │       ├── education_service.go
│   │       ├── experience_service.go
│   │       ├── mocks
│   │       │   ├── AchievementService.go
│   │       │   ├── CertificateService.go
│   │       │   ├── EducationService.go
│   │       │   ├── ExperienceService.go
│   │       │   ├── ProjectService.go
│   │       │   ├── Service.go
│   │       │   └── UserLoginServive.go
│   │       ├── project_service.go
│   │       ├── service.go
│   │       ├── tests
│   │       │   ├── achievement_service_test.go
│   │       │   ├── certificate_service_test.go
│   │       │   ├── education_service_test.go
│   │       │   ├── experience_service_test.go
│   │       │   ├── project_service_test.go
│   │       │   ├── service_test.go
│   │       │   └── user_login_service_test.go
│   │       └── user_login_service.go
│   ├── cron-job
│   │   └── app.go
│   ├── db
│   │   ├── migrate.go
│   │   └── migrations
│   │       ├── 000001_initial.down.sql
│   │       ├── 000001_initial.up.sql
│   │       ├── 000002_add_colm_josh_joining_date.down.sql
│   │       └── 000002_add_colm_josh_joining_date.up.sql
│   ├── pkg
│   │   ├── constants
│   │   │   └── app.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── helpers
│   │   │   ├── app.go
│   │   │   └── helpers.go
│   │   ├── jwt_token
│   │   │   └── createToken.go
│   │   ├── middleware
│   │   │   ├── auth.go
│   │   │   ├── response_writer.go
│   │   │   ├── tests
│   │   │   │   ├── auth_test.go
│   │   │   │   └── verify_jwt_token_test.go
│   │   │   └── verify_jwt_token.go
│   │   └── specs
│   │       ├── achievement.go
│   │       ├── api.go
│   │       ├── certificate.go
│   │       ├── education.go
│   │       ├── experience.go
│   │       ├── profile.go
│   │       ├── project.go
│   │       └── user_login.go
│   └── repository
│       ├── achievement_repository.go
│       ├── certificate_repository.go
│       ├── config.go
│       ├── education_repository.go
│       ├── experience_repository.go
│       ├── init.go
│       ├── mocks
│       │   ├── AchievementStorer.go
│       │   ├── CertificateStorer.go
│       │   ├── EducationStorer.go
│       │   ├── ExperienceStorer.go
│       │   ├── ProfileStorer.go
│       │   ├── ProjectStorer.go
│       │   ├── RepositoryTrasanctions.go
│       │   ├── Trasanctions.go
│       │   └── UserStorer.go
│       ├── model.go
│       ├── profile_repository.go
│       ├── project_repository.go
│       ├── repo.go
│       └── user_login_repository.go
├── Makefile
├── README.md
└── swagger.yaml

22 directories, 90 files

```
