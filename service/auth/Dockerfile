# Base image
FROM node:20-alpine 

# App directory
WORKDIR /app 

# Dependency files
COPY package*.json ./
RUN npm install

COPY . .

# Expose port
EXPOSE 3000

# Start app
CMD ["node", "index.js"]