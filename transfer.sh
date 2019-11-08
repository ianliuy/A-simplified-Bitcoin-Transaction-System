rm *.db
go run *.go
go run *.go send aaa bbb 10 eee "aaa->bbb:10"
##go run *.go printChain
##bash checkBalance.sh
go run *.go send aaa ccc 20 eee "aaa->ccc:20"
#go run *.go printChain
#bash checkBalance.sh
go run *.go send ccc bbb 2 eee "ccc->bbb:20"
#bash checkBalance.sh
go run *.go send ccc bbb 3 eee "ccc->bbb:20"
#bash checkBalance.sh
go run *.go send ccc aaa 5 eee "ccc->aaa:20"
#bash checkBalance.sh
go run *.go send bbb ddd 14 eee "bbb->ddd:20"
#bash checkBalance.sh
bash checkBalance.sh
