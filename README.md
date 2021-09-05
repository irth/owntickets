# OwnTickets
A personal ticketing system that's simple to use and deploy.

It is limited to a single admin user by default. Ticket creation can be either public, or limited by a password.

**NOT YET IN A WORKING STATE**

## Setup
1. Download the binary from the Releases page (TODO: put binaries there).
2. Make sure it is executable (`chmod +x owntickets`)
2. Hash the admin password using `./owntickets hash-password`
3. If you want the ticket creation page to be password protected, hash that password too
4. Create a config (see below)
5. Run `./owntickets serve` (or `./owntickets serve -c=config.json`)

## Configuration
There are two ways to configure OwnTickets: it can either take it's configuration from a JSON file, or from the environment. If a value is specified both in JSON and environment, the one from the environment is used.

JSON config:

```json
{
    "database": "<path to the sqlite3 database file>",
    "passwordHash": "<bcrypt admin password hash>",
    "requirePasswordForTicketCreation": true,
    "ticketCreationPasswordHash": "<bcrypt ticket creation password hash>"
}
```

environment:
```bash
OWNTICKETS_DATABASE="./sqlite3.db"
OWNTICKETS_PASSWORD_HASH="hash"
OWNTICKETS_REQUIRE_PASSWORD="yes/no" # require password for ticket creation
OWNTICKETS_TICKET_PASSWORD_HASH="hash"
```