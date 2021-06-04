HTMLhouse
=========

Publish HTML quickly. 

[![HTMLhouse screenshot](https://lh3.googleusercontent.com/-iy91ChUsBOx0h0-AsXo-uSAwZOk-L2QSbLxxRUoTO8mL1ArJH8qiWfxFZ8OhOUnC9B8tzGHQA=s640-h400-e365-rw)](https://html.house)

HTMLhouse uses [ACE editor](https://ace.c9.io/) for modifying HTML and shows a live preview of what you've created in an `iframe` alongside the source code.

No user signup is required -- authorization to modify an individual published page is saved on the creator's device in local storage as an ID and [JWT](http://jwt.io/) generated by the server.

It is also available as a [Chrome extension](https://chrome.google.com/webstore/detail/htmlhouse/aljfacibkadhobknaalpfbclcpoepopi) and browsable on [Android](https://play.google.com/store/apps/details?id=com.abunchtell.htmlhouse).

[![#writeas on freenode](https://img.shields.io/badge/freenode-%23writeas-blue.svg)](http://webchat.freenode.net/?channels=writeas) [![Public Slack discussion](http://slack.write.as/badge.svg)](http://slack.write.as/)

## Development

### Requirements

* Go
* Node.js
* MySQL

### Setup

1. Clone the repo
2. Run `go get -d` to get necessary dependencies
3. Run `make install` to install LESS compiler locally and generate the CSS files
4. Run the queries in `init.sql` to set up the database
5. _Optional_. Run `./keys.sh prod` to create a new keypair

### Running the server

* Run `go run cmd/htmlhouse/main.go` in the top level directory, optionally by creating a simple run script
```bash
#!/bin/bash

DB_USER=dbuser DB_PASSWORD=pass DB_DB=htmlhouse PRIVATE_KEY=keys/dev PUBLIC_KEY=keys/dev.pub go run main/main.go
```
* Open your browser to http://localhost:8080

#### Environment Variables

| Variable | What it is | Default value |
| -------- | ---------- | ------------- |
| `DB_USER` | Database user | None. **Required** |
| `DB_PASSWORD` | Database password | None. **Required** |
| `DB_DB` | Database name | None. **Required** |
| `DB_HOST` | Database host | `localhost` |
| `PRIVATE_KEY` | Generated private key | None. **Required** |
| `PUBLIC_KEY` | Generated public key | None. **Required** |
| `PORT` | Port to run app on | `8080` |
| `STATIC_DIR` | Relative dir where static files are stored | `static` |
| `ALLOW_PUBLISH` | Allow users to publish posts | true |
| `AUTO_APPROVE` | Automatically approves public posts | false |
| `PREVIEWS_HOST` | Fully-qualified URL (without trailing slash) of screenshot server | None. |
| `ADMIN_PASS` | Password to perform admin functions via API | `uhoh` |
| `BROWSE_ITEMS` | Number of items to show on Browse page | 10 |
| `BLACKLIST_TERMS` | Comma-separated list of terms to prevent a post from being made public | None. |
| `TWITTER_KEY` | Twitter consumer key | `notreal` |
| `TWITTER_SECRET` | Twitter consumer secret | `notreal` |
| `TWITTER_TOKEN` | Twitter access token of the posting Twitter account | `notreal` |
| `TWITTER_TOKEN_SECRET` | Twitter access token secret of the posting Twitter account | `notreal` |
| `WF_MODE` | Run CSShorse, not HTMLhouse — for customizing WriteFreely blogs | `false` |

### Notes

**Changing CSS**. Run `make` after all changes to update the stylesheets.

**When you don't need to reload the app**. When you make changes to any files in `static/` you can simply refresh the resource without restarting the app.

**When to reload the app**. If you change any of the templates in `templates/` or any `.go` file, you'll need to re-run the app.
