# woodblock

Rotate your macOS desktop wallpaper through Unsplash photos from accounts you trust.

By default it shows Japanese woodblock prints from museum and library accounts
(Art Institute of Chicago, Cleveland Museum of Art, The Met, and others). Any
Unsplash search and any account allowlist work.

## Install

```sh
mise use -g packslip:github.com/nateberkopec/woodblock
```

Create an Unsplash app at <https://unsplash.com/oauth/applications> and give
woodblock its access key, either in the config file or as `UNSPLASH_ACCESS_KEY`.

## Usage

```
woodblock                         show the current wallpaper's title, credit, and links
woodblock open [source|unsplash]  open the source page (e.g. the museum), or the Unsplash page
woodblock next                    set a new wallpaper
woodblock back                    go back to the previous wallpaper
```

Run `woodblock next` on a schedule (launchd, cron) to rotate daily.

Each pick searches a random query, jumps to a random page of results, and
picks a random photo from an allowed account that is not in recent history.

## Configuration

`~/.config/woodblock/config.toml`, every key optional:

```toml
access_key = "..."
queries = ["woodblock print", "japanese print"]
accounts = ["artchicago", "clevelandart", "metropolitanmuseum"]
orientation = "landscape"
history_size = 100
```

An empty `accounts` list allows every account.

History lives in `~/.local/state/woodblock/history.json` and images in
`~/Library/Caches/woodblock`.
