package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/jmoiron/sqlx"
)

type apiHandler struct {
	lock sync.Mutex
	db   *sqlx.DB
	stmt *sqlx.Stmt
}

const getAllDinos = `
SELECT
	d.name,
	w.name as world,
	c1.name as class_name,
	d.dino_id1|d.dino_id2 as dino_id,
	level_wild,
	level_tamed,
	level_total,
	is_cryo,
	c2.name as parent_class,
	parent_name,
	x, y, z,
	color0, color1, color2, color3, color4, color5,
	health_current, stamina_current, torpor_current, oxygen_current, food_current, weight_current, melee_current, speed_current,
	health_wild, stamina_wild, torpor_wild, oxygen_wild, food_wild, weight_wild, melee_wild, speed_wild,
	health_tamed, stamina_tamed, torpor_tamed, oxygen_tamed, food_tamed, weight_tamed, melee_tamed, speed_tamed,
	health_total, stamina_total, torpor_total, oxygen_total, food_total, weight_total, melee_total, speed_total
FROM
	dinos d
LEFT JOIN worlds w ON d.world == w.id
LEFT JOIN classes c1 ON d.class == c1.id
LEFT JOIN classes c2 ON d.parent_class == c2.id
`

type dinoResult struct {
	Name        string `json:"name" db:"name"`
	World       string `json:"world" db:"world"`
	Class       string `json:"class_name" db:"class_name"`
	DinoId      int    `json:"dino_id" db:"dino_id"`
	LevelsWild  int    `json:"levels_wild" db:"level_wild"`
	LevelsTamed int    `json:"levels_tamed" db:"level_tamed"`
	LevelsTotal int    `json:"levels_total" db:"level_total"`

	IsCryopod   bool    `json:"is_cryo" db:"is_cryo"`
	ParentClass *string `json:"parent_class" db:"parent_class"`
	ParentName  *string `json:"parent_name" db:"parent_name"`

	X float64 `json:"x" db:"x"`
	Y float64 `json:"y" db:"y"`
	Z float64 `json:"z" db:"z"`

	Color0 int `json:"color0" db:"color0"`
	Color1 int `json:"color1" db:"color1"`
	Color2 int `json:"color2" db:"color2"`
	Color3 int `json:"color3" db:"color3"`
	Color4 int `json:"color4" db:"color4"`
	Color5 int `json:"color5" db:"color5"`

	HealthCurrent  float64 `json:"health_current" db:"health_current"`
	StaminaCurrent float64 `json:"stamina_current" db:"stamina_current"`
	TorporCurrent  float64 `json:"torpor_current" db:"torpor_current"`
	OxygenCurrent  float64 `json:"oxygen_current" db:"oxygen_current"`
	FoodCurrent    float64 `json:"food_current" db:"food_current"`
	WeightCurrent  float64 `json:"weight_current" db:"weight_current"`
	MeleeCurrent   float64 `json:"melee_current" db:"melee_current"`
	SpeedCurrent   float64 `json:"speed_current" db:"speed_current"`

	HealthWild  int64 `json:"health_wild" db:"health_wild"`
	StaminaWild int64 `json:"stamina_wild" db:"stamina_wild"`
	TorporWild  int64 `json:"torpor_wild" db:"torpor_wild"`
	OxygenWild  int64 `json:"oxygen_wild" db:"oxygen_wild"`
	FoodWild    int64 `json:"food_wild" db:"food_wild"`
	WeightWild  int64 `json:"weight_wild" db:"weight_wild"`
	MeleeWild   int64 `json:"melee_wild" db:"melee_wild"`
	SpeedWild   int64 `json:"speed_wild" db:"speed_wild"`

	HealthTamed  int64 `json:"health_tamed" db:"health_tamed"`
	StaminaTamed int64 `json:"stamina_tamed" db:"stamina_tamed"`
	TorporTamed  int64 `json:"torpor_tamed" db:"torpor_tamed"`
	OxygenTamed  int64 `json:"oxygen_tamed" db:"oxygen_tamed"`
	FoodTamed    int64 `json:"food_tamed" db:"food_tamed"`
	WeightTamed  int64 `json:"weight_tamed" db:"weight_tamed"`
	MeleeTamed   int64 `json:"melee_tamed" db:"melee_tamed"`
	SpeedTamed   int64 `json:"speed_tamed" db:"speed_tamed"`

	HealthTotal  int64 `json:"health_total" db:"health_total"`
	StaminaTotal int64 `json:"stamina_total" db:"stamina_total"`
	TorporTotal  int64 `json:"torpor_total" db:"torpor_total"`
	OxygenTotal  int64 `json:"oxygen_total" db:"oxygen_total"`
	FoodTotal    int64 `json:"food_total" db:"food_total"`
	WeightTotal  int64 `json:"weight_total" db:"weight_total"`
	MeleeTotal   int64 `json:"melee_total" db:"melee_total"`
	SpeedTotal   int64 `json:"speed_total" db:"speed_total"`
}

func (ah *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.lock.Lock()
	defer ah.lock.Unlock()

	if ah.stmt == nil {
		stmt, err := ah.db.Preparex(getAllDinos)
		if err != nil {
			log.Printf("Error preparing SQL: %s", err)
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		ah.stmt = stmt
	}

	var result []dinoResult
	err := ah.stmt.Select(&result)
	if err != nil {
		log.Printf("Error querying database: %s", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling result: %s", err)
		http.Error(w, "JSON Error", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
