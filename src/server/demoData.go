package main

var demoUser = &User{
	Username: "demo",
	Secret:   []byte("demo"),
}

var demoEvents = []*Event{demoEvent1, demoEvent2, demoEvent3, demoEvent4, demoEvent5, demoEvent6, demoEvent7, demoEvent8, demoEvent9, demoEvent10}

var demoEvent1 = &Event{
	ID:          1,
	Name:        "Create Demo Events",
	StartTime:   "2018-02-28 19:00",
	EndTime:     "2018-02-28 20:00",
	Description: "Create demo events using time.Now() and time.Add() to create events that start from the time the app is run for the first time.",
	CreatedBy:   "demo",
}

var demoEvent2 = &Event{
	ID:          2,
	Name:        "Create db stmts create, read, update, delete.",
	StartTime:   "2018-02-28 22:00",
	EndTime:     "2018-02-28 23:50",
	Description: "Once demo events are added to the db test CRUD functionality with db stmts.",
	CreatedBy:   "demo",
}

var demoEvent3 = &Event{
	ID:          3,
	Name:        "Install docker-compose.",
	StartTime:   "2018-03-01 12:00",
	EndTime:     "2018-03-01 13:00",
	Description: "Add docker-compose to start mysql and app with one command.",
	CreatedBy:   "demo",
}

var demoEvent4 = &Event{
	ID:          4,
	Name:        "Create GRPC app.",
	StartTime:   "2018-03-02 12:00",
	EndTime:     "2018-03-02 01:00",
	Description: "GRPC app will have front end javascript.",
	CreatedBy:   "demo",
}

var demoEvent5 = &Event{
	ID:          5,
	Name:        "Create Docker/Postgres/Django app.",
	StartTime:   "2018-03-06 05:00",
	EndTime:     "2018-03-06 08:00",
	Description: "App should use Nginx and UWSGI.",
	CreatedBy:   "demo",
}

var demoEvent6 = &Event{
	ID:          6,
	Name:        "Create personal website.",
	StartTime:   "2018-03-10 15:00",
	EndTime:     "2018-03-10 18:00",
	Description: "Publish website.",
	CreatedBy:   "demo",
}

var demoEvent7 = &Event{
	ID:          7,
	Name:        "Gather portfolio.",
	StartTime:   "2018-03-11 13:00",
	EndTime:     "2018-03-11 19:00",
	Description: "Collect links to all work and publish on social media.",
	CreatedBy:   "demo",
}

var demoEvent8 = &Event{
	ID:          8,
	Name:        "Rewrite Resume.",
	StartTime:   "2018-03-11 02:00",
	EndTime:     "2018-03-11 06:00",
	Description: "Update resume to contain all of the latest work and skills.",
	CreatedBy:   "demo",
}

var demoEvent9 = &Event{
	ID:          9,
	Name:        "Submit Portfolio and Resume",
	StartTime:   "2018-02-12 04:00",
	EndTime:     "2018-02-12 12:00",
	Description: "Apply for attractive positions using newly finished portfolio.",
	CreatedBy:   "demo",
}

var demoEvent10 = &Event{
	ID:          10,
	Name:        "Do what I'm passionate about.",
	StartTime:   "2018-04-01 13:00",
	EndTime:     "2018-04-01 17:00",
	Description: "Programming, coding, configuring, learning, expanding, growing.",
	CreatedBy:   "demo",
}
