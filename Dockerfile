FROM debian:12.5-slim


#TODO add build and version info

RUN mkdir -p /opt/app/bundles

COPY bin/pr-bot /opt/app/
COPY config/* /opt/app/

RUN useradd -ms /bin/bash pr-bot
RUN chown -R pr-bot /opt/app
USER pr-bot

WORKDIR /opt/app
EXPOSE 9090

ENTRYPOINT ["/opt/app/pr-bot"]
CMD ["-config", "/opt/app/dev.yaml"]
