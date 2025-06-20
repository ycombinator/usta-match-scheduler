## Installation

```
npm install
```

## Running

1. Generate public assets whenever there are changes to files.
```
fswatch -o -r ./ | xargs -n1 -I% npm run build
```

2. (In a different window) Start server (API proxy + public assets).
```
npm run serve
```

Public assets are served at http://localhost:3000
API proxy is served at http://localhost:3000/api/, proxying to http://localhost:8000/