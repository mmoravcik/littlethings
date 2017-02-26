"""
Usage:

`python streak.py [streak_amount]`
Default streak_amount is 4
"""

import requests
import re
import sys
from bs4 import BeautifulSoup

url = "http://www.sportsnet.ca/hockey/nhl/standings/"

min_streak = int(sys.argv[1]) if len(sys.argv) > 1 else 4

r = requests.get(url)

html = r.text

soup = BeautifulSoup(html, 'html.parser')

table_data = []

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
        team = data[1]
        streak = data[-1]

        streak_count = int(re.findall('[0-9]+', streak)[0])
        if streak_count >= min_streak:
            results.append({'team': team, 'streak': streak})


for result in results:
    print result['team'], result['streak']