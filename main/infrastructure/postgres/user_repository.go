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

func (repository *UserRepository) FindUserTasksByAddress(address string) (*trading.User, error) {
    user, err := getUser(repository.client, address)
    if err != nil {
        return user, err
    }

    taskRecords, err := repository.getTaskRecord(user, address)
    if err != nil {
        return user, err
    }
    user.SetTaskRecords(taskRecords)
    return user, nil
}

func (repository *UserRepository) getTaskRecord(user *trading.User, address string) ([]*trading.TaskRecord, error) {
    query := `
    SELECT tr.id as id, t.name as name, tr.status as status, tr.amount as amount, tr.points as points FROM TaskRecord tr
    LEFT JOIN User0 u ON tr.user_id = u.id
    LEFT JOIN Task t ON tr.task_id = t.id
    WHERE u.address = $1
    Limit 100
    `
    rows, err := repository.client.Query(query, address)
    if err != nil {
        log.Printf("Query getTaskRecord fail: %v", err)
        return nil, err
    }
    defer rows.Close()

    var taskRecords []*trading.TaskRecord
    for rows.Next() {
        taskRecordEntity := &TaskRecordEntity{}
        taskEntity := &TaskEntity{}

        err = rows.Scan(
            &taskRecordEntity.Id, &taskEntity.Name, &taskRecordEntity.Status, &taskRecordEntity.Amount, &taskRecordEntity.Points)
        if err != nil {
            log.Printf("Parse getTaskRecord entity fail: %v", err)
            return nil, err
        }
        taskRecords = append(taskRecords, ToTaskRecord(taskRecordEntity, user, ToTask(taskEntity)))
    }

    return taskRecords, nil
}

func getUser(client *sql.DB, address string) (*trading.User, error) {
    query := `
    SELECT id, address, amount, points FROM User0
    WHERE address = $1
    `
    rows, err := client.Query(query, address)
    if err != nil {
        log.Printf("Query getUser fail: %v", err)
        return nil, err
    }
    defer rows.Close()
    userEntity := &UserEntity{}

    for rows.Next() {
        err = rows.Scan(&userEntity.Id, &userEntity.Address, &userEntity.Amount, &userEntity.Points)
        if err != nil {
            log.Printf("Parse getUser entity fail: %v", err)
            return nil, err
        }
    }

    if userEntity.Address == "" {
        return nil, nil
    }

    return ToUser(userEntity), nil
}

func (repository *UserRepository) FindAllUserTasks() ([]*trading.User, error) {
    query := `
    SELECT u.id as uid, u.address as address, u.amount as amount, u.points as points, t.name as name, tr.id as trid, tr.status as task_status, tr.amount as task_amount, tr.points as task_points FROM TaskRecord tr
    LEFT JOIN User0 u ON tr.user_id = u.id
    LEFT JOIN Task t ON tr.task_id = t.id
    `

    rows, err := repository.client.Query(query)
    if err != nil {
        log.Printf("Query findAllUserTasks fail: %v", err)
        return nil, err
    }
    defer rows.Close()

    usersMap := make(map[string]*trading.User)

    for rows.Next() {
        userEntity := &UserEntity{}
        taskRecordEntity := &TaskRecordEntity{}
        taskEntity := &TaskEntity{}

        err = rows.Scan(&userEntity.Id, &userEntity.Address, &userEntity.Amount, &userEntity.Points, &taskEntity.Name, &taskRecordEntity.Id, &taskRecordEntity.Status, &taskRecordEntity.Amount, &taskRecordEntity.Points)
        if err != nil {
            log.Printf("Parse findAllUserTasks entity fail: %v", err)
            return nil, err
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

    return users, nil
}

func (repository *UserRepository) SaveAllUser(users []*trading.User) error {
    tx, err := repository.client.Begin()
    if err != nil {
        log.Printf("Transaction saveAllUser fail: %v", err)
        return err
    }

    query := `
    SELECT id, name FROM Task
    `

    rows, err := repository.client.Query(query)
    if err != nil {
        log.Printf("Query saveAllUser fail: %v", err)
        return err
    }

    tasksMap := make(map[string]*TaskEntity)

    for rows.Next() {
        task := &TaskEntity{}

        err = rows.Scan(&task.Id, &task.Name)
        if err != nil {
            log.Printf("Parse saveAllUser entity fail: %v", err)
            return err
        }
        tasksMap[task.Name] = task
    }

    err = mergeUser(tx, users)
    if err != nil {
        return err
    }
    err = mergeTaskRecord(tx, users, tasksMap)
    if err != nil {
        return err
    }

    err = tx.Commit()
    if err != nil {
        log.Printf("Commit saveAllUser fail: %v", err)
        return err
    }
    return nil
}

func mergeUser(tx *sql.Tx, users []*trading.User) error {
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
        log.Printf("Prepare mergeUser fail: %v", err)
        newErr := tx.Rollback()
        if newErr != nil {
            log.Printf("Rollback mergeUser fail: %v", err)
            return newErr
        }
    }
    defer stmt.Close()

    for _, user := range users {
        entity := ToUserEntity(user)
        _, err = stmt.Exec(entity.Id, entity.Address, entity.Amount, entity.Points, entity.CreatedAt, entity.UpdatedAt)
        if err != nil {
            log.Printf("Parse mergeUser entity fail: %v", err)
            newErr := tx.Rollback()
            if newErr != nil {
                log.Printf("Rollback mergeUser entity fail: %v", err)
                return newErr
            }
        }
    }
    return nil
}

func mergeTaskRecord(tx *sql.Tx, users []*trading.User, tasksMap map[string]*TaskEntity) error {
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
        log.Printf("Prepare mergeTaskRecord fail: %v", err)
        newErr := tx.Rollback()
        if newErr != nil {
            log.Printf("Rollback mergeTaskRecord fail: %v", err)
            return newErr
        }
    }
    defer stmt.Close()

    for _, user := range users {
        for _, taskRecord := range user.GetTaskRecords() {
            entity := ToTaskRecordEntity(taskRecord, user, tasksMap)
            _, err = stmt.Exec(entity.Id, entity.UserId, entity.TaskId, entity.Status, entity.Amount, entity.Points, entity.CreatedAt, entity.UpdatedAt)
            if err != nil {
                log.Printf("Parse mergeTaskRecord entity fail: %v", err)
                newErr := tx.Rollback()
                if newErr != nil {
                    log.Printf("Rollback mergeTaskRecord entity fail: %v", err)
                    return newErr
                }
            }
        }
    }
    return err
}
