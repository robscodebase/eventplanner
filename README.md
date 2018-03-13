# Docker Golang Mysql Event Planner

Docker Golang Mysql Event Planner is created by Robert Reyna.

Unless otherwise noted, these source files are distributed under the
BSD-style license found in the LICENSE file.

### Clone Repo

You must have Docker and Docker-Compose installed on your machine.
* Clone repo.
  * `git clone https://github.com/robscodebase/eventplanner.git`


### Run Docker-Compose

* Navigate to /eventplanner
  * `sudo docker-compose up`

### View App

* Open browser and navigate to localhost:8081/

### Demo App Creditials

* Username: `demo`
* Password: `demo`

### Run Tests
* `git checkout test`
* `./runtest.sh`
* This will start a mysql instance and run tests on the server.

#### Non-standard libraries.
    * `github.com/gorilla/mux
    github.com/gorilla/handlers
    github.com/go-sql-driver/mysql
    golang.org/x/crypto/bcrypt
    github.com/nu7hatch/gouuid`
