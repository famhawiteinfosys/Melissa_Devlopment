FROM golang

WORKDIR /go/src/github.com/famhawiteinfosys/Melissa

COPY . .

ENTRYPOINT ["go", "run", "."]
