FROM ubuntu

WORKDIR /projects

ADD rapper_be /projects
ADD config.ini /projects

EXPOSE 8080

CMD ["./rapper_be"]