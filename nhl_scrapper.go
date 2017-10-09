package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"
)

type TeamRecord struct {
    Streak struct {
        StreakCode   string `json:"streakCode"`
        StreakNumber int    `json:"streakNumber"`
        StreakType   string `json:"streakType"`
    } `json:"streak"`
    Team struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    } `json:"team"`
}


type Standings struct {
	Records   []struct {
		TeamRecords  []TeamRecord
	} `json:"records"`
}



type Schedule struct {
	Copyright string `json:"copyright"`
	Dates     []struct {
		Date   string        `json:"date"`
		Games  []struct {
			GameDate string `json:"gameDate"`
			GameType string `json:"gameType"`
			Teams struct {
				Away struct {
					Team  struct {
						Abbreviation string `json:"abbreviation"`
						ID              int    `json:"id"`
						Name            string `json:"name"`
						ShortName       string `json:"shortName"`
						TeamName        string `json:"teamName"`
					} `json:"team"`
				} `json:"away"`
				Home struct {
					Team  struct {
						Abbreviation string `json:"abbreviation"`
						ID              int    `json:"id"`
						Name            string `json:"name"`
						ShortName       string `json:"shortName"`
						TeamName        string `json:"teamName"`
					} `json:"team"`
				} `json:"home"`
			} `json:"teams"`
		} `json:"games"`
		Matches      []interface{} `json:"matches"`
		TotalEvents  int           `json:"totalEvents"`
		TotalGames   int           `json:"totalGames"`
		TotalItems   int           `json:"totalItems"`
		TotalMatches int           `json:"totalMatches"`
	} `json:"dates"`
	TotalEvents  int `json:"totalEvents"`
	TotalGames   int `json:"totalGames"`
	TotalItems   int `json:"totalItems"`
	TotalMatches int `json:"totalMatches"`
	Wait         int `json:"wait"`
}


func getBodyFromJson(url string) []byte {
    client := http.Client{
        Timeout: time.Second * 2, // Maximum of 2 secs
    }

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        log.Fatal(err)
    }

    req.Header.Set("User-Agent", "go-api-test")

    res, getErr := client.Do(req)
    if getErr != nil {
        log.Fatal(getErr)
    }

    body, readErr := ioutil.ReadAll(res.Body)
    if readErr != nil {
        log.Fatal(readErr)
    }

    return body
}


const MINIMAL_STREAK = 2

func main() {


    schedule_url := fmt.Sprintf(
        "https://statsapi.web.nhl.com/api/v1/schedule?startDate=%s&endDate=%s&expand=schedule.teams,schedule.teams.team&site=en_nhl&teamId=&gameType=&timecode=",
        time.Now().Local().Format("2006-01-02"),
        time.Now().Local().Format("2006-01-02"),
    )
    standings_url := "https://statsapi.web.nhl.com/api/v1/standings"

    var standings_body = getBodyFromJson(standings_url)
    var standings Standings

    jsonStandingsErr := json.Unmarshal(standings_body, &standings)
    if jsonStandingsErr != nil {
        log.Fatal(jsonStandingsErr)
    }

    var streak_teams []TeamRecord

    for _, division := range standings.Records {
        for _, team_record := range division.TeamRecords {
            //fmt.Println(team_record.Team.Name)
            //fmt.Println(team_record.Streak.StreakNumber)
            if team_record.Streak.StreakNumber == MINIMAL_STREAK {
                streak_teams = append(streak_teams, team_record)
            }


            //fmt.Println(team_record.Streak.StreakNumber == MINIMAL_STREAK)
            //fmt.Println("\n")
        }
    }

    for _, team_record := range streak_teams {
        fmt.Println(team_record)
    }

    var schedule_body = getBodyFromJson(schedule_url)

    var schedule Schedule

    jsonScheduleErr := json.Unmarshal(schedule_body, &schedule)
    if jsonScheduleErr != nil {
        log.Fatal(jsonScheduleErr)
    }

    for _, date := range schedule.Dates {

        for _, game := range(date.Games) {
            var _ = fmt.Sprintf(
                "%s vs %s on %s",
                game.Teams.Away.Team.Abbreviation,
                game.Teams.Home.Team.Abbreviation,
                date.Date,
            )
            //fmt.Println(s)
        }
    }
}
