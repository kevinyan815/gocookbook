package tts

import (
	"bytes"
	"xxx/oss"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
)


var (
	accessKeyId     = "xxxxx"
	accessKeySecret = "xxx"
	appKey          = "xxxxxxx"

	OSSEndpoint = "xxxx"
	OSSAccessKey = "xxx"
	OSSKeySecrete = "xxx"

)

var clientConfig *nls.ConnectionConfig

func getClient() *nls.ConnectionConfig {
	return clientConfig
}

type userParam struct {
	audioStream bytes.Buffer //转换成功回调时把音频流保存在这里
}

func init() {

	var ak, as string
	if os.Getenv("env") == "prod" {
		ak = accessKeyIdProd
		as = accessKeySecretProd
	} else {
		ak = accessKeyId
		as = accessKeySecret
	}
	cf, err := nls.NewConnectionConfigWithAKInfoDefault(nls.DEFAULT_URL, appKey, ak, as)
	if err != nil {
		panic(err)
	}
	clientConfig = cf
}

// 转换任务错误执行
func onTaskFailed(text string, param interface{}) {
	log.Print("TtsTextToVoiceTaskFailedError", text)
}

// 转换任务成功回调音频流
func onSynthesisResult(data []byte, param interface{}) {
	p, ok := param.(*userParam)
	if !ok {
		log.Print("TtsTextToVoiceResultUserParamUnExpected", param)
		return
	}
	p.audioStream.Write(data)
}

// 执行完毕
func onCompleted(text string, param interface{}) {
	log.Print("TtsTextToVoiceOnCompleted:", text)
}

func onClose(param interface{}) {
	log.Print("TtsTextToVoiceOnClosed")
}

func waitReady(ch chan bool, logger *nls.NlsLogger) error {
	select {
	case done := <-ch:
		{
			if !done {
				logger.Println("Wait failed")
				return errors.New("wait failed")
			}
			logger.Println("Wait done")
		}
	case <-time.After(60 * time.Second):
		{
			logger.Println("Wait timeout")
			return errors.New("wait timeout")
		}
	}
	return nil
}

func TextToVoice(text string) (audioFileUrl string, err error) {
	clientConf := getClient()
	param := nls.DefaultSpeechSynthesisParam()
	param.Voice = "yuer"  // 儿童剧女声
	strId := "ai-" // tts日志的前缀
	logger := nls.NewNlsLogger(os.Stderr, strId, log.LstdFlags|log.Lmicroseconds)
	logger.SetLogSil(false)
	logger.SetDebug(true)
	userParam := new(userParam)
	//第三个参数控制是否请求长文本语音合成，false为短文本语音合成
	tts, err := nls.NewSpeechSynthesis(clientConf, logger, false,
		onTaskFailed, onSynthesisResult, nil,
		onCompleted, onClose, userParam)
	if err != nil {
		logger.Fatalln(err)
		return
	}
	ch, err := tts.Start(text, param, nil)
	if err != nil {
		tts.Shutdown()
		dlog.Error("TtsTextToVoiceStartError", err)
		return
	}
	err = waitReady(ch, logger)
	if err != nil {
		tts.Shutdown()
		dlog.Error("TtsTextToVoiceWaitError", err)
		return
	}
	tts.Shutdown() //tts关闭后再把文件流上传到OSS
	bucket, err := oss.GetBucket("qschou", OSSEndpoint, OSSAccessKey, OSSKeySecrete)
	if err != nil {
		log.Print("func", "SaveImgToOSS", "msg", "GetBucket Error", "err", err)
		return
	}
	uniqId := Md5(fmt.Sprintf("%d", time.Now().UnixNano()))
	fileName := fmt.Sprintf("files/ai-toolbox/text-to-audio/%s.%s", uniqId, "wav")
	err = bucket.PutObject(fileName, bytes.NewReader(userParam.audioStream.Bytes()))
	if err != nil {
		log.Print("func", "SaveImgToOSS", "msg", "PutObject Error", "err", err)
		return "", err
	}

	return "https://thumb.qschou.com" + "/" + fileName, nil
	//本地测试用
	//fname := "ttsdump.wav"
	//fout, err := os.OpenFile(fname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	//fout.Write(userParam.audioStream.Bytes())
	//return "", nil
}


func Md5(str string) string {
	tmp := md5.Sum([]byte(str))
	return hex.EncodeToString(tmp[:])
}
