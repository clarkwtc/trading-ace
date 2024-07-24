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

    taskRecords := getTaskRecord(repository.client, address)
    user.Tasks = taskRecords
    return user
}

func getTaskRecord(client *sql.DB, address string) []*trading.TaskRecord {
    query := `
    SELECT tr.id as id, t.name as name, tr.status as status, tr.amount as amount, tr.points as points FROM TaskRecord tr
    LEFT JOIN User0 u ON tr.user_id = u.id
    LEFT JOIN Task t ON tr.task_id = t.id
    WHERE u.address = $1
    Limit 100
    `
    rows, err := client.Query(query, address)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var taskRecords []*trading.TaskRecord
    for rows.Next() {
        taskRecord := &TaskRecordEntity{}
        task := &TaskEntity{}

        err := rows.Scan(
            &taskRecord.Id, &task.Name, &taskRecord.Status, &taskRecord.Amount, &taskRecord.Points)
        if err != nil {
            log.Fatal(err)
        }
        taskRecords = append(taskRecords, ToTaskRecord(taskRecord, task))
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

func (repository *UserRepository) FindUserRewardByAddress(address string) *trading.User {
    user := getUser(repository.client, address)
    if user == nil {
        return nil
    }

    rewardRecord := getRewardRecord(repository.client, address)
    user.PointHistory = rewardRecord
    return user

    return user
}

func getRewardRecord(client *sql.DB, address string) []*trading.RewardRecord {
    query := `
    SELECT rr.id as id, t.name as name, rr.points as points, rr.created_at as CreatedAt FROM RewardRecord rr
    LEFT JOIN User0 u ON rr.user_id = u.id
    LEFT JOIN Task t ON rr.task_id = t.id
    WHERE u.address = $1
    Limit 100
    `
    rows, err := client.Query(query, address)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var rewardRecords []*trading.RewardRecord

    for rows.Next() {
        rewardRecord := &RewardRecordEntity{}
        task := &TaskEntity{}

        err := rows.Scan(
            &rewardRecord.Id, &task.Name, &rewardRecord.Points, &rewardRecord.CreatedAt)
        if err != nil {
            log.Fatal(err)
        }
        rewardRecords = append(rewardRecords, ToRewardRecord(rewardRecord, task))
    }
    return rewardRecords
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
        usersMap[userEntity.Address].AddTask(ToTaskRecord(taskRecordEntity, taskEntity))
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
    insertRewardRecord(tx, users, tasksMap)

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
        for _, task := range user.Tasks {
            entity := ToTaskRecordEntity(user, task, tasksMap)
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

func insertRewardRecord(tx *sql.Tx, users []*trading.User, tasksMap map[string]*TaskEntity) {
    query := `
    INSERT INTO rewardrecord (id, user_id, task_id, points, created_at)
    VALUES ($1, $2, $3, $4, $5)
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
        for _, reward := range user.PointHistory {
            entity := ToRewardRecordEntity(user, reward, tasksMap)
            _, err = stmt.Exec(entity.Id, entity.UserId, entity.TaskId, entity.Points, entity.CreatedAt)
            if err != nil {
                newErr := tx.Rollback()
                if newErr != nil {
                    log.Fatal(err)
                }
            }
        }
    }
}
