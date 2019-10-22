FROM debian
COPY ./csrgo /csrgo
ENTRYPOINT /csrgo
