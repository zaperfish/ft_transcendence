FROM node:24-alpine
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build
EXPOSE 7771
ENV PORT=7771
ENV NODE_ENV=production
CMD ["npm", "start"]
