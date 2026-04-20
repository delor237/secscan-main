FROM node:20-alpine
WORKDIR /app

# Bağımlılıkları yükle
COPY package.json package-lock.json* ./
RUN npm install

# Kaynak kodları kopyala ve build al
COPY . .
RUN npm run build

# Portu aç ve başlat
EXPOSE 3000
ENV PORT 3000

CMD ["npm", "start"]
