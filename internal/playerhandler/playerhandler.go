package playerhandler

import (
    "database/sql"
    "log"
    "errors"
    "github.com/mattn/go-sqlite3"
)

var (
    ErrDuplicate    = errors.New("record already exists")
    ErrNotExists    = errors.New("row not exists")
    ErrUpdateFailed = errors.New("update failed")
    ErrDeleteFailed = errors.New("delete failed")
)

type Database struct {
    inst *sql.DB
}

type Player struct {
    Name string
    Inventory []string
}

func LoadDatabase(filename string) (*Database) {
    db, err := sql.Open("sqlite3", filename + ".db")
    if err != nil {
        log.Fatal(err)
    }

    tmpDB := new(Database) 
    tmpDB.inst = db

    err = tmpDB.InitTable()
    if err != nil {
        log.Fatal(err)
    }
    return tmpDB
}

func (db *Database) InitTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS players(
        name TEXT NOT NULL UNIQUE
    );
    `

    _, err :=  db.inst.Exec(query)
    return err
}

func (db *Database) LoadPlayer(playername string) (*Player,error) {
    row := db.inst.QueryRow("SELECT * FROM players WHERE name=?", playername)

    var player Player
    if err := row.Scan(&player.Name); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotExists
        }
        return nil, err
    }
    return &player, nil
}

func (db *Database) CreatePlayer(name string) error {
    _, err := db.inst.Exec("INSERT INTO players (name) VALUES(?)", name)
    if err != nil {
        var sqliteErr sqlite3.Error
        if errors.As(err, &sqliteErr) {
            if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
                return ErrDuplicate
            }
        }
        return err
    }

    return nil
}
