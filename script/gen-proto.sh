# execute to generate .pb.go and .pb.micro.go files
cd ..
cd proto
protoc --micro_out=./ --go_out=./ userinfo/userinfo.proto
cd .. && cd ..