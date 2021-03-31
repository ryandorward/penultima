FROM node:latest

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY package.json .
COPY Gruntfile.js .

RUN npm install

CMD ["npm","start"]
