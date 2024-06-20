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

2. Run following command to Make migrations of DB

```bash
make migrate
```

## APIs

1. <b>Login API</b> : `POST http://localhost:1925/login`
2. <b>Create Profile</b> : `POST http://localhost:1925/profiles`
3. <b>Create Educations for Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/educations`
4. <b>Create Projects for Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/projects`
5. <b>Create Experiences for Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/experiences`
6. <b>Create Certificates for Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/certificates`
7. <b>Create Achievements for Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/achievements`
8. <b>List Skills</b> : `POST http://localhost:1925/skills`
9. <b>Get Profile of Specific ID</b> : `POST http://localhost:1925/profiles/{profile_id}`
10. <b>Get Educations of Specific Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/educations`
11. <b>Get Projects of Specific Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/projects`
12. <b>Get Experiences of Specific Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/experiences`
13. <b>Get Certificates of Specific Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/certificates`
14. <b>Get Achievements of Specific Profile</b> : `POST http://localhost:1925/profiles/{profile_id}/achievements`

## Postman Collection

[here](postman_collection.json)

## Project Structure

```
jspnlp@unispab:~/JOSH/Profile Builder Backend$ tree
.
├── cmd
│   └── main.go
├── coverage.out
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
│   │   │   └── project_handler.go
│   │   └── router.go
│   ├── app
│   │   ├── dependencies.go
│   │   ├── mocks
│   │   │   └── Service.go
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
│   │       │   └── Service.go
│   │       ├── project_service.go
│   │       └── service.go
│   ├── db
│   │   ├── migrate.go
│   │   └── migrations
│   │       ├── postgres_DOWN.sql
│   │       └── postrges_UP.sql
│   ├── pkg
│   │   ├── constants
│   │   │   └── app.go
│   │   ├── specs
│   │   │   ├── achievement.go
│   │   │   ├── api.go
│   │   │   ├── certificate.go
│   │   │   ├── education.go
│   │   │   ├── experience.go
│   │   │   ├── profile.go
│   │   │   └── project.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── helpers
│   │   │   └── app.go
│   │   └── middleware
│   │       └── response_writer.go
│   └── repository
│       ├── achievement_repository.go
│       ├── certificate_repository.go
│       ├── education_repository.go
│       ├── experience_repository.go
│       ├── init.go
│       ├── mocks
│       │   ├── AchievementStorer.go
│       │   ├── CertificateStorer.go
│       │   ├── EducationStorer.go
│       │   ├── ExperienceStorer.go
│       │   ├── ProfileStorer.go
│       │   └── ProjectStorer.go
│       ├── model.go
│       ├── profile_repository.go
│       └── project_repository.go
└── Makefile

18 directories, 55 files

```
