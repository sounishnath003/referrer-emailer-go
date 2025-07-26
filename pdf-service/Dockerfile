FROM node:23-slim

RUN apt-get update && apt-get install -y build-essential chromium \
    --no-install-recommends && rm -fr /var/lib/apt/lists/*

ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
ENV PORT=3001

WORKDIR /app

RUN mkdir -p /usr/share/fonts/truetype/
COPY ./fonts/Inter.ttf /usr/share/fonts/truetype/inter/
RUN fc-cache -f -v

COPY package.json package-lock.json* /app
RUN npm install --production

COPY . .
EXPOSE 3001

ENTRYPOINT ["node", "index.js"]