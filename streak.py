"""
Usage:

`python streak.py [streak_amount]`
Default streak_amount is 4
"""

import json
import requests
import re
import sys
from bs4 import BeautifulSoup

# This table has a good streak info, including overtime wins and loses counted
# as normal wins and loses, e.g. LL OT1 will be L3. NHL.com will list OT1
streak_table_url = "http://www.sportsnet.ca/hockey/nhl/standings/"
nhl_standings_api = "https://statsapi.web.nhl.com/api/v1/standings?expand=standings.team,standings.league,team.schedule.next&season=20162017"

min_streak = int(sys.argv[1]) if len(sys.argv) > 1 else 4

r = requests.get(streak_table_url)

html = r.text

soup = BeautifulSoup(html, 'html.parser')

table_data = []

team_data_json = requests.get(nhl_standings_api).content
all_team_data = json.loads(team_data_json)
team_data = {}
for x in xrange(0, 4):  # 4 conferences
    for team in all_team_data['records'][x]['teamRecords']:
        team_data[team['team']['name']] = {
            'next_game': team['team']['nextGameSchedule']['dates'][0]['date']
        }


for link in soup.findChildren("div", {"id": "leaguestandings"}):
    for table in link.find_all("table"):
        rows = table.find_all("tr")
        for row in rows:
            cols = row.find_all('td')
            cols = [el.text.strip() for el in cols]
            table_data.append([el for el in cols if el])

results = []

for data in table_data:
    if data:
        team_name = data[1]
        streak = data[-1]

        for team, data in team_data.iteritems():
            if team_name in team:
                next_game = data['next_game']
                break

        streak_count = int(re.findall('[0-9]+', streak)[0])
        if streak_count >= min_streak:
            results.append({
                'team': team_name,
                'streak': streak,
                'next_game': next_game
            })


for result in results:
    print result['team'], result['streak'], result['next_game']