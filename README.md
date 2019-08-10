# Pocket Gopher!

> Pocket gophers, commonly referred to as just gophers, are burrowing rodents of the family Geomyidae.

**(Work In Progress)**

Pocket is a neat little web library to help you write cleaner HTTP request handlers!

It's inspired by how a few other languages provide handler APIs, mainly [Rocket](https://rocket.rs/) (hence the name!)

## Example

Here's a quick example:

```go
func handler(props struct {
    MethodGet
    ParamUserID string
    BodyJSON    J
}) error {
    // Do stuff with `props`!
    // simply return an error! no more writing to w!
}
```

## Why?

Go has a great ecosystem for HTTP middleware thanks to the standard Request/ResponseWriter interface. This pattern is
great, until it comes to writing your actual business logic handlers. The Request/ResponseWriter is just too low level
for doing quick things like getting query parameters or grabbing some JSON.

Sure, you can write some helpers for this but then if you have a few projects, you end up copying/importing and calling
this boilerplate all over the place. Doing it declaratively is much nicer!

Note: This package makes generous use of reflection and, while all efforts are made to cache unnecessary reflection
calls and perform most of at handler generation-time, this reflection _may_ impact extremely high-traffic web servers.
There are currently no benchmarks but this is something I plan to work on in future.

## Some Ideas

_A scratchpad for ideas for this library!_

- `func(string)`/`func() string` signatures for quick and easy text requests/responses.
- `func(T)`/`func(T)` signatures for quick JSON endpoints.
- nice default for `func(...) error` signatures
- Provide some configuration for the handler generator:
  - mapping from `error` types to statuses, for those `func(...) error` sigs.
  - wrap JSON responses in some default body
- generate index/sitemap/documentation based on signatures
