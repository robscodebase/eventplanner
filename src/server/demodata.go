// Copyright (c) 2018 Robert Reyna. All rights reserved.
// License BSD 3-Clause https://github.com/robscodebase/eventplanner/blob/master/LICENSE.md
// demodata.go contains the data for populating the db on startup.
package main

var demoUser = &User{
	Username: "demo@demo.com",
	Secret:   []byte("demo"),
}

var demoEvents = []*Event{demoEvent1, demoEvent2, demoEvent3, demoEvent4, demoEvent5, demoEvent6, demoEvent7, demoEvent8, demoEvent9, demoEvent10}

var demoEvent1 = &Event{
	ID:          1,
	Name:        "Create Demo Events",
	StartTime:   "2018-02-28 19:00",
	EndTime:     "2018-02-28 20:00",
	Description: "Create demo events using time and add to create events that start from the time the app is run for the first time.",
	UserID:   1,
}

var demoEvent2 = &Event{
	ID:          2,
	Name:        "Create db stmts.",
	StartTime:   "2018-02-28 22:00",
	EndTime:     "2018-02-28 23:50",
	Description: "Once demo events are added to the db test CRUD functionality with db stmts. Add view-events.html and add-events.html.",
	UserID:   1,
}

var demoEvent3 = &Event{
	ID:          3,
	Name:        "Install docker-compose.",
	StartTime:   "2018-03-01 12:00",
	EndTime:     "2018-03-01 13:00",
	Description: "Add docker-compose to start mysql and app with one command. Include statements for creating docker network for mysql conn.",
	UserID:   1,
}

var demoEvent4 = &Event{
	ID:          4,
	Name:        "Create GRPC app.",
	StartTime:   "2018-03-02 12:00",
	EndTime:     "2018-03-02 01:00",
	Description: "GRPC app will have front end javascript. Using the mapbox api, canvas draw, and webpack to produce a usable front end.",
	UserID:   1,
}

var demoEvent5 = &Event{
	ID:          5,
	Name:        "Create Postgres/Django.",
	StartTime:   "2018-03-06 05:00",
	EndTime:     "2018-03-06 08:00",
	Description: "App should use Nginx and UWSGI. Postgres and Django container will be connected using docker network and docker-compose.",
	UserID:   1,
}

var demoEvent6 = &Event{
	ID:          6,
	Name:        "Create personal website.",
	StartTime:   "2018-03-10 15:00",
	EndTime:     "2018-03-10 18:00",
	Description: "Publish website. Include professional and personal information. Create links to github profile page and examples of work.",
	UserID:   1,
}

var demoEvent7 = &Event{
	ID:          7,
	Name:        "Gather portfolio.",
	StartTime:   "2018-03-11 13:00",
	EndTime:     "2018-03-11 19:00",
	Description: "Collect links to all work and publish on social media. Also, collect links to all social media and list them on sites.",
	UserID:   1,
}

var demoEvent8 = &Event{
	ID:          8,
	Name:        "Rewrite Resume.",
	StartTime:   "2018-03-11 02:00",
	EndTime:     "2018-03-11 06:00",
	Description: "Update resume to contain all of the latest work and skills. Update resume on Linked-In and other job posting sites.",
	UserID:   1,
}

var demoEvent9 = &Event{
	ID:          9,
	Name:        "Submit Resume",
	StartTime:   "2018-02-12 04:00",
	EndTime:     "2018-02-12 12:00",
	Description: "Apply for attractive positions using newly finished portfolio. Write targeted cover letters for each employer.",
	UserID:   1,
}

var demoEvent10 = &Event{
	ID:          10,
	Name:        "Do what I love.",
	StartTime:   "2018-04-01 13:00",
	EndTime:     "2018-04-01 17:00",
	Description: "Programming, coding, configuring, learning, expanding, growing. Practice daily, think of new projects, work hard.",
	UserID:   1,
}
