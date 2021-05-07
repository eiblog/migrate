# migrate

迁移 eiblog 数据：

* 博客版本相同，迁移存储后端，如mongodb -> sqlite
* 博客版本升级，迁移数据格式，如v1 -> v2

配置请参考 app.yml:
```
from:
  version: v1
  driver: mongodb
  source: mongodb://127.0.0.1:27018
to:
  version: v2
  driver: mongodb
  source: mongodb://127.0.0.1:27017
```
