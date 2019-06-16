# PIU

Pull images up to date.

Do not bother to recreate your services, PIU can help you
pull the latest images and update your containers.

## How it works

- Scan all the containers and find their images
- Watch the changes of the images
- When the digest changed
  - pull the latest image
  - recreate the container
  - keep everything except the image
