#根据数据库表生成对应的go model文件

modeldir=./genModel

# 数据库配置
host=127.0.0.1
port=3306
dbname=big_market
username=root
passwd=Huangji1

# 获取数据库中所有表的名字
tables=$(mysql -h ${host} -P ${port} -u ${username} -p${passwd} -D ${dbname} -e "SHOW TABLES;" | awk 'NR > 1')

# 为每个表生成模型文件
for table in $tables
do
    echo "开始创建库：$dbname 的表：$table"
    tabledir="${modeldir}/${table}" # 为每个表创建一个目录
    mkdir -p "$tabledir"            # 确保目录存在，如果不存在则创建
    goctl model mysql datasource -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" -table="${table}" -dir="${tabledir}" -cache=true --style=goZero
done
