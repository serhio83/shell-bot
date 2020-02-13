FROM scratch

ENV PORT 9999
EXPOSE $PORT

COPY shell-bot /
CMD ["/shell-bot"]
