# Library

This repository serves as a library of community-made resources 
(blog posts, videos, presentations, ...). The content is indexed, versioned, 
and searchable from [`vale.sh/docs`](https://vale.sh/docs/).

## Submitting a resource

If you'd like to submit your own resource, open a PR that adds it to the `library.json`
file at the root of this repository.

## Searching

We support a query-string syntax powered by [Bleve](https://blevesearch.com/):

- [x] Faceted search: `date:>2021` or `author:jdkato` (`date`, `title`, `author`, `text`, `type`).
- [x] `+foo` (required have)
- [x] `-foo` (required to not have)
- [x] `foo~2` (fuzzy)
- [x] `"foo bar"` (phrase)
- [x] `foo bar` (OR)

See the [Bleve docs](https://blevesearch.com/docs/Query-String-Query/) for more information.