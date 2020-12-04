# pd-tools-bot
![Go](https://github.com/jaceklubzinski/pd-tools-bot/workflows/Go/badge.svg?branch=master)
![golangci-lint](https://github.com/jaceklubzinski/pd-tools-bot/workflows/golangci-lint/badge.svg?branch=master)
![security scan](https://github.com/jaceklubzinski/pd-tools-bot/workflows/security%20scan/badge.svg?branch=master)
![docker build](https://github.com/jaceklubzinski/pd-tools-bot/workflows/docker%20build/badge.svg?branch=latest)

# Deployment
Public docker image available on docker hub https://hub.docker.com/repository/docker/jlubzinski/pd-tools-bot

docker compose [deployments/docker-compose.yml](deployments/docker-compose.yml) requires environment variables to work
```
# required
PDBOT_PAGER_DUTY_AUTH_TOKEN:
PDBOT_SLACK_AUTH_TOKEN:
PDBOT_DUTY_PAY="holiday:xx,workday:xx,weekend:xx"
# optional
## to use short names instead of PD ids
PDBOT_PAGER_DUTY_TEAM_ID: "teamName:teamID,teamName:teamID"
PDBOT_PAGER_DUTY_SCHEDULE_ID: "scheduleName:scheduleID,scheduleName:scheduleID"
## channel authorized to set maintenance mode
PDBOT_SLACK_AUTHORIZED_CHANNELS:
```
# Supported slack command
```
help - help
oncall today - PagerDuty today oncall user
Example: oncall today
oncall month lts current|next|prev - PagerDuty oncall current/next/prev month summary with profit
Example: oncall month lts
incident list pdteam - PagerDuty list of triggered and acknowledged incident for specific team
Example: incident list lts
incident duty pdteam pdhour - PagerDuty list all incident incident for specific team and since defined hours
Example: incident duty lts 24h
schedule list - PagerDuty schedule list
Example: schedule list
team list - PagerDuty team list
Example: team list
service list pdteam - PagerDuty list of services assigned to specific team
Example: service list lts
maintenance list pdteam - PagerDuty list of service maintenance to specific team
Example: maintenace list lts
maintenance create pdservice pdhour - PagerDuty create maintenance window for specific service from current time + given duration *
Example: maintenace create lts 4h
secret - Very secret stuff *
* Authorized users only
```