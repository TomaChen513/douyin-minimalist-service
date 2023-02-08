package rabbitmq

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/streadway/amqp"
)

//01和02.这里是rabbitmq最简单的两个模式：simple模式以及work模式。
//simple模式也就是由生产者将消息送到队列里，然后由消费者到队列里取出来消费。
//另外这里的代码work模式也是相同的，也是可以得用的。两个的差别是：work模式在simple模式的基础上多了消费者而已。

// 创建简单模式下的实例，只需要queueName这个参数，其中exchange是默认的，key则不需要。
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	//获取参数connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "连接connection失败")
	//获取channel参数
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel参数失败")
	return rabbitmq

}

// 直接模式,生产者.
func (r *RabbitMQ) PublishSimple(message string) {
	//第一步，申请队列，如不存在，则自动创建之，存在，则路过。
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("创建连接队列失败：%s", err)
	}

	//第二步，发送消息到队列中
	r.channel.Publish(
		r.ExChange,
		r.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 直接模式，消费者
func (r *RabbitMQ) ConsumeSimple() {
	//第一步,申请队列,如果队列不存在则自动创建,存在则跳过
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	//第二步,接收消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",   //用来区分多个消费者
		true, //是否自动应答,告诉我已经消费完了
		false,
		false, //若设置为true,则表示为不能将同一个connection中发送的消息传递给这个connection中的消费者.
		false, //消费队列是否设计阻塞
		nil,
	)
	if err != nil {
		fmt.Printf("消费者接收消息出现问题:%s", err)
	}

	forever := make(chan bool)
	switch r.QueueName {
	case "like_add": //启用协程处理消息
		go func() {
			for d := range msgs {
				params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
				userId, _ := strconv.ParseInt(params[0], 10, 64)
				videoId, _ := strconv.ParseInt(params[1], 10, 64)
				log.Printf("userId:%v\n", userId)
				log.Printf("videoId:%v\n", videoId)
				model.InsertFavourite(model.Like{UserId: userId + 100, VideoId: videoId})
			}
		}()

	case "like_del":
		//启用协程处理消息
		go func() {
			for d := range msgs {
				params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
				userId, _ := strconv.ParseInt(params[0], 10, 64)
				videoId, _ := strconv.ParseInt(params[1], 10, 64)
				log.Printf("userId:%v\n", userId)
				log.Printf("videoId:%v\n", videoId)
				model.DeleteFavourite(model.Like{UserId: userId + 100, VideoId: videoId})
			}
		}()

	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// rabbitmq := rabbitmq.NewRabbitMQSimple("like_del")
var LikeAdd *RabbitMQ
var LikeDel *RabbitMQ

func InitLikeRabbitMQ() {
	LikeAdd = NewRabbitMQSimple("like_add")
	go LikeAdd.ConsumeSimple()

	LikeDel = NewRabbitMQSimple("like_del")
	go LikeDel.ConsumeSimple()
}

// type LikeMQ struct {
// 	RabbitMQ
// 	channel   *amqp.Channel
// 	queueName string
// 	exchange  string
// }

// // NewLikeRabbitMQ 获取likeMQ的对应队列。
// func NewLikeRabbitMQ(queueName string) *LikeMQ {
// 	likeMQ := &LikeMQ{
// 		RabbitMQ:  *Rmq,
// 		queueName: queueName,
// 	}
// 	cha, err := likeMQ.conn.Channel()
// 	likeMQ.channel = cha
// 	Rmq.failOnErr(err, "获取通道失败")
// 	return likeMQ
// }

// // Publish like操作的发布配置。
// func (l *LikeMQ) Publish(message string) {

// 	_, err := l.channel.QueueDeclare(
// 		l.queueName,
// 		//是否持久化
// 		false,
// 		//是否为自动删除
// 		false,
// 		//是否具有排他性
// 		false,
// 		//是否阻塞
// 		false,
// 		//额外属性
// 		nil,
// 	)
// 	if err != nil {
// 		fmt.Printf("创建连接队列失败：%s", err)
// 	}

// 	l.channel.Publish(
// 		l.exchange,
// 		l.queueName,
// 		false,
// 		false,
// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        []byte(message),
// 		})

// }

// // Consumer like关系的消费逻辑。
// func (l *LikeMQ) Consumer() {

// 	_, err := l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	//2、接收消息
// 	messages, err1 := l.channel.Consume(
// 		l.queueName,
// 		//用来区分多个消费者
// 		"",
// 		//是否自动应答
// 		true,
// 		//是否具有排他性
// 		false,
// 		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
// 		false,
// 		//消息队列是否阻塞
// 		false,
// 		nil,
// 	)
// 	if err1 != nil {
// 		fmt.Printf("消费者接收消息出现问题:%s", err)

// 		// panic(err1)
// 	}

// 	forever := make(chan bool)
// 	switch l.queueName {
// 	case "like_add":
// 		//点赞消费队列
// 		go l.consumerLikeAdd(messages)
// 	case "like_del":
// 		//取消赞消费队列
// 		go l.consumerLikeDel(messages)

// 	}

// 	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

// 	<-forever

// }

// // consumerLikeAdd 赞关系添加的消费方式。
// func (l *LikeMQ) consumerLikeAdd(messages <-chan amqp.Delivery) {
// 	for d := range messages {
// 		// 参数解析。
// 		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
// 		userId, _ := strconv.ParseInt(params[0], 10, 64)
// 		videoId, _ := strconv.ParseInt(params[1], 10, 64)
// 		//最多尝试操作数据库的次数
// 		for i := 0; i < 5; i++ {
// 			flag := false //默认无问题
// 			//如果查询没有数据，用来生成该条点赞信息，存储在likeData中
// 			var likeData model.Like
// 			//先查询是否有这条数据
// 			likeInfo, err := model.GetLike(userId, videoId)
// 			//如果有问题，说明查询数据库失败，打印错误信息err:"get likeInfo failed"
// 			if err != nil {
// 				log.Printf(err.Error())
// 				flag = true //出现问题
// 			} else {
// 				if likeInfo == (model.Like{}) { //没查到这条数据，则新建这条数据；
// 					likeData.UserId = userId   //插入userId
// 					likeData.VideoId = videoId //插入videoId
// 					likeData.Cancel = 0        //插入点赞cancel=0
// 					//如果有问题，说明插入数据库失败，打印错误信息err:"insert data fail"
// 					if !model.InsertFavourite(likeData) {
// 						flag = true //出现问题
// 					}
// 				} else { //查到这条数据,更新即可;
// 					//如果有问题，说明插入数据库失败，打印错误信息err:"update data fail"
// 					if err := model.UpdateLike(userId, videoId, 0); err != nil {
// 						log.Printf(err.Error())
// 						flag = true //出现问题
// 					}
// 				}
// 				//一遍流程下来正常执行了，那就打断结束，不再尝试
// 				if flag == false {
// 					break
// 				}
// 			}
// 		}
// 	}
// }

// // consumerLikeDel 赞关系删除的消费方式。
// func (l *LikeMQ) consumerLikeDel(messages <-chan amqp.Delivery) {
// 	for d := range messages {
// 		// 参数解析。
// 		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
// 		userId, _ := strconv.ParseInt(params[0], 10, 64)
// 		videoId, _ := strconv.ParseInt(params[1], 10, 64)
// 		//最多尝试操作数据库的次数
// 		for i := 0; i < 5; i++ {
// 			flag := false //默认无问题
// 			//取消赞行为，只有当前状态是点赞状态才会发起取消赞行为，所以如果查询到，必然是cancel==0(点赞)
// 			//先查询是否有这条数据
// 			likeInfo, err := model.GetLike(userId, videoId)
// 			//如果有问题，说明查询数据库失败，返回错误信息err:"get likeInfo failed"
// 			if err != nil {
// 				log.Printf(err.Error())
// 				flag = true //出现问题
// 			} else {
// 				if likeInfo == (model.Like{}) { //只有当前是点赞状态才能取消点赞这个行为
// 					// 所以如果查询不到数据则返回错误信息:"can't find data,this action invalid"
// 					log.Printf(errors.New("can't find data,this action invalid").Error())
// 				} else {
// 					//如果查询到数据，则更新为取消赞状态
// 					//如果有问题，说明插入数据库失败，打印错误信息err:"update data fail"
// 					if err := model.UpdateLike(userId, videoId, 1); err != nil {
// 						log.Printf(err.Error())
// 						flag = true
// 					}
// 				}
// 			}
// 			//一遍流程下来正常执行了，那就打断结束，不再尝试
// 			if flag == false {
// 				break
// 			}
// 		}
// 	}
// }
