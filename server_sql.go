package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"time"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var ch_sql = make(chan Task, 50)

type FileId struct {
	gorm.Model
	//Делаем ссылку уникальной, иначе gorm clause работать не будет
	Link string `gorm:"unique"`
	Code string
}

func server_init() {
	fill_base_solved()
	fill_base_empty()
	queue(Task{},3)
	go update_sql()
}

func update_sql(){
	defer wg.Done()
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		logmes(9,"Can't open base", err.Error())
		return
	}
	db.AutoMigrate(&FileId{})
	for {
		task, ok := <-ch_sql
		if !ok {
			return
		}
		if task.Link == "" {
			continue
		}
		if task.Code == "" {
			db.Model(&FileId{}).Where("Link = ?", task.Link).Update("code", nil)
		} else {
			db.Model(&FileId{}).Where("Link = ?", task.Link).Update("code", task.Code)
		}
	}
}

func get_link_sql() (link_result string) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		//Меняем уровень логера, тк при создании 200+ строк выскакивают предупреждения о скорости
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		logmes(9,"Can't open base", err.Error())
		return
	}
	db.AutoMigrate(&FileId{})
	var result []FileId
	db.Where("code = ?", "").Order("updated_at").First(&result)
	if len(result) == 0 {
		return
	}
	link_result = result[0].Link

	sqlDB, _ := db.DB()
	sqlDB.Close()
	return
}


func fill_base_solved() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		//Меняем уровень логера, тк при создании 200+ строк выскакивают предупреждения о скорости
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		logmes(9,"Can't open base", err.Error())
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.AutoMigrate(&FileId{})

	file, err := os.Open("txtfiles\\links.txt")
	if err != nil {
		logfatal("Can't open file with links", err)
	}
	defer file.Close()

	var files = []FileId{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		ntext := strings.TrimSuffix(text, "@2x.webp")
		ntext = strings.TrimSuffix(ntext, ".webp")
		ntext = strings.TrimSuffix(ntext, ".jpg")
		if text != ntext {
			files = append(files, FileId{Link: ntext[:len(ntext)-6], Code: ntext[len(ntext)-5:]})
		}
	}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "link"}},
		DoUpdates: clause.AssignmentColumns([]string{"code"}),
	}).CreateInBatches(&files, 1000)
}

func fill_base_empty() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		//Меняем уровень логера, тк при создании 200+ строк выскакивают предупреждения о скорости
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		logmes(9,"Can't open base", err.Error())
		return
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.AutoMigrate(&FileId{})

	file, err := os.Open("txtfiles\\check.txt")
	if err != nil {
		logfatal("Can't open file with links", err)
	}
	defer file.Close()

	var files = []FileId{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		files = append(files, FileId{Link: text})
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&files, 1000)
}

func check_link(info Task) (result bool) {
	conn := &http.Client{
		//	Timeout:   4 * time.Second,
	}
	link, _ := get_suffix_rune(&info.Link)
	if link == "" {
		return
	}
	link = "https://cdn.townofsins.com/media/assets/images/" + info.Link + "_" + info.Code + link
	for {
		res, err := conn.Head(link)
		if err == nil && res != nil && res.StatusCode == 200 {
			return true
		} else if err == nil && res != nil && res.StatusCode == 500 {
			break
		}
	}
	return
}

// mode:
// 1 - создание; 2 - удаление, 3 - заполнение карты
func queue(info Task, mode uint8) {
	db, err := gorm.Open(sqlite.Open("queue.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		logmes(9,"Can't open base with queue", err.Error())
		return
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.AutoMigrate(&FileId{})
	if  mode == 1 {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&FileId{Link: info.Link, Code: info.Secret})
	} else if mode == 2 {
		var result []FileId
		db.Where("Link = ?", info.Link).Delete(&result)
	} else if mode == 3 {
		ctime := time.Now().Add(time.Duration(-24) * time.Hour)
		db.Where("updated_at < ?", ctime).Delete(&FileId{})

		rows, _ := db.Model(&FileId{}).Rows()
		defer rows.Close()
		for rows.Next() {
			var result FileId
			db.ScanRows(rows, &result)
			TaskList[result.Link] = result.Code	
		}
	}
	return
}


