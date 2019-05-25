FROM aoepeople/scratch-go-env

ADD VERSION .
ADD rating /rating
ADD config /config
ADD static /static
ADD templates /templates
ADD sql /sql

ENTRYPOINT ["/rating"]
CMD ["serve"]
