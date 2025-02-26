# Obsidian Share

Self-hosted alternative to the paid [Obsidian Publish](https://obsidian.md/publish). Features:

- Admin page to manage and pull files
- Fast and memory efficient (executable is 12.5 MB!)

## Development

Simple golang application -> `go run .`

## Prod

### Build with Docker

```bash
docker build -t obsidian-share .
```

### Run with Docker

Create an `.env` file:

```
SSH_PRIVATE_KEY=
REPO_URL=
MONGODB_URI=
ADMIN_PASS=
PORT=
```

Then simply run the following command with Docker:

```bash
docker run --rm --env-file=.env --mount type=bind,source=/home/mikey/.ssh/id_ed25519,target=/id_ed25519 --mount type=bind,source=/home/mikey/.ssh/known_hosts,target=/.ssh/known_hosts -p 5000:8080 obsidian-share:latest
```

Replace `/home/mikey` with your own home directory path. My secret key is called `id_ed25519`; change this to match your github secret key. Make sure to specify this in your `.env` file too. For the ports, remember that it's `<HOST_PORT>:<CONTAINER_PORT>`. `CONTAINER_PORT` should match what you define in your env file.

## Usage

The `/admin` route will have the admin page, which can be logged in using the `ADMIN_PASS` password set in the environment. Once you have added files via the `Add New...` button, you must click the `Pull and compile files` button, which will actually do a git clone/fetch and compile the markdown into HTML. Then, simply go to your defined slug and open the file!
