
    用来做备份的一个工具，目前支持localfile备份 ssh远程备份 pgsql的备份，并能打包成一个压缩文件，提供HTTP下载
只需要一个配置文件就可以在本地，远程，或者将pqgsql 的数据全量备份下来，非常的方便，可以当成CMD单次备份使用，也可以当成一个服务进行定时备使用，后面将支持MYSQL等数据的备份。

# 设置配置文件
```
        #备份的文件夹
        backupdir=D:/backup
        #是否以服务方式启动，1 为以服务启动，0为以CMD运行
        service = 1
        #service = 1 才起效 m:分钟 s:秒
        sleeptime = 5/m
        #保存多少天,超过时间的删除
        savedays = 2


        ###备份本地文件
        [local]
        mode = file
        dir = D:/work/test/edge/edgelog

        ###备份SSH远程文件
        [edgelog]
        mode = ssh
        host = 127.0.0.1
        port = 22
        user = test
        passwd = test
        dir = /var/log/data

        ##备份Pgsql文件
        [pgsql]
        mode = pgsql
        dsn="host=127.0.0.1 user=postgres password=123456 dbname=test port=30388  sslmode=disable TimeZone=Asia/Shanghai"
        schema=hvac

```
#启动效果
```
DownLoading D:/work/test/edge-services/edgelog [##################################################] 100% 63/63 Done!
:: Begin Execute Task 【local2】::
DownLoading D:/work/test/edge-services/edgexport [##################################################] 100% 76/76 Done!
:: Begin Execute Task 【edgelog】::
DownLoad /edge-services/edgelog/data:[##################################################] 100% 1/1 Done!
:: Begin Compress BackupDataDict 【D:/backup/tmp000696D3ACD7BA78】::
################## Compress Backup Data Is Success! ################
```
##并可以提供http下载，访问，http://{ip}:17894,列表显示打包文件,并可以进行下载。
```
-rw-rw-rw-   2021-10-25@BackupFile30C1DD94D2C6C2B8.tar.gz  10843377 byte  2021-10-25 10:16:08.4752685 +0800 CST
-rw-rw-rw-   2021-10-25@BackupFile617A4B500AC5728D.tar.gz  10843360 byte  2021-10-25 10:05:40.2144843 +0800 CST
-rw-rw-rw-   2021-10-25@BackupFileCF3D28798CF3C0F2.tar.gz  10843416 byte  2021-10-25 10:10:54.1849612 +0800 CST

```
