# brandhang-portfolio

Personal portfolio site for Brandon Hang — Platform Engineer, Sys Admin, Solutions Architect.
Live at [brandhang.net](https://brandhang.net).

Server-rendered Go, no JavaScript build step.

| Layer | Choice |
| --- | --- |
| HTTP | [Echo](https://echo.labstack.com/) v4 |
| Templates | [templ](https://templ.guide/) (compiled to Go) |
| Interactivity | [HTMX](https://htmx.org/) (vendored, no bundler) |
| Styling | [Tailwind CSS](https://tailwindcss.com/) v4 via the standalone CLI |
| Contact form | Discord bot (`bwmarrin/discordgo`) |
| Deploy | Docker image pushed to Docker Hub by GitHub Actions |

## Layout

```
cmd/api/main.go        entrypoint, graceful shutdown
cmd/web/               templ views + embedded static assets
internal/server/       HTTP server config and routes
internal/handler/      portfolio data loading, contact-form delivery
```

Portfolio entries live in `internal/handler/portfolio.json` and are embedded at
compile time — add a project by editing that file, no Go changes needed.

## Getting started

```bash
cp .env.example .env   # then fill in the Discord values
make build             # templ generate -> tailwind -> go build
make run
```

Then open <http://localhost:8080>.

## Environment

| Variable | Required | Notes |
| --- | --- | --- |
| `PORT` | no | Defaults to `8080`. |
| `APP_ENV` | no | `local` or `production`. |
| `DISCORD_BOT_TOKEN` | for contact form | Bot must be able to post in the channel. |
| `DISCORD_CHANNEL_ID` | for contact form | Where submissions are relayed. |

If the Discord values are missing the site still runs; the contact form just
reports an error to the visitor.

## Make targets

| Target | Does |
| --- | --- |
| `make build` | Generate templ, build CSS, compile the binary |
| `make run` | Run without building assets |
| `make test` | `go test ./... -v` |
| `make watch` | Live reload via [air](https://github.com/air-verse/air) |
| `make clean` | Remove the built binary |

## Gotcha: `output.css` is committed

The Dockerfile runs `templ generate` but **not** Tailwind, so production serves
whatever `cmd/web/assets/css/output.css` is committed to the repo. After
changing `input.css` or any class in a `.templ` file, run `make build` and commit
the regenerated CSS. CI fails the build if the committed CSS is stale.

## Deployment

Pushes to `master` run the `test` job (vet, tests, gofmt, tidy check, CSS drift
check); if it passes, the `docker` job builds and pushes
`:latest` and `:<commit-sha>` to Docker Hub. Pull requests run `test` only.

Required repo config: `DOCKERHUB_USERNAME` and `DOCKERIMAGE_NAME` variables, and
a `DOCKERHUB_TOKEN` secret.

```bash
docker build -t brandhang-portfolio .
docker run --rm -p 8080:8080 --env-file .env brandhang-portfolio
```

## License

MIT — see [LICENSE](LICENSE).
