
    用来做备份的一个工具，目前支持localfile备份 ssh远程备份 pgsql的备份，并能打包成一个压缩文件，提供HTTP下载
只需要一个配置文件就可以在本地，远程，或者将pqgsql 的数据全量备份下来，非常的方便，可以当成CMD单次备份使用，也可以当成一个服务进行定时备使用，后面将支持MYSQL等数据的备份。


```
        #备份的文件夹
        backupdir=D:/backup
        #是否以服务方式启动
        service = 1
        #service = 1 才起效 m:分钟 s:秒
        sleeptime = 5/m
        #保存多少天,超过时间的删除
        savedays = 2

        [local]
        mode = file
        dir = D:/work/test/edge/edgelog


        [edgelog]
        mode = ssh
        host = 127.0.0.1
        port = 22
        user = test
        passwd = test
        dir = /var/log/data

        [pgsql]
        mode = pgsql
        dsn="host=127.0.0.1 user=postgres password=123456 dbname=test port=30388  sslmode=disable TimeZone=Asia/Shanghai"
        schema=hvac

```
