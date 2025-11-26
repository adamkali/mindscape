# Frontend Development Dockerfile
FROM node:22-alpine

# Install pnpm
RUN npm install -g pnpm

# Set working directory
WORKDIR /app/web

# Copy package files
COPY web/package.json web/pnpm-lock.yaml ./

# Install dependencies
RUN CI=true pnpm install

# Copy the rest of the web directory
COPY web/ ./

# Expose the dev server port
EXPOSE 5173

# Start the development server
CMD ["pnpm", "dev", "--host", "0.0.0.0"]
