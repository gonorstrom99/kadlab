FROM archlinux:latest AS builder 
WORKDIR /up 
RUN pacman -Syu --noconfirm && pacman -S go gcc --noconfirm 
COPY . ./
RUN ./compile.sh

FROM alpine:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab
WORKDIR /app

# TODO: kör funktion för korrekt tilldelande av port

#fil måste vara kompilerad innan man kör
# görs med
# ./compile.sh
# om error, kör först
# sudo chmod +x ./compile.sh


COPY --from=builder /up/d7024e ./

# sudo docker logs --follow [namn på nod]
# skriver vi saker i nodens log bör detta isf skrivas ut i terminalen

ENTRYPOINT ./d7024e