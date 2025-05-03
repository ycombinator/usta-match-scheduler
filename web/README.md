## Installation

```
npm install
```

## Running

1. Generate public assets.
```
npm run build
```

2. Start server (API proxy + public assets).
```
npm run start
```

Public assets are served at http://localhost:3000
API proxy is served at http://localhost:3000/api/, proxying to http://localhost:8000/