FROM debian
COPY ./csrgo /csrgo
EXPOSE 80
ENTRYPOINT /csrgo
