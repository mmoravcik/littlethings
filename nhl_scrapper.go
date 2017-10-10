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


type Game struct {
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
}


type Schedule struct {
	Dates     []struct {
		Date   string        `json:"date"`
		Games  []Game
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
}


func getBodyFromJson(url string) []byte {
    client := http.Client{
        Timeout: time.Second * 3,
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


const MINIMAL_STREAK = 4

func main() {
    var start_date = time.Now().AddDate(0, 0, -1).Local().Format("2006-01-02")
    var end_date = time.Now().AddDate(0, 0, 1).Local().Format("2006-01-02")

    schedule_url := fmt.Sprintf(
        "https://statsapi.web.nhl.com/api/v1/schedule?startDate=%s&endDate=%s&expand=schedule.teams,schedule.teams.team&site=en_nhl&teamId=&gameType=&timecode=",
        start_date,
        end_date,
    )
    standings_url := "https://statsapi.web.nhl.com/api/v1/standings"

    var standings_body = getBodyFromJson(standings_url)
    var standings Standings

    jsonStandingsErr := json.Unmarshal(standings_body, &standings)
    if jsonStandingsErr != nil {
        log.Fatal(jsonStandingsErr)
    }

    // This is where we store teams with desired streak
    var streak_team_records []TeamRecord
    // This is where we store games with teams on streak
    var streak_games []Game

    // Go through all divisions and their teams to populate streak_teams
    for _, division := range standings.Records {
        for _, team_record := range division.TeamRecords {
            if team_record.Streak.StreakNumber == MINIMAL_STREAK {
                streak_team_records = append(streak_team_records, team_record)
            }
        }
    }

    var schedule_body = getBodyFromJson(schedule_url)

    var schedule Schedule

    jsonScheduleErr := json.Unmarshal(schedule_body, &schedule)
    if jsonScheduleErr != nil {
        log.Fatal(jsonScheduleErr)
    }

    fmt.Println(fmt.Sprintf("Games for %s - %s date range with minimal streak of %d", start_date, end_date, MINIMAL_STREAK))

    // Go through the schedule and store games with streak teams
    for _, date := range schedule.Dates {
        for _, game := range(date.Games) {
            // TODO this is ineffective
            for _, team_record := range streak_team_records {
                // If any of the teams are on the streak, print them
                if game.Teams.Away.Team.ID == team_record.Team.ID || game.Teams.Home.Team.ID == team_record.Team.ID {
                    streak_games = append(streak_games, game)
                    var s = fmt.Sprintf(
                        "%s: %s vs %s - %s with %s streak",
                        date.Date,
                        game.Teams.Away.Team.Abbreviation,
                        game.Teams.Home.Team.Abbreviation,
                        team_record.Team.Name,
                        team_record.Streak.StreakCode,
                    )
                    fmt.Println(s)
                }
            }
        }
    }

    if len(streak_games) == 0 {
        fmt.Println("No games matching the criteria")
    }

    fmt.Println("FINISHED")

}
