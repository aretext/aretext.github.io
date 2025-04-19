aretext.github.io
-----------------

Website for the [aretext](https://github.com/aretext/aretext) editor.

Setup
=====

1. Install [uv](https://docs.astral.sh/uv/getting-started/installation/)
2. `uv tool install mkdocs --with mkdocs-material`


Building the site
=================

1. `git submodule update --init` to check out the aretext repo (so we can pull in the documentation).
2. `make build`

To serve the site locally, run `make server`.

Updating the aretext documentation
==================================

```
cd ./aretext
git fetch
git checkout $COMMIT_SHA
cd ..
git add aretext
git commit -m "Updated aretext"
```
