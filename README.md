#### 功能
- 通过beat框架(v5.6.8)实现发送的所有告警信息录入elastic或其他beat支持的output端
- 通过发送给NagiosAPI服务(nagios被动检测方式)实现业务告警通过nagios实现告警分发，告警策略配置等

#### alertbeate的API

- 业务API
```
curl -XPOST http://127.0.0.1:8989/v1/t8t -d'{
	"labels":{	   
		"alertname":"10 test content",
		"env":"idc",
		"from":"kapacitor",	
		"proj":"test-abc-abc",	
		"type":"java"
	},
	"annotations":{
		"count":"2",
		"domain":"test-abc-abc",
		"duration":"0",
		"host":"10.10.10.11",
		"interface":"com.to8to.rpc.utils.LogUtils.logByLevel",
		"level":"INFO",
		"time":"4:00PM 03/26"
	}
}'
```

- 另一种业务格式API
```
curl http://127.0.0.1:8989/v1/basic -d'{
	"alarmid":"2001",
	"content":"test cotent"
}'
```

- 测试查询
```
curl http://127.0.0.1:8989/v1/
```
