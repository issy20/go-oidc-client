ARG VARIANT="1.19-bullseye" 
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT} as dev

ENV PATH /go/bin:$PATH
RUN ["chmod", "+w", "/go/src"]