# rabbitmq-bench



## 用法

### 声明exchange，queue，routingKey

如果不声明exchange，queue，routingKey，对于生产者也是可以工作的，使用`--declare`选项时，生产者和消费者都会生成默认的依赖项，详见命令行帮助信息

- 声明默认的exchage-生产者

  ```sh
  ./rabbitmq-bench producer --log=debug --declare -n 100 e-name q-name  00message00
  
  # partial output:
  INFO[0000] declaring exchange                            autoDelete=false durable=false internal=false name=e-name noWait=false type=direct    
  ```

- 声明默认的exchage，queue-消费者

```

```



- 一次声明exchange，queue，routingKey

```sh
$ ./rabbitmq-bench prepare all --log=info   e-name q-name
# output：
INFO[0000] declaring exchange                            autoDelete=false durable=false internal=false name=e-name noWait=false type=direct
INFO[0000] declaring queue                               autoDelete=false durable=false internal=false name=q-name noWait=false
INFO[0000] queue state                                   consumers=0 message=0 name=q-name
INFO[0000] declaring binding                             exchange=e-name noWait=false queue=q-name routingKey=q-name
INFO[0000] Exiting...

```

如果queue的名称不指定，那么与routingKey是一致的

- 
- 

### 生产者发送消息

- 作为生产者发送一条msg，显示详细日志，如果exchange不存在则创建默认的

```
./rabbitmq-bench producer --log=debug --declare  -n 1 e-name q-name 00message00
```

- 作为生产者发送大量msg(100000)，显示发送的统计信息

```
$ ./rabbitmq-bench producer --declare  -n 100000 e-name q-name 00message00
```



消费者接收消息

```
$ ./rabbitmq-bench consumer --log=debug --declare --stat.interval=1 e-name q-name
```

