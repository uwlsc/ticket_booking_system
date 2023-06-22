FROM mysql:8.0

COPY ./deployments/docker/custom.cnf /etc/mysql/conf.d/custom.cnf