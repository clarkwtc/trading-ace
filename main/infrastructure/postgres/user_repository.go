package postgres

import (
    "database/sql"
    "log"
    "tradingACE/main/trading"
)

type UserRepository struct {
    client *sql.DB
}

func NewUserRepository(client *sql.DB) *UserRepository {
    return &UserRepository{client}
}

func (repository *UserRepository) FindUserTasksByAddress(address string) *trading.User {
    user := getUser(repository.client, address)
    if user == nil {
        return nil
    }

    taskRecords := repository.getTaskRecord(user, address)
    user.SetTaskRecords(taskRecords)
    return user
}

func (repository *UserRepository) getTaskRecord(user *trading.User, address string) []*trading.TaskRecord {
    query := `
    SELECT tr.id as id, t.name as name, tr.status as status, tr.amount as amount, tr.points as points FROM TaskRecord tr
    LEFT JOIN User0 u ON tr.user_id = u.id
    LEFT JOIN Task t ON tr.task_id = t.id
    WHERE u.address = $1
    Limit 100
    `
    rows, err := repository.client.Query(query, address)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var taskRecords []*trading.TaskRecord
    for rows.Next() {
        taskRecordEntity := &TaskRecordEntity{}
        taskEntity := &TaskEntity{}

        err := rows.Scan(
            &taskRecordEntity.Id, &taskEntity.Name, &taskRecordEntity.Status, &taskRecordEntity.Amount, &taskRecordEntity.Points)
        if err != nil {
            log.Fatal(err)
        }
        taskRecords = append(taskRecords, ToTaskRecord(taskRecordEntity, user, ToTask(taskEntity)))
    }

    return taskRecords
}

func getUser(client *sql.DB, address string) *trading.User {
    query := `
    SELECT id, address, amount, points FROM User0
    WHERE address = $1
    `
    rows, err := client.Query(query, address)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    userEntity := &UserEntity{}

    for rows.Next() {
        err := rows.Scan(&userEntity.Id, &userEntity.Address, &userEntity.Amount, &userEntity.Points)
        if err != nil {
            log.Fatal(err)
        }
    }

    if userEntity.Address == "" {
        return nil
    }

    return ToUser(userEntity)
}

func (repository *UserRepository) FindAllUserTasks() []*trading.User {
    query := `
    SELECT u.id as uid, u.address as address, u.amount as amount, u.points as points, t.name as name, tr.id as trid, tr.status as task_status, tr.amount as task_amount, tr.points as task_points FROM TaskRecord tr
    LEFT JOIN User0 u ON tr.user_id = u.id
    LEFT JOIN Task t ON tr.task_id = t.id
    `

    rows, err := repository.client.Query(query)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    usersMap := make(map[string]*trading.User)

    for rows.Next() {
        userEntity := &UserEntity{}
        taskRecordEntity := &TaskRecordEntity{}
        taskEntity := &TaskEntity{}

        err := rows.Scan(&userEntity.Id, &userEntity.Address, &userEntity.Amount, &userEntity.Points, &taskEntity.Name, &taskRecordEntity.Id, &taskRecordEntity.Status, &taskRecordEntity.Amount, &taskRecordEntity.Points)
        if err != nil {
            log.Fatal(err)
        }

        _, exists := usersMap[userEntity.Address]
        if !exists {
            usersMap[userEntity.Address] = ToUser(userEntity)
        }
        usersMap[userEntity.Address].AcceptTask(ToTaskRecord(taskRecordEntity, usersMap[userEntity.Address], ToTask(taskEntity)))
    }

    var users []*trading.User
    for _, value := range usersMap {
        users = append(users, value)
    }

    return users
}

func (repository *UserRepository) SaveAllUser(users []*trading.User) {
    tx, err := repository.client.Begin()
    if err != nil {
        log.Fatal(err)
    }

    query := `
    SELECT id, name FROM Task
    `

    rows, err := repository.client.Query(query)
    if err != nil {
        log.Fatal(err)
    }

    tasksMap := make(map[string]*TaskEntity)

    for rows.Next() {
        task := &TaskEntity{}

        err := rows.Scan(&task.Id, &task.Name)
        if err != nil {
            log.Fatal(err)
        }
        tasksMap[task.Name] = task
    }

    mergeUser(tx, users)
    mergeTaskRecord(tx, users, tasksMap)

    err = tx.Commit()
    if err != nil {
        log.Fatal(err)
    }
}

func mergeUser(tx *sql.Tx, users []*trading.User) {
    query := `
    INSERT INTO user0 (id, address, amount, points, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (address)
    DO UPDATE SET
        amount = EXCLUDED.amount,
        points = EXCLUDED.points,
        updated_at = EXCLUDED.updated_at
    `

    stmt, err := tx.Prepare(query)
    if err != nil {
        newErr := tx.Rollback()
        if newErr != nil {
            log.Fatal(err)
        }
    }
    defer stmt.Close()

    for _, user := range users {
        entity := ToUserEntity(user)
        _, err = stmt.Exec(entity.Id, entity.Address, entity.Amount, entity.Points, entity.CreatedAt, entity.UpdatedAt)
        if err != nil {
            newErr := tx.Rollback()
            if newErr != nil {
                log.Fatal(err)
            }
        }
    }
}

func mergeTaskRecord(tx *sql.Tx, users []*trading.User, tasksMap map[string]*TaskEntity) {
    query := `
    INSERT INTO taskrecord (id, user_id, task_id, status, amount, points, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    ON CONFLICT (id)
    DO UPDATE SET
        status = EXCLUDED.status,
        amount = EXCLUDED.amount,
        points = EXCLUDED.points,
        updated_at = EXCLUDED.updated_at
    `
    stmt, err := tx.Prepare(query)
    if err != nil {
        newErr := tx.Rollback()
        if newErr != nil {
            log.Fatal(err)
        }
    }
    defer stmt.Close()

    for _, user := range users {
        for _, taskRecord := range user.GetTaskRecords() {
            entity := ToTaskRecordEntity(taskRecord, user, tasksMap)
            _, err = stmt.Exec(entity.Id, entity.UserId, entity.TaskId, entity.Status, entity.Amount, entity.Points, entity.CreatedAt, entity.UpdatedAt)
            if err != nil {
                newErr := tx.Rollback()
                if newErr != nil {
                    log.Fatal(err)
                }
            }
        }
    }
}
