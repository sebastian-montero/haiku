# Use an official Node.js image as the base image
FROM node:18-alpine

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application files
COPY . .

# Build the app for production
RUN npm run build

# Install `serve` to serve the built application
RUN npm install -g serve

# Expose the port that the app will run on
EXPOSE 3000

# Command to run the app
CMD ["npm", "run", "start"]
