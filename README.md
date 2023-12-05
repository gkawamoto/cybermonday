# Cybermonday

Markdown server without complications.

## Overview

Simple golang + nginx server to host Markdown files without too much complication.
This project was originally based on [https://github.com/russross/blackfriday](https://github.com/russross/blackfriday) (hence the name), but for the sake of simplicity and compliance, it was rewritten to use [https://gitlab.com/golang-commonmark/markdown](https://gitlab.com/golang-commonmark/markdown).

## Quickstart

```bash
docker run --name cybermonday -v $(pwd):/usr/share/nginx/html:ro gkawamoto/cybermonday:stable
```

## Volume

This containers uses a single volume mountpoint to host all the files (markdown or not): `/usr/share/nginx/html`.

## Environment variables

| Variable | Description |
|:-|-:|
| `CYBERMONDAY_TITLE` | Page title in the upper left corner
| `CYBERMONDAY_BOOTSTRAP_REF` | Bootstrap CSS CDN path

## Contributing

Issues and merge-requests are welcome.

## License

MIT License

---
Copyright (c) 2023, Gustavo Kawamoto
