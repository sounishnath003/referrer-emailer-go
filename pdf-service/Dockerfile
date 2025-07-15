FROM node:22-slim

# Install necessary dependencies for Puppeteer/Chrome
RUN apt-get update && \
    apt-get install -y wget gnupg && \
    apt-get install -y \
        libxss1 \
        libasound2 \
        libnss3 \
        libatk-bridge2.0-0 \
        libgtk-3-0 \
        libx11-xcb1 \
        libxcomposite1 \
        libxdamage1 \
        libxrandr2 \
        libgbm1 \
    rm -rf /var/lib/apt/lists/*

# We don't need the standalone Chromium
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD true

WORKDIR /app
COPY . .

RUN npm install

# Install Chrome for Puppeteer to use
RUN npx puppeteer browsers install chrome

EXPOSE 3001

CMD ["node", "index.js"]