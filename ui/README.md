# Starship Api-server management Web UI

To run management Web UI locally:
```bash
# First install the build toolchain
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
nvm install node
npm install yarn

yarn install
yarn run dev
# Open localhost:8000, if port 8000 is already taken，8001 will be used, and so on
```

To change the endpoint of management Web UI's backend server on API server,
update configurations in `config/proxy.ts`:

```
test: {
  // localhost:8000/api/** -> https://preview.pro.ant.design/api/**
  '/api/': {
    target: 'http://ec2-3-93-75-222.compute-1.amazonaws.com:8080',
    changeOrigin: true,
    pathRewrite: { '^': '' },
  },
},
```

To run nginx with local build:
```
yarn run build
sudo nginx -c  docker/nginx_test.conf
sudo systemctl restart nginx
```

```bash
# Install node dependencies
npm install

# Or using yarn
yarn

# To start local dev server
npm start

# Build the project
npm run build

# Check code style
npm run lint

# You can also use script to auto fix some lint error:
npm run lint:fix

# Test code
npm test
```

To build docker image
```
# Build docker image
yarn install
yarn run build
yarn run docker:image -- TAG=${tag}
```
