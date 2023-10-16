package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var url = "https://tickets.fc-zenit.ru/football/tickets/"

type Event struct {
	HomeTeam     string
	VisitingTeam string
	DateAndTime  string
	Championship string
	Place        string
}

func main() {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var htmlContent string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		log.Fatal(err)
	}

	var teams, championships, places, dates []*cdp.Node

	err = chromedp.Run(ctx,
		chromedp.Nodes(
			`.member-name--team.member-name--event-big`, &teams,
			chromedp.ByQueryAll,
		),
		chromedp.Nodes(
			".additional-block__item.additional-block__item--first", &championships,
			chromedp.ByQueryAll,
		),
		chromedp.Nodes(
			".additional-block__item.additional-block__item--with-dot", &places,
			chromedp.ByQueryAll,
		),
		chromedp.Nodes(
			".date-block--event-big", &dates, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Fatal(err)
	}

	var events []Event

	for i := 1; i < len(teams); i += 2 {
		events = append(events, Event{
			HomeTeam:     teams[i-1].Children[0].NodeValue,
			VisitingTeam: teams[i].Children[0].NodeValue,
			Championship: championships[i].Children[0].NodeValue,
			Place:        places[(i-1)/2].Children[0].NodeValue,
			DateAndTime:  dates[i].Children[0].NodeValue,
		})
	}

	for _, event := range events {
		fmt.Printf("Date: %s \nChampionship: %s \n%s - %s\nPlace: %s\n",
			event.DateAndTime,
			event.Championship,
			event.HomeTeam,
			event.VisitingTeam,
			event.Place,
		)
	}

}
