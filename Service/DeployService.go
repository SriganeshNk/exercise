package Service

import (
	"errors"
	"log"
	"os"
	"slackProject/Database"
	"slackProject/Models"
	"strconv"
	"time"
)

type DeployService struct {
	deploys []Models.Deploys
	deployDB *Database.DeployDB
}


func PopulateDeploys(file os.File) (*DeployService, error) {
	db, err := Database.NewDataBase(file.Name())
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select * from deploys")
	if err != nil {
		return nil, err
	}
	var result []Models.Deploys
	for rows.Next() {
		d := Models.Deploys{}
		var t int64
		rows.Scan(&d.Id, &d.Sha, &t, &d.Action, &d.Engineer)
		d.Date = time.Unix(t, 0)
		result = append(result, d)
	}
	return &DeployService{result, db}, nil
}


func (d * DeployService) GetDistinctEngineers() ([]string, error) {
	rows, err := d.deployDB.Query("select distinct engineer from deploys")
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for rows.Next() {
		var eng string
		rows.Scan(&eng)
		result = append(result, eng)
	}
	return result, nil
}


func (d* DeployService) GetActionsOfEng(eng string) ([][]string, error) {
	rows, err := d.deployDB.Query("select date, action from deploys where engineer=\"" + eng +"\"")
	if err != nil {
		return nil, err
	}
	result := make([][]string, 0)
	for rows.Next() {
		var eng string
		var t int64
		rows.Scan(&t, &eng)
		var temp []string
		temp = append(temp, time.Unix(t, 0).String(), eng)
		result = append(result, temp)
	}
	return result, nil
}


func (d *DeployService) GetEventsDuring(from int64, to int64) ([]map[string]interface{}, error) {
	q := "select engineer, action, date from deploys where date >= "+
		strconv.Itoa(int(from)) + " and date <= "  + strconv.Itoa(int(to))
	log.Print(q)
	rows, err := d.deployDB.Query(q)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		var eng string
		var t int64
		var action string
		rows.Scan(&eng, &action, &t)
		temp := make(map[string]interface{})
		temp["action"] = action
		temp["engineer"] = eng
		temp["Date"] = time.Unix(t, 0)
		result = append(result, temp)
	}
	return result, nil
}


func (d *DeployService) GetEngineers() ([]string, error) {
	engineers := make(map[string]bool)
	result := make([]string, 0)
	if len(d.deploys) == 0 {
		return nil, errors.New("No records available")
	}
	for _, ele := range d.deploys {
		_, isPresent := engineers[ele.Engineer]
		if !isPresent {
			engineers[ele.Engineer] = true
			result = append(result, ele.Engineer)
		}
	}
	return result, nil
}


func (d *DeployService) GetActionsOfEngineer(eng string) ([][]string, error) {
	result := make([][]string, 0)
	if len(d.deploys) == 0 {
		return nil, errors.New("No records available")
	}
	for _, ele := range d.deploys {
		if ele.Engineer == eng {
			temp := make([]string, 2)
			temp = append(temp, ele.Action)
			temp = append(temp, ele.Date.String())
			result = append(result, temp)
		}
	}
	return result, nil
}

func (d *DeployService) GetEvents(from int64, to int64) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	start, end := time.Unix(from, 0), time.Unix(to, 0)
	if len(d.deploys) == 0 {
		return nil, errors.New("No records available")
	}
	for _, ele := range d.deploys {
		if (start.Before(ele.Date) || start.Equal(ele.Date)) && (end.After(ele.Date) || end.Equal(ele.Date)) {
			temp := make(map[string]interface{})
			temp["action"] = ele.Action
			temp["engineer"] = ele.Engineer
			temp["Date"] = ele.Date
			result = append(result, temp)
		}
	}
	return result, nil
}


func (d *DeployService) GetEventStats() (map[string]map[string]int64, error) {
	result := make(map[string]map[string]int64, 0)
	if len(d.deploys) == 0 {
		return nil, errors.New("No records available")
	}
	for _, ele := range d.deploys {
		key := ele.Date.Format("2006-01-02")
		_, isPresent := result[key]
		if !isPresent {
			result[key] = make(map[string]int64)
			result[key][ele.Action] = 1
		} else {
			_, isPresent := result[key][ele.Action]
			if !isPresent {
				result[key][ele.Action] = 1
			} else {
				result[key][ele.Action] += 1
			}
		}
	}
	return result, nil
}