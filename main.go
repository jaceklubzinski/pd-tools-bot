package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/extensions"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/incident"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/maintenance"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/oncall"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/schedule"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/services"
	team "github.com/jaceklubzinski/pd-tools-bot/pkg/teams"
	"github.com/kelseyhightower/envconfig"
	"github.com/shomali11/slacker"
)

type envConfig struct {
	PagerDutyAuthToken      string            `required:"true" split_words:"true"`
	SlackAuthToken          string            `required:"true" split_words:"true"`
	PagerDutyTeamID         map[string]string `split_words:"true"`
	PagerDutyScheduleID     map[string]string `split_words:"true"`
	SlackAuthorizedChannels string            `split_words:"true"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("pdbot", &env); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(env.PagerDutyScheduleID["LTS"])
	fmt.Println(env.PagerDutyTeamID["LTS"])

	pdclient := pagerduty.NewClient(env.PagerDutyAuthToken)
	conn := client.NewAPIClient(pdclient)

	bot := slacker.NewClient(env.SlackAuthToken)

	bot.Init(func() {
		log.Println("Connected!")
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	bot.DefaultCommand(func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Say what? try `help` command")
	})

	bot.DefaultEvent(func(event interface{}) {
		fmt.Println(event)
	})

	authorizedChannels := []string{env.SlackAuthorizedChannels}

	oncallDuty := &slacker.CommandDefinition{
		Description: "PagerDuty today oncall user",
		Example:     "oncall today",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			onCall := &oncall.AdminOnDutyList{OnCall: conn}
			schedule := &schedule.NewSchedule{Schedule: conn}
			today := time.Now()
			for _, pdschedule := range schedule.GetAll() {
				onCall.UsersOnCallOptions(today.String(), today.String(), pdschedule.APIObject.ID)
				if err := onCall.UsersOnCall(today, today); err != nil {
					response.ReportError(errors.New("something went wrong during processing bot command"))
				}
				oncall := onCall.PrintTodayDuty(pdschedule.Name)
				response.Typing()
				time.Sleep(time.Second)
				response.Reply(oncall)
			}
		},
	}

	oncallMonth := &slacker.CommandDefinition{
		Description: "PagerDuty oncall current month summary with profit",
		Example:     "oncall month lts",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			onCall := &oncall.AdminOnDutyList{OnCall: conn}
			pdschedule := request.StringParam("pdschedule", "lts")
			envPDSchedule, ok := env.PagerDutyScheduleID[strings.ToUpper(pdschedule)]
			if ok {
				pdschedule = envPDSchedule
			}
			today := time.Now()
			start := extensions.BeginningOfMonth(today)
			end := extensions.EndOfMonth(today)
			onCall.UsersOnCallOptions(start.String(), end.String(), pdschedule)
			if err := onCall.UsersOnCall(start, end); err != nil {
				response.ReportError(errors.New("something went wrong during processing bot command"))
			}
			oncall := onCall.PrintDutySummary(true)
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(oncall)
		},
	}

	incidentListTeam := &slacker.CommandDefinition{
		Description: "PagerDuty list of triggered and acknowledged incident for specific team",
		Example:     "incident list lts",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			incident := incident.Incidents{Incident: conn}
			pdTeam := request.StringParam("pdteam", "lts")
			envPDTeam, ok := env.PagerDutyTeamID[strings.ToUpper(pdTeam)]
			if ok {
				pdTeam = envPDTeam
			}
			pdTeamList := []string{}
			pdTeamList = append(pdTeamList, pdTeam)
			incidentOutls, err := incident.GetTeam(pdTeamList)
			if err != nil {
				response.ReportError(errors.New("something went wrong during processing bot command"))
			}
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(incidentOutls)
		},
	}

	incidentListTeamDuty := &slacker.CommandDefinition{
		Description: "PagerDuty list all incident incident for specific team and since defined hours",
		Example:     "incident duty lts 24h",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			incident := incident.Incidents{Incident: conn}
			toHour := request.StringParam("pdhour", "24h")
			pdTeam := request.StringParam("pdteam", "lts")
			envPDTeam, ok := env.PagerDutyTeamID[strings.ToUpper(pdTeam)]
			if ok {
				pdTeam = envPDTeam
			}
			pdTeamList := []string{}
			pdTeamList = append(pdTeamList, pdTeam)
			incidentOutls, err := incident.GetTeamDuty(pdTeamList, toHour)
			if err != nil {
				response.ReportError(errors.New("something went wrong during processing bot command"))
			}
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(incidentOutls)
		},
	}

	serviceListTeam := &slacker.CommandDefinition{
		Description: "PagerDuty list of services assigned to specific team",
		Example:     "service list lts",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			service := &services.Services{Service: conn}
			pdTeam := request.StringParam("pdteam", "lts")
			envPDTeam, ok := env.PagerDutyTeamID[strings.ToUpper(pdTeam)]
			if ok {
				pdTeam = envPDTeam
			}
			pdTeamList := []string{}
			pdTeamList = append(pdTeamList, pdTeam)
			serviceOutls, err := service.GetTeam(pdTeamList)
			if err != nil {
				response.ReportError(errors.New("something went wrong during processing bot command"))
			}
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(serviceOutls)
		},
	}

	maintenanceListTeam := &slacker.CommandDefinition{
		Description: "PagerDuty list of service maintenance to specific team",
		Example:     "maintenace list lts",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			maintenance := &maintenance.Maintenances{Maintenance: conn}
			pdTeam := request.StringParam("pdteam", "lts")
			envPDTeam, ok := env.PagerDutyTeamID[strings.ToUpper(pdTeam)]
			if ok {
				pdTeam = envPDTeam
			}
			pdTeamList := []string{}
			pdTeamList = append(pdTeamList, pdTeam)
			maintenanceOutls, err := maintenance.Get(pdTeamList)
			if err != nil {
				response.ReportError(errors.New("something went wrong during processing bot command"))
			}
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(maintenanceOutls)
		},
	}

	maintenanceCreateTeam := &slacker.CommandDefinition{
		Description: "PagerDuty create maintenance window for specific service from current time + given duration",
		Example:     "maintenace create PDSERVICEID 4h",
		AuthorizationFunc: func(botCtx slacker.BotContext, request slacker.Request) bool {
			return contains(authorizedChannels, botCtx.Event().Channel)
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			maintenance := &maintenance.Maintenances{Maintenance: conn}
			pdServiceID := request.StringParam("pdservice", "PDSERVICEID")
			toHour := request.StringParam("pdhour", "4h")
			maintenanceCreateOutls, err := maintenance.Create(pdServiceID, toHour)
			if err != nil {
				response.ReportError(errors.New("something went wrong during processing bot command"))
			}
			response.Reply(maintenanceCreateOutls)
		},
	}

	scheduleList := &slacker.CommandDefinition{
		Description: "PagerDuty schedule list",
		Example:     "schedule list",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			schedule := &schedule.NewSchedule{Schedule: conn}
			scheduleOutls := schedule.PrintSchedules()
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(scheduleOutls)
		},
	}

	teamList := &slacker.CommandDefinition{
		Description: "PagerDuty team list",
		Example:     "team list",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			team := &team.Team{Team: conn}
			teamOutls := team.PrintTeams()
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(teamOutls)
		},
	}

	bot.Command("oncall today", oncallDuty)
	bot.Command("oncall month <pdschedule>", oncallMonth)
	bot.Command("incident list <pdteam>", incidentListTeam)
	bot.Command("incident duty <pdteam> <pdhour>", incidentListTeamDuty)
	bot.Command("schedule list", scheduleList)
	bot.Command("team list", teamList)
	bot.Command("service list <pdteam>", serviceListTeam)
	bot.Command("maintenance list <pdteam>", maintenanceListTeam)
	bot.Command("maintenance create <pdservice> <pdhour>", maintenanceCreateTeam)

	authorizedDefinition := &slacker.CommandDefinition{
		Description: "Very secret stuff",
		AuthorizationFunc: func(botCtx slacker.BotContext, request slacker.Request) bool {
			return contains(authorizedChannels, botCtx.Event().Channel)
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("You are authorized!")
		},
	}

	bot.Command("secret", authorizedDefinition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func contains(list []string, element string) bool {
	for _, value := range list {
		if value == element {
			return true
		}
	}
	return false
}
