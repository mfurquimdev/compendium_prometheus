FROM alpine

ENV OLA_MUNDO "Olá do valor padrão"

ADD printvar.sh /
RUN chmod a+x /printvar.sh

ENTRYPOINT ["/printvar.sh"]
