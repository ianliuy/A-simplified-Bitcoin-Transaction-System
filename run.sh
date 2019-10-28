rm *.db
go run *.go
go run *.go printChain
go run *.go getBalance --address genesis
go run *.go send genesis aaa 10 minerbbb "1st transaction"