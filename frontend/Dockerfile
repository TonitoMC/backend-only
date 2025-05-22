
# Use the official Nginx base image (lightweight Alpine variant)
FROM nginx:alpine

# Copy all files from current directory to Nginx's html directory
COPY . /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# No CMD is needed since the base image already starts Nginx
