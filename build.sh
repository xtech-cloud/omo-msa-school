export GO111MODULE=on
export GOSUMDB=off
export GOPROXY=https://mirrors.aliyun.com/goproxy/
go install omo.msa.school
mkdir _build
mkdir _build/bin

cp -rf /root/go/bin/omo.msa.school _build/bin/
cp -rf conf _build/
cd _build
tar -zcf msa.school.tar.gz ./*
mv msa.school.tar.gz ../
cd ../
rm -rf _build
