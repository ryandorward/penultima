### PART I ###

# pull official base image
FROM node:13.12.0-alpine AS builder

# take the following two lines out for production
RUN apk update && apk add bash 
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# set working directory
RUN mkdir /app
WORKDIR /app

#RUN mkdir -p /usr/src/app
#WORKDIR /usr/src/app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install app dependencies
COPY package.json ./
# COPY package-lock.json ./
RUN npm install

# add app
COPY . ./

#RUN npm run build

# start app
CMD ["npm", "start"] 
#CMD ["npm","run","build"]


# no, instead just build?
#CMD ["npm", "BUILD"] 
