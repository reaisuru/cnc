# cnc

Leaked command and control server for a DDoS botnet.
It supports API, daily attacks, image rendering, termfx, etc..
Code is garbage.

## How to run?

 - Clone this repository.
 - Create a new PostgreSQL user (you'll find tutorials for this on Google)
 - Create a database named `botnet`.
 - Configure database in [database.go](https://github.com/reaisuru/cnc/blob/main/internal/database/database.go)
 - Execute the binary. `go run cmd/main.go`

### Requirements
 - [Go](https://go.dev) 1.22.4 or above.
 - [PostgreSQL](https://www.postgresql.org/)

##  TODO

- [x] Captcha
- Users
  - [x] Add user
  - [x] Remove user
  - [x] Edit user
  - [x] Add group to user
  - [x] Remove group from user
  - [x] Sessions view
  - [x] Kick sessions
  - [x] Password changíng
- Floods
  - [x] Attack logs
  - [x] Broadcast stats
  - [x] Target blacklist
  - [x] User cooldown for attacks
  - [x] Enable/disable attacks
  - [x] Daily attack limit
  - [x] Daily attack log removal
  - [x] Global cooldown for attacks
  - [x] Expiry check
  - [x] Duration check
  - [x] Logs command
  - [x] Kill all attacks
