package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/mrmt1204/server-side-application/httputil"
	"github.com/mrmt1204/server-side-application/model"
	"github.com/gin-gonic/gin"
)

// Message is controller for requests to messages
type Message struct {
	DB              *sql.DB
	SimpleBotStream chan *model.Message
	GachaBotStream  chan *model.Message
}

// All は全てのメッセージを取得してJSONで返します
func (m *Message) All(c *gin.Context) {
	msgs, err := model.MessagesAll(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if len(msgs) == 0 {
		c.JSON(http.StatusOK, make([]*model.Message, 0))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": msgs,
		"error":  nil,
	})
}

// GetByID はパラメーターで受け取ったidのメッセージを取得してJSONで返します
func (m *Message) GetByID(c *gin.Context) {
	msg, err := model.MessageByID(m.DB, c.Param("id"))

	switch {
	case err == sql.ErrNoRows:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusNotFound, resp)
		return
	case err != nil:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": msg,
		"error":  nil,
	})
}

// Create は新しいメッセージ保存し、作成したメッセージをJSONで返します
func (m *Message) Create(c *gin.Context) {
	var msg model.Message

	if c.Request.ContentLength == 0 {
		resp := httputil.NewErrorResponse(errors.New("body is missing"))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.BindJSON(&msg); err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	inserted, err := msg.Insert(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// bot対応
	m.SimpleBotStream <- inserted
	m.GachaBotStream <- inserted

	c.JSON(http.StatusCreated, gin.H{
		"result": inserted,
		"error":  nil,
	})
}

// UpdateByID は...
func (m *Message) UpdateByID(c *gin.Context) {
	// Mission 1-1. メッセージを編集しよう
	//

	var msg model.Message
	if c.Request.ContentLength == 0 {
		resp := httputil.NewErrorResponse(errors.New("body is missin    g"))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.BindJSON(&msg); err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	oldMsg, err := model.MessageByID(m.DB, c.Param("id"))
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	oldMsg.Body = msg.Body

	newMsg, err := oldMsg.Update(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": newMsg,
		"error":  nil,
	})
}

// DeleteByID は...
func (m *Message) DeleteByID(c *gin.Context) {
	// Mission 1-2. メッセージを削除しよう
	// ...
	msg, err := model.MessageByID(m.DB, c.Param("id"))

	switch {
	case err == sql.ErrNoRows:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusNotFound, resp)
		return
	case err != nil:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	err = msg.Delete(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": nil,
		"error":  nil,
	})
}
