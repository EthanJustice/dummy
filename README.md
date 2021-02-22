# dummy

a dead-simple json mocking api

```typescript
fetch('/new', {
    method: 'post', // /new endpoint only accepts POST requests
    headers: {
        'Content-Type': 'application/json',
        'X-Route': 'route/to/thing' // this is the route to the new path, and must be set
    },
    body: JSON.stringify({
        test: 'example string'
    })
}).then((resp) => resp.json()).then(v => console.log(v)); // { "status": "success" }
```

## Roadmap

+ random input generation using `rsg`
+ stop codes (to remove endpoints)
+ more backends

## Development

1. Download the source
2. Run `go run main.go` to start the server
3. With /static as the working directory, run `yarn run build` or `npm run build` to bundle the frontend

## Usage

### Creating an Endpoint

To create a new route, make a POST request against `/new`. Make sure you set the `X-Route`
    header to the route you want to use (do not include a slash at the beginning).

### Getting JSON from an Endpoint

To get a response from an endpoint already created, make a GET request against
    `/whatever/path/you/set/as/the/X-Route`.

## Endpoints

```text
GET /
serves the default root page

POST /new
    X-Route | string | the route to use, do not add a leading slash (e.g. this/is/a/route rather than /this/is/a/route)
    body | JSON | JSON to serve at the new endpoint

GET /{any}
    body | JSON | the body that was set when /new was POSTed
```
