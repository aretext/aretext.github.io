aretext.github.io
-----------------

Website for the [aretext](github.com/aretext/aretext) editor.

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
