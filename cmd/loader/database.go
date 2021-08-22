package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/justinian/ark"
)

var schema = []string{`
	CREATE TABLE worlds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);`,
	`CREATE TABLE classes (
		id INTEGER PRIMARY KEY,
		class TEXT,
		name TEXT
	);`,
	`CREATE TABLE dinos (
		id INTEGER,
		list INTEGER,
		world INTEGER,
		class INTEGER,
		name TEXT,
		level_wild INTEGER,
		level_tamed INTEGER,
		dino_id1 INTEGER,
		dino_id2 INTEGER,
		is_cryo BOOLEAN,
		parent_class INTEGER,
		parent_name TEXT,
		x FLOAT,
		y FLOAT,
		z FLOAT,

		color0 INTEGER,
		color1 INTEGER,
		color2 INTEGER,
		color3 INTEGER,
		color4 INTEGER,
		color5 INTEGER,

		health_current FLOAT,
		stamina_current FLOAT,
		torpor_current FLOAT,
		oxygen_current FLOAT,
		food_current FLOAT,
		weight_current FLOAT,
		melee_current FLOAT,
		speed_current FLOAT,

		health_wild INTEGER,
		stamina_wild INTEGER,
		torpor_wild INTEGER,
		oxygen_wild INTEGER,
		food_wild INTEGER,
		weight_wild INTEGER,
		melee_wild INTEGER,
		speed_wild INTEGER,

		health_tamed INTEGER,
		stamina_tamed INTEGER,
		torpor_tamed INTEGER,
		oxygen_tamed INTEGER,
		food_tamed INTEGER,
		weight_tamed INTEGER,
		melee_tamed INTEGER,
		speed_tamed INTEGER,

		level_total INTEGER AS (level_wild+level_tamed),

		health_total INTEGER AS (health_wild+health_tamed),
		stamina_total INTEGER AS (stamina_wild+stamina_tamed),
		torpor_total INTEGER AS (torpor_wild+torpor_tamed),
		oxygen_total INTEGER AS (oxygen_wild+oxygen_tamed),
		food_total INTEGER AS (food_wild+food_tamed),
		weight_total INTEGER AS (weight_wild+weight_tamed),
		melee_total INTEGER AS (melee_wild+melee_tamed),
		speed_total INTEGER AS (speed_wild+speed_tamed),

		PRIMARY KEY (id, list, world)
	);`,
}

func openDatabase(filename string) (*sql.DB, error) {
	err := os.Rename(filename, fmt.Sprintf("%s.bak", filename))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("Could not move old db file:\n%w", err)
		}
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("Could not open db file:\n%w", err)
	}

	for _, table := range schema {
		_, err = db.Exec(table)
		if err != nil {
			return nil, fmt.Errorf("Could not create SQL schema:\n%w", err)
		}
	}

	return db, nil
}

func processSaves(db *sql.DB, saves []*ark.SaveGame, classNames map[string]string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Could not begin SQL transaction:\n%w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO classes VALUES (?,?,?)")
	if err != nil {
		return fmt.Errorf("Could not prepare SQL class insert:\n%w", err)
	}

	classCount := 1 // leave 0 for "none"
	classes := make(map[string]int, 300)
	for _, save := range saves {
		for _, objlist := range save.Objects {
			for _, obj := range objlist {
				className := obj.ClassName.Name
				if _, ok := classes[className]; !ok {
					classes[className] = classCount
					name, ok := classNames[className]
					if !ok {
						// chomp the _C suffix
						name = className[:len(className)-2]
					}
					_, err = stmt.Exec(classCount, className, name)
					if err != nil {
						return fmt.Errorf("Inserting class: (%d, %s):\n%w", classCount, name, err)
					}
					classCount++
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Could not commit SQL transaction:\n%w", err)
	}

	baseStmt, err := db.Prepare(`INSERT INTO dinos VALUES (
									?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,
									?,?,?,?,?,?,
									?,?,?,?,?,?,?,?,
									?,?,?,?,?,?,?,?,
									?,?,?,?,?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("Could not prepare SQL insert:\n%w", err)
	}

	for _, save := range saves {
		tx, err = db.Begin()
		if err != nil {
			return fmt.Errorf("Could not begin SQL transaction:\n%w", err)
		}

		worldName := save.DataFiles[0]
		if strings.HasSuffix(worldName, "_P") {
			worldName = worldName[:len(worldName)-2]
		}

		res, err := tx.Exec("INSERT INTO worlds (name) VALUEs (?)", worldName)
		if err != nil {
			return fmt.Errorf("Could not insert world name:\n%w", err)
		}

		worldId, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("Could not get world id:\n%w", err)
		}

		stmt := tx.Stmt(baseStmt)

		for listNum, objlist := range save.Objects {
			for i, obj := range objlist {
				server := obj.Properties.Get("TamedOnServerName", 0)
				if server == nil {
					continue
				}

				name := obj.Properties.GetString("TamedName", 0)
				statsCurrent := make([]float64, 12)
				pointsWild := make([]int64, 12)
				pointsTamed := make([]int64, 12)
				var levelWild int64
				var levelTamed int64

				loc := obj.Location
				parentClass := 0
				parentName := ""
				if obj.Parent != nil {
					loc = obj.Parent.Location
					parentClass = classes[obj.Parent.ClassName.Name]
					parentName = obj.Parent.Properties.GetString("BoxName", 0)
					if parentName == "" {
						parentName = obj.Parent.Properties.GetString("PlayerName", 0)
					}
				}

				cscProp := obj.Properties.GetTyped("MyCharacterStatusComponent", 0, ark.ObjectPropertyType)
				if cscProp != nil {
					cscId := cscProp.(*ark.ObjectProperty).Id
					csc := objlist[cscId]

					for index := 0; index < 12; index++ {
						statsCurrent[index] = csc.Properties.GetFloat("CurrentStatusValues", index)
						pointsWild[index] = csc.Properties.GetInt("NumberOfLevelUpPointsApplied", index)
						pointsTamed[index] = csc.Properties.GetInt("NumberOfLevelUpPointsAppliedTamed", index)
					}

					levelWild = csc.Properties.GetInt("BaseCharacterLevel", 0)
					levelTamed = csc.Properties.GetInt("ExtraCharacterLevel", 0)
				}

				dinoId1 := obj.Properties.GetInt("DinoID1", 0)
				dinoId2 := obj.Properties.GetInt("DinoID2", 0)

				colors := make([]int64, 6)
				for i := range colors {
					colors[i] = obj.Properties.GetInt("ColorSetIndices", i)
				}

				classId := classes[obj.ClassName.Name]
				_, err = stmt.Exec(
					i,
					listNum,
					worldId,
					classId,
					name,
					levelWild,
					levelTamed,
					dinoId1,
					dinoId2,
					obj.IsCryopod,
					parentClass,
					parentName,

					loc.X, loc.Y, loc.Z,

					colors[0], colors[1], colors[2],
					colors[3], colors[4], colors[5],

					statsCurrent[0], statsCurrent[1], statsCurrent[2], statsCurrent[3],
					statsCurrent[4], statsCurrent[7], statsCurrent[8], statsCurrent[9],

					pointsWild[0], pointsWild[1], pointsWild[2], pointsWild[3],
					pointsWild[4], pointsWild[7], pointsWild[8], pointsWild[9],

					pointsTamed[0], pointsTamed[1], pointsTamed[2], pointsTamed[3],
					pointsTamed[4], pointsTamed[7], pointsTamed[8], pointsTamed[9],
				)

				if err != nil {
					return fmt.Errorf("Could not insert object %d:\n%w", i, err)
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("Could not commit SQL transaction:\n%w", err)
		}
	}

	return nil
}
