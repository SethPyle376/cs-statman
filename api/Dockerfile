FROM node:14.10.1-alpine3.10
COPY . /app
WORKDIR /app
RUN npm install
ENTRYPOINT ["npm", "run", "start"]