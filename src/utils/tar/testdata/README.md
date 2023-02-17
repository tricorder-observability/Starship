# testdata

#### test.tar.gz
Is is a compressed file, The directory structure is
```
➜  testdata git:(spec/linux-headers) ✗ tar -zxvf test.tar.gz 
./hello.txt
```

#### wrong_file_format.tar.gz
it's fake tar.gz file, In fact is a text file
```
➜  testdata git:(spec/linux-headers) ✗ cat wrong_file_format.tar.gz 
222
```
