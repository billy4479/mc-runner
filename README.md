# MC Server Runner

Each docker container runs just one server.

## Frontend

Served on `/`.

Svelte static page.

## Backend

Backend on `/api`, Go with Echo.

Legend:
- (R): restricted, for automation, not for users
- (U): authenticated users
- (A): admin only

### Hooks

On `/api/hook` (R).

- `/api/hooks/mc`: get player and server lifetime events from the companion mod, 
stores into `players`.

- `/api/hooks/tg`: receive updates from Telegram from the companion bot
    - commands:
        - `/status` -> returns list of players
        - `/link <link_tg_token>` -> links a username to account
    - send msg on
        - server open (by who)
        - server close (how much time has it been open)
        - player joins or leaves

### Authentication

On `/api/auth`, just a cookie (lasts forever, until removed from db).

No password, ever.

- `/api/auth/register?t=<register_token>`: the token is a random string.
Registration flow as follows:
    - Admin invites user with `/api/admin/invite`
    - A `registration_token` gets added to `token_store` table
    - User clicks the link and chooses username (no password)
    - Server removes `registration_token` from `token_store`, adds user to `users` and the
    refresh cookie to `token_store`
    - Server returns the cookie

- `/api/auth/add-device` (U): generate a `login_token` to add to `token_store` and returns
it to user to be used in `/api/auth/login`. Expires after 5 minutes.

- `/api/auth/login?t=<login_token>`: removes `login_token` from `token_store`, 
generates `refresh_token` and adds it to `token_store`, return `refresh_token`

- `/api/auth/logout`: unset the `refresh_token`, removes it from `token_store`

- `/api/auth/link-tg` (U): generates `link_tg_token`

- `/api/auth/claim-mc` (U): sets `mc_name` to the provided one TODO: idk how to verify authenticity.

### Server Management

On `/api/management`.
Server closes only after a timeout.

- `/api/management/status` (U): returns if on or off and a list of players
- `/api/management/start` (U): starts the server if off
- `/api/management/logs` (U): websocket for logs

## Database

### users

- id
- name
- mc_name (nullable)
- tg_id (nullable)

### players

- id
- in_game_name
- is_online
- is_bot (use [this](https://carpet.tis.world/docs/rules#fakeplayernameprefix))

### token_store

- token
- expires
- type (enum of `refresh_token`, `login_token`, `register_token`, `link_tg_token`)
- user_id
