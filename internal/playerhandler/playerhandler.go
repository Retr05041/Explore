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

// Loads the database corresponding to the map
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

// If the tables don't exist, create them
func (db *Database) InitTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS players(
        player_name TEXT NOT NULL UNIQUE
    );
    CREATE TABLE IF NOT EXISTS inventory(
        inventory_id INTEGER PRIMARY KEY AUTOINCREMENT,
        player_name TEXT NOT NULL UNIQUE,
        item TEXT NOT NULL UNIQUE,
        FOREIGN KEY (player_name) REFERENCES players(player_name)
    );
    `

    _, err :=  db.inst.Exec(query)
    return err
}

// Populate the player struct given a player name (each one is uniue)
func (db *Database) LoadPlayer(playername string) (*Player,error) {
    // Load Player
    row := db.inst.QueryRow("SELECT * FROM players WHERE player_name=?", playername)

    var player Player
    if err := row.Scan(&player.Name); err != nil {
        // this is where we catch out errors for selecting a player
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotExists
        }
        return nil, err
    }

    // Load Inv
    rows, err := db.inst.Query("SELECT item FROM inventory WHERE player_name=?", playername)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var item string
        if err := rows.Scan(&item); err != nil {
            return nil, err
        }
        player.Inventory = append(player.Inventory, item)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return &player, nil
}

// Create a player if they don't exist
func (db *Database) CreatePlayer(name string) error {
    _, err := db.inst.Exec("INSERT INTO players(player_name) VALUES(?)", name)
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

func (p *Player) AddToInv(item string) {
   p.Inventory = append(p.Inventory, item) 
}

func (db *Database) SavePlayerInfo(p *Player) error {
    for _, item := range p.Inventory {
        _, err := db.inst.Exec("INSERT INTO inventory(player_name, item) VALUES(?,?)", p.Name, item)
        if err != nil {
            return err
        }
    }
    return nil
}
