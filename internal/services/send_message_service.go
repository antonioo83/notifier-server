package services

import (
	"fmt"
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/models"
	"github.com/antonioo83/notifier-server/internal/repositories/interfaces"
	interfaces2 "github.com/antonioo83/notifier-server/internal/services/interfaces"
	"github.com/antonioo83/notifier-server/internal/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type sendMessageService struct {
	cfg  config.Config
	rep  interfaces.MessageRepository
	jRep interfaces.JournalRepository
}

func NewMessageSenderService(cfg config.Config, rep interfaces.MessageRepository, jRep interfaces.JournalRepository) interfaces2.SendMessageService {
	return &sendMessageService{cfg, rep, jRep}
}

func (s sendMessageService) Run() {
	wg := &sync.WaitGroup{}
	messages, err := s.rep.FindAll(s.cfg.Sender.MaxAttempts, s.cfg.Sender.ItemCount, 0)
	if err != nil {
		fmt.Printf("can't get messages: " + err.Error())
	}

	for _, message := range *messages {
		log.Println("checking", message)
		wg.Add(1)
		go s.sendMessage(message, wg)
	}

	// в отдельной горутине ждём завершения всех healthCheck
	// после этого закрываем канал errCh — больше записей не будет
	go func() {
		wg.Wait()
		log.Println("wait!")
		time.Sleep(s.cfg.Sender.LoadInterval)
		s.Run()
	}()

	log.Println("successful!")
}

func (s sendMessageService) sendMessage(message models.Message, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	isSent := false
	status, content, err := s.sendRequest(message)
	if err != nil {
		log.Printf("can't send message to the resource: %s", err.Error())
	}
	if status == message.SuccessHttpStatus {
		isSent = true
	}

	if isSent {
		err = s.rep.MarkSent(message.Code)
		if err != nil {
			log.Println("can't mark a message as sent!")
		}
		log.Println("the message was sent!")
	} else {
		err = s.rep.MarkUnSent(message.Code, message.AttemptCount+1)
		if err != nil {
			log.Println("can't mark a message as not sent!")
		}
		log.Println("the message wasn't sent")
	}

	if len(content) > 300 {
		content = content[:300]
	}

	journal := models.Journal{
		UserId:          message.User.ID,
		ResourceId:      message.ResourceId,
		MessageId:       message.ID,
		ResponseStatus:  status,
		ResponseContent: content,
	}
	err = s.jRep.Save(journal)
	if err != nil {
		log.Println("can't write a record to the journal!")
	}

	return
}

func (s sendMessageService) sendRequest(message models.Message) (status int, content string, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: s.cfg.Sender.Timeout,
	}

	req, err := http.NewRequest(strings.ToUpper(message.Command), message.Resource.URL, strings.NewReader(message.Content))
	if err != nil {
		return 0, "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", message.User.AuthToken)
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	defer utils.ResourceClose(resp.Body)

	return resp.StatusCode, string(respBody), nil
}
