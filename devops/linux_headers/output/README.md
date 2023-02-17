# Output

Linux header archives.


## timeconst
It is generated from the following command:
```
wget https://raw.githubusercontent.com/torvalds/linux/35728b8209ee7d25b6241a56304ee926469bd154/kernel/time/timeconst.bc
echo 100 | bc timeconst.bc > timeconst_100.h
echo 250 | bc timeconst.bc > timeconst_250.h
echo 300 | bc timeconst.bc > timeconst_300.h 
echo 1000 | bc timeconst.bc > timeconst_1000.h 
```