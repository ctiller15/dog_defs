# DogDefs
- A project to create and share definitions for any and all dog-related terminology.

## Running the container

```bash
# macos m1+
docker buildx build --platform linux/amd64 -t dog-defs .
docker run -p 8080:8080 dog-defs

docker build -t dog-defs .
docker run -p 8080:8080 dog-defs
```