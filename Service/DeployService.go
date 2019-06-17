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

/**
Has a cache of the data as well as the pointer to the DB for queries
 */

type DeployService struct {
	deploys []Models.Deploys
	deployDB *Database.DeployDB
}


func PopulateDeploys(file os.File) (*DeployService, error) {
	/**
	Initializes the Deploy Service. The db object to query the DB or a cache which holds the data
	 */
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
		if err := rows.Scan(&d.Id, &d.Sha, &t, &d.Action, &d.Engineer); err != nil {
			log.Printf("Something went wrong in reading data to cache %v", err)
		}
		d.Date = time.Unix(t, 0)
		result = append(result, d)
	}
	return &DeployService{result, db}, nil
}


func (d * DeployService) GetDistinctEngineers() ([]string, error) {
	/**
	Query based method to get the distinct set of engineers from the table
	 */
	rows, err := d.deployDB.Query("select distinct engineer from deploys")
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for rows.Next() {
		var eng string
		if err := rows.Scan(&eng); err != nil {
			log.Printf("Something went wrong in reading data from DB %v", err)
		}
		result = append(result, eng)
	}
	return result, nil
}


func (d* DeployService) GetActionsOfEng(eng string) ([][]string, error) {
	/**
	Query based method to get the actions of the engineer
	Returns the time and the action that was performed at that time by the engineer
	 */
	rows, err := d.deployDB.Query("select date, action from deploys where engineer=\"" + eng +"\"")
	if err != nil {
		return nil, err
	}
	result := make([][]string, 0)
	for rows.Next() {
		var eng string
		var t int64
		if err := rows.Scan(&t, &eng); err != nil {
			log.Printf("Something went wrong in reading data from DB %v", err)
		}
		var temp []string
		temp = append(temp, time.Unix(t, 0).String(), eng)
		result = append(result, temp)
	}
	return result, nil
}


func validateInput(fromString string, toString string) (int64, int64, error) {
	/**
	Internal validation method to check for correctness in the query input (for date ranges)
	 */
	formatMessage := "i only support unix seconds format"
	rangeMessage := "invalid range, please choose a proper range"
	if len(fromString) != 10 || len(toString) != 10 {
		return -1, -1, errors.New(formatMessage)
	} else {
		from, err := strconv.Atoi(fromString)
		if err != nil {
			return -1, -1, errors.New(formatMessage)
		}
		to, err := strconv.Atoi(toString)
		if err != nil {
			return -1, -1, errors.New(formatMessage)
		}
		if from > to {
			return -1, -1, errors.New(rangeMessage)
		}
		return int64(from), int64(to), nil
	}
}


func (d *DeployService) GetEventsDuring(start string, stop string) ([]map[string]interface{}, error) {
	/**
	Query based method to get the consolidated event summary during the given time range
	 */
	from, to, err := validateInput(start, stop)
	if err != nil {
		return nil, err
	}
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
		if err := rows.Scan(&eng, &action, &t); err != nil {
			log.Printf("Something went wrong in reading data from DB %v", err)
		}
		temp := make(map[string]interface{})
		temp["action"] = action
		temp["engineer"] = eng
		temp["Date"] = time.Unix(t, 0)
		result = append(result, temp)
	}
	return result, nil
}


func (d *DeployService) GetEngineers() ([]string, error) {
	/**
	Cache retrieval of the distinct engineers in the DB
	 */
	engineers := make(map[string]bool)
	result := make([]string, 0)
	if len(d.deploys) == 0 {
		return nil, errors.New("no records available")
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
	/**
	Cache retrieval of the actions and dates the actions was performed during that date by the engineer
	 */
	result := make([][]string, 0)
	if len(d.deploys) == 0 {
		return nil, errors.New("no records available")
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


func (d *DeployService) GetEvents(begin string, end string) ([]map[string]interface{}, error) {
	/**
	Cache retrieval of the events that took place between a specified time range
	 */
	from, to, err := validateInput(begin, end)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	start, stop := time.Unix(from, 0), time.Unix(to, 0)
	if len(d.deploys) == 0 {
		return nil, errors.New("no records available")
	}
	for _, ele := range d.deploys {
		if (start.Before(ele.Date) || start.Equal(ele.Date)) && (stop.After(ele.Date) || stop.Equal(ele.Date)) {
			temp := make(map[string]interface{})
			temp["action"] = ele.Action
			temp["engineer"] = ele.Engineer
			temp["Date"] = ele.Date
			result = append(result, temp)
		}
	}
	return result, nil
}


func (d * DeployService) getDate(date *string, all bool) error {
	if len(*date) != 10 && !all {
		return errors.New("not a valid date, please provide unix 10 digit timestamp")
	}
	if !all {
		inputDate, err := strconv.Atoi(*date)
		if err != nil {
			return errors.New("not a valid date, please provide unix 10 digit timestamp")
		}
		*date = time.Unix(int64(inputDate), 0).Format("2006-01-02")
	}
	return nil
}


func (d * DeployService) getEventStatsHelper(date string, all bool) (map[string]map[string]int64, error) {
	/**
	Retrieves the Stats of the events that was performed on a given date or a consolidated date
	 */
	if len(d.deploys) == 0 {
		return nil, errors.New("no records available")
	}
	if err := d.getDate(&date, all); err != nil {
		return nil, err
	}
	result := make(map[string]map[string]int64)
	for _, ele := range d.deploys {
		key := ele.Date.Format("2006-01-02")
		if key == date || all {
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
	}
	return result, nil
}


func (d *DeployService) GetEventStats(date string) (map[string]map[string]int64, error) {
	if date != "" {
		return d.getEventStatsHelper(date, false)
	} else {
		return d.getEventStatsHelper(date, true)
	}
}