FROM node:8-alpine

WORKDIR /slackkubewatch
COPY package*.json ./
COPY slackkubewatch.js ./
RUN npm install

ENTRYPOINT [ "node", "slackkubewatch.js" ]

